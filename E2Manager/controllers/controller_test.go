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
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
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

func TestX2SetupSuccess(t *testing.T) {
/*	log := initLog(t)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	config := configuration.ParseConfiguration()

	header := http.Header{}
	header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)

	httpRequest := tests.GetHttpRequest()
	httpRequest.Header = header
	controller.X2SetupHandler(writer, httpRequest)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)*/
}

func TestShutdownHandlerRnibError(t *testing.T) {
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

	rnibErr := &common.ResourceNotFoundError{}
	var nbIdentityList []*entities.NbIdentity
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibErr)

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)
	controller.ShutdownHandler(writer, tests.GetHttpRequest())

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, e2managererrors.NewRnibDbError().Code)
}

func TestHeaderValidationFailed(t *testing.T) {
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

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)

	header := &http.Header{}

	controller.handleRequest(writer, header, httpmsghandlerprovider.ShutdownRequest, nil, true)

	var errorResponse = parseJsonRequest(t, writer.Body)
	err := e2managererrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestShutdownStatusNoContent(t *testing.T) {
	log := initLog(t)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	config := configuration.ParseConfiguration()

	var rnibError error
	nbIdentityList := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibError)

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)
	controller.ShutdownHandler(writer, tests.GetHttpRequest())

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestHandleInternalError(t *testing.T) {
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

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)
	err := e2managererrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleCommandAlreadyInProgressError(t *testing.T) {
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
	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)
	err := e2managererrors.NewCommandAlreadyInProgressError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestValidateHeaders(t *testing.T) {
	log := initLog(t)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	config := configuration.ParseConfiguration()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, ranSetupManager)

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
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	config := configuration.ParseConfiguration()
	rmrService := getRmrService(rmrMessengerMock, log)

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, rmrService, readerProvider, writerProvider, config, ranSetupManager)

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
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	config := configuration.ParseConfiguration()
	rmrService := getRmrService(rmrMessengerMock, log)

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, rmrService, readerProvider, writerProvider, config, ranSetupManager)

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
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rmrMessengerMock := &mocks.RmrMessengerMock{}

	config := configuration.ParseConfiguration()
	rmrService := getRmrService(rmrMessengerMock, log)

	writer := httptest.NewRecorder()
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, rmrService, readerProvider, writerProvider, config, ranSetupManager)

	// Invalid json: attribute name without quotes (should be "cause":).
	b := strings.NewReader("{cause:\"protocol:transfer-syntax-error\"")
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2ResetHandler(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	_, ok := rmrService.E2sessions[ranName]
	assert.False(t, ok)

}

func TestHandleErrorResponse(t *testing.T) {
	log := initLog(t)

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rmrMessengerMock := &mocks.RmrMessengerMock{}

	config := configuration.ParseConfiguration()
	rmrService := getRmrService(rmrMessengerMock, log)
	ranSetupManager := managers.NewRanSetupManager(log, getRmrService(rmrMessengerMock, log), rNibWriter.GetRNibWriter)
	controller := NewController(log, rmrService, readerProvider, writerProvider, config, ranSetupManager)

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
