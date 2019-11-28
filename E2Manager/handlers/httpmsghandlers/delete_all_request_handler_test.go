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


package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func setupTest(t *testing.T) (*logger.Logger, *configuration.Configuration, *mocks.RnibReaderMock, *mocks.RnibWriterMock, services.RNibDataService, *mocks.RmrMessengerMock) {
	log := initLog(t)
	config := configuration.ParseConfiguration()

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rnibDataService := services.NewRnibDataService(log, config, readerProvider, writerProvider)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	return log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock
}

func TestHandleBeforeTimerGetListNodebIdsFailedFlow(t *testing.T) {
	log, config, readerMock, _, rnibDataService, rmrMessengerMock := setupTest(t)

	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	rnibErr := &common.ResourceNotFoundError{}
	var nbIdentityList []*entities.NbIdentity
	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibErr)

	expected := &e2managererrors.RnibDbError{}
	_, actual := handler.Handle(nil)
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func TestHandleAfterTimerGetListNodebIdsFailedFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1

	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	rnibErr := &common.ResourceNotFoundError{}
	//Before timer: Disconnected->ShutDown, ShuttingDown->Ignore, Connected->ShuttingDown
	nbIdentityList := createIdentityList()

	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil).Return(nbIdentityList, rnibErr)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3).Return(nil)

	expected := &e2managererrors.RnibDbError{}
	_, actual := handler.Handle(nil)

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func TestHandleSuccessFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	//Before timer: Disconnected->ShutDown, ShuttingDown->Ignore, Connected->ShuttingDown
	nbIdentityList := createIdentityList()
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3).Return(nil)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	nb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3AfterTimer, nil)

	updatedNb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb2AfterTimer).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3AfterTimer).Return(nil)

	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, nil)

	_, actual := handler.Handle(nil)

	assert.Nil(t, actual)
}

func TestHandleSuccessGetNextStatusFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(nil)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, nil)

	updatedNb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1AfterTimer).Return(nil)

	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, nil)

	_, actual := handler.Handle(nil)

	assert.Nil(t, actual)
}

func TestHandleShuttingDownStatusFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, nil)

	updatedNb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1AfterTimer).Return(nil)

	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, nil)

	_, actual := handler.Handle(nil)

	assert.Nil(t, actual)
}

func TestHandleGetNodebFailedFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	//Before timer: Disconnected->ShutDown(will fail), ShuttingDown->Ignore, Connected->ShuttingDown
	nbIdentityList := createIdentityList()
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	errRnib := &common.ResourceNotFoundError{}
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, errRnib)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(errRnib)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3).Return(nil)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	nb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, errRnib)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3AfterTimer, nil)

	updatedNb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb2AfterTimer).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3AfterTimer).Return(nil)

	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, nil)

	_, actual := handler.Handle(nil)

	assert.Nil(t, actual)
}

func TestHandleSaveFailedFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	//Before timer: Disconnected->ShutDown, ShuttingDown->Ignore, Connected->ShuttingDown(will fail)
	nbIdentityList := createIdentityList()
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	errRnib := &common.ResourceNotFoundError{}
	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3).Return(errRnib)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown(will fail)
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	nb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3AfterTimer, nil)

	updatedNb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb2AfterTimer).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3AfterTimer).Return(errRnib)

	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, nil)

	_, actual := handler.Handle(nil)

	assert.Nil(t, actual)
}

func TestHandleSendRmrFailedFlow(t *testing.T) {
	log, config, readerMock, writerMock, rnibDataService, rmrMessengerMock := setupTest(t)

	config.BigRedButtonTimeoutSec = 1
	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	//Before timer: Disconnected->ShutDown, ShuttingDown->Ignore, Connected->ShuttingDown(will fail)
	nbIdentityList := createIdentityList()
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb1).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3).Return(nil)

	//after timer: ShutDown->Ignore, ShuttingDown->ShutDown(will fail)
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1AfterTimer := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	nb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	nb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2AfterTimer, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3AfterTimer, nil)

	updatedNb2AfterTimer := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb3AfterTimer := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("SaveNodeb", mock.Anything, updatedNb2AfterTimer).Return(nil)
	writerMock.On("SaveNodeb", mock.Anything, updatedNb3AfterTimer).Return(nil)

	expected := e2managererrors.NewRmrError()
	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize, "RanName", &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf))).Return(mbuf, expected)

	_, actual := handler.Handle(nil)

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func TestHandleGetListEnbIdsEmptyFlow(t *testing.T) {
	log, config, readerMock, _, rnibDataService, rmrMessengerMock := setupTest(t)

	handler := NewDeleteAllRequestHandler(log, getRmrSender(rmrMessengerMock, log), config, rnibDataService)

	var rnibError error
	nbIdentityList := []*entities.NbIdentity{}

	readerMock.On("GetListNodebIds").Return(nbIdentityList, rnibError)

	_, actual := handler.Handle(nil)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 0)
	assert.Nil(t, actual)
}

func createIdentityList() []*entities.NbIdentity {
	nbIdentity1 := entities.NbIdentity{InventoryName: "RanName_1"}
	nbIdentity2 := entities.NbIdentity{InventoryName: "RanName_2"}
	nbIdentity3 := entities.NbIdentity{InventoryName: "RanName_3"}

	var nbIdentityList []*entities.NbIdentity
	nbIdentityList = append(nbIdentityList, &nbIdentity1)
	nbIdentityList = append(nbIdentityList, &nbIdentity2)
	nbIdentityList = append(nbIdentityList, &nbIdentity3)

	return nbIdentityList
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}

func getRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, &rmrMessenger)
}
