//
// Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.
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
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"encoding/xml"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

const E2ResetRequestLogInfoElapsedTime = "#E2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms"

var (
	resetRequestEmptyTagsToReplaceToSelfClosingTags = []string{"reject", "ignore", "protocolIEs", "procedureCode", "ResetResponse", "ResetResponseIEs", "id", "criticality", "TransactionID"}
)

type E2ResetRequestNotificationHandler struct {
	logger                            *logger.Logger
	rnibDataService                   services.RNibDataService
	config                            *configuration.Configuration
	rmrSender                         *rmrsender.RmrSender
	ranResetManager                   *managers.RanResetManager
	changeStatusToConnectedRanManager *managers.ChangeStatusToConnectedRanManager
}

func NewE2ResetRequestNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, config *configuration.Configuration, rmrSender *rmrsender.RmrSender, ranResetManager *managers.RanResetManager, changeStatusToConnectedRanManager *managers.ChangeStatusToConnectedRanManager) *E2ResetRequestNotificationHandler {
	return &E2ResetRequestNotificationHandler{
		logger:                            logger,
		rnibDataService:                   rnibDataService,
		config:                            config,
		rmrSender:                         rmrSender,
		ranResetManager:                   ranResetManager,
		changeStatusToConnectedRanManager: changeStatusToConnectedRanManager,
	}
}

func (e *E2ResetRequestNotificationHandler) Handle(request *models.NotificationRequest) {

	e.logger.Infof("#E2ResetRequestNotificationHandler.Handle - RAN name: %s - received E2_Reset. Payload: %x", request.RanName, request.Payload)

	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - RIC_E2_Node_Reset parsed successfully ")

	nodebInfo, err := e.getNodebInfo(request.RanName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to retrieve nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
		e.logger.Infof(E2ResetRequestLogInfoElapsedTime, utils.ElapsedTime(request.StartTime))
		return
	}

	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - nodeB entity retrieved. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

	nodebInfo.ConnectionStatus = entities.ConnectionStatus_UNDER_RESET

	ranName := request.RanName
	isResetDone, err := e.ranResetManager.ResetRan(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to update and notify connection status of nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
	} else {
		if isResetDone {
			nodebInfoupdated, err1 := e.getNodebInfo(request.RanName)
			if err1 != nil {
				e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to get updated nodeB entity. RanName: %s. Error: %s", request.RanName, err1.Error())
			}
			e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - Reset Done Successfully ran: %s , Connection status updated : %s", ranName, nodebInfoupdated.ConnectionStatus)
		} else {
			e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - Reset Failed")
		}
	}

	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to update connection status of nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
	}

	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - nodeB entity under reset state. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

	e.logger.Infof(E2ResetRequestLogInfoElapsedTime, utils.ElapsedTime(request.StartTime))

	e.waitfortimertimeout(request)

	resetRequest, err := e.parseE2ResetMessage(request.Payload)
	if err != nil {
		e.logger.Errorf(err.Error())
		return
	}
	e.logger.Infof("#E2ResetRequestNotificationHandler.Handle - RIC_RESET_REQUEST has been parsed successfully %+v", resetRequest)
	e.handleSuccessfulResponse(ranName, request, resetRequest)

	isConnectedStatus, err := e.changeStatusToConnectedRanManager.ChangeStatusToConnectedRan(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to update and notify connection status of nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
	} else {
		if isConnectedStatus {
			nodebInfoupdated, err1 := e.getNodebInfo(request.RanName)
			if err1 != nil {
				e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to get updated nodeB entity. RanName: %s. Error: %s", request.RanName, err1.Error())
			}
			e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - Connection status Set Successfully ran: %s , Connection status updated : %s", ranName, nodebInfoupdated.ConnectionStatus)
		} else {
			e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - Connection status Setting Failed")
		}
	}

	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - nodeB entity connected state. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

}

func (e *E2ResetRequestNotificationHandler) getNodebInfo(ranName string) (*entities.NodebInfo, error) {

	nodebInfo, err := e.rnibDataService.GetNodeb(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to retrieve nodeB entity. RanName: %s. Error: %s", ranName, err.Error())
		return nil, err
	}
	return nodebInfo, err
}

func (e *E2ResetRequestNotificationHandler) waitfortimertimeout(request *models.NotificationRequest) {
	timeout := e.config.E2ResetTimeOutSec
	for {
		timeElapsed := utils.ElapsedTime(request.StartTime)
		e.logger.Infof(E2ResetRequestLogInfoElapsedTime, utils.ElapsedTime(request.StartTime))
		if int(timeElapsed) > timeout {
			break
		}
		time.Sleep(time.Duration(timeout/100) * time.Millisecond)
	}
}

func (e *E2ResetRequestNotificationHandler) parseE2ResetMessage(payload []byte) (*models.E2ResetRequestMessage, error) {
	e2resetMessage := models.E2ResetRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload), &(e2resetMessage.E2APPDU))

	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - error in parsing request message: %+v", err)
		return nil, err
	}
	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - Unmarshalling is successful %v", e2resetMessage.E2APPDU.InitiatingMessage.ProcedureCode)
	return &e2resetMessage, nil
}

func (h *E2ResetRequestNotificationHandler) handleSuccessfulResponse(ranName string, req *models.NotificationRequest, resetRequest *models.E2ResetRequestMessage) {

	successResponse := models.NewE2ResetResponseMessage(resetRequest)
	h.logger.Debugf("#E2ResetRequestNotificationHandler.handleSuccessfulResponse - E2_RESET_RESPONSE has been built successfully %+v", successResponse)

	responsePayload, err := xml.Marshal(&successResponse.E2ApPdu)
	if err != nil {
		h.logger.Warnf("#E2ResetRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_RESET_RESP. Payload: %s", ranName, responsePayload)
	}

	responsePayload = utils.ReplaceEmptyTagsWithSelfClosing(responsePayload, resetRequestEmptyTagsToReplaceToSelfClosingTags)

	h.logger.Infof("#E2ResetRequestNotificationHandler.handleSuccessfulResponse - payload: %s", responsePayload)

	msg := models.NewRmrMessage(rmrCgo.RIC_E2_RESET_RESP, ranName, responsePayload, req.TransactionId, req.GetMsgSrc())
	h.logger.Infof("#E2ResetRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - RIC_E2_RESET_RESP message has been built successfully. Message: %x", ranName, msg)
	err = h.rmrSender.Send(msg)
	if err != nil {
		h.logger.Errorf("#E2ResetRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - Error sending e2 success response %+v", ranName, msg)
	}
}
