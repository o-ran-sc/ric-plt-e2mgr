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
	"e2mgr/models"
	"e2mgr/services"
	"e2mgr/utils"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

const E2ResetRequestLogInfoElapsedTime = "#E2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms"

type E2ResetRequestNotificationHandler struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	config          *configuration.Configuration
}

func NewE2ResetRequestNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, config *configuration.Configuration) *E2ResetRequestNotificationHandler {
	return &E2ResetRequestNotificationHandler{
		logger:          logger,
		rnibDataService: rnibDataService,
		config:          config,
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

	err = e.rnibDataService.UpdateNodebInfoAndPublish(nodebInfo)

	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to update connection status of nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
	}

	e.logger.Debugf("#E2ResetRequestNotificationHandler.Handle - nodeB entity under reset state. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

	e.logger.Infof(E2ResetRequestLogInfoElapsedTime, utils.ElapsedTime(request.StartTime))
	timeout := e.config.E2ResetTimeOutSec

	for {
		timeElapsed := utils.ElapsedTime(request.StartTime)
		e.logger.Infof(E2ResetRequestLogInfoElapsedTime, utils.ElapsedTime(request.StartTime))
		if int(timeElapsed) > timeout {
			break
		}
		time.Sleep(time.Duration(timeout/100) * time.Millisecond)
	}
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED

	err = e.rnibDataService.UpdateNodebInfoAndPublish(nodebInfo)

	if err != nil {
		e.logger.Errorf("#E2ResetRequestNotificationHandler.Handle - failed to update connection status of nodeB entity. RanName: %s. Error: %s", request.RanName, err.Error())
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
