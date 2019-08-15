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
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers"
	"e2mgr/rNibWriter"
	"e2mgr/tests"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	var messageChannel chan<- *models.NotificationResponse

	rnibErr := &common.RNibError{}
	var nbIdentityList []*entities.NbIdentity
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibErr)

	writer := httptest.NewRecorder()
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)
	controller.ShutdownHandler(writer, tests.GetHttpRequest(), nil)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, e2managererrors.NewRnibDbError().Err.Code)
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
	var messageChannel chan<- *models.NotificationResponse

	writer := httptest.NewRecorder()

	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)

	header := &http.Header{}

	controller.handleRequest(writer, header, providers.ShutdownRequest, nil, true, http.StatusNoContent)

	var errorResponse = parseJsonRequest(t, writer.Body)
	err := e2managererrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Err.Code)
	assert.Equal(t, errorResponse.Message, err.Err.Message)
}

func TestShutdownStatusNoContent(t *testing.T){
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
	var messageChannel chan<- *models.NotificationResponse

	var rnibError common.IRNibError
	nbIdentityList := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibError)

	writer := httptest.NewRecorder()
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)
	controller.ShutdownHandler(writer, tests.GetHttpRequest(), nil)

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
	var messageChannel chan<- *models.NotificationResponse

	writer := httptest.NewRecorder()
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)
	err := e2managererrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Err.Code)
	assert.Equal(t, errorResponse.Message, err.Err.Message)
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
	var messageChannel chan<- *models.NotificationResponse
	writer := httptest.NewRecorder()
	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)
	err := e2managererrors.NewCommandAlreadyInProgressError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Err.Code)
	assert.Equal(t, errorResponse.Message, err.Err.Message)
}

func TestValidateHeaders(t *testing.T){
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
	var messageChannel chan<- *models.NotificationResponse

	controller := NewController(log, getRmrService(rmrMessengerMock, log), readerProvider, writerProvider, config, messageChannel)

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
