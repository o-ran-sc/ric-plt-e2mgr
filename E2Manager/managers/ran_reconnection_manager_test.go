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
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/sessions"
	"e2mgr/tests"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func initRanLostConnectionTest(t *testing.T) (*logger.Logger,*mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *RanReconnectionManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrService := getRmrService(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	ranSetupManager := NewRanSetupManager(logger,rmrService, rnibWriterProvider)
	ranReconnectionManager := NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibReaderProvider, rnibWriterProvider, ranSetupManager)
	return logger,rmrMessengerMock, readerMock, writerMock, ranReconnectionManager
}

func TestRanReconnectionGetNodebFailure(t *testing.T) {
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestShutdownRanReconnection(t *testing.T) {
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
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
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
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

func TestConnectingRanWithMaxAttemptsReconnection(t *testing.T) {
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, ConnectionAttempts: 20}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestUnconnectableRanUpdateNodebInfoFailure(t *testing.T) {
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
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
	_,rmrMessengerMock, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	updatedNodebInfo.ConnectionAttempts++
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	rmrMessengerMock.On("SendMsg",mock.Anything, mock.AnythingOfType("int")).Return(&rmrCgo.MBuf{},nil)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestConnectedRanExecuteSetupFailure(t *testing.T) {
	_,_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
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

// TODO: should extract to test_utils
func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, make(sessions.E2Sessions), messageChannel)
}
