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
package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/rnibBuilders"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestSetupHandleNewRanSave_Error(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)

	ranName := "RanName"
	rnibErr := &common.ResourceNotFoundError{}
	sr := models.SetupRequest{"127.0.0.1", 8080, ranName,}

	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, rnibErr)

	vErr := &common.ValidationError{}
	updatedNb, _ := rnibBuilders.CreateInitialNodeInfo(&sr, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(vErr)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, Ip: sr.RanIp, Port: uint32(sr.RanPort), ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	actual := handler.Handle(sr)
	expected := &e2managererrors.RnibDbError{}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}

	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestSetupHandleNewRan_Success(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)

	ranName := "RanName"
	rnibErr := &common.ResourceNotFoundError{}
	sr := models.SetupRequest{"127.0.0.1", 8080, ranName,}

	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, rnibErr)

	updatedNb, _ := rnibBuilders.CreateInitialNodeInfo(&sr, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, Ip: sr.RanIp, Port: uint32(sr.RanPort), ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	actual := handler.Handle(sr)

	assert.Nil(t, actual)

	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestEndcSetupHandleRmr_Error(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)

	ranName := "RanName"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	var nbDisconnected = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	writerMock.On("UpdateNodebInfo", nbDisconnected).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrErr := &e2managererrors.RmrError{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, rmrErr)

	sr := models.SetupRequest{"127.0.0.1", 8080, ranName,}
	actual := handler.Handle(sr)

	if reflect.TypeOf(actual) != reflect.TypeOf(rmrErr) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, rmrErr)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestEndcSetupHandleExistingDisconnectedRan_Success(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)

	ranName := "RanName"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	sr := models.SetupRequest{"127.0.0.1", 8080, ranName,}
	actual := handler.Handle(sr)

	assert.Nil(t, actual)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestX2SetupHandleExistingConnectedRan_Success(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	ranName := "RanName"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	sr := models.SetupRequest{"127.0.0.1", 8080, ranName,}
	actual := handler.Handle(sr)

	assert.Nil(t, actual)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestX2SetupHandleRnibGet_Error(t *testing.T) {
	readerMock, _, handler,rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	rnibErr := &common.ValidationError{}
	nb := &entities.NodebInfo{RanName: "RanName", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName").Return(nb, rnibErr)

	sr := models.SetupRequest{"127.0.0.1", 8080, "RanName",}
	actual := handler.Handle(sr)

	expected := &e2managererrors.RnibDbError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestX2SetupHandleShuttingDownRan_Error(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	nb := &entities.NodebInfo{RanName: "RanName", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName").Return(nb, nil)

	sr := models.SetupRequest{"127.0.0.1", 8080, "RanName",}
	actual := handler.Handle(sr)

	expected := &e2managererrors.WrongStateError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestX2SetupHandleNoPort_Error(t *testing.T) {
	_, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{"127.0.0.1", 0, "RanName",}
	actual := handler.Handle(sr)

	expected := &e2managererrors.RequestValidationError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestX2SetupHandleNoRanName_Error(t *testing.T) {
	_, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{}
	sr.RanPort = 8080
	sr.RanIp = "127.0.0.1"

	actual := handler.Handle(sr)

	expected := &e2managererrors.RequestValidationError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestX2SetupHandleNoIP_Error(t *testing.T) {
	_, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{}
	sr.RanPort = 8080
	sr.RanName = "RanName"

	actual := handler.Handle(sr)

	expected := &e2managererrors.RequestValidationError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestX2SetupHandleInvalidIp_Error(t *testing.T) {
	_, writerMock, handler, rmrMessengerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{}
	sr.RanPort = 8080
	sr.RanName = "RanName"
	sr.RanIp = "invalid ip"

	actual := handler.Handle(sr)

	expected := &e2managererrors.RequestValidationError{}
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func initSetupRequestTest(t *testing.T, protocol entities.E2ApplicationProtocol)(*mocks.RnibReaderMock, *mocks.RnibWriterMock, *SetupRequestHandler, *mocks.RmrMessengerMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrService := getRmrService(rmrMessengerMock, log)

	rnibDataService := services.NewRnibDataService(log, config, readerProvider, writerProvider)

	ranSetupManager := managers.NewRanSetupManager(log, rmrService, rnibDataService)
	handler := NewSetupRequestHandler(log, rnibDataService, ranSetupManager, protocol)

	return readerMock, writerMock, handler, rmrMessengerMock
}