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
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequestController(t *testing.T) {
	rnibReaderProvider := func() reader.RNibReader {
		return &mocks.RnibReaderMock{}
	}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	assert.NotNil(t, NewNodebController(&logger.Logger{}, &services.RmrService{}, rnibReaderProvider, rnibWriterProvider))
}

func TestHandleHealthCheckRequest(t *testing.T) {
	rc := NewNodebController(nil, nil, nil, nil)
	writer := httptest.NewRecorder()
	rc.HandleHealthCheckRequest(writer, nil)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
}

func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, make(sessions.E2Sessions), messageChannel)
}

func executeGetNodeb(logger *logger.Logger, writer *httptest.ResponseRecorder, rnibReaderProvider func() reader.RNibReader) {
	req, _ := http.NewRequest("GET", "/nodeb", nil)
	req = mux.SetURLVars(req, map[string]string{"ranName": "testNode"})

	NewNodebController(logger, nil, rnibReaderProvider, nil).GetNodeb(writer, req)
}

func TestNodebController_GetNodeb_Success(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodeb_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	var rnibError error
	rnibReaderMock.On("GetNodeb", "testNode").Return(&entities.NodebInfo{}, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodeb(log, writer, rnibReaderProvider)

	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
}

func TestNodebController_GetNodeb_NotFound(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodeb_NotFound - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}
	rnibError := common.NewResourceNotFoundErrorf("#reader.GetNodeb - responding node %s not found", "testNode")
	var nodebInfo *entities.NodebInfo
	rnibReaderMock.On("GetNodeb", "testNode").Return(nodebInfo, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodeb(log, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusNotFound)
}

func TestNodebController_GetNodeb_InternalError(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodeb_InternalError - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	rnibError := common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error"))
	var nodebInfo *entities.NodebInfo
	rnibReaderMock.On("GetNodeb", "testNode").Return(nodebInfo, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodeb(log, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusInternalServerError)
}

func executeGetNodebIdList(logger *logger.Logger, writer *httptest.ResponseRecorder, rnibReaderProvider func() reader.RNibReader) {
	req, _ := http.NewRequest("GET", "/nodeb-ids", nil)
	NewNodebController(logger, nil, rnibReaderProvider, nil ).GetNodebIdList(writer,req)
}

func TestNodebController_GetNodebIdList_Success(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}
	var rnibError error

	nbList := []*entities.NbIdentity{
		{InventoryName:"test1", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId1",NbId: "nbId1"}},
		{InventoryName:"test2", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId2",NbId: "nbId2"}},
		{InventoryName:"test3", GlobalNbId: &entities.GlobalNbId{PlmnId:"",NbId: ""}},
	}
	rnibReaderMock.On("GetListNodebIds").Return(nbList, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
	bodyBytes, err := ioutil.ReadAll(writer.Body)
	assert.Equal(t, "[{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}},{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}},{\"inventoryName\":\"test3\",\"globalNbId\":{}}]",string(bodyBytes) )
}

func TestNodebController_GetNodebIdList_EmptyList(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_EmptyList - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	var rnibError error
	nbList := []*entities.NbIdentity{}
	rnibReaderMock.On("GetListNodebIds").Return(nbList, rnibError)


	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(log, writer, rnibReaderProvider)

	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
	bodyBytes, err := ioutil.ReadAll(writer.Body)
	assert.Equal(t, "[]", string(bodyBytes))
}

func TestNodebController_GetNodebIdList_InternalError(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_InternalError - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	rnibError := common.NewInternalError(errors.New("#reader.GetEnbIdList - Internal Error"))
	var nbList []*entities.NbIdentity
	rnibReaderMock.On("GetListNodebIds").Return(nbList, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusInternalServerError)
}