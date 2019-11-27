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
//

package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type IRanReconnectionManager interface {
	ReconnectRan(inventoryName string) error
}

type RanReconnectionManager struct {
	logger          *logger.Logger
	config          *configuration.Configuration
	rnibDataService services.RNibDataService
	ranSetupManager *RanSetupManager
}

func NewRanReconnectionManager(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, ranSetupManager *RanSetupManager) *RanReconnectionManager {
	return &RanReconnectionManager{
		logger:          logger,
		config:          config,
		rnibDataService: rnibDataService,
		ranSetupManager: ranSetupManager,
	}
}

func (m *RanReconnectionManager) ReconnectRan(inventoryName string) error {
	nodebInfo, rnibErr := m.rnibDataService.GetNodeb(inventoryName)

	if rnibErr != nil {
		m.logger.Errorf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Failed fetching RAN from rNib. Error: %v", inventoryName, rnibErr)
		return rnibErr
	}

	m.logger.Infof("#RanReconnectionManager.ReconnectRan - RAN name: %s - RAN's connection status: %s, RAN's connection attempts: %d", nodebInfo.RanName, nodebInfo.ConnectionStatus, nodebInfo.ConnectionAttempts)

	if !m.canReconnectRan(nodebInfo) {
		return m.setConnectionStatusOfUnconnectableRan(nodebInfo)
	}

	err := m.ranSetupManager.ExecuteSetup(nodebInfo, entities.ConnectionStatus_CONNECTING)

	if err != nil {
		m.logger.Errorf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Failed executing setup. Error: %v", inventoryName, err)
		return err
	}

	return nil
}

func (m *RanReconnectionManager) canReconnectRan(nodebInfo *entities.NodebInfo) bool {
	connectionStatus := nodebInfo.GetConnectionStatus()
	return connectionStatus != entities.ConnectionStatus_SHUT_DOWN && connectionStatus != entities.ConnectionStatus_SHUTTING_DOWN &&
		int(nodebInfo.GetConnectionAttempts()) < m.config.MaxConnectionAttempts
}

func (m *RanReconnectionManager) updateNodebInfoStatus(nodebInfo *entities.NodebInfo, connectionStatus entities.ConnectionStatus) error {
	if nodebInfo.ConnectionStatus == connectionStatus {
		return nil
	}

	nodebInfo.ConnectionStatus = connectionStatus;
	err := m.rnibDataService.UpdateNodebInfo(nodebInfo)

	if err != nil {
		m.logger.Errorf("#RanReconnectionManager.updateNodebInfoStatus - RAN name: %s - Failed updating RAN's connection status to %s in rNib. Error: %v", nodebInfo.RanName, connectionStatus, err)
		return err
	}

	m.logger.Infof("#RanReconnectionManager.updateNodebInfoStatus - RAN name: %s - Successfully updated rNib. RAN's current connection status: %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)
	return nil
}

func (m *RanReconnectionManager) setConnectionStatusOfUnconnectableRan(nodebInfo *entities.NodebInfo) error {
	connectionStatus := nodebInfo.GetConnectionStatus()

	if connectionStatus == entities.ConnectionStatus_SHUT_DOWN {
		m.logger.Warnf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Cannot reconnect RAN. Reason: connection status is SHUT_DOWN", nodebInfo.RanName)
		return nil
	}

	if connectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		m.logger.Warnf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Cannot reconnect RAN. Reason: connection status is SHUTTING_DOWN", nodebInfo.RanName)
		return m.updateNodebInfoStatus(nodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	}

	if int(nodebInfo.GetConnectionAttempts()) >= m.config.MaxConnectionAttempts {
		m.logger.Warnf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Cannot reconnect RAN. Reason: RAN's connection attempts exceeded the limit (%d)", nodebInfo.RanName, m.config.MaxConnectionAttempts)
		return m.updateNodebInfoStatus(nodebInfo, entities.ConnectionStatus_DISCONNECTED)
	}

	return nil
}
