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

package controllers

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupControllerTest(t *testing.T) (*Controller, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock){
	log := initLog(t)
	config := configuration.ParseConfiguration()

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rnibDataService := services.NewRnibDataService(log, config, readerProvider, writerProvider)

	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rnibDataService)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), rnibDataService, config, ranSetupManager)
	return controller, readerMock, writerMock, rmrMessengerMock
}

func TestX2SetupInvalidBody(t *testing.T) {

	controller, _, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest, _ := http.NewRequest("POST", "http://localhost:3800/v1/nodeb/x2-setup", strings.NewReader("{}{}"))
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.X2SetupHandler(writer, httpRequest)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	assert.Equal(t, e2managererrors.NewInvalidJsonError().Code, errorResponse.Code)
}

func TestX2SetupSuccess(t *testing.T) {

	controller, readerMock, writerMock, rmrMessengerMock := setupControllerTest(t)

	ranName := "test"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest := tests.GetHttpRequest()
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.X2SetupHandler(writer, httpRequest)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestEndcSetupSuccess(t *testing.T) {

	controller, readerMock, writerMock, rmrMessengerMock := setupControllerTest(t)

	ranName := "test"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)

	var nbUpdated = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", nbUpdated).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(msg, nil)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest := tests.GetHttpRequest()
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.EndcSetupHandler(writer, httpRequest)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestShutdownHandlerRnibError(t *testing.T) {
	controller, readerMock, _, _:= setupControllerTest(t)

	rnibErr := &common.ResourceNotFoundError{}
	var nbIdentityList []*entities.NbIdentity
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibErr)

	writer := httptest.NewRecorder()

	controller.ShutdownHandler(writer, tests.GetHttpRequest())

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, e2managererrors.NewRnibDbError().Code)
}

func TestHeaderValidationFailed(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()

	header := &http.Header{}

	controller.handleRequest(writer, header, httpmsghandlerprovider.ShutdownRequest, nil, true)

	var errorResponse = parseJsonRequest(t, writer.Body)
	err := e2managererrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestShutdownStatusNoContent(t *testing.T) {
	controller, readerMock, _, _ := setupControllerTest(t)

	var rnibError error
	nbIdentityList := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibError)

	writer := httptest.NewRecorder()
	controller.ShutdownHandler(writer, tests.GetHttpRequest())

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestHandleInternalError(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := e2managererrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleCommandAlreadyInProgressError(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	err := e2managererrors.NewCommandAlreadyInProgressError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestValidateHeaders(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	result := controller.validateRequestHeader(&header)

	assert.Nil(t, result)
}

func parseJsonRequest(t *testing.T, r io.Reader) models.ErrorResponse {

	var errorResponse models.ErrorResponse
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error cannot deserialize json request")
	}
	json.Unmarshal(body, &errorResponse)

	return errorResponse
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}

func TestX2ResetHandleSuccessfulRequestedCause(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock := setupControllerTest(t)

	ranName := "test1"
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	data4Req := map[string]interface{}{"cause": "protocol:transfer-syntax-error"}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data4Req)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2ResetHandler(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)

}

func TestX2ResetHandleSuccessfulRequestedDefault(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock := setupControllerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	// no body
	b := new(bytes.Buffer)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2ResetHandler(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)

}

func TestX2ResetHandleFailureInvalidBody(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)

	ranName := "test1"

	writer := httptest.NewRecorder()

	// Invalid json: attribute name without quotes (should be "cause":).
	b := strings.NewReader("{cause:\"protocol:transfer-syntax-error\"")
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2ResetHandler(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

}

func TestHandleErrorResponse(t *testing.T) {
	controller, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRnibDbError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewCommandAlreadyInProgressError(), writer)
	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewHeaderValidationError(), writer)
	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewWrongStateError("", ""), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRequestValidationError(), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRmrError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewResourceNotFoundError(), writer)
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(fmt.Errorf("ErrorError"), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}