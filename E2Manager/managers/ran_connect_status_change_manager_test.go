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
	"e2mgr/mocks"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const EventChannelForTest = "RAN_CONNECTION_STATUS_CHANGE"

func initRanConnectStatusChangeManagerTest(t *testing.T) (*mocks.RnibWriterMock, *mocks.RanListManagerMock, *mocks.RanAlarmServiceMock, *RanConnectStatusChangeManager) {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize log, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel: EventChannelForTest,
		},
	}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	ranListManagerMock := &mocks.RanListManagerMock{}
	ranAlarmServiceMock := &mocks.RanAlarmServiceMock{}
	ranConnectStatusChangeManager := NewRanConnectStatusChangeManager(log, rnibDataService, ranListManagerMock, ranAlarmServiceMock)
	return writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager
}

func TestChangeStatusSuccessNewRan(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_UNKNOWN_CONNECTION_STATUS}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, EventChannelForTest, RanName+"_"+CONNECTED_RAW_EVENT).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	ranAlarmServiceMock.On("SetConnectivityChangeAlarm", &updatedNodebInfo).Return(nil)
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_CONNECTED)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusSuccessEventNone1(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusSuccessEventNone2(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusSuccessEventConnected(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, EventChannelForTest, RanName+"_"+CONNECTED_RAW_EVENT).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	ranAlarmServiceMock.On("SetConnectivityChangeAlarm", &updatedNodebInfo).Return(nil)
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_CONNECTED)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusSuccessEventDisconnected(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, EventChannelForTest, RanName+"_"+DISCONNECTED_RAW_EVENT).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	ranAlarmServiceMock.On("SetConnectivityChangeAlarm", &updatedNodebInfo).Return(nil)
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_DISCONNECTED)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusRnibErrorEventNone(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	assert.NotNil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusRnibErrorEventConnected(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, EventChannelForTest, RanName+"_"+CONNECTED_RAW_EVENT).Return(common.NewInternalError(errors.New("Error")))
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_CONNECTED)
	assert.NotNil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusRanListManagerError(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}

func TestChangeStatusRanAlarmServiceErrorEventConnected(t *testing.T) {
	writerMock, ranListManagerMock, ranAlarmServiceMock, ranConnectStatusChangeManager := initRanConnectStatusChangeManagerTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, EventChannelForTest, RanName+"_"+CONNECTED_RAW_EVENT).Return(nil)
	ranListManagerMock.On("UpdateRanState", &updatedNodebInfo).Return(nil)
	ranAlarmServiceMock.On("SetConnectivityChangeAlarm", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranConnectStatusChangeManager.ChangeStatus(origNodebInfo, entities.ConnectionStatus_CONNECTED)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
	ranAlarmServiceMock.AssertExpectations(t)
}
