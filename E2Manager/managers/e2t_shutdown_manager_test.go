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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package managers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const E2TAddress3 = "10.10.2.17:9800"

func initE2TShutdownManagerTest(t *testing.T) (*E2TShutdownManager, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock, *mocks.RmrMessengerMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, E2TInstanceDeletionTimeoutMs: 15000}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	associationManager := NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, log)
	ranSetupManager := NewRanSetupManager(log, rmrSender, rnibDataService)
	shutdownManager := NewE2TShutdownManager(log, config, rnibDataService, e2tInstancesManager, associationManager, ranSetupManager)

	return shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock
}

func TestShutdownSuccess1OutOf3Instances(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2,e2tInstance3}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test1", "test5")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	readerMock.On("GetE2TInstance", E2TAddress2).Return(e2tInstance2, nil)
	e2tInstance2updated := *e2tInstance2
	e2tInstance2updated.AssociatedRanList = append(e2tInstance2updated.AssociatedRanList, "test1", "test5")
	writerMock.On("SaveE2TInstance", &e2tInstance2updated).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb1new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)
	nodeb5new := *nodeb5
	nodeb5new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb5new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb5new).Return(nil)

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction1 := []byte("test1")
	msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test1", &payload, &xaction1)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg1, nil)
	xaction5 := []byte("test5")
	msg5 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test5", &payload, &xaction5)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg5, nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestShutdownSuccess1InstanceWithoutRans(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownSuccess1Instance2Rans(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1", "test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	e2tInstance1updated := *e2tInstance1
	e2tInstance1updated.State = entities.ToBeDeleted
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance1updated}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)
	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = ""
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}



func TestShutdownE2tInstanceAlreadyBeingDeleted(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance1.AssociatedRanList = []string{"test1"}
	e2tInstance1.DeletionTimestamp = time.Now().UnixNano()

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownFailureMarkInstanceAsToBeDeleted(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(e2managererrors.NewRnibDbError())

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownFailureReassociatingInMemoryNodebNotFound(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	var nodeb1 *entities.NodebInfo
	readerMock.On("GetNodeb", "test1").Return(nodeb1, common.NewResourceNotFoundError("for tests"))
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test2")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2}).Return(nil)

	readerMock.On("GetE2TInstance", E2TAddress2).Return(e2tInstance2, nil)
	e2tInstance2updated := *e2tInstance2
	e2tInstance2updated.AssociatedRanList = append(e2tInstance2updated.AssociatedRanList, "test2")
	writerMock.On("SaveE2TInstance", &e2tInstance2updated).Return(nil)

	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb2new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	nodeb2connected := *nodeb2
	nodeb2connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb2connected.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2connected).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction2 := []byte("test2")
	msg2 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test2", &payload, &xaction2)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg2, nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestShutdownFailureReassociatingInMemoryGetNodebError(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	var nodeb1 *entities.NodebInfo
	readerMock.On("GetNodeb", "test1").Return(nodeb1, common.NewInternalError(fmt.Errorf("for tests")))

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownFailureRoutingManagerError(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2,e2tInstance3}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test1", "test5")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)

	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, nil)
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.RoutingManagerFailure })).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownFailureInClearNodebsAssociation(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1", "test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	e2tInstance1updated := *e2tInstance1
	e2tInstance1updated.State = entities.ToBeDeleted
	readerMock.On("GetE2TInstances", []string{E2TAddress}).Return([]*entities.E2TInstance{&e2tInstance1updated}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(common.NewInternalError(fmt.Errorf("for tests")))

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownFailureInRmr(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2,e2tInstance3}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test1", "test5")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	readerMock.On("GetE2TInstance", E2TAddress2).Return(e2tInstance2, nil)
	e2tInstance2updated := *e2tInstance2
	e2tInstance2updated.AssociatedRanList = []string{"test3", "test1", "test5"}
	writerMock.On("SaveE2TInstance", &e2tInstance2updated).Return(nil)

	nodeb1reassigned := *nodeb1
	nodeb1reassigned.AssociatedE2TInstanceAddress = E2TAddress2
	writerMock.On("UpdateNodebInfo", &nodeb1reassigned).Return(nil)
	nodeb5reassigned := *nodeb5
	nodeb5reassigned.AssociatedE2TInstanceAddress = E2TAddress2
	writerMock.On("UpdateNodebInfo", &nodeb5reassigned).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb1new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)
	nodeb5new := *nodeb5
	nodeb5new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb5new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb5new).Return(nil)

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	nodeb1connected.ConnectionAttempts = 0
	//nodeb1connected.E2ApplicationProtocol = entities.E2ApplicationProtocol_X2_SETUP_REQUEST
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	nodeb5connected.ConnectionAttempts = 0
	//nodeb5connected.E2ApplicationProtocol = entities.E2ApplicationProtocol_X2_SETUP_REQUEST
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction1 := []byte("test1")
	msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test1", &payload, &xaction1)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg1, common.NewInternalError(fmt.Errorf("for test")))
	xaction2 := []byte("test5")
	msg2 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test5", &payload, &xaction2)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg2, common.NewInternalError(fmt.Errorf("for test")))

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestShutdownFailureDbErrorInAsociateAndSetupNodebs(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2,e2tInstance3}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test1", "test5")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	readerMock.On("GetE2TInstance", E2TAddress2).Return(e2tInstance2, nil)
	e2tInstance2updated := *e2tInstance2
	e2tInstance2updated.AssociatedRanList = []string{"test3", "test1", "test5"}
	writerMock.On("SaveE2TInstance", &e2tInstance2updated).Return(nil)

	nodeb1reassigned := *nodeb1
	nodeb1reassigned.AssociatedE2TInstanceAddress = E2TAddress2
	writerMock.On("UpdateNodebInfo", &nodeb1reassigned).Return(common.NewInternalError(fmt.Errorf("for tests")))

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestShutdownSuccess1OutOf3InstancesStateIsRoutingManagerFailure(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock, rmrMessengerMock := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress)
	e2tInstance1.State = entities.RoutingManagerFailure
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	readerMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{e2tInstance2,e2tInstance3}, nil)

	e2tDataList := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress2, "test1", "test5")}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, e2tDataList)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	readerMock.On("GetE2TInstance", E2TAddress2).Return(e2tInstance2, nil)
	e2tInstance2updated := *e2tInstance2
	e2tInstance2updated.AssociatedRanList = []string{"test3", "test1", "test5"}
	writerMock.On("SaveE2TInstance", &e2tInstance2updated).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb1new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)
	nodeb5new := *nodeb5
	nodeb5new.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5new.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodeb5new.ConnectionAttempts = 1
	writerMock.On("UpdateNodebInfo", &nodeb5new).Return(nil)

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	//nodeb1connected.E2ApplicationProtocol = entities.E2ApplicationProtocol_X2_SETUP_REQUEST
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = E2TAddress2
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	//nodeb5connected.E2ApplicationProtocol = entities.E2ApplicationProtocol_X2_SETUP_REQUEST
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	xaction1 := []byte("test1")
	msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test1", &payload, &xaction1)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg1, nil)
	xaction5 := []byte("test5")
	msg5 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), "test5", &payload, &xaction5)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg5, nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}