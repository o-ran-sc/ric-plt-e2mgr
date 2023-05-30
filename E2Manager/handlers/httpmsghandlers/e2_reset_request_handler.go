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

package httpmsghandlers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type E2ResetRequestHandler struct {
	rNibDataService                   services.RNibDataService
	rmrSender                         *rmrsender.RmrSender
	logger                            *logger.Logger
	ranResetManager                   *managers.RanResetManager
	changeStatusToConnectedRanManager *managers.ChangeStatusToConnectedRanManager
}

func NewE2ResetRequestHandler(logger *logger.Logger, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService, ranResetManager *managers.RanResetManager, changeStatusToConnectedRanManager *managers.ChangeStatusToConnectedRanManager) *E2ResetRequestHandler {
	return &E2ResetRequestHandler{
		rNibDataService:                   rNibDataService,
		rmrSender:                         rmrSender,
		logger:                            logger,
		ranResetManager:                   ranResetManager,
		changeStatusToConnectedRanManager: changeStatusToConnectedRanManager,
	}
}

func (e *E2ResetRequestHandler) Handle(request models.Request) error {
	resetRequest := request.(models.ResetRequest)
	e.logger.Infof("#E2ResetRequestHandler.Handle - Ran name: %s", resetRequest.RanName)
	nodebInfo, err := e.getNodebInfo(resetRequest.RanName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to get status of RAN: %s from RNIB. Error: %s", resetRequest.RanName, err.Error())
		_, ok := err.(*common.ResourceNotFoundError)
		if ok {
			return e2managererrors.NewResourceNotFoundError()
		}
		return e2managererrors.NewRnibDbError()
	}

	e.logger.Debugf("#E2ResetRequestHandler.Handle - nodeB entity retrieved. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

	ranName := resetRequest.RanName
	isResetDone, err := e.ranResetManager.ResetRan(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to update and notify connection status of nodeB entity. RanName: %s. Error: %s", resetRequest.RanName, err.Error())
		return err
	} else {
		if isResetDone {
			nodebInfoupdated, err1 := e.getNodebInfo(resetRequest.RanName)
			if err1 != nil {
				e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to get updated nodeB entity. RanName: %s. Error: %s", resetRequest.RanName, err1.Error())
				return err1
			}
			e.logger.Debugf("#E2ResetRequestHandler.Handle - Reset Done Successfully ran: %s , Connection status updated : %s", ranName, nodebInfoupdated.ConnectionStatus)
		} else {
			e.logger.Debugf("#E2ResetRequestHandler.Handle - Reset Failed")
		}
	}

	//Todo: add timer
	time.Sleep(5 * time.Second)

	isConnectedStatus, err := e.changeStatusToConnectedRanManager.ChangeStatusToConnectedRan(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to update and notify connection status of nodeB entity. RanName: %s. Error: %s", resetRequest.RanName, err.Error())
		return err
	} else {
		if isConnectedStatus {
			nodebInfoupdated, err1 := e.getNodebInfo(resetRequest.RanName)
			if err1 != nil {
				e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to get updated nodeB entity. RanName: %s. Error: %s", resetRequest.RanName, err1.Error())
				return err1
			}
			e.logger.Debugf("#E2ResetRequestHandler.Handle - Connection status Set Successfully ran: %s , Connection status updated : %s", ranName, nodebInfoupdated.ConnectionStatus)
		} else {
			e.logger.Debugf("#E2ResetRequestHandler.Handle - Connection status Setting Failed")
		}
	}

	e.logger.Debugf("#E2ResetRequestHandler.Handle - nodeB entity connected state. RanName %s, ConnectionStatus %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)

	return nil
}

func (e *E2ResetRequestHandler) getNodebInfo(ranName string) (*entities.NodebInfo, error) {

	nodebInfo, err := e.rNibDataService.GetNodeb(ranName)
	if err != nil {
		e.logger.Errorf("#E2ResetRequestHandler.Handle - failed to retrieve nodeB entity. RanName: %s. Error: %s", ranName, err.Error())
		return nil, err
	}
	return nodebInfo, nil
}
