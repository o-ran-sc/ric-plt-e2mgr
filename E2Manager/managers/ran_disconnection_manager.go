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

package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type IRanDisconnectionManager interface {
	DisconnectRan(inventoryName string) error
}

type RanDisconnectionManager struct {
	logger                        *logger.Logger
	config                        *configuration.Configuration
	rnibDataService               services.RNibDataService
	e2tAssociationManager         *E2TAssociationManager
	ranConnectStatusChangeManager IRanConnectStatusChangeManager
}

func NewRanDisconnectionManager(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, e2tAssociationManager *E2TAssociationManager, ranConnectStatusChangeManager IRanConnectStatusChangeManager) *RanDisconnectionManager {
	return &RanDisconnectionManager{
		logger:                        logger,
		config:                        config,
		rnibDataService:               rnibDataService,
		e2tAssociationManager:         e2tAssociationManager,
		ranConnectStatusChangeManager: ranConnectStatusChangeManager,
	}
}

func (m *RanDisconnectionManager) DisconnectRan(inventoryName string) error {
	nodebInfo, err := m.rnibDataService.GetNodeb(inventoryName)

	if err != nil {
		m.logger.Errorf("#RanDisconnectionManager.DisconnectRan - RAN name: %s - Failed fetching RAN from rNib. Error: %v", inventoryName, err)
		return err
	}

	connectionStatus := nodebInfo.GetConnectionStatus()
	m.logger.Infof("#RanDisconnectionManager.DisconnectRan - RAN name: %s - RAN's connection status: %s", nodebInfo.RanName, connectionStatus)

	if connectionStatus == entities.ConnectionStatus_SHUT_DOWN {
		m.logger.Warnf("#RanDisconnectionManager.DisconnectRan - RAN name: %s - quit. RAN's connection status is SHUT_DOWN", nodebInfo.RanName)
		return nil
	}

	if connectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		_, err = m.ranConnectStatusChangeManager.ChangeStatus(nodebInfo, entities.ConnectionStatus_SHUT_DOWN)
		return err
	}

	_, err = m.ranConnectStatusChangeManager.ChangeStatus(nodebInfo, entities.ConnectionStatus_DISCONNECTED)

	if err != nil {
		return err
	}

	e2tAddress := nodebInfo.AssociatedE2TInstanceAddress
	return m.e2tAssociationManager.DissociateRan(e2tAddress, nodebInfo.RanName)
}
