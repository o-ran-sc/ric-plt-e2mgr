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
	"fmt"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

//type FunctionChange int

const (
	RAN_FUNCTIONS_ADDED int = 10 + iota
	RAN_FUNCTIONS_DELETED
	RAN_FUNCTIONS_MODIFIED
)

type functionDetails struct {
	functionChange     int
	functionId         uint32
	functionDefinition string
	functionRevision   uint32
	functionOID        string
}

type RicServiceUpdateHandler struct {
	logger                  *logger.Logger
	rmrSender               *rmrsender.RmrSender
	rNibDataService         services.RNibDataService
	ranListManager          managers.RanListManager
	RicServiceUpdateManager managers.IRicServiceUpdateManager
}

func NewRicServiceUpdateHandler(logger *logger.Logger, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService, ranListManager managers.RanListManager, RicServiceUpdateManager managers.IRicServiceUpdateManager) *RicServiceUpdateHandler {
	return &RicServiceUpdateHandler{
		logger:                  logger,
		rmrSender:               rmrSender,
		rNibDataService:         rNibDataService,
		ranListManager:          ranListManager,
		RicServiceUpdateManager: RicServiceUpdateManager,
	}
}

func (h *RicServiceUpdateHandler) Handle(request *models.NotificationRequest) {
	ranName := request.RanName
	h.logger.Infof("#RicServiceUpdateHandler.Handle - RAN name: %s - received RIC_SERVICE_UPDATE. Payload: %s", ranName, request.Payload)

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			h.logger.Errorf("#RicServiceUpdateHandler.Handle - failed to get nodeB entity for ran name: %v due to RNIB Error: %s", ranName, err)
		} else {
			h.logger.Errorf("#RicServiceUpdateHandler.Handle - nobeB entity of RanName:%s absent in RNIB. Error: %s", ranName, err)
		}
		return
	}

	ricServiceUpdate, err := h.parseSetupRequest(request.Payload)
	if err != nil {
		h.logger.Errorf(err.Error())
		return
	}
	h.logger.Infof("#RicServiceUpdateHandler.Handle - RIC_SERVICE_UPDATE has been parsed successfully %+v", ricServiceUpdate)
	h.RicServiceUpdateManager.StoreExistingRanFunctions(ranName)
	h.logger.Infof("#RicServiceUpdate.Handle - Getting the ranFunctions before we do the RIC ServiceUpdate handling")

	if len(ricServiceUpdate.E2APPDU.InitiatingMessage.Value.RICServiceUpdate.ProtocolIEs.RICServiceUpdateIEs) == 0 {
		h.logger.Errorf("#RicServiceUpdateHandler.Handle - RAN name: %s - RICServiceUpdateIEs empty", ranName)
		return
	}

	ackFunctionIds := h.updateFunctions(ricServiceUpdate.E2APPDU.InitiatingMessage.Value.RICServiceUpdate.ProtocolIEs.RICServiceUpdateIEs, nodebInfo)
	if len(ricServiceUpdate.E2APPDU.InitiatingMessage.Value.RICServiceUpdate.ProtocolIEs.RICServiceUpdateIEs) > 1 {
		err = h.rNibDataService.UpdateNodebInfoAndPublish(nodebInfo)
		if err != nil {
			h.logger.Errorf("#RicServiceUpdateHandler.Handle - RAN name: %s - Failed at UpdateNodebInfoAndPublish. error: %s", nodebInfo.RanName, err)
			return
		}
	}

	oldNbIdentity, newNbIdentity := h.ranListManager.UpdateHealthcheckTimeStampReceived(nodebInfo.RanName)
	err = h.ranListManager.UpdateNbIdentities(nodebInfo.NodeType, []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity})
	if err != nil {
		h.logger.Errorf("#RicServiceUpdate.Handle - failed to Update NbIdentities: %s", err)
		return
	}

	updateAck := models.NewServiceUpdateAck(ackFunctionIds, ricServiceUpdate.E2APPDU.InitiatingMessage.Value.RICServiceUpdate.ProtocolIEs.RICServiceUpdateIEs[0].Value.TransactionID)
	err = h.sendUpdateAck(updateAck, nodebInfo, request)
	if err != nil {
		h.logger.Errorf("#RicServiceUpdate.Handle - failed to send RIC_SERVICE_UPDATE_ACK message to RMR: %s", err)
		return
	}

	h.logger.Infof("#RicServiceUpdate.Handle - Completed successfully")
	models.UpdateProcedureType(ranName, models.RicServiceUpdateCompleted)
	h.logger.Debugf("#RicServiceUpdateHandler.Handle  - updating the enum value to RicServiceUpdateCompleted completed")
}

func (h *RicServiceUpdateHandler) sendUpdateAck(updateAck models.RicServiceUpdateAckE2APPDU, nodebInfo *entities.NodebInfo, request *models.NotificationRequest) error {
	payLoad, err := xml.Marshal(updateAck)
	if err != nil {
		h.logger.Errorf("#RicServiceUpdate.sendUpdateAck - RAN name: %s - Error marshalling RIC_SERVICE_UPDATE_ACK. Payload: %s", nodebInfo.RanName, payLoad)
	}

	toReplaceTags := []string{"reject", "ignore", "procedureCode", "id", "RANfunctionID-Item", "RANfunctionsID-List"}
	payLoad = utils.ReplaceEmptyTagsWithSelfClosing(payLoad, toReplaceTags)

	h.logger.Infof("#RicServiceUpdate.sendUpdateAck - Sending RIC_SERVICE_UPDATE_ACK to RAN name: %s with payload %s", nodebInfo.RanName, payLoad)
	msg := models.NewRmrMessage(rmrCgo.RIC_SERVICE_UPDATE_ACK, nodebInfo.RanName, payLoad, request.TransactionId, request.GetMsgSrc())
	err = h.rmrSender.Send(msg)
	return err
}

func (h *RicServiceUpdateHandler) updateFunctions(RICServiceUpdateIEs []models.RICServiceUpdateIEs, nodebInfo *entities.NodebInfo) []models.RicServiceAckRANFunctionIDItem {
	ranFunctions := nodebInfo.GetGnb().RanFunctions
	RanFIdtoIdxMap := make(map[uint32]int)
	var acceptedFunctionIds []models.RicServiceAckRANFunctionIDItem
	functionsToBeDeleted := make(map[int]bool)

	for index, ranFunction := range ranFunctions {
		RanFIdtoIdxMap[ranFunction.RanFunctionId] = index
	}

	for _, ricServiceUpdateIE := range RICServiceUpdateIEs {
		functionDetails, err := h.getFunctionDetails(ricServiceUpdateIE)
		if err != nil {
			h.logger.Errorf("#RicServiceUpdate.updateFunctions- GetFunctionDetails returned err: %s", err)
		}

		for _, functionDetail := range functionDetails {
			functionChange, functionId, functionDefinition, functionRevision, functionOID := functionDetail.functionChange,
				functionDetail.functionId, functionDetail.functionDefinition, functionDetail.functionRevision, functionDetail.functionOID
			ranFIndex, ok := RanFIdtoIdxMap[functionId]
			if !ok {
				switch functionChange {
				case RAN_FUNCTIONS_ADDED, RAN_FUNCTIONS_MODIFIED:
					ranFunctions = append(ranFunctions, &entities.RanFunction{RanFunctionId: functionId,
						RanFunctionDefinition: functionDefinition, RanFunctionRevision: functionRevision, RanFunctionOid: functionOID})
				case RAN_FUNCTIONS_DELETED:
					//Do nothing
				}
			} else {
				switch functionChange {
				case RAN_FUNCTIONS_ADDED, RAN_FUNCTIONS_MODIFIED:
					ranFunctions[ranFIndex].RanFunctionDefinition = functionDefinition
					ranFunctions[ranFIndex].RanFunctionRevision = functionRevision
				case RAN_FUNCTIONS_DELETED:
					functionsToBeDeleted[ranFIndex] = true
				}
			}
			serviceupdateAckFunctionId := models.RicServiceAckRANFunctionIDItem{RanFunctionID: functionId, RanFunctionRevision: functionRevision}
			acceptedFunctionIds = append(acceptedFunctionIds, serviceupdateAckFunctionId)
		}
	}
	finalranFunctions := h.remove(ranFunctions, functionsToBeDeleted)
	nodebInfo.GetGnb().RanFunctions = finalranFunctions
	return acceptedFunctionIds
}

func (h *RicServiceUpdateHandler) remove(ranFunctions []*entities.RanFunction, functionsToBeDeleted map[int]bool) []*entities.RanFunction {
	if len(functionsToBeDeleted) == 0 {
		return ranFunctions
	}
	var finalranFunctions []*entities.RanFunction
	for i := 0; i < len(ranFunctions); i++ {
		_, ok := functionsToBeDeleted[i]
		if !ok {
			finalranFunctions = append(finalranFunctions, ranFunctions[i])
		}
	}
	return finalranFunctions
}

func (h *RicServiceUpdateHandler) getFunctionDetails(ricServiceUpdateIE models.RICServiceUpdateIEs) ([]functionDetails, error) {
	functionChange := ricServiceUpdateIE.ID
	switch functionChange {
	case RAN_FUNCTIONS_ADDED, RAN_FUNCTIONS_MODIFIED:
		return h.getFunctionsAddedModifiedHandler(ricServiceUpdateIE)
	case RAN_FUNCTIONS_DELETED:
		return h.getFunctionsDeleteHandler(ricServiceUpdateIE)
	default:
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - Unknown change type %v", functionChange))
	}
	return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - Internal Error"))
}

func (h *RicServiceUpdateHandler) getFunctionsAddedModifiedHandler(ricServiceUpdateIE models.RICServiceUpdateIEs) ([]functionDetails, error) {
	functionChange := ricServiceUpdateIE.ID
	ranFunctionsIEList := ricServiceUpdateIE.Value.RANfunctionsList.RANfunctionsItemProtocolIESingleContainer
	const MaxRanFunctions = 256
	if len(ranFunctionsIEList) == 0 {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - function change type is %v but Functions list is empty", functionChange))
	}
	if len(ranFunctionsIEList) > MaxRanFunctions {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - Functions list size %d exceeds maximum %d", len(ranFunctionsIEList), MaxRanFunctions))
	}

	functionDetailsList := make([]functionDetails, len(ranFunctionsIEList))
	for index, ranFunctionIE := range ranFunctionsIEList {
		ranFunction := ranFunctionIE.Value.RANfunctionItem
		functionDetailsList[index] = functionDetails{functionChange: functionChange, functionId: ranFunction.RanFunctionID,
			functionDefinition: ranFunction.RanFunctionDefinition, functionRevision: ranFunction.RanFunctionRevision, functionOID: ranFunction.RanFunctionOID}
	}
	return functionDetailsList, nil
}

func (h *RicServiceUpdateHandler) getFunctionsDeleteHandler(ricServiceUpdateIE models.RICServiceUpdateIEs) ([]functionDetails, error) {
	functionChange := ricServiceUpdateIE.ID
	ranFunctionIdIEsList := ricServiceUpdateIE.Value.RANfunctionsIDList.RANfunctionsItemIDProtocolIESingleContainer
	const MaxRanFunctions = 256
	if len(ranFunctionIdIEsList) == 0 {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - function change type is %v but FunctionIds list is empty", functionChange))
	}
	if len(ranFunctionIdIEsList) > MaxRanFunctions {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdate.getFunctionDetails - FunctionIds list size %d exceeds maximum %d", len(ranFunctionIdIEsList), MaxRanFunctions))
	}

	functionDetailsList := make([]functionDetails, len(ranFunctionIdIEsList))
	for index, ranFunctionIdIE := range ranFunctionIdIEsList {
		ranFunctionId := ranFunctionIdIE.Value.RANfunctionIDItem
		functionDetailsList[index] = functionDetails{functionChange: functionChange, functionId: ranFunctionId.RanFunctionID,
			functionDefinition: "", functionRevision: ranFunctionId.RanFunctionRevision}
	}
	return functionDetailsList, nil
}

func (h *RicServiceUpdateHandler) parseSetupRequest(payload []byte) (*models.RICServiceUpdateMessage, error) {
	pipInd := bytes.IndexByte(payload, '|')
	if pipInd < 0 {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdateHandler.parseSetupRequest - Error parsing RIC SERVICE UPDATE failed extract Payload: no | separator found"))
	}

	ricServiceUpdate := &models.RICServiceUpdateMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload[pipInd+1:]), &ricServiceUpdate.E2APPDU)
	if err != nil {
		return nil, common.NewInternalError(fmt.Errorf("#RicServiceUpdateHandler.parseSetupRequest - Error unmarshalling RIC SERVICE UPDATE payload: %x", payload))
	}
	return ricServiceUpdate, nil
}
