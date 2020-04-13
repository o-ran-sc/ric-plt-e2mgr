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
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
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
	"strconv"
	"strings"
)

type E2SetupRequestNotificationHandler struct {
	logger                 *logger.Logger
	config                 *configuration.Configuration
	e2tInstancesManager    managers.IE2TInstancesManager
	rmrSender              *rmrsender.RmrSender
	rNibDataService       services.RNibDataService
	e2tAssociationManager *managers.E2TAssociationManager
}

func NewE2SetupRequestNotificationHandler(logger *logger.Logger, config *configuration.Configuration, e2tInstancesManager managers.IE2TInstancesManager, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService, e2tAssociationManager *managers.E2TAssociationManager) E2SetupRequestNotificationHandler {
	return E2SetupRequestNotificationHandler{
		logger:                 logger,
		config:                 config,
		e2tInstancesManager:    e2tInstancesManager,
		rmrSender: rmrSender,
		rNibDataService: rNibDataService,
		e2tAssociationManager: e2tAssociationManager,
	}
}

func (h E2SetupRequestNotificationHandler) Handle(request *models.NotificationRequest){
	ranName := request.RanName
	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - received E2_SETUP_REQUEST. Payload: %x", ranName, request.Payload)

	setupRequest, e2tIpAddress, err := h.parseSetupRequest(request.Payload)
	if err != nil {
		h.logger.Errorf(err.Error())
		return
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - E2T Address: %s - handling E2_SETUP_REQUEST", e2tIpAddress)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.Handle - E2_SETUP_REQUEST has been parsed successfully %+v", setupRequest)

	_, err = h.e2tInstancesManager.GetE2TInstance(e2tIpAddress)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
		return
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil{
		if _, ok := err.(*common.ResourceNotFoundError); ok{
			nbIdentity := h.buildNbIdentity(ranName, setupRequest)
			nodebInfo, err = h.buildNodebInfo(ranName, e2tIpAddress, setupRequest)
			if err != nil{
				h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to build nodebInfo entity. Error: %s", ranName, err)
				return
			}
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
		nodebInfo.GetGnb().RanFunctions, err = setupRequest.GetExtractRanFunctionsList()
		if err != nil{
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to update nodebInfo entity. Error: %s", ranName, err)
			return
		}
	}
	err = h.e2tAssociationManager.AssociateRan(e2tIpAddress, nodebInfo)
	if err != nil{
		if _, ok := err.(*e2managererrors.RoutingManagerError); ok{
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to associate E2T to nodeB entity. Error: %s", ranName, err)
			h.handleUnsuccessfulResponse(nodebInfo, request)
		}
		return
	}
	h.handleSuccessfulResponse(ranName, request, setupRequest)
}

func (h E2SetupRequestNotificationHandler) handleUnsuccessfulResponse(nodebInfo *entities.NodebInfo, req *models.NotificationRequest){
	failureResponse := models.NewE2SetupFailureResponseMessage()
	h.logger.Debugf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - E2_SETUP_RESPONSE has been built successfully %+v", failureResponse)

	responsePayload, err := xml.Marshal(&failureResponse.E2APPDU)
	if err != nil{
		h.logger.Warnf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_SETUP_RESP. Payload: %s", nodebInfo.RanName, responsePayload)
	}

	responsePayload = replaceEmptyTagsWithSelfClosing(responsePayload)

	msg := models.NewRmrMessage(rmrCgo.RIC_E2_SETUP_FAILURE, nodebInfo.RanName, responsePayload, req.TransactionId, req.GetMsgSrc())
	h.logger.Infof("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - RAN name: %s - RIC_E2_SETUP_RESP message has been built successfully. Message: %x", nodebInfo.RanName, msg)
	_ = h.rmrSender.WhSend(msg)

}

func (h E2SetupRequestNotificationHandler) handleSuccessfulResponse(ranName string, req *models.NotificationRequest, setupRequest *models.E2SetupRequestMessage){

	ricNearRtId, err := convertTo20BitString(h.config.GlobalRicId.RicNearRtId)
	if err != nil{
		h.logger.Errorf("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - failed to convert RicNearRtId value %s to 20 bit string . Error: %s", ranName, h.config.GlobalRicId.RicNearRtId, err)
		return
	}
	successResponse := models.NewE2SetupSuccessResponseMessage(h.config.GlobalRicId.PlmnId, ricNearRtId,setupRequest)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - E2_SETUP_RESPONSE has been built successfully %+v", successResponse)

	responsePayload, err := xml.Marshal(&successResponse.E2APPDU)
	if err != nil{
		h.logger.Warnf("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_SETUP_RESP. Payload: %s", ranName, responsePayload)
	}

	responsePayload = replaceEmptyTagsWithSelfClosing(responsePayload)

	msg := models.NewRmrMessage(rmrCgo.RIC_E2_SETUP_RESP, ranName, responsePayload, req.TransactionId, req.GetMsgSrc())
	h.logger.Infof("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - RIC_E2_SETUP_RESP message has been built successfully. Message: %x", ranName, msg)
	_ = h.rmrSender.Send(msg)
}


func replaceEmptyTagsWithSelfClosing(responsePayload []byte) []byte {
	responseString := strings.NewReplacer(
		"<reject></reject>", "<reject/>",
		"<ignore></ignore>", "<ignore/>",
		"<transport-resource-unavailable></transport-resource-unavailable>", "<transport-resource-unavailable/>",
		).Replace(string(responsePayload))
	return []byte(responseString)
}

func convertTo20BitString(ricNearRtId string) (string, error){
	r, err := strconv.ParseUint(ricNearRtId, 16, 32)
	if err != nil{
		return "", err
	}
	return fmt.Sprintf("%020b", r)[:20], nil
}

func (h E2SetupRequestNotificationHandler) parseSetupRequest(payload []byte)(*models.E2SetupRequestMessage, string, error){

	pipInd := bytes.IndexByte(payload, '|')
	if pipInd < 0 {
		return nil, "", errors.New( "#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
	}

	e2tIpAddress := string(payload[:pipInd])
	if len(e2tIpAddress) == 0 {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Empty E2T Address received")
	}
	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(payload[pipInd + 1:], &setupRequest.E2APPDU)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: %x", payload))
	}

	return setupRequest, e2tIpAddress, nil
}

func (h E2SetupRequestNotificationHandler) buildNodebInfo(ranName string, e2tAddress string, request *models.E2SetupRequestMessage) (*entities.NodebInfo, error){
	var err error
	nodebInfo := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tAddress,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		RanName: ranName,
		NodeType: entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
		GlobalNbId:  &entities.GlobalNbId{
			PlmnId: request.GetPlmnId(),
			NbId:   request.GetNbId(),
		},
	}
	nodebInfo.GetGnb().RanFunctions, err = request.GetExtractRanFunctionsList()
	return nodebInfo, err
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