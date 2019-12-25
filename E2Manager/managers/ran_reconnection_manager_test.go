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
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *RanReconnectionManager, *mocks.HttpClientMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)
	ranSetupManager := NewRanSetupManager(logger, rmrSender, rnibDataService)
	httpClient := &mocks.HttpClientMock{}
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClient)
	e2tAssociationManager := NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient)
	ranReconnectionManager := NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, ranSetupManager, e2tAssociationManager)
	return logger, rmrMessengerMock, readerMock, writerMock, ranReconnectionManager, httpClient
}

func TestRanReconnectionGetNodebFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestShutdownRanReconnection(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestShuttingdownRanReconnection(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnectingRanWithMaxAttemptsReconnectionDissociateSucceeds(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, httpClient:= initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, ConnectionAttempts: 20, AssociatedE2TInstanceAddress: E2TAddress}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList:[]string{ranName}}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	e2tInstanceToSave := * e2tInstance
	e2tInstanceToSave .AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &e2tInstanceToSave).Return(nil)
	mockHttpClient(httpClient, clients.DissociateRanE2TInstanceApiSuffix, true)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
}

func TestConnectingRanWithMaxAttemptsReconnectionDissociateFails(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	e2tAddress := "10.0.2.15"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, ConnectionAttempts: 20, AssociatedE2TInstanceAddress: e2tAddress}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address:e2tAddress, AssociatedRanList:[]string{ranName}}
	readerMock.On("GetE2TInstance",e2tAddress).Return(e2tInstance, common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	writerMock.AssertNotCalled(t, "SaveE2TInstance", )
}

func TestUnconnectableRanUpdateNodebInfoFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnectedRanExecuteSetupSuccess(t *testing.T) {
	_, rmrMessengerMock, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	updatedNodebInfo.ConnectionAttempts++
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, nil)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestConnectedRanExecuteSetupFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	updatedNodebInfo.ConnectionAttempts++
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestNoSetConnectionStatus(t *testing.T) {
	_, _, _, _, ranReconnectionManager, _ := initRanLostConnectionTest(t)
	nodebInfo := &entities.NodebInfo{RanName: "ranName", GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	err := ranReconnectionManager.updateUnconnectableRan(nodebInfo)
	assert.Nil(t, err)
}

func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}
