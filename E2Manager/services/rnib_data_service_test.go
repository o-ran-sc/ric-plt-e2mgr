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

package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

func setupRnibDataServiceTest(t *testing.T) (*rNibDataService, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	return setupRnibDataServiceTestWithMaxAttempts(t, 3)
}

func setupRnibDataServiceTestWithMaxAttempts(t *testing.T, maxAttempts int) (*rNibDataService, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	DebugLevel := int8(4)
	logger, err := logger.InitLogger(DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: maxAttempts, RnibWriter: configuration.RnibWriterConfig{RanManipulationMessageChannel: "RAN_MANIPULATION", StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE"}}

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := NewRnibDataService(logger, config, readerMock, writerMock)
	assert.NotNil(t, rnibDataService)

	return rnibDataService, readerMock, writerMock
}

func TestSuccessfulSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("SaveNodeb", nodebInfo).Return(nil)

	rnibDataService.SaveNodeb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveNodeb", nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 3)
}

func TestNonConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection failure")}
	writerMock.On("SaveNodeb", nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestSuccessfulUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(nil)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnFailureUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(mockErr)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 3)
}

func TestSuccessfulUpdateNodebInfoAndPublish(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateNodebInfoAndPublish", nodebInfo).Return(nil)

	rnibDataService.UpdateNodebInfoAndPublish(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoAndPublish", 1)
}

func TestConnFailureUpdateNodebInfoAndPublish(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNodebInfoAndPublish", nodebInfo).Return(mockErr)

	rnibDataService.UpdateNodebInfoAndPublish(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoAndPublish", 3)
}

func TestSuccessfulSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(nil)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 1)
}

func TestConnFailureSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(mockErr)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 3)
}

func TestSuccessfulGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	nodebInfo := &entities.NodebInfo{}
	readerMock.On("GetNodeb", invName).Return(nodebInfo, nil)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
	assert.Equal(t, nodebInfo, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeb, res)
}

func TestSuccessfulGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	nodeIds := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.Equal(t, nodeIds, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)
}

func TestConnFailureTwiceGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)

	res2, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeb, res2)
}

func TestConnFailureWithAnotherConfig(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTestWithMaxAttempts(t, 5)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 5)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)
}

func TestGetGeneralConfigurationConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var config *entities.GeneralConfiguration = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetGeneralConfiguration").Return(config, mockErr)

	res, err := rnibDataService.GetGeneralConfiguration()
	readerMock.AssertNumberOfCalls(t, "GetGeneralConfiguration", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetGeneralConfigurationOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	config := &entities.GeneralConfiguration{}
	readerMock.On("GetGeneralConfiguration").Return(config, nil)

	res, err := rnibDataService.GetGeneralConfiguration()
	readerMock.AssertNumberOfCalls(t, "GetGeneralConfiguration", 1)
	assert.Equal(t, config, res)
	assert.Nil(t, err)
}

func TestGetGeneralConfigurationOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var config *entities.GeneralConfiguration = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetGeneralConfiguration").Return(config, mockErr)

	res, err := rnibDataService.GetGeneralConfiguration()
	readerMock.AssertNumberOfCalls(t, "GetGeneralConfiguration", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestSaveGeneralConfigurationConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	config := &entities.GeneralConfiguration{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveGeneralConfiguration", config).Return(mockErr)

	err := rnibDataService.SaveGeneralConfiguration(config)
	writerMock.AssertNumberOfCalls(t, "SaveGeneralConfiguration", 3)
	assert.NotNil(t, err)
}

func TestSaveGeneralConfigurationOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	config := &entities.GeneralConfiguration{}
	writerMock.On("SaveGeneralConfiguration", config).Return(nil)

	err := rnibDataService.SaveGeneralConfiguration(config)
	writerMock.AssertNumberOfCalls(t, "SaveGeneralConfiguration", 1)
	assert.Nil(t, err)
}

func TestSaveGeneralConfigurationOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	config := &entities.GeneralConfiguration{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("SaveGeneralConfiguration", config).Return(mockErr)

	err := rnibDataService.SaveGeneralConfiguration(config)
	writerMock.AssertNumberOfCalls(t, "SaveGeneralConfiguration", 1)
	assert.NotNil(t, err)
}

func TestRemoveServedCellsConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	cellIds := []*entities.ServedCellInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("RemoveServedCells", ranName, cellIds).Return(mockErr)

	err := rnibDataService.RemoveServedCells(ranName, cellIds)
	writerMock.AssertNumberOfCalls(t, "RemoveServedCells", 3)
	assert.NotNil(t, err)
}

func TestRemoveServedCellsOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	cellIds := []*entities.ServedCellInfo{}
	writerMock.On("RemoveServedCells", ranName, cellIds).Return(nil)

	err := rnibDataService.RemoveServedCells(ranName, cellIds)
	writerMock.AssertNumberOfCalls(t, "RemoveServedCells", 1)
	assert.Nil(t, err)
}

func TestRemoveServedCellsOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	var ranName string = "abcd"
	cellIds := []*entities.ServedCellInfo{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("RemoveServedCells", ranName, cellIds).Return(mockErr)

	err := rnibDataService.RemoveServedCells(ranName, cellIds)
	writerMock.AssertNumberOfCalls(t, "RemoveServedCells", 1)
	assert.NotNil(t, err)
}

func TestUpdateEnbConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	cellIds := []*entities.ServedCellInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateEnb", nodebInfo, cellIds).Return(mockErr)

	err := rnibDataService.UpdateEnb(nodebInfo, cellIds)
	writerMock.AssertNumberOfCalls(t, "UpdateEnb", 3)
	assert.NotNil(t, err)
}

func TestUpdateEnbOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	cellIds := []*entities.ServedCellInfo{}
	writerMock.On("UpdateEnb", nodebInfo, cellIds).Return(nil)

	err := rnibDataService.UpdateEnb(nodebInfo, cellIds)
	writerMock.AssertNumberOfCalls(t, "UpdateEnb", 1)
	assert.Nil(t, err)
}

func TestUpdateEnbOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	cellIds := []*entities.ServedCellInfo{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("UpdateEnb", nodebInfo, cellIds).Return(mockErr)

	err := rnibDataService.UpdateEnb(nodebInfo, cellIds)
	writerMock.AssertNumberOfCalls(t, "UpdateEnb", 1)
	assert.NotNil(t, err)
}

func TestAddEnbConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("AddEnb", nodebInfo).Return(mockErr)

	err := rnibDataService.AddEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "AddEnb", 3)
	assert.NotNil(t, err)
}

func TestAddEnbOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("AddEnb", nodebInfo).Return(nil)

	err := rnibDataService.AddEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "AddEnb", 1)
	assert.Nil(t, err)
}

func TestAddEnbOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("AddEnb", nodebInfo).Return(mockErr)

	err := rnibDataService.AddEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "AddEnb", 1)
	assert.NotNil(t, err)
}

func TestUpdateNbIdentityConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeId := &entities.NbIdentity{}
	newNodeId := &entities.NbIdentity{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNbIdentities", gnbType, []*entities.NbIdentity{oldNodeId},
		[]*entities.NbIdentity{newNodeId}).Return(mockErr)

	err := rnibDataService.UpdateNbIdentity(gnbType, oldNodeId, newNodeId)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 3)
	assert.NotNil(t, err)
}

func TestUpdateNbIdentityOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeId := &entities.NbIdentity{}
	newNodeId := &entities.NbIdentity{}
	writerMock.On("UpdateNbIdentities", gnbType, []*entities.NbIdentity{oldNodeId},
		[]*entities.NbIdentity{newNodeId}).Return(nil)

	err := rnibDataService.UpdateNbIdentity(gnbType, oldNodeId, newNodeId)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 1)
	assert.Nil(t, err)
}

func TestUpdateNbIdentityOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeId := &entities.NbIdentity{}
	newNodeId := &entities.NbIdentity{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("UpdateNbIdentities", gnbType, []*entities.NbIdentity{oldNodeId},
		[]*entities.NbIdentity{newNodeId}).Return(mockErr)

	err := rnibDataService.UpdateNbIdentity(gnbType, oldNodeId, newNodeId)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 1)
	assert.NotNil(t, err)
}

func TestUpdateNbIdentitiesConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeIds := []*entities.NbIdentity{}
	newNodeIds := []*entities.NbIdentity{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNbIdentities", gnbType, oldNodeIds, newNodeIds).Return(mockErr)

	err := rnibDataService.UpdateNbIdentities(gnbType, oldNodeIds, newNodeIds)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 3)
	assert.NotNil(t, err)
}

func TestUpdateNbIdentitiesOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeIds := []*entities.NbIdentity{}
	newNodeIds := []*entities.NbIdentity{}
	writerMock.On("UpdateNbIdentities", gnbType, oldNodeIds, newNodeIds).Return(nil)

	err := rnibDataService.UpdateNbIdentities(gnbType, oldNodeIds, newNodeIds)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 1)
	assert.Nil(t, err)
}

func TestUpdateNbIdentitiesOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	gnbType := entities.Node_GNB
	oldNodeIds := []*entities.NbIdentity{}
	newNodeIds := []*entities.NbIdentity{}
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	writerMock.On("UpdateNbIdentities", gnbType, oldNodeIds, newNodeIds).Return(mockErr)

	err := rnibDataService.UpdateNbIdentities(gnbType, oldNodeIds, newNodeIds)
	writerMock.AssertNumberOfCalls(t, "UpdateNbIdentities", 1)
	assert.NotNil(t, err)
}

func TestPingRnibConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.False(t, res)
}

func TestPingRnibOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}

func TestPingRnibOkOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}

func TestSuccessfulUpdateNodebInfoOnConnectionStatusInversion(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)
	event := "event"

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", nodebInfo, event).Return(nil)

	rnibDataService.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, event)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 1)
}

func TestConnFailureUpdateNodebInfoOnConnectionStatusInversion(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)
	event := "event"

	nodebInfo := &entities.NodebInfo{}
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", nodebInfo, event).Return(mockErr)

	rnibDataService.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, event)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfoOnConnectionStatusInversion", 3)
}

func TestGetE2TInstanceConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstanceOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	e2tInstance := &entities.E2TInstance{}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, nil)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstanceOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstance(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstanceNoLogsConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstanceNoLogs(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstanceNoLogsOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	e2tInstance := &entities.E2TInstance{}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, nil)

	res, err := rnibDataService.GetE2TInstanceNoLogs(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstanceNoLogsOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	address := "10.10.5.20:3200"
	var e2tInstance *entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TInstance", address).Return(e2tInstance, mockErr)

	res, err := rnibDataService.GetE2TInstanceNoLogs(address)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstance", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstancesConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	var e2tInstances []*entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstances, mockErr)

	res, err := rnibDataService.GetE2TInstances(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstancesOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	e2tInstance := []*entities.E2TInstance{}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstance, nil)

	res, err := rnibDataService.GetE2TInstances(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 1)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstancesOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	var e2tInstances []*entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstances, mockErr)

	res, err := rnibDataService.GetE2TInstances(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstancesNoLogsConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	var e2tInstances []*entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstances, mockErr)

	res, err := rnibDataService.GetE2TInstancesNoLogs(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TInstancesNoLogsOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	e2tInstance := []*entities.E2TInstance{}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstance, nil)

	res, err := rnibDataService.GetE2TInstancesNoLogs(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 1)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstancesNoLogsOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	var e2tInstances []*entities.E2TInstance = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TInstances", addresses).Return(e2tInstances, mockErr)

	res, err := rnibDataService.GetE2TInstancesNoLogs(addresses)
	readerMock.AssertNumberOfCalls(t, "GetE2TInstances", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TAddressesConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var addresses []string = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TAddresses").Return(addresses, mockErr)

	res, err := rnibDataService.GetE2TAddresses()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TAddressesOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	readerMock.On("GetE2TAddresses").Return(addresses, nil)

	res, err := rnibDataService.GetE2TAddresses()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 1)
	assert.Nil(t, err)
	assert.Equal(t, addresses, res)
}

func TestGetE2TAddressesOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var addresses []string = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TAddresses").Return(addresses, mockErr)

	res, err := rnibDataService.GetE2TAddresses()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TAddressesNoLogsConnFailure(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var addresses []string = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetE2TAddresses").Return(addresses, mockErr)

	res, err := rnibDataService.GetE2TAddressesNoLogs()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 3)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestGetE2TAddressesNoLogsOkNoError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	readerMock.On("GetE2TAddresses").Return(addresses, nil)

	res, err := rnibDataService.GetE2TAddressesNoLogs()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 1)
	assert.Nil(t, err)
	assert.Equal(t, addresses, res)
}

func TestGetE2TAddressesNoLogsOtherError(t *testing.T) {
	rnibDataService, readerMock, _ := setupRnibDataServiceTest(t)

	var addresses []string = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetE2TAddresses").Return(addresses, mockErr)

	res, err := rnibDataService.GetE2TAddressesNoLogs()
	readerMock.AssertNumberOfCalls(t, "GetE2TAddresses", 1)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestSaveE2TInstanceConnFailure(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("SaveE2TInstance", e2tInstance).Return(mockErr)

	err := rnibDataService.SaveE2TInstance(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 3)
	assert.NotNil(t, err)
}

func TestSaveE2TInstanceOkNoError(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("SaveE2TInstance", e2tInstance).Return(nil)

	err := rnibDataService.SaveE2TInstance(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	assert.Nil(t, err)
}

func TestSaveE2TInstanceOtherError(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("SaveE2TInstance", e2tInstance).Return(mockErr)

	err := rnibDataService.SaveE2TInstance(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)

	assert.NotNil(t, err)
}

func TestSaveE2TInstanceNoLogsConnFailure(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("SaveE2TInstance", e2tInstance).Return(mockErr)

	err := rnibDataService.SaveE2TInstanceNoLogs(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 3)
	assert.NotNil(t, err)
}

func TestSaveE2TInstanceNoLogsOkNoError(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("SaveE2TInstance", e2tInstance).Return(nil)

	err := rnibDataService.SaveE2TInstanceNoLogs(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)
	assert.Nil(t, err)
}

func TestSaveE2TInstanceNoLogsOtherError(t *testing.T) {
	e2tInstance := &entities.E2TInstance{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("SaveE2TInstance", e2tInstance).Return(mockErr)

	err := rnibDataService.SaveE2TInstanceNoLogs(e2tInstance)
	writerMock.AssertNumberOfCalls(t, "SaveE2TInstance", 1)

	assert.NotNil(t, err)
}

func TestSaveE2TAddressesConnFailure(t *testing.T) {
	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("SaveE2TAddresses", addresses).Return(mockErr)

	err := rnibDataService.SaveE2TAddresses(addresses)
	writerMock.AssertNumberOfCalls(t, "SaveE2TAddresses", 3)
	assert.NotNil(t, err)
}

func TestSaveE2TAddressesOkNoError(t *testing.T) {
	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("SaveE2TAddresses", addresses).Return(nil)

	err := rnibDataService.SaveE2TAddresses(addresses)
	writerMock.AssertNumberOfCalls(t, "SaveE2TAddresses", 1)
	assert.Nil(t, err)
}

func TestSaveE2TAddressesOtherError(t *testing.T) {
	addresses := []string{"10.10.5.20:3200", "10.10.5.21:3200"}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("SaveE2TAddresses", addresses).Return(mockErr)

	err := rnibDataService.SaveE2TAddresses(addresses)
	writerMock.AssertNumberOfCalls(t, "SaveE2TAddresses", 1)

	assert.NotNil(t, err)
}

func TestRemoveE2TInstanceConnFailure(t *testing.T) {
	address := "10.10.5.20:3200"
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("RemoveE2TInstance", address).Return(mockErr)

	err := rnibDataService.RemoveE2TInstance(address)
	writerMock.AssertNumberOfCalls(t, "RemoveE2TInstance", 3)
	assert.NotNil(t, err)
}

func TestRemoveE2TInstanceOkNoError(t *testing.T) {
	address := "10.10.5.20:3200"
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("RemoveE2TInstance", address).Return(nil)

	err := rnibDataService.RemoveE2TInstance(address)
	writerMock.AssertNumberOfCalls(t, "RemoveE2TInstance", 1)
	assert.Nil(t, err)
}

func TestRemoveE2TInstanceOtherError(t *testing.T) {
	address := "10.10.5.20:3200"
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("RemoveE2TInstance", address).Return(mockErr)

	err := rnibDataService.RemoveE2TInstance(address)
	writerMock.AssertNumberOfCalls(t, "RemoveE2TInstance", 1)

	assert.NotNil(t, err)
}

func TestAddNbIdentityConnFailure(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("AddNbIdentity", gnbType, nbIdentity).Return(mockErr)

	err := rnibDataService.AddNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "AddNbIdentity", 3)
	assert.NotNil(t, err)
}

func TestAddNbIdentityOkNoError(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("AddNbIdentity", gnbType, nbIdentity).Return(nil)

	err := rnibDataService.AddNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "AddNbIdentity", 1)
	assert.Nil(t, err)
}

func TestAddNbIdentityOtherError(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("AddNbIdentity", gnbType, nbIdentity).Return(mockErr)

	err := rnibDataService.AddNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "AddNbIdentity", 1)

	assert.NotNil(t, err)
}

func TestRemoveNbIdentityConnFailure(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("RemoveNbIdentity", gnbType, nbIdentity).Return(mockErr)

	err := rnibDataService.RemoveNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "RemoveNbIdentity", 3)
	assert.NotNil(t, err)
}

func TestRemoveNbIdentityOkNoError(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("RemoveNbIdentity", gnbType, nbIdentity).Return(nil)

	err := rnibDataService.RemoveNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "RemoveNbIdentity", 1)
	assert.Nil(t, err)
}

func TestRemoveNbIdentityOtherError(t *testing.T) {
	gnbType := entities.Node_GNB
	nbIdentity := &entities.NbIdentity{}
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("RemoveNbIdentity", gnbType, nbIdentity).Return(mockErr)

	err := rnibDataService.RemoveNbIdentity(gnbType, nbIdentity)
	writerMock.AssertNumberOfCalls(t, "RemoveNbIdentity", 1)
	assert.NotNil(t, err)
}

func TestRemoveServedNrCellsConnFailure(t *testing.T) {
	var ranName string = "abcd"
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}

	writerMock.On("RemoveServedNrCells", ranName, servedNrCells).Return(mockErr)

	err := rnibDataService.RemoveServedNrCells(ranName, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "RemoveServedNrCells", 3)
	assert.NotNil(t, err)
}

func TestRemoveServedNrCellsOkNoError(t *testing.T) {
	var ranName string = "abcd"
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	writerMock.On("RemoveServedNrCells", ranName, servedNrCells).Return(nil)

	err := rnibDataService.RemoveServedNrCells(ranName, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "RemoveServedNrCells", 1)
	assert.Nil(t, err)
}

func TestRemoveServedNrCellsOtherError(t *testing.T) {
	var ranName string = "abcd"
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}

	writerMock.On("RemoveServedNrCells", ranName, servedNrCells).Return(mockErr)

	err := rnibDataService.RemoveServedNrCells(ranName, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "RemoveServedNrCells", 1)
	assert.NotNil(t, err)
}

func TestRemoveEnbConnFailure(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	nodebInfo := &entities.NodebInfo{}
	writerMock.On("RemoveEnb", nodebInfo).Return(mockErr)

	err := rnibDataService.RemoveEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "RemoveEnb", 3)
	assert.NotNil(t, err)
}

func TestRemoveEnbOkNoError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("RemoveEnb", nodebInfo).Return(nil)

	err := rnibDataService.RemoveEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "RemoveEnb", 1)
	assert.Nil(t, err)
}

func TestRemoveEnbOtherError(t *testing.T) {
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	nodebInfo := &entities.NodebInfo{}
	writerMock.On("RemoveEnb", nodebInfo).Return(mockErr)

	err := rnibDataService.RemoveEnb(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "RemoveEnb", 1)
	assert.NotNil(t, err)
}

func TestUpdateGnbCellsConnFailure(t *testing.T) {
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateGnbCells", nodebInfo, servedNrCells).Return(mockErr)

	err := rnibDataService.UpdateGnbCells(nodebInfo, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "UpdateGnbCells", 3)
	assert.NotNil(t, err)
}

func TestUpdateGnbCellsOkNoError(t *testing.T) {
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateGnbCells", nodebInfo, servedNrCells).Return(nil)

	err := rnibDataService.UpdateGnbCells(nodebInfo, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "UpdateGnbCells", 1)
	assert.Nil(t, err)
}

func TestUpdateGnbCellsOtherError(t *testing.T) {
	var servedNrCells []*entities.ServedNRCell
	rnibDataService, _, writerMock := setupRnibDataServiceTest(t)

	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateGnbCells", nodebInfo, servedNrCells).Return(mockErr)

	err := rnibDataService.UpdateGnbCells(nodebInfo, servedNrCells)
	writerMock.AssertNumberOfCalls(t, "UpdateGnbCells", 1)
	assert.NotNil(t, err)
}
