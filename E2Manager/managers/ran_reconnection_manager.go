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
	"e2mgr/rNibWriter"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type IRanReconnectionManager interface {
	ReconnectRan(inventoryName string) error
}

type RanReconnectionManager struct {
	logger             *logger.Logger
	config             *configuration.Configuration
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
	ranSetupManager    *RanSetupManager
}

func NewRanReconnectionManager(logger *logger.Logger, config *configuration.Configuration, rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter, rmrService *services.RmrService) *RanReconnectionManager {
	return &RanReconnectionManager{
		logger:             logger,
		config:             config,
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
		ranSetupManager:    NewRanSetupManager(logger,rmrService,rnibWriterProvider),
	}
}

func (m *RanReconnectionManager) ReconnectRan(inventoryName string) error {
	nodebInfo, rnibErr := m.rnibReaderProvider().GetNodeb(inventoryName)

	if rnibErr != nil {
		m.logger.Errorf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Failed fetching RAN from rNib. Error: %v", inventoryName, rnibErr)
		return rnibErr
	}

	if !m.canReconnectRan(nodebInfo) {
		m.logger.Warnf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Cannot reconnect RAN", inventoryName)
		return m.setConnectionStatusOfUnconnectableRan(nodebInfo)
	}

	err := m.ranSetupManager.ExecuteSetup(nodebInfo)

	if err != nil {
		m.logger.Errorf("#RanReconnectionManager.ReconnectRan - RAN name: %s - Failed executing setup. Error: %v", inventoryName, err)
		return err
	}

	m.logger.Infof("#RanReconnectionManager.ReconnectRan - RAN name: %s - Successfully done executing setup. RAN's connection attempts: %d", inventoryName, nodebInfo.ConnectionAttempts)
	return nil
}

func (m *RanReconnectionManager) canReconnectRan(nodebInfo *entities.NodebInfo) bool {
	connectionStatus := nodebInfo.GetConnectionStatus()
	return connectionStatus != entities.ConnectionStatus_SHUT_DOWN && connectionStatus != entities.ConnectionStatus_SHUTTING_DOWN &&
		int(nodebInfo.GetConnectionAttempts()) < m.config.MaxConnectionAttempts
}

func (m *RanReconnectionManager) updateNodebInfoStatus(nodebInfo *entities.NodebInfo, connectionStatus entities.ConnectionStatus) common.IRNibError {
	if nodebInfo.ConnectionStatus == connectionStatus {
		return nil
	}

	nodebInfo.ConnectionStatus = connectionStatus;
	err := m.rnibWriterProvider().UpdateNodebInfo(nodebInfo)

	if err != nil {
		m.logger.Errorf("#RanReconnectionManager.updateNodebInfoStatus - RAN name: %s - Failed updating RAN's connection status to %s in rNib. Error: %v", nodebInfo.RanName, connectionStatus, err)
		return err
	}

	m.logger.Infof("#RanReconnectionManager.updateNodebInfoStatus - RAN name: %s - Successfully updated RAN's connection status to %s in rNib", nodebInfo.RanName, connectionStatus)
	return nil
}

func (m *RanReconnectionManager) setConnectionStatusOfUnconnectableRan(nodebInfo *entities.NodebInfo) common.IRNibError {
	connectionStatus := nodebInfo.GetConnectionStatus()
	m.logger.Warnf("#RanReconnectionManager.setConnectionStatusOfUnconnectableRan - RAN name: %s, RAN's connection status: %s, RAN's connection attempts: %d", nodebInfo.RanName, nodebInfo.ConnectionStatus, nodebInfo.ConnectionAttempts)

	if connectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		return m.updateNodebInfoStatus(nodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	}

	if int(nodebInfo.GetConnectionAttempts()) >= m.config.MaxConnectionAttempts {
		m.logger.Warnf("#RanReconnectionManager.setConnectionStatusOfUnconnectableRan - RAN name: %s - RAN's connection attempts are greater than %d", nodebInfo.RanName, m.config.MaxConnectionAttempts)
		return m.updateNodebInfoStatus(nodebInfo, entities.ConnectionStatus_DISCONNECTED)
	}

	return nil
}
