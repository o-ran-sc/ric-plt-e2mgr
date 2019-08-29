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
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestExecuteSetupConnectingX2Setup(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var rnibErr common.IRNibError
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err != nil {
		t.Errorf("want: success, got: error: %s", err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupConnectingEndcX2Setup(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var rnibErr common.IRNibError
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err != nil {
		t.Errorf("want: success, got: error: %s", err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupDisconnected(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	var rnibErr common.IRNibError
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(rnibErr)

	payload := []byte{0}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, fmt.Errorf("send failure"))
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err == nil {
		t.Errorf("want: failure, got: success")
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupConnectingRnibError(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	var rnibErr = common.NewInternalError(fmt.Errorf("DB error"))
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(rnibErr)

	payload := []byte{0}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, fmt.Errorf("send failure"))
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err == nil {
		t.Errorf("want: failure, got: success")
	} else {
		assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestExecuteSetupDisconnectedRnibError(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	var rnibErr common.IRNibError
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(common.NewInternalError(fmt.Errorf("DB error")))

	payload := []byte{0}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, fmt.Errorf("send failure"))
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err == nil {
		t.Errorf("want: failure, got: success")
	} else {
		assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupUnsupportedProtocol(t *testing.T) {
	log := initLog(t)

	ranName := "test1"

	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_UNKNOWN_E2_APPLICATION_PROTOCOL}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_UNKNOWN_E2_APPLICATION_PROTOCOL, ConnectionAttempts: 1}
	var rnibErr common.IRNibError
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)
	rmrService := getRmrService(rmrMessengerMock, log)

	mgr := NewRanSetupManager(log, rmrService, writerProvider)
	if err := mgr.ExecuteSetup(initialNodeb); err == nil {
		t.Errorf("want: error, got: success")
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}
