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
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	//        "github.com/pkg/errors"
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

	for _, ranName := range ranNameList {
		nodebInfo, err := h.rNibDataService.GetNodeb(ranName) //This method is needed for getting RAN functions with later commits
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

		}
	}
	if isAtleastOneRanConnected == false {
		return nil, e2managererrors.NewNoConnectedRanError()
	}

	return nil, nil
}

func (h *HealthCheckRequestHandler) getRanNameList(request models.Request) []string {
	healthCheckRequest := request.(models.HealthCheckRequest)
	if request != nil && len(healthCheckRequest.RanList) != 0 {
		return healthCheckRequest.RanList
	}
	nodeIds := h.ranListManager.GetNbIdentityList()

	var ranNameList []string
	for _, nbIdentity := range nodeIds {
		if nbIdentity.ConnectionStatus == entities.ConnectionStatus_CONNECTED {
			ranNameList = append(ranNameList, nbIdentity.InventoryName)
		}
	}
	return ranNameList
}
