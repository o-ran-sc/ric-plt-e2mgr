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

package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type ranAlarmServiceInstance struct {
	logger *logger.Logger
	config *configuration.Configuration
}

type RanAlarmService interface {
	SetConnectivityChangeAlarm(nodebInfo *entities.NodebInfo) error
}

func NewRanAlarmService(logger *logger.Logger, config *configuration.Configuration) RanAlarmService {
	return &ranAlarmServiceInstance{
		logger: logger,
		config: config,
	}
}

func (m *ranAlarmServiceInstance) SetConnectivityChangeAlarm(nodebInfo *entities.NodebInfo) error {
	m.logger.Infof("#ranAlarmServiceInstance.SetConnectivityChangeAlarm - RAN name: %s - Connectivity state was changed to %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)
	return nil
}
