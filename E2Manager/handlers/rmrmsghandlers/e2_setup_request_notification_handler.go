//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type E2SetupRequestNotificationHandler struct {
	logger                 *logger.Logger
	e2tInstancesManager    managers.IE2TInstancesManager
	rmrSender              *rmrsender.RmrSender
	rNibDataService       services.RNibDataService
	e2tAssociationManager *managers.E2TAssociationManager
}

func NewE2SetupRequestNotificationHandler(logger *logger.Logger, e2tInstancesManager managers.IE2TInstancesManager, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService, e2tAssociationManager *managers.E2TAssociationManager) E2SetupRequestNotificationHandler {
	return E2SetupRequestNotificationHandler{
		logger:                 logger,
		e2tInstancesManager:    e2tInstancesManager,
		rmrSender: rmrSender,
		rNibDataService: rNibDataService,
		e2tAssociationManager: e2tAssociationManager,
	}
}

func (h E2SetupRequestNotificationHandler) Handle(request *models.NotificationRequest){
	ranName := request.RanName
	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - received E2 Setup Request. Payload: %x", ranName, request.Payload)

	setupRequest, e2tIpAddress, err := h.parseSetupRequest(request.Payload)
	if err != nil {
		h.logger.Errorf(err.Error())
		return
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - E2T Address: %s - handling E2_SETUP_REQUEST", e2tIpAddress)

	_, err = h.e2tInstancesManager.GetE2TInstance(e2tIpAddress)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
		return
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil{
		if _, ok := err.(*common.ResourceNotFoundError); ok{
			nbIdentity := h.buildNbIdentity(ranName, setupRequest)
			nodebInfo = h.buildNodebInfo(ranName, e2tIpAddress, setupRequest)
			err = h.rNibDataService.SaveNodeb(nbIdentity, nodebInfo)
			if err != nil{
				h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to save nodebInfo entity. Error: %s", ranName, err)
				return
			}
		} else{
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to retrieve nodebInfo entity. Error: %s", ranName, err)
			return
		}

	} else {
		if nodebInfo.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in incorrect state", nodebInfo.RanName, nodebInfo.ConnectionStatus)
			h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling setup request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
			return
		}
		h.updateNodeBFunctions(nodebInfo, setupRequest)
	}
	err = h.e2tAssociationManager.AssociateRan(e2tIpAddress, nodebInfo)
	if err != nil{
		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to associate E2T to nodeB entity. Error: %s", ranName, err)
		return
	}
	successResponse := &models.E2SetupSuccessResponseMessage{}
	successResponse.SetPlmnId(setupRequest.GetPlmnId())
	successResponse.SetNbId("&" + fmt.Sprintf("%020b", 0xf0))
	responsePayload, err := xml.Marshal(successResponse)
	if err != nil{
		h.logger.Warnf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - Error marshalling E2 Setup Response. Response: %x", ranName, responsePayload)
	}
	msg := models.NewRmrMessage(rmrCgo.RIC_E2_SETUP_RESP, ranName, responsePayload, request.TransactionId)
	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - E2 Setup Request has been built. Message: %x", ranName, msg)
	//TODO err = h.rmrSender.Send(msg)

}

func (h E2SetupRequestNotificationHandler) parseSetupRequest(payload []byte)(*models.E2SetupRequestMessage, string, error){

	colonInd := bytes.IndexByte(payload, ':')
	if colonInd < 0 {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request, failed extract E2T IP Address: no ':' separator found")
	}

	e2tIpAddress := string(payload[:colonInd])
	if len(e2tIpAddress) == 0 {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Empty E2T Address received")
	}

	pipInd := bytes.IndexByte(payload, '|')
	if pipInd < 0 {
		return nil, "", errors.New( "#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
	}

	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(payload[pipInd + 1:], &setupRequest)
	if err != nil {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: %s")
	}

	return setupRequest, e2tIpAddress, nil
}

func (h E2SetupRequestNotificationHandler) updateNodeBFunctions(nodeB *entities.NodebInfo, request *models.E2SetupRequestMessage){
	//TODO the function should be implemented in the scope of the US 192 "Save the entire Setup request in RNIB"
}

func (h E2SetupRequestNotificationHandler) buildNodebInfo(ranName string, e2tAddress string, request *models.E2SetupRequestMessage) *entities.NodebInfo{
	nodebInfo := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tAddress,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		RanName: ranName,
		NodeType: entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}
	h.updateNodeBFunctions(nodebInfo, request)
	return nodebInfo
}

func (h E2SetupRequestNotificationHandler) buildNbIdentity(ranName string, setupRequest *models.E2SetupRequestMessage)*entities.NbIdentity{
	return &entities.NbIdentity{
		InventoryName:ranName,
		GlobalNbId: &entities.GlobalNbId{
			PlmnId: setupRequest.GetPlmnId(),
			NbId:   setupRequest.GetNbId(),
		},
	}
}