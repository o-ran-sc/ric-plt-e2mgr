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

package managers

import (
	// "e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/services"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type IChangeStatusToConnectedRanManager interface {
	ChangeStatusToConnectedRan(inventoryName string) error
}

type ChangeStatusToConnectedRanManager struct {
	logger                        *logger.Logger
	rnibDataService               services.RNibDataService
	ranConnectStatusChangeManager IRanConnectStatusChangeManager
}

func NewChangeStatusToConnectedRanManager(logger *logger.Logger, rnibDataService services.RNibDataService, ranConnectStatusChangeManager IRanConnectStatusChangeManager) *ChangeStatusToConnectedRanManager {
	return &ChangeStatusToConnectedRanManager{
		logger:                        logger,
		rnibDataService:               rnibDataService,
		ranConnectStatusChangeManager: ranConnectStatusChangeManager,
	}
}

func (m *ChangeStatusToConnectedRanManager) ChangeStatusToConnectedRan(inventoryName string) (bool, error) {
	nodebInfo, err := m.rnibDataService.GetNodeb(inventoryName)

	if err != nil {
		m.logger.Errorf("#ChangeStatusToConnectedRanManager.ChangeStatusToConnectedRan - RAN name: %s - Failed fetching RAN from rNib. Error: %v", inventoryName, err)
		return false, err
	}
	connectionStatus := nodebInfo.GetConnectionStatus()
	m.logger.Infof("#ChangeStatusToConnectedRanManager.ChangeStatusToConnectedRan - RAN name: %s - RAN's connection status: %s", nodebInfo.RanName, connectionStatus)

	ranConnectStatusChange, err := m.ranConnectStatusChangeManager.ChangeStatus(nodebInfo, entities.ConnectionStatus_CONNECTED)

	if err != nil {
		return ranConnectStatusChange, err
	}
	m.logger.Infof("#ChangeStatusToConnectedRanManager.ChangeStatusToConnectedRan Status changed \n")
	return ranConnectStatusChange, nil
}
