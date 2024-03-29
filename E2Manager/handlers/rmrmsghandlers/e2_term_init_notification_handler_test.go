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

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

const (
	e2tInitPayload     = "{\"address\":\"10.0.2.15\", \"fqdn\":\"\"}"
	e2tInstanceAddress = "10.0.2.15"
	podName            = "podNAme_test"
)

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, E2TermInitNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock) {

	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	routingManagerClientMock := &mocks.RoutingManagerClientMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)

	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}

	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := &mocks.RanAlarmServiceMock{}
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock, ranConnectStatusChangeManager)

	ranDisconnectionManager := managers.NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager)
	handler := NewE2TermInitNotificationHandler(logger, ranDisconnectionManager, e2tInstancesManagerMock, routingManagerClientMock)

	return logger, handler, readerMock, writerMock, e2tInstancesManagerMock, routingManagerClientMock
}

func initRanLostConnectionTestWithRealE2tInstanceManager(t *testing.T) (*logger.Logger, *configuration.Configuration, E2TermInitNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock, managers.RanListManager) {

	logger := initLog(t)
	config := configuration.ParseConfiguration()

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}
	httpClientMock := &mocks.HttpClientMock{}

	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClientMock)
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)

	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager)
	ranDisconnectionManager := managers.NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager)
	handler := NewE2TermInitNotificationHandler(logger, ranDisconnectionManager, e2tInstancesManager, routingManagerClient)
	return logger, config, handler, readerMock, writerMock, httpClientMock, ranListManager
}

func TestE2TermInitUnmarshalPayloadFailure(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte("asd")}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitEmptyE2TAddress(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte("{\"address\":\"\"}")}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitGetE2TInstanceFailure(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	var e2tInstance *entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("internal error")))
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitGetE2TInstanceDbFailure(t *testing.T) {
	_, _, handler, readerMock, writerMock, _, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	var e2tInstance *entities.E2TInstance
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("internal error")))
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestE2TermInitNewE2TInstance(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	var e2tInstance *entities.E2TInstance

	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewResourceNotFoundError("not found"))
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.AddE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	var e2tAddresses []string
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, common.NewResourceNotFoundError(""))

	e2tAddresses = append(e2tAddresses, e2tInstanceAddress)
	writerMock.On("SaveE2TAddresses", e2tAddresses).Return(nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)

	httpClientMock.AssertCalled(t, "Post", url, mock.Anything, mock.Anything)
	writerMock.AssertCalled(t, "SaveE2TInstance", mock.Anything)
	writerMock.AssertCalled(t, "SaveE2TAddresses", e2tAddresses)
}

func TestE2TermInitNewE2TInstance__RoutingManagerError(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	var e2tInstance *entities.E2TInstance

	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewResourceNotFoundError("not found"))

	url := config.RoutingManager.BaseUrl + clients.AddE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("error"))

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 0)
}

func TestE2TermInitExistingE2TInstanceNoAssociatedRans(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock, routingManagerClientMock := initRanLostConnectionTest(t)
	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	var rtmgrErr error
	e2tInstancesManagerMock.On("ResetKeepAliveTimestamp", e2tInstanceAddress).Return(nil)
	routingManagerClientMock.On("AddE2TInstance", e2tInstanceAddress).Return(rtmgrErr, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceAddress)
}

func TestE2TermInitHandlerSuccessOneRan(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	oldNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}

	var rnibErr error
	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, "test_DISCONNECTED").Return(nil)

	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", argNodeb.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	var disconnectedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(disconnectedNodeb, rnibErr)

	var updatedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

func TestE2TermInitHandlerSuccessOneRan_RoutingManagerError(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	oldNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}
	var rnibErr error
	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, "test_DISCONNECTED").Return(nil)

	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", argNodeb.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	var disconnectedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(disconnectedNodeb, rnibErr)

	var updatedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{}, errors.New("error"))

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

func TestE2TermInitHandlerSuccessOneRanShuttingdown(t *testing.T) {
	_, _, handler, readerMock, writerMock, _, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	var rnibErr error

	oldNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}
	var initialNodeb = &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", argNodeb.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestE2TermInitHandlerSuccessOneRan_ToBeDeleted(t *testing.T) {
	_, _, handler, readerMock, writerMock, httpClientMock, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	argNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.State = entities.ToBeDeleted
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)

	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	httpClientMock.AssertNotCalled(t, "Post", mock.Anything, mock.Anything, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestE2TermInitHandlerSuccessTwoRans(t *testing.T) {

	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	test2 := "test2"
	oldNbIdentity1 := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	oldNbIdentity2 := &entities.NbIdentity{InventoryName: test2, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}}
	oldNbIdentityList := []*entities.NbIdentity{oldNbIdentity1, oldNbIdentity2}
	readerMock.On("GetListNodebIds").Return(oldNbIdentityList, nil)

	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}
	var rnibErr error

	//First RAN
	var firstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	var disconnectedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(firstRan, rnibErr).Return(disconnectedFirstRan, rnibErr)

	var updatedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", updatedFirstRan.GetNodeType(), []*entities.NbIdentity{oldNbIdentity1}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	var updatedDisconnectedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedDisconnectedFirstRan.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	//Second RAN
	var secondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	var disconnectedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", test2).Return(secondRan, rnibErr).Return(disconnectedSecondRan, rnibErr)

	var updatedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	updatedSecondRan.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	newNbIdentity2 := &entities.NbIdentity{InventoryName: test2, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}}
	writerMock.On("UpdateNbIdentities", updatedFirstRan.GetNodeType(), []*entities.NbIdentity{oldNbIdentity2}, []*entities.NbIdentity{newNbIdentity2}).Return(nil)

	var updatedDisconnectedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: ""}
	updatedDisconnectedSecondRan.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, test2)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 4)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 2)
	httpClientMock.AssertNumberOfCalls(t, "Post", 2)
}

func TestE2TermInitHandlerSuccessTwoRansSecondRanShutdown(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	oldNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}
	var rnibErr error
	test2 := "test2"

	//First RAN
	var firstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	var disconnectedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(firstRan, rnibErr).Return(disconnectedFirstRan, rnibErr)

	var updatedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", updatedFirstRan.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	var updatedDisconnectedFirstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedDisconnectedFirstRan.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	//Second RAN
	var secondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", test2).Return(secondRan, rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

func TestE2TermInitHandlerSuccessTwoRansFirstNotFoundFailure(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	test2 := "test2"
	oldNbIdentity := &entities.NbIdentity{InventoryName: test2, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}
	var rnibErr error

	//First RAN
	var firstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(firstRan, common.NewResourceNotFoundError("not found"))

	//Second RAN
	var secondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	var disconnectedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", test2).Return(secondRan, rnibErr).Return(disconnectedSecondRan, rnibErr)

	var updatedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	var updatedDisconnectedSecondRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: test2, AssociatedE2TInstanceAddress: ""}
	updatedDisconnectedSecondRan.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	newNbIdentity := &entities.NbIdentity{InventoryName: test2, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", updatedSecondRan.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, test2)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

func TestE2TermInitHandlerSuccessTwoRansFirstRnibInternalErrorFailure(t *testing.T) {
	_, _, handler, readerMock, writerMock, httpClientMock, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	test2 := "test2"

	//First RAN
	var firstRan = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(firstRan, common.NewInternalError(fmt.Errorf("internal error")))

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, test2)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 0)
	httpClientMock.AssertNumberOfCalls(t, "Post", 0)
}

func TestE2TermInitHandlerSuccessZeroRans(t *testing.T) {
	_, handler, _, writerMock, e2tInstancesManagerMock, routingManagerClientMock := initRanLostConnectionTest(t)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	var rtmgrErr error
	e2tInstancesManagerMock.On("ResetKeepAliveTimestamp", e2tInstanceAddress).Return(nil)
	routingManagerClientMock.On("AddE2TInstance", e2tInstanceAddress).Return(rtmgrErr, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestE2TermInitHandlerFailureGetNodebInternalError(t *testing.T) {
	_, handler, readerMock, writerMock, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)

	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", "test1").Return(nodebInfo, common.NewInternalError(fmt.Errorf("internal error")))

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, "test1")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)

	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestE2TermInitHandlerOneRanNoRanInNbIdentityMap(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, _ := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	var rnibErr error
	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	argNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, "test_DISCONNECTED").Return(nil)

	var disconnectedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(disconnectedNodeb, rnibErr)

	var updatedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

func TestE2TermInitHandlerOneRanUpdateNbIdentitiesFailure(t *testing.T) {
	_, config, handler, readerMock, writerMock, httpClientMock, ranListManager := initRanLostConnectionTestWithRealE2tInstanceManager(t)

	oldNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{oldNbIdentity}, nil)
	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}

	var rnibErr error
	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, "test_DISCONNECTED").Return(nil)

	newNbIdentity := &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	writerMock.On("UpdateNbIdentities", argNodeb.GetNodeType(), []*entities.NbIdentity{oldNbIdentity}, []*entities.NbIdentity{newNbIdentity}).Return(common.NewInternalError(fmt.Errorf("internal error")))

	var disconnectedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: "10.0.2.15"}
	readerMock.On("GetNodeb", RanName).Return(disconnectedNodeb, rnibErr)

	var updatedNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNodeb.StatusUpdateTimeStamp = uint64(time.Now().UnixNano())
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress, podName)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil).Return(e2tInstance, nil)
	writerMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	url := config.RoutingManager.BaseUrl + clients.DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, mock.Anything, mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	httpClientMock.AssertNumberOfCalls(t, "Post", 1)
}

// TODO: extract to test_utils
func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	InfoLevel := int8(3)
	log, err := logger.InitLogger(InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
