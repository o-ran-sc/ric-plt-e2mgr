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
	"testing"
)

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *RanReconnectionManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	rmrService := getRmrService(&mocks.RmrMessengerMock{}, logger)

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	ranReconnectionManager := NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibReaderProvider, rnibWriterProvider, rmrService)
	return logger, readerMock, writerMock, ranReconnectionManager
}

func TestLostConnectionFetchingNodebFailure(t *testing.T) {
	_, readerMock, _, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
}

func TestLostConnectionUpdatingNodebForUnconnectableRanFailure(t *testing.T) {
	_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr common.IRNibError
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.NotNil(t, err)
}

func TestLostConnectionOfConnectedRanWithMaxAttempts(t *testing.T) {
	_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED, ConnectionAttempts: 20}
	var rnibErr common.IRNibError
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
}

func TestLostConnectionOfShutdownRan(t *testing.T) {
	_, readerMock, _, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	var rnibErr common.IRNibError
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
}

func TestLostConnectionOfShuttingdownRan(t *testing.T) {
	_, readerMock, writerMock, ranReconnectionManager := initRanLostConnectionTest(t)
	ranName := "test"
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr common.IRNibError
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	err := ranReconnectionManager.ReconnectRan(ranName)
	assert.Nil(t, err)
}

// TODO: should extract to test_utils
func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, make(sessions.E2Sessions), messageChannel)
}
