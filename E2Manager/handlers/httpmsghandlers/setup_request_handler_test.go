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
package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

const E2TAddress = "10.0.2.15:8989"
const RanName = "test"

func initSetupRequestTest(t *testing.T, protocol entities.E2ApplicationProtocol) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *SetupRequestHandler, *mocks.E2TInstancesManagerMock, *mocks.RanSetupManagerMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	ranSetupManagerMock := &mocks.RanSetupManagerMock{}
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	handler := NewSetupRequestHandler(log, rnibDataService, ranSetupManagerMock, protocol, e2tInstancesManagerMock)

	return readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock
}

func TestX2SetupHandleNoPortError(t *testing.T) {
	readerMock, _, handler, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{"127.0.0.1", 0, RanName,}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleNoRanNameError(t *testing.T) {
	readerMock, _, handler, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanIp: "127.0.0.1"}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleNoIpError(t *testing.T) {
	readerMock, _, handler, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanName: RanName}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleInvalidIpError(t *testing.T) {
	readerMock, _, handler, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanName: RanName, RanIp: "invalid ip"}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleGetNodebFailure(t *testing.T) {
	readerMock, _, handler, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	rnibErr := &common.ValidationError{}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)

	sr := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
}

func TestSetupNewRanSelectE2TInstancesDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewRnibDbError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "AssociateRan")
	writerMock.AssertNotCalled(t, "SaveNodeb")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanSelectE2TInstancesNoInstances(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewE2TInstanceAbsenceError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "AssociateRan")
	writerMock.AssertNotCalled(t, "SaveNodeb")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanAssociateRanFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(e2managererrors.NewRnibDbError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	writerMock.AssertNotCalled(t, "SaveNodeb")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanSaveNodebFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST, E2TAddress)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(common.NewInternalError(fmt.Errorf("")))
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanSetupDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST, E2TAddress)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", nodebInfo, entities.ConnectionStatus_CONNECTING).Return(e2managererrors.NewRnibDbError())
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
}

func TestSetupNewRanSetupRmrError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST, E2TAddress)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", nodebInfo, entities.ConnectionStatus_CONNECTING).Return(e2managererrors.NewRmrError())
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RmrError{}, err)
}

func TestSetupNewRanSetupSuccess(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST, E2TAddress)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", nodebInfo, entities.ConnectionStatus_CONNECTING).Return(nil)
	_, err := handler.Handle(setupRequest)
	assert.Nil(t, err)
}

func TestX2SetupExistingRanShuttingDown(t *testing.T) {
	readerMock, _, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}, nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.WrongStateError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestEndcSetupExistingRanShuttingDown(t *testing.T) {
	readerMock, _, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}, nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.WrongStateError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewRnibDbError())
	updatedNb := *nb
	updatedNb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectNoInstanceError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewE2TInstanceAbsenceError())
	updatedNb := *nb
	updatedNb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectNoInstanceErrorUpdateFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewE2TInstanceAbsenceError())
	updatedNb := *nb
	updatedNb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectErrorDisconnected(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:"", ConnectionStatus:entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, e2managererrors.NewE2TInstanceAbsenceError())
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceAssociateRanFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(e2managererrors.NewRnibDbError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceAssociateRanSucceedsUpdateNodebFails(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionAttempts = 0
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceExecuteSetupFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionAttempts = 0
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(e2managererrors.NewRnibDbError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
}

func TestSetupExistingRanWithoutAssocE2TInstanceSuccess(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:""}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AssociateRan", RanName, E2TAddress).Return(nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionAttempts = 0
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(nil)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.Nil(t, err)
}

func TestSetupExistingRanWithAssocE2TInstanceUpdateNodebFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AssociateRan")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingConnectedRanWithAssocE2TInstanceSuccess(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", RanName).Return(nb , nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTED).Return(nil)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.Nil(t, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AssociateRan")
}