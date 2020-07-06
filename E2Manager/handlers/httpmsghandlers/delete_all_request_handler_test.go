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
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const E2TAddress = "10.0.2.15:8989"
const BaseRMUrl = "http://10.10.2.15:12020/routingmanager"

func setupDeleteAllRequestHandlerTest(t *testing.T) (*DeleteAllRequestHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.HttpClientMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibWriter: configuration.RnibWriterConfig{ StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE"}}
	config.BigRedButtonTimeoutSec = 1
	config.RoutingManager.BaseUrl = BaseRMUrl

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := getRmrSender(rmrMessengerMock, log)

	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)

	ranListManager := managers.NewRanListManager(log)
	ranAlarmService := services.NewRanAlarmService(log, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(log, rnibDataService,ranListManager, ranAlarmService)

	handler := NewDeleteAllRequestHandler(log, rmrSender, config, rnibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager)
	return handler, readerMock, writerMock, rmrMessengerMock, httpClientMock
}

func mapE2TAddressesToE2DataList(e2tAddresses []string) models.RoutingManagerE2TDataList {
	e2tDataList := make(models.RoutingManagerE2TDataList, len(e2tAddresses))

	for i, v := range e2tAddresses {
		e2tDataList[i] = models.NewRoutingManagerE2TData(v)
	}

	return e2tDataList
}

func mockHttpClientDissociateAllRans(httpClientMock *mocks.HttpClientMock, e2tAddresses []string, ok bool) {
	data := mapE2TAddressesToE2DataList(e2tAddresses)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := BaseRMUrl + clients.DissociateRanE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))

	var status int
	if ok {
		status = http.StatusOK
	} else {
		status = http.StatusBadRequest
	}
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: status, Body: respBody}, nil)
}

func TestGetE2TAddressesFailure(t *testing.T) {
	h, readerMock, _, _, _ := setupDeleteAllRequestHandlerTest(t)
	readerMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
}

func TestOneRanGetE2TAddressesEmptyList(t *testing.T) {
	h, readerMock, writerMock, _, _ := setupDeleteAllRequestHandlerTest(t)
	readerMock.On("GetE2TAddresses").Return([]string{}, nil)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb1).Return(nil)
	_, err := h.Handle(nil)
	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestTwoRansGetE2TAddressesEmptyListOneGetNodebFailure(t *testing.T) {
	h, readerMock, writerMock, _, _ := setupDeleteAllRequestHandlerTest(t)
	readerMock.On("GetE2TAddresses").Return([]string{}, nil)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}, {InventoryName: "RanName_2"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1
	updatedNb1.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNb1, "RAN_CONNECTION_STATUS_CHANGE", "RanName_1_DISCONNECTED").Return(nil)

	nodeb1NotAssociated := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	writerMock.On("UpdateNodebInfo", nodeb1NotAssociated).Return(nil)

	var nb2 *entities.NodebInfo
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t,&e2managererrors.RnibDbError{}, err)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 0)
	readerMock.AssertExpectations(t)
}

func TestUpdateNodebInfoOnConnectionStatusInversionFailure(t *testing.T) {
	h, readerMock, writerMock, _, _ := setupDeleteAllRequestHandlerTest(t)
	readerMock.On("GetE2TAddresses").Return([]string{}, nil)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}, {InventoryName: "RanName_2"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1
	updatedNb1.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNb1, "RAN_CONNECTION_STATUS_CHANGE", "RanName_1_DISCONNECTED").Return(common.NewInternalError(errors.New("error")))

	_, err := h.Handle(nil)

	assert.IsType(t,&e2managererrors.RnibDbError{}, err)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
	readerMock.AssertExpectations(t)
}

func TestTwoRansGetE2TAddressesEmptyListOneUpdateNodebInfoFailure(t *testing.T) {
	h, readerMock, writerMock, _, _ := setupDeleteAllRequestHandlerTest(t)
	readerMock.On("GetE2TAddresses").Return([]string{}, nil)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}, {InventoryName: "RanName_2"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb1).Return(nil)

	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	updatedNb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb2).Return(common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t,&e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanDissociateSucceedsTryShuttingDownFailure(t *testing.T) {
	h, readerMock, writerMock, _, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanDissociateFailsTryShuttingDownFailure(t *testing.T) {
	h, readerMock, writerMock, _, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}

	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, false)
	nbIdentityList := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanWithStateShutDown(t *testing.T) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)

	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1"}}
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, nil)

	_, err := h.Handle(nil)

	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanShutDown(t *testing.T) {
	h, readerMock, writerMock, _, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	nodeb1NotAssociated := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("UpdateNodebInfo", nodeb1NotAssociated).Return(nil)

	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)

	_, err := h.Handle(nil)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanTryShuttingDownSucceedsClearFails(t *testing.T) {
	h, readerMock, writerMock, _, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1
	updatedNb1.ConnectionStatus = entities.ConnectionStatus_SHUTTING_DOWN
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNb1, "RAN_CONNECTION_STATUS_CHANGE", "RanName_1_DISCONNECTED").Return(nil)

	nodeb1NotAssociated := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", nodeb1NotAssociated).Return(nil)

	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{}, common.NewInternalError(errors.New("error")))
	_, err := h.Handle(nil)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanTryShuttingDownUpdateNodebError(t *testing.T) {
	h, readerMock, writerMock, _, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1
	updatedNb1.ConnectionStatus = entities.ConnectionStatus_SHUTTING_DOWN
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNb1, "RAN_CONNECTION_STATUS_CHANGE", "RanName_1_DISCONNECTED").Return(nil)

	nodeb1NotAssociated := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", nodeb1NotAssociated).Return(common.NewInternalError(errors.New("error")))

	_, err := h.Handle(nil)

	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestOneRanTryShuttingDownSucceedsClearSucceedsRmrSendFails(t *testing.T) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1
	updatedNb1.ConnectionStatus = entities.ConnectionStatus_SHUTTING_DOWN
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNb1, "RAN_CONNECTION_STATUS_CHANGE", "RanName_1_DISCONNECTED").Return(nil)

	nodeb1NotAssociated := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", nodeb1NotAssociated).Return(nil)

	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1"}}
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, e2managererrors.NewRmrError())
	_, err := h.Handle(nil)
	assert.IsType(t, &e2managererrors.RmrError{}, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func testTwoRansTryShuttingDownSucceedsClearSucceedsRmrSucceedsAllRansAreShutdown(t *testing.T, partial bool) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, !partial)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}, {InventoryName: "RanName_2"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1", "RanName_2"}}
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, nil)
	resp, err := h.Handle(nil)
	assert.Nil(t, err)

	if partial {
		assert.IsType(t, &models.RedButtonPartialSuccessResponseModel{}, resp)
	} else {
		assert.Nil(t, resp)
	}

	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestTwoRansTryShuttingDownSucceedsClearSucceedsRmrSucceedsAllRansAreShutdownSuccess(t *testing.T) {
	testTwoRansTryShuttingDownSucceedsClearSucceedsRmrSucceedsAllRansAreShutdown(t, false)
}

func TestTwoRansTryShuttingDownSucceedsClearSucceedsRmrSucceedsAllRansAreShutdownPartialSuccess(t *testing.T) {
	testTwoRansTryShuttingDownSucceedsClearSucceedsRmrSucceedsAllRansAreShutdown(t, true)
}

func TestOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDownUpdateFailure(t *testing.T) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	updatedNb1 := *nb1 //&entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, AssociatedE2TInstanceAddress: E2TAddress}
	writerMock.On("UpdateNodebInfo", &updatedNb1).Return(nil)

	nodeb1NotAssociated := *nb1
	nodeb1NotAssociated.AssociatedE2TInstanceAddress = ""
	nodeb1NotAssociated.ConnectionStatus = entities.ConnectionStatus_SHUTTING_DOWN
	writerMock.On("UpdateNodebInfo", &nodeb1NotAssociated).Return(nil)

	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1"}}
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, nil)

	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	readerMock.On("GetNodeb", "RanName_1").Return(updatedNb1, nil)

	updatedNb2 := *nb1//&entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	updatedNb2.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	updatedNb2.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb2).Return(common.NewInternalError(errors.New("error")))

	_, err := h.Handle(nil)

	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func testOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDown(t *testing.T, partial bool) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, !partial)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb1).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1"}}
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, nil)

	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	readerMock.On("GetNodeb", "RanName_1").Return(updatedNb1, nil)
	updatedNb2 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb2).Return(nil)
	_, err := h.Handle(nil)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 3)
}

func TestOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDownSuccess (t *testing.T) {
	testOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDown(t, false)
}

func TestOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDownPartialSuccess (t *testing.T) {
	testOneRanTryShuttingDownSucceedsClearSucceedsRmrSucceedsRanStatusIsShuttingDown(t, true)
}

func TestSuccessTwoE2TInstancesSixRans(t *testing.T) {
	h, readerMock, writerMock, rmrMessengerMock, httpClientMock := setupDeleteAllRequestHandlerTest(t)
	e2tAddresses := []string{E2TAddress, E2TAddress2}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	mockHttpClientDissociateAllRans(httpClientMock, e2tAddresses, true)
	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1"}, {InventoryName: "RanName_2"}, {InventoryName: "RanName_3"}, {InventoryName: "RanName_4"}, {InventoryName: "RanName_5"}, {InventoryName: "RanName_6"}}
	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)

	updatedNb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb1).Return(nil)
	updatedNb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb2).Return(nil)
	updatedNb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb3).Return(nil)
	updatedNb4 := &entities.NodebInfo{RanName: "RanName_4", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb4).Return(nil)
	updatedNb5 := &entities.NodebInfo{RanName: "RanName_5", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb5).Return(nil)
	updatedNb6 := &entities.NodebInfo{RanName: "RanName_6", ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	writerMock.On("UpdateNodebInfo", updatedNb6).Return(nil)

	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"RanName_1", "RanName_2", "RanName_3"}}
	e2tInstance2 := entities.E2TInstance{Address: E2TAddress2, AssociatedRanList: []string{"RanName_4", "RanName_5", "RanName_6"}}
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{&e2tInstance, &e2tInstance2}, nil)
	updatedE2tInstance := e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	updatedE2tInstance2 := e2tInstance2
	updatedE2tInstance2.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)
	writerMock.On("SaveE2TInstance", &updatedE2tInstance2).Return(nil)

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	mbuf := rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, nil)

	readerMock.On("GetListNodebIds").Return(nbIdentityList, nil)
	readerMock.On("GetNodeb", "RanName_1").Return(updatedNb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(updatedNb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(updatedNb3, nil)
	readerMock.On("GetNodeb", "RanName_4").Return(updatedNb4, nil)
	readerMock.On("GetNodeb", "RanName_5").Return(updatedNb5, nil)
	readerMock.On("GetNodeb", "RanName_6").Return(updatedNb6, nil)

	updatedNb1AfterTimer := *updatedNb1
	updatedNb1AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb1AfterTimer).Return(nil)
	updatedNb2AfterTimer := *updatedNb2
	updatedNb2AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb2AfterTimer).Return(nil)
	updatedNb3AfterTimer := *updatedNb3
	updatedNb3AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb3AfterTimer).Return(nil)
	updatedNb4AfterTimer := *updatedNb4
	updatedNb4AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb4AfterTimer).Return(nil)
	updatedNb5AfterTimer := *updatedNb5
	updatedNb5AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb5AfterTimer).Return(nil)
	updatedNb6AfterTimer := *updatedNb6
	updatedNb6AfterTimer.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNb6AfterTimer).Return(nil)
	_, err := h.Handle(nil)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf, true)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 18)
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}

func getRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}
