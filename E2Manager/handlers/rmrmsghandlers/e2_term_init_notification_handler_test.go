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
package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/mock"
	"testing"
)

const e2tInstanceAddress = "10.0.2.15"
const e2tInitPayload = "{\"address\":\"10.0.2.15\", \"fqdn\":\"\"}"

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, E2TermInitNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock) {

	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	routingManagerClientMock := &mocks.RoutingManagerClientMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)

	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}

	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, ranSetupManager, e2tInstancesManagerMock)
	handler := NewE2TermInitNotificationHandler(logger, ranReconnectionManager, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock)

	return logger, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock
}

func initRanLostConnectionTestWithRealE2tInstanceManager(t *testing.T) (*logger.Logger, E2TermInitNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock) {

	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)

	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, ranSetupManager, e2tInstancesManager)
	handler := NewE2TermInitNotificationHandler(logger, ranReconnectionManager, rnibDataService, e2tInstancesManager, routingManagerClientMock)
	return logger, handler, readerMock, writerMock, rmrMessengerMock
}

func TestE2TermInitUnmarshalPayloadFailure(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte("asd")}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitEmptyE2TAddress(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte("{\"address\":\"\"}")}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitGetE2TInstanceFailure(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var e2tInstance *entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("internal error")))
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance")
}

func TestE2TermInitGetE2TInstanceDbFailure(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTestWithRealE2tInstanceManager(t)
	var e2tInstance *entities.E2TInstance
	readerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewInternalError(fmt.Errorf("internal error")))
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2TermInitNewE2TInstance(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, routingManagerClient := initRanLostConnectionTest(t)
	var e2tInstance *entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewResourceNotFoundError("not found"))
	e2tInstance = entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstancesManagerMock.On("AddE2TInstance", e2tInstanceAddress).Return(nil)
	routingManagerClient.On("AddE2TInstance", e2tInstanceAddress).Return(nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	routingManagerClient.AssertCalled(t, "AddE2TInstance", e2tInstanceAddress)
	e2tInstancesManagerMock.AssertCalled(t, "AddE2TInstance", e2tInstanceAddress)
}

func TestE2TermInitNewE2TInstance_RoutingManagerError(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, routingManagerClient := initRanLostConnectionTest(t)
	var e2tInstance *entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, common.NewResourceNotFoundError("not found"))
	e2tInstance = entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstancesManagerMock.On("AddE2TInstance", e2tInstanceAddress).Return(nil)
	routingManagerClient.On("AddE2TInstance", e2tInstanceAddress).Return(fmt.Errorf("error"))
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	routingManagerClient.AssertCalled(t, "AddE2TInstance", e2tInstanceAddress)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddE2TInstance", e2tInstanceAddress)
}

func TestE2TermInitExistingE2TInstanceNoAssociatedRans(t *testing.T) {
	_, handler, _, _, _, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceAddress)
}

func TestE2TermInitHandlerSuccessOneRan(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TermInitHandlerSuccessOneRanShuttingdown(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2TermInitHandlerSuccessOneRan_ToBeDeleted(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.State = entities.ToBeDeleted
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)

	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2TermInitHandlerSuccessOneRan_RoutingManagerFailure(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.State = entities.RoutingManagerFailure
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)

	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	e2tInstancesManagerMock.On("ActivateE2TInstance", e2tInstance).Return(nil)
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TermInitHandlerSuccessOneRan_RoutingManagerFailure_Error(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.State = entities.RoutingManagerFailure
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)

	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	e2tInstancesManagerMock.On("ActivateE2TInstance", e2tInstance).Return(fmt.Errorf(" Error "))
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestE2TermInitHandlerSuccessTwoRans(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error
	var initialNodeb0 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", "test2").Return(initialNodeb1, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName, "test2")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestE2TermInitHandlerSuccessTwoRansSecondRanShutdown(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	var rnibErr error
	var initialNodeb0 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", "test2").Return(initialNodeb1, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName, "test2")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
	writerMock.AssertExpectations(t)
}

func TestE2TermInitHandlerSuccessThreeRansFirstRmrFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, rnibErr)
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Fail = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Fail).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg0, fmt.Errorf("RMR Error"))

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, "test1", "test2", "test3")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	//test1 (before send +1, after failure +1), test2 (0) test3 (0)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	//test1 failure (+1), test2  (0). test3 (0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TermInitHandlerSuccessThreeRansSecondNotFoundFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, common.NewResourceNotFoundError("not found"))
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Success = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Success).Return(rnibErr)

	var argNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb2Success = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb2).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb2Success).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg0, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, "test1", "test2", "test3")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	//test1 (+1), test2 failure (0) test3 (+1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	//test1 success (+1), test2  (0). test3 (+1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestE2TermInitHandlerSuccessThreeRansSecondRnibInternalErrorFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, common.NewInternalError(fmt.Errorf("internal error")))
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Success = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Success).Return(rnibErr)

	var argNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb2Success = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb2).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb2Success).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg0, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, "test1", "test2", "test3")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 2)
	//test1 (+1), test2 failure (0) test3 (0)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	//test1 success (+1), test2  (0). test3 (+1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TermInitHandlerSuccessZeroRans(t *testing.T) {
	_, handler, _, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _  := initRanLostConnectionTest(t)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2TermInitHandlerFailureGetNodebInternalError(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)

	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", "test1").Return(nodebInfo, common.NewInternalError(fmt.Errorf("internal error")))

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, "test1")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}
	handler.Handle(notificationRequest)

	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2TermInitHandlerSuccessTwoRansSecondIsDisconnected(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _ := initRanLostConnectionTest(t)
	var rnibErr error
	var initialNodeb0 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", "test2").Return(initialNodeb1, rnibErr)

	var argNodeb1 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb1).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	e2tInstance := entities.NewE2TInstance(e2tInstanceAddress)
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName, "test2")
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceAddress).Return(e2tInstance, nil)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(e2tInitPayload)}

	handler.Handle(notificationRequest)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}


// TODO: extract to test_utils
func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
