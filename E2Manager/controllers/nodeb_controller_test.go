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
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const x2SetupMessageType = "x2-setup"

func TestNewRequestController(t *testing.T) {
	rnibReaderProvider := func() reader.RNibReader {
		return &mocks.RnibReaderMock{}
	}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	assert.NotNil(t, NewNodebController(&logger.Logger{}, &services.RmrService{}, rnibReaderProvider, rnibWriterProvider))
}

func TestHandleRequestSuccess(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#nodeb_controller_test.TestHandleRequestSuccess - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock :=&mocks.RmrMessengerMock{}
	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize,"RanName", &tests.DummyPayload, &tests.DummyXAction)

	rmrMessengerMock.On("SendMsg",
		mock.AnythingOfType(fmt.Sprintf("%T", mbuf)),
		tests.MaxMsgSize).Return(mbuf, nil)

	writer := httptest.NewRecorder()

	handleRequest(writer, log, rmrMessengerMock, tests.GetHttpRequest(), x2SetupMessageType)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
}

func TestHandleRequestFailure_InvalidRequestDetails(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#nodeb_controller_test.TestHandleRequestFailure - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	var mbuf *rmrCgo.MBuf

	rmrMessengerMock.On("SendMsg",
		mock.AnythingOfType(fmt.Sprintf("%T", mbuf)),
		tests.MaxMsgSize).Return(mbuf, errors.New("test failure"))

	writer := httptest.NewRecorder()

	handleRequest(writer, log, rmrMessengerMock, tests.GetInvalidRequestDetails(), x2SetupMessageType)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
}

func TestHandleRequestFailure_InvalidMessageType(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#nodeb_controller_test.TestHandleRequestFailure - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	var mbuf *rmrCgo.MBuf

	rmrMessengerMock.On("SendMsg",
		mock.AnythingOfType(fmt.Sprintf("%T", mbuf)),
		tests.MaxMsgSize).Return(mbuf, errors.New("test failure"))

	writer := httptest.NewRecorder()

	handleRequest(writer, log, rmrMessengerMock, tests.GetInvalidMessageType(), "dummy")
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)
}

func TestHandleHealthCheckRequest(t *testing.T) {
	rc := NewNodebController(nil, nil, nil, nil)
	writer := httptest.NewRecorder()
	rc.HandleHealthCheckRequest(writer, nil, nil)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
}

func handleRequest(writer *httptest.ResponseRecorder, log *logger.Logger, rmrMessengerMock *mocks.RmrMessengerMock,
	request *http.Request, messageType string) {
	rmrService := getRmrService(rmrMessengerMock, log)
	params := []httprouter.Param{{Key: "messageType", Value: messageType}}

	var nodebInfo *entities.NodebInfo
	var nbIdentity *entities.NbIdentity

	rnibWriterMock := mocks.RnibWriterMock{}
	rnibWriterMock.On("SaveNodeb",
		mock.AnythingOfType(fmt.Sprintf("%T", nbIdentity)),
		mock.AnythingOfType(fmt.Sprintf("%T", nodebInfo))).Return(nil)

	rnibReaderProvider := func() reader.RNibReader {
		return &mocks.RnibReaderMock{}
	}

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &rnibWriterMock
	}

	NewNodebController(log, rmrService, rnibReaderProvider, rnibWriterProvider).HandleRequest(writer, request, params)
}

func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rnibReaderProvider := func() reader.RNibReader {
		return &mocks.RnibReaderMock{}
	}

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}
	nManager := managers.NewNotificationManager(rnibReaderProvider, rnibWriterProvider)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	messageChannel := make(chan *models.NotificationResponse)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, E2Sessions, nManager, messageChannel)
}

func executeGetNodeb(logger *logger.Logger, writer *httptest.ResponseRecorder, rnibReaderProvider func() reader.RNibReader) {
	req, _ := http.NewRequest("GET", "/nodeb", nil)

	params := []httprouter.Param{{Key: "ranName", Value: "testNode"}}

	NewNodebController(logger, nil, rnibReaderProvider, nil).GetNodeb(writer, req, params)
}

func TestNodebController_GetNodeb_Success(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err != nil {
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodeb_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	var rnibError common.IRNibError
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
	rnibError := common.NewResourceNotFoundError(errors.Errorf("#reader.GetNodeb - responding node %s not found", "testNode"))
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
	NewNodebController(logger, nil, rnibReaderProvider, nil ).GetNodebIdList(writer,req,nil)
}

func TestNodebController_GetNodebIdList_Success(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)

	if err!=nil{
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}
	var rnibError common.IRNibError

	enbList := []*entities.NbIdentity{&entities.NbIdentity{InventoryName:"test1", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId1",NbId: "nbId1"}}}
	gnbList := []*entities.NbIdentity{&entities.NbIdentity{InventoryName:"test2", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId2",NbId: "nbId2"}}}

	rnibReaderMock.On("GetListEnbIds").Return(&enbList, rnibError)
	rnibReaderMock.On("GetListGnbIds").Return(&gnbList, rnibError)


	rnibReaderProvider:= func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
	bodyBytes, err := ioutil.ReadAll(writer.Body)
	assert.Equal(t, "[{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}},{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}}]",string(bodyBytes) )
}

func TestNodebController_GetNodebIdList_EmptyList(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)

	if err!=nil{
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_EmptyList - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}

	var rnibError common.IRNibError
	enbList := []*entities.NbIdentity{}
	gnbList := []*entities.NbIdentity{}

	rnibReaderMock.On("GetListEnbIds").Return(&enbList, rnibError)
	rnibReaderMock.On("GetListGnbIds").Return(&gnbList, rnibError)


	rnibReaderProvider:= func() reader.RNibReader {
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
	var enbList *[]*entities.NbIdentity
	rnibReaderMock.On("GetListEnbIds").Return(enbList, rnibError)

	rnibReaderProvider := func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusInternalServerError)
}

func TestNodebController_GetNodebIdList_Success_One(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)

	if err!=nil{
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}
	var rnibError common.IRNibError

	enbList := []*entities.NbIdentity{}
	gnbList := []*entities.NbIdentity{&entities.NbIdentity{InventoryName:"test2", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId2",NbId: "nbId2"}}}

	rnibReaderMock.On("GetListEnbIds").Return(&enbList, rnibError)
	rnibReaderMock.On("GetListGnbIds").Return(&gnbList, rnibError)


	rnibReaderProvider:= func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
	bodyBytes, err := ioutil.ReadAll(writer.Body)
	assert.Equal(t, "[{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}}]",string(bodyBytes) )
}

func TestNodebController_GetNodebIdList_Success_Many(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)

	if err!=nil{
		t.Errorf("#nodeb_controller_test.TestNodebController_GetNodebIdList_Success - failed to initialize logger, error: %s", err)
	}

	writer := httptest.NewRecorder()

	rnibReaderMock := mocks.RnibReaderMock{}
	var rnibError common.IRNibError

	enbList := []*entities.NbIdentity{&entities.NbIdentity{InventoryName:"test1", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId1",NbId: "nbId1"}}}
	gnbList := []*entities.NbIdentity{&entities.NbIdentity{InventoryName:"test2", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId2",NbId: "nbId2"}}, {InventoryName:"test3", GlobalNbId: &entities.GlobalNbId{PlmnId:"plmnId3",NbId: "nbId3"}}}

	rnibReaderMock.On("GetListEnbIds").Return(&enbList, rnibError)
	rnibReaderMock.On("GetListGnbIds").Return(&gnbList, rnibError)


	rnibReaderProvider:= func() reader.RNibReader {
		return &rnibReaderMock
	}

	executeGetNodebIdList(logger, writer, rnibReaderProvider)
	assert.Equal(t, writer.Result().StatusCode, http.StatusOK)
	bodyBytes, err := ioutil.ReadAll(writer.Body)
	assert.Equal(t, "[{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}},{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}},{\"inventoryName\":\"test3\",\"globalNbId\":{\"plmnId\":\"plmnId3\",\"nbId\":\"nbId3\"}}]",string(bodyBytes) )
}