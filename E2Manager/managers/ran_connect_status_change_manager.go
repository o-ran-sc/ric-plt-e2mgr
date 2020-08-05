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
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

const (
	CONNECTED_RAW_EVENT    = "CONNECTED"
	DISCONNECTED_RAW_EVENT = "DISCONNECTED"
	NONE_RAW_EVENT         = "NONE"
)

type IRanConnectStatusChangeManager interface {
	ChangeStatus(nodebInfo *entities.NodebInfo, nextStatus entities.ConnectionStatus) (bool, error)
}

type RanConnectStatusChangeManager struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	ranListManager  RanListManager
	ranAlarmService services.RanAlarmService
}

func NewRanConnectStatusChangeManager(logger *logger.Logger, rnibDataService services.RNibDataService, ranListManager RanListManager, ranAlarmService services.RanAlarmService) *RanConnectStatusChangeManager {
	return &RanConnectStatusChangeManager{
		logger:          logger,
		rnibDataService: rnibDataService,
		ranListManager:  ranListManager,
		ranAlarmService: ranAlarmService,
	}
}

func (m *RanConnectStatusChangeManager) ChangeStatus(nodebInfo *entities.NodebInfo, nextStatus entities.ConnectionStatus) (bool, error) {
	m.logger.Infof("#RanConnectStatusChangeManager.ChangeStatus - RAN name: %s, currentStatus: %s, nextStatus: %s", nodebInfo.RanName, nodebInfo.GetConnectionStatus(), nextStatus)

	var ranStatusChangePublished bool

	// set the proper event
	event := m.setEvent(nodebInfo, nextStatus)
	isConnectivityEvent := event != NONE_RAW_EVENT

	// only after determining event we set next status
	nodebInfo.ConnectionStatus = nextStatus;
	if !isConnectivityEvent {
		err := m.updateNodebInfo(nodebInfo)
		if err != nil {
			return ranStatusChangePublished, err
		}
	} else {
		err := m.updateNodebInfoOnConnectionStatusInversion(nodebInfo, event)
		if err != nil {
			return ranStatusChangePublished, err
		}
		ranStatusChangePublished = true
	}

	// in any case, update RanListManager
	connectionStatus := nodebInfo.GetConnectionStatus()
	m.logger.Infof("#RanConnectStatusChangeManager.ChangeStatus - RAN name: %s, updating RanListManager... status: %s", nodebInfo.RanName, connectionStatus)
	err := m.ranListManager.UpdateNbIdentityConnectionStatus(nodebInfo.GetNodeType(), nodebInfo.RanName, connectionStatus)
	if err != nil {
		m.logger.Errorf("#RanConnectStatusChangeManager.ChangeStatus - RAN name: %s - Failed updating RAN's connection status by RanListManager. Error: %v", nodebInfo.RanName, err)
		// log and proceed...
	}

	if isConnectivityEvent {
		m.logger.Infof("#RanConnectStatusChangeManager.ChangeStatus - RAN name: %s, setting alarm at RanAlarmService... event: %s", nodebInfo.RanName, event)
		err := m.ranAlarmService.SetConnectivityChangeAlarm(nodebInfo)
		if err != nil {
			m.logger.Errorf("#RanConnectStatusChangeManager.ChangeStatus - RAN name: %s - Failed setting an alarm by RanAlarmService. Error: %v", nodebInfo.RanName, err)
			// log and proceed...
		}
	}

	return ranStatusChangePublished, nil
}

func (m *RanConnectStatusChangeManager) updateNodebInfoOnConnectionStatusInversion(nodebInfo *entities.NodebInfo, event string) error {

	err := m.rnibDataService.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, event)

	if err != nil {
		m.logger.Errorf("#RanConnectStatusChangeManager.updateNodebInfoOnConnectionStatusInversion - RAN name: %s - Failed updating RAN's connection status in rNib. Error: %v", nodebInfo.RanName, err)
		return err
	}

	m.logger.Infof("#RanConnectStatusChangeManager.updateNodebInfoOnConnectionStatusInversion - RAN name: %s - Successfully updated rNib.", nodebInfo.RanName)
	return nil
}

func (m *RanConnectStatusChangeManager) updateNodebInfo(nodebInfo *entities.NodebInfo) error {

	err := m.rnibDataService.UpdateNodebInfo(nodebInfo)

	if err != nil {
		m.logger.Errorf("#RanConnectStatusChangeManager.updateNodebInfo - RAN name: %s - Failed updating RAN's connection status in rNib. Error: %v", nodebInfo.RanName, err)
		return err
	}

	m.logger.Infof("#RanConnectStatusChangeManager.updateNodebInfo - RAN name: %s - Successfully updated rNib.", nodebInfo.RanName)
	return nil
}

func (m *RanConnectStatusChangeManager) setEvent(nodebInfo *entities.NodebInfo, nextState entities.ConnectionStatus) string {
	currentConnectionStatus := nodebInfo.GetConnectionStatus()

	var event string
	if currentConnectionStatus != entities.ConnectionStatus_CONNECTED && nextState == entities.ConnectionStatus_CONNECTED {
		event = nodebInfo.RanName + "_" + CONNECTED_RAW_EVENT
	} else if currentConnectionStatus == entities.ConnectionStatus_CONNECTED && nextState != entities.ConnectionStatus_CONNECTED {
		event = nodebInfo.RanName + "_" + DISCONNECTED_RAW_EVENT
	} else {
		event = NONE_RAW_EVENT
	}

	m.logger.Infof("#RanConnectStatusChangeManager.setEvent - Connectivity Event for RAN %s is: %s", nodebInfo.RanName, event)
	return event
}
