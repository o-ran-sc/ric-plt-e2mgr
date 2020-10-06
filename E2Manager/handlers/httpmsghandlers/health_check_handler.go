//
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package httpmsghandlers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"encoding/xml"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

var(
	healthCheckSuccessResponse          = "Request Accepted"
	healthCheckEmptyTagsToReplaceToSelfClosingTags = []string{"reject", "ignore", "protocolIEs", "procedureCode"}
)

type HealthCheckRequestHandler struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
	ranListManager  managers.RanListManager
	rmrsender       *rmrsender.RmrSender
}

func NewHealthCheckRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService, ranListManager managers.RanListManager, rmrsender *rmrsender.RmrSender) *HealthCheckRequestHandler {
	return &HealthCheckRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
		ranListManager:  ranListManager,
		rmrsender:       rmrsender,
	}
}

func (h *HealthCheckRequestHandler) Handle(request models.Request) (models.IResponse, error) {
	ranNameList := h.getRanNameList(request)
	isAtleastOneRanConnected := false

	nodetypeToNbIdentityMapOld := make(map[entities.Node_Type][]*entities.NbIdentity)
	nodetypeToNbIdentityMapNew := make(map[entities.Node_Type][]*entities.NbIdentity)

	for _, ranName := range ranNameList {
		nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
		if err != nil {
			_, ok := err.(*common.ResourceNotFoundError)
			if !ok {
				h.logger.Errorf("#HealthCheckRequest.Handle - failed to get nodeBInfo entity for ran name: %v from RNIB. Error: %s", ranName, err)
				return nil, e2managererrors.NewRnibDbError()
			}
			continue
		}

		if nodebInfo.ConnectionStatus == entities.ConnectionStatus_CONNECTED {
			isAtleastOneRanConnected = true

			err := h.sendRICServiceQuery(nodebInfo)
			if err != nil {
				return nil,err
			}

			oldnbIdentity, newnbIdentity := h.ranListManager.UpdateHealthcheckTimeStampSent(ranName)
			nodetypeToNbIdentityMapOld[nodebInfo.NodeType] = append(nodetypeToNbIdentityMapOld[nodebInfo.NodeType], oldnbIdentity)
			nodetypeToNbIdentityMapNew[nodebInfo.NodeType] = append(nodetypeToNbIdentityMapNew[nodebInfo.NodeType], newnbIdentity)
		}
	}

	for k, _ := range nodetypeToNbIdentityMapOld {
		err := h.ranListManager.UpdateNbIdentities(k, nodetypeToNbIdentityMapOld[k], nodetypeToNbIdentityMapNew[k])
		if err != nil {
			return nil,err
		}
	}

	if isAtleastOneRanConnected == false {
		return nil, e2managererrors.NewNoConnectedRanError()
	}

	h.logger.Infof("#HealthcheckRequest.Handle - HealthcheckTimeStampSent Update completed to RedisDB")

	return models.NewHealthCheckSuccessResponse(healthCheckSuccessResponse), nil
}

func (h *HealthCheckRequestHandler) sendRICServiceQuery(nodebInfo *entities.NodebInfo) error {

	serviceQuery := models.NewRicServiceQueryMessage(nodebInfo.GetGnb().RanFunctions)
	payLoad, err := xml.Marshal(serviceQuery.E2APPDU)
	if err != nil {
		h.logger.Errorf("#HealthCHeckRequest.Handle- RAN name: %s - Error marshalling RIC_SERVICE_QUERY. Payload: %s", nodebInfo.RanName, payLoad)
		//return nil, e2managererrors.NewInternalError()
	}

	payLoad = utils.ReplaceEmptyTagsWithSelfClosing(payLoad,healthCheckEmptyTagsToReplaceToSelfClosingTags)

	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := models.NewRmrMessage(rmrCgo.RIC_SERVICE_QUERY, nodebInfo.RanName, payLoad, xAction, msgSrc)

	err = h.rmrsender.Send(msg)

	if err != nil {
		h.logger.Errorf("#HealthCHeckRequest.Handle - failed to send RIC_SERVICE_QUERY message to RMR for %s. Error: %s", nodebInfo.RanName, err)
		//return nil, e2managererrors.NewRmrError()
	} else {
		h.logger.Infof("#HealthCHeckRequest.Handle - RAN name : %s - Successfully built and sent RIC_SERVICE_QUERY. Message: %x", nodebInfo.RanName, msg)
	}

	return nil
}

func (h *HealthCheckRequestHandler) getRanNameList(request models.Request) []string {
	healthCheckRequest := request.(models.HealthCheckRequest)
	if request != nil && len(healthCheckRequest.RanList) != 0 {
		return healthCheckRequest.RanList
	}

	h.logger.Infof("#HealthcheckRequest.getRanNameList - Empty request sent, fetching all connected NbIdentitylist")

	nodeIds := h.ranListManager.GetNbIdentityList()
	var ranNameList []string

	for _, nbIdentity := range nodeIds {
		if nbIdentity.ConnectionStatus == entities.ConnectionStatus_CONNECTED {
			ranNameList = append(ranNameList, nbIdentity.InventoryName)
		}
	}

	return ranNameList
}
