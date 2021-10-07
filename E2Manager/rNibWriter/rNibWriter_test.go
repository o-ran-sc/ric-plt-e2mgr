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

package rNibWriter

import (
	"e2mgr/configuration"
	"e2mgr/mocks"
	"encoding/json"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var namespace = common.GetRNibNamespace()

const (
	RanName = "test"
	PlmnId  = "02f829"
	NbId    = "4a952a0a"
)

func initSdlMock() (w RNibWriter, sdlMock *mocks.MockSdlSyncStorage) {
	sdlMock = new(mocks.MockSdlSyncStorage)
	w = GetRNibWriter(sdlMock, configuration.RnibWriterConfig{StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE", RanManipulationMessageChannel: "RAN_MANIPULATION"})
	return
}

func generateNodebInfo(inventoryName string, nodeType entities.Node_Type, plmnId string, nbId string) *entities.NodebInfo {
	nodebInfo := &entities.NodebInfo{
		RanName:          inventoryName,
		GlobalNbId:       &entities.GlobalNbId{PlmnId: plmnId, NbId: nbId},
		NodeType:         nodeType,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
	}

	if nodeType == entities.Node_ENB {
		nodebInfo.Configuration = &entities.NodebInfo_Enb{
			Enb: &entities.Enb{},
		}
	} else if nodeType == entities.Node_GNB {
		nodebInfo.Configuration = &entities.NodebInfo_Gnb{
			Gnb: &entities.Gnb{},
		}
	}

	return nodebInfo
}

func generateServedNrCells(cellIds ...string) []*entities.ServedNRCell {

	var servedNrCells []*entities.ServedNRCell

	for i, v := range cellIds {
		servedNrCells = append(servedNrCells, &entities.ServedNRCell{ServedNrCellInformation: &entities.ServedNRCellInformation{
			CellId: v,
			ChoiceNrMode: &entities.ServedNRCellInformation_ChoiceNRMode{
				Fdd: &entities.ServedNRCellInformation_ChoiceNRMode_FddInfo{},
			},
			NrMode:      entities.Nr_FDD,
			NrPci:       uint32(i + 1),
			ServedPlmns: []string{"whatever"},
		}})
	}

	return servedNrCells
}

func generateServedCells(cellIds ...string) []*entities.ServedCellInfo {

	var servedCells []*entities.ServedCellInfo

	for i, v := range cellIds {
		servedCells = append(servedCells, &entities.ServedCellInfo{
			CellId: v,
			ChoiceEutraMode: &entities.ChoiceEUTRAMode{
				Fdd: &entities.FddInfo{},
			},
			Pci:            uint32(i + 1),
			BroadcastPlmns: []string{"whatever"},
		})
	}

	return servedCells
}

func generateServedCellInfos(cellIds ...string) []*entities.ServedCellInfo {

	servedCells := []*entities.ServedCellInfo{}

	for i, v := range cellIds {
		servedCells = append(servedCells, &entities.ServedCellInfo{
			CellId: v,
			Pci:    uint32(i + 1),
		})
	}

	return servedCells
}

func TestRemoveServedNrCellsSuccess(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedNrCellsToRemove := generateServedNrCells("whatever1", "whatever2")
	sdlMock.On("Remove", namespace, buildServedNRCellKeysToRemove(RanName, servedNrCellsToRemove)).Return(nil)
	err := w.RemoveServedNrCells(RanName, servedNrCellsToRemove)
	assert.Nil(t, err)
}

func TestRemoveServedNrCellsFailure(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedNrCellsToRemove := generateServedNrCells("whatever1", "whatever2")
	sdlMock.On("Remove", namespace, buildServedNRCellKeysToRemove(RanName, servedNrCellsToRemove)).Return(errors.New("expected error"))
	err := w.RemoveServedNrCells(RanName, servedNrCellsToRemove)
	assert.IsType(t, &common.InternalError{}, err)
}

func TestUpdateGnbCellsInvalidNodebInfoFailure(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := &entities.NodebInfo{}
	sdlMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestAddNbIdentitySuccess(t *testing.T) {
	w, sdlMock := initSdlMock()

	nbIdentity := &entities.NbIdentity{InventoryName: RanName, GlobalNbId: &entities.GlobalNbId{PlmnId: PlmnId, NbId: NbId}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Fatalf("#rNibWriter_test.TestAddNbIdentitySuccess - Failed to marshal NodeB Identity entity. Error: %v", err)
	}

	sdlMock.On("AddMember", namespace, "ENB", []interface{}{nbIdData}).Return(nil)
	rNibErr := w.AddNbIdentity(entities.Node_ENB, nbIdentity)
	assert.Nil(t, rNibErr)
}

func TestAddNbIdentityMarshalNilFailure(t *testing.T) {
	w, _ := initSdlMock()

	rNibErr := w.AddNbIdentity(entities.Node_ENB, nil)
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	assert.Equal(t, expectedErr, rNibErr)
}

func TestAddNbIdentitySdlFailure(t *testing.T) {
	w, sdlMock := initSdlMock()

	nbIdentity := &entities.NbIdentity{InventoryName: RanName, GlobalNbId: &entities.GlobalNbId{PlmnId: PlmnId, NbId: NbId}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Fatalf("#rNibWriter_test.TestAddNbIdentitySdlFailure - Failed to marshal NodeB Identity entity. Error: %v", err)
	}

	sdlMock.On("AddMember", namespace, "ENB", []interface{}{nbIdData}).Return(errors.New("expected error"))
	rNibErr := w.AddNbIdentity(entities.Node_ENB, nbIdentity)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestUpdateGnbCellsInvalidCellFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedNrCells := []*entities.ServedNRCell{{ServedNrCellInformation: &entities.ServedNRCellInformation{}}}
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	sdlMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func getUpdateEnbCellsSetExpected(t *testing.T, nodebInfo *entities.NodebInfo, servedCells []*entities.ServedCellInfo) []interface{} {

	nodebInfoData, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Fatalf("#rNibWriter_test.getUpdateEnbCellsSetExpected - Failed to marshal NodeB entity. Error: %s", err)
	}

	nodebNameKey, _ := common.ValidateAndBuildNodeBNameKey(nodebInfo.RanName)
	nodebIdKey, _ := common.ValidateAndBuildNodeBIdKey(nodebInfo.NodeType.String(), nodebInfo.GlobalNbId.PlmnId, nodebInfo.GlobalNbId.NbId)
	setExpected := []interface{}{nodebNameKey, nodebInfoData, nodebIdKey, nodebInfoData}

	for _, cell := range servedCells {

		cellEntity := entities.Cell{Type: entities.Cell_LTE_CELL, Cell: &entities.Cell_ServedCellInfo{ServedCellInfo: cell}}
		cellData, err := proto.Marshal(&cellEntity)

		if err != nil {
			t.Fatalf("#rNibWriter_test.getUpdateEnbCellsSetExpected - Failed to marshal cell entity. Error: %s", err)
		}

		nrCellIdKey, _ := common.ValidateAndBuildCellIdKey(cell.GetCellId())
		cellNamePciKey, _ := common.ValidateAndBuildCellNamePciKey(nodebInfo.RanName, cell.GetPci())
		setExpected = append(setExpected, nrCellIdKey, cellData, cellNamePciKey, cellData)
	}

	return setExpected
}

func getUpdateGnbCellsSetExpected(t *testing.T, nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) []interface{} {

	nodebInfoData, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Fatalf("#rNibWriter_test.getUpdateGnbCellsSetExpected - Failed to marshal NodeB entity. Error: %s", err)
	}

	nodebNameKey, _ := common.ValidateAndBuildNodeBNameKey(nodebInfo.RanName)
	nodebIdKey, _ := common.ValidateAndBuildNodeBIdKey(nodebInfo.NodeType.String(), nodebInfo.GlobalNbId.PlmnId, nodebInfo.GlobalNbId.NbId)
	setExpected := []interface{}{nodebNameKey, nodebInfoData, nodebIdKey, nodebInfoData}

	for _, v := range servedNrCells {

		cellEntity := entities.Cell{Type: entities.Cell_NR_CELL, Cell: &entities.Cell_ServedNrCell{ServedNrCell: v}}
		cellData, err := proto.Marshal(&cellEntity)

		if err != nil {
			t.Fatalf("#rNibWriter_test.getUpdateGnbCellsSetExpected - Failed to marshal cell entity. Error: %s", err)
		}

		nrCellIdKey, _ := common.ValidateAndBuildNrCellIdKey(v.GetServedNrCellInformation().GetCellId())
		cellNamePciKey, _ := common.ValidateAndBuildCellNamePciKey(nodebInfo.RanName, v.GetServedNrCellInformation().GetNrPci())
		setExpected = append(setExpected, nrCellIdKey, cellData, cellNamePciKey, cellData)
	}

	return setExpected
}

func TestUpdateGnbCellsSdlFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	setExpected := getUpdateGnbCellsSetExpected(t, nodebInfo, servedNrCells)
	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(errors.New("expected error"))
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestUpdateGnbCellsRnibKeyValidationError(t *testing.T) {
	//Empty RAN name fails RNIB validation
	inventoryName := ""
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, _ := initSdlMock()
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells

	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateGnbCellsSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	setExpected := getUpdateGnbCellsSetExpected(t, nodebInfo, servedNrCells)
	var e error
	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(e)
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNodebInfoAndPublishSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoAndPublish(nodebInfo)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNodebInfoMissingInventoryNameFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	nodebInfo := &entities.NodebInfo{}
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)

	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateNodebInfoMissingGlobalNbId(t *testing.T) {
	inventoryName := "name"
	w, sdlMock := initSdlMock()
	nodebInfo := &entities.NodebInfo{}
	nodebInfo.RanName = inventoryName
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	setExpected = append(setExpected, nodebNameKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)

	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoSdlSetFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	e := errors.New("expected error")
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestSaveEnb(t *testing.T) {
	ranName := "RAN:" + RanName
	w, sdlMock := initSdlMock()
	nb := entities.NodebInfo{
		RanName:          RanName,
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   "4a952a0a",
			PlmnId: "02f829",
		},
	}

	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaff", Pci: 3}
	cellEntity := entities.Cell{Type: entities.Cell_LTE_CELL, Cell: &entities.Cell_ServedCellInfo{ServedCellInfo: cell}}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	data, err := proto.Marshal(&nb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error

	cellData, err := proto.Marshal(&cellEntity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal Cell entity. Error: %v", err)
	}
	var setExpected []interface{}
	setExpected = append(setExpected, ranName, data)
	setExpected = append(setExpected, "ENB:02f829:4a952a0a", data)
	setExpected = append(setExpected, fmt.Sprintf("CELL:%s", cell.GetCellId()), cellData)
	setExpected = append(setExpected, fmt.Sprintf("PCI:%s:%02x", RanName, cell.GetPci()), cellData)

	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)
	rNibErr := w.SaveNodeb(&nb)
	assert.Nil(t, rNibErr)
}

func TestSaveEnbCellIdValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{}
	nb.RanName = "name"
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.SaveNodeb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveEnbInventoryNameValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   "4a952a0a",
			PlmnId: "02f829",
		},
	}
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaa", Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.SaveNodeb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNodeBNameKey - an empty inventory name received", rNibErr.Error())
}

func TestSaveEnbGlobalNbIdPlmnValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{
		RanName:          RanName,
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId: "4a952a0a",
			//Empty PLMNID fails RNIB validation
			PlmnId: "",
		},
	}
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaa", Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.SaveNodeb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNodeBIdKey - an empty plmnId received", rNibErr.Error())
}

func TestSaveGnbCellIdValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{}
	nb.RanName = "name"
	nb.NodeType = entities.Node_GNB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	gnb := entities.Gnb{}
	cellInfo := &entities.ServedNRCellInformation{NrPci: 2}
	cell := &entities.ServedNRCell{ServedNrCellInformation: cellInfo}
	gnb.ServedNrCells = []*entities.ServedNRCell{cell}
	nb.Configuration = &entities.NodebInfo_Gnb{Gnb: &gnb}

	rNibErr := w.SaveNodeb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNrCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveGnb(t *testing.T) {
	ranName := "RAN:" + RanName
	w, sdlMock := initSdlMock()
	nb := entities.NodebInfo{
		RanName:          RanName,
		NodeType:         entities.Node_GNB,
		ConnectionStatus: 1,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   "4a952a0a",
			PlmnId: "02f829",
		},
		Ip:   "localhost",
		Port: 5656,
	}

	gnb := entities.Gnb{}
	cellInfo := &entities.ServedNRCellInformation{NrPci: 2, CellId: "ccdd"}
	cell := &entities.ServedNRCell{ServedNrCellInformation: cellInfo}
	cellEntity := entities.Cell{Type: entities.Cell_NR_CELL, Cell: &entities.Cell_ServedNrCell{ServedNrCell: cell}}
	gnb.ServedNrCells = []*entities.ServedNRCell{cell}
	nb.Configuration = &entities.NodebInfo_Gnb{Gnb: &gnb}
	data, err := proto.Marshal(&nb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveGnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error

	cellData, err := proto.Marshal(&cellEntity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveGnb - Failed to marshal Cell entity. Error: %v", err)
	}
	var setExpected []interface{}
	setExpected = append(setExpected, ranName, data)
	setExpected = append(setExpected, "GNB:02f829:4a952a0a", data)
	setExpected = append(setExpected, fmt.Sprintf("NRCELL:%s", cell.GetServedNrCellInformation().GetCellId()), cellData)
	setExpected = append(setExpected, fmt.Sprintf("PCI:%s:%02x", RanName, cell.GetServedNrCellInformation().GetNrPci()), cellData)

	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)
	rNibErr := w.SaveNodeb(&nb)
	assert.Nil(t, rNibErr)
}

func TestSaveRanLoadInformationSuccess(t *testing.T) {
	inventoryName := "name"
	loadKey, validationErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformationSuccess - Failed to build ran load infromation key. Error: %v", validationErr)
	}

	w, sdlMock := initSdlMock()

	ranLoadInformation := generateRanLoadInformation()
	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformation - Failed to marshal RanLoadInformation entity. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveRanLoadInformation(inventoryName, ranLoadInformation)
	assert.Nil(t, rNibErr)
}

func TestSaveRanLoadInformationMarshalNilFailure(t *testing.T) {
	inventoryName := "name2"
	w, _ := initSdlMock()

	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	err := w.SaveRanLoadInformation(inventoryName, nil)
	assert.Equal(t, expectedErr, err)
}

func TestSaveRanLoadInformationEmptyInventoryNameFailure(t *testing.T) {
	inventoryName := ""
	w, _ := initSdlMock()

	err := w.SaveRanLoadInformation(inventoryName, nil)
	assert.NotNil(t, err)
	assert.IsType(t, &common.ValidationError{}, err)
}

func TestSaveRanLoadInformationSdlFailure(t *testing.T) {
	inventoryName := "name2"

	loadKey, validationErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformationSuccess - Failed to build ran load infromation key. Error: %v", validationErr)
	}

	w, sdlMock := initSdlMock()

	ranLoadInformation := generateRanLoadInformation()
	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformation - Failed to marshal RanLoadInformation entity. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(expectedErr)

	rNibErr := w.SaveRanLoadInformation(inventoryName, ranLoadInformation)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func generateCellLoadInformation() *entities.CellLoadInformation {
	cellLoadInformation := entities.CellLoadInformation{}

	cellLoadInformation.CellId = "123"

	ulInterferenceOverloadIndication := entities.UlInterferenceOverloadIndication_HIGH_INTERFERENCE
	cellLoadInformation.UlInterferenceOverloadIndications = []entities.UlInterferenceOverloadIndication{ulInterferenceOverloadIndication}

	ulHighInterferenceInformation := entities.UlHighInterferenceInformation{
		TargetCellId:                 "456",
		UlHighInterferenceIndication: "xxx",
	}

	cellLoadInformation.UlHighInterferenceInfos = []*entities.UlHighInterferenceInformation{&ulHighInterferenceInformation}

	cellLoadInformation.RelativeNarrowbandTxPower = &entities.RelativeNarrowbandTxPower{
		RntpPerPrb:                       "xxx",
		RntpThreshold:                    entities.RntpThreshold_NEG_4,
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V1_ANT_PRT,
		PB:                               1,
		PdcchInterferenceImpact:          2,
		EnhancedRntp: &entities.EnhancedRntp{
			EnhancedRntpBitmap:     "xxx",
			RntpHighPowerThreshold: entities.RntpThreshold_NEG_2,
			EnhancedRntpStartTime:  &entities.StartTime{StartSfn: 500, StartSubframeNumber: 5},
		},
	}

	cellLoadInformation.AbsInformation = &entities.AbsInformation{
		Mode:                             entities.AbsInformationMode_ABS_INFO_FDD,
		AbsPatternInfo:                   "xxx",
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V2_ANT_PRT,
		MeasurementSubset:                "xxx",
	}

	cellLoadInformation.InvokeIndication = entities.InvokeIndication_ABS_INFORMATION

	cellLoadInformation.ExtendedUlInterferenceOverloadInfo = &entities.ExtendedUlInterferenceOverloadInfo{
		AssociatedSubframes:                       "xxx",
		ExtendedUlInterferenceOverloadIndications: cellLoadInformation.UlInterferenceOverloadIndications,
	}

	compInformationItem := &entities.CompInformationItem{
		CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "789", CompHypothesis: "xxx"}},
		BenefitMetric:      50,
	}

	cellLoadInformation.CompInformation = &entities.CompInformation{
		CompInformationItems:     []*entities.CompInformationItem{compInformationItem},
		CompInformationStartTime: &entities.StartTime{StartSfn: 123, StartSubframeNumber: 456},
	}

	cellLoadInformation.DynamicDlTransmissionInformation = &entities.DynamicDlTransmissionInformation{
		State:             entities.NaicsState_NAICS_ACTIVE,
		TransmissionModes: "xxx",
		PB:                2,
		PAList:            []entities.PA{entities.PA_DB_NEG_3},
	}

	return &cellLoadInformation
}

func generateRanLoadInformation() *entities.RanLoadInformation {
	ranLoadInformation := entities.RanLoadInformation{}

	ranLoadInformation.LoadTimestamp = uint64(time.Now().UnixNano())

	cellLoadInformation := generateCellLoadInformation()
	ranLoadInformation.CellLoadInfos = []*entities.CellLoadInformation{cellLoadInformation}

	return &ranLoadInformation
}

func TestSaveNilEntityFailure(t *testing.T) {
	w, _ := initSdlMock()
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	actualErr := w.SaveNodeb(nil)
	assert.Equal(t, expectedErr, actualErr)
}

func TestSaveUnknownTypeEntityFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := &entities.NodebInfo{}
	nb.Port = 5656
	nb.Ip = "localhost"
	actualErr := w.SaveNodeb(nb)
	assert.IsType(t, &common.ValidationError{}, actualErr)
}

func TestSaveEntitySetFailure(t *testing.T) {
	name := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"

	w, sdlMock := initSdlMock()
	gnb := entities.NodebInfo{
		RanName:          name,
		NodeType:         entities.Node_GNB,
		ConnectionStatus: 1,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   nbId,
			PlmnId: plmnId,
		},
		Ip:   "localhost",
		Port: 5656,
	}
	data, err := proto.Marshal(&gnb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEntityFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	setExpected := []interface{}{"RAN:" + name, data}
	setExpected = append(setExpected, "GNB:"+plmnId+":"+nbId, data)
	expectedErr := errors.New("expected error")
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(expectedErr)
	rNibErr := w.SaveNodeb(&gnb)
	assert.NotEmpty(t, rNibErr)
}

func TestSaveEntitySetAndPublishFailure(t *testing.T) {
	name := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"

	w, sdlMock := initSdlMock()
	enb := entities.NodebInfo{
		RanName:          name,
		NodeType:         entities.Node_ENB,
		ConnectionStatus: 1,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   nbId,
			PlmnId: plmnId,
		},
		Ip:   "localhost",
		Port: 5656,
	}
	data, err := proto.Marshal(&enb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEntityFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	setExpected := []interface{}{"RAN:" + name, data}
	setExpected = append(setExpected, "ENB:"+plmnId+":"+nbId, data)
	expectedErr := errors.New("expected error")
	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", name + "_" + RanAddedEvent}, []interface{}{setExpected}).Return(expectedErr)
	rNibErr := w.AddEnb(&enb)
	assert.NotEmpty(t, rNibErr)
}

func TestGetRNibWriter(t *testing.T) {
	received, _ := initSdlMock()
	assert.NotEmpty(t, received)
}

func TestSaveE2TInstanceSuccess(t *testing.T) {
	address := "10.10.2.15:9800"
	loadKey, validationErr := common.ValidateAndBuildE2TInstanceKey(address)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSuccess - Failed to build E2T Instance key. Error: %v", validationErr)
	}

	w, sdlMock := initSdlMock()

	e2tInstance := generateE2tInstance(address)
	data, err := json.Marshal(e2tInstance)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSuccess - Failed to marshal E2tInstance entity. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveE2TInstance(e2tInstance)
	assert.Nil(t, rNibErr)
}

func TestSaveE2TInstanceNullE2tInstanceFailure(t *testing.T) {
	w, _ := initSdlMock()
	var address string
	e2tInstance := entities.NewE2TInstance(address, "test")
	err := w.SaveE2TInstance(e2tInstance)
	assert.NotNil(t, err)
	assert.IsType(t, &common.ValidationError{}, err)
}

func TestSaveE2TInstanceSdlFailure(t *testing.T) {
	address := "10.10.2.15:9800"
	loadKey, validationErr := common.ValidateAndBuildE2TInstanceKey(address)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSdlFailure - Failed to build E2T Instance key. Error: %v", validationErr)
	}

	w, sdlMock := initSdlMock()

	e2tInstance := generateE2tInstance(address)
	data, err := json.Marshal(e2tInstance)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSdlFailure - Failed to marshal E2tInstance entity. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(expectedErr)

	rNibErr := w.SaveE2TInstance(e2tInstance)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func generateE2tInstance(address string) *entities.E2TInstance {
	e2tInstance := entities.NewE2TInstance(address, "pod test")

	e2tInstance.AssociatedRanList = []string{"test1", "test2"}

	return e2tInstance
}

func TestSaveE2TAddressesSuccess(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlMock := initSdlMock()

	e2tAddresses := []string{address}
	data, err := json.Marshal(e2tAddresses)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInfoListSuccess - Failed to marshal E2TInfoList. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, E2TAddressesKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveE2TAddresses(e2tAddresses)
	assert.Nil(t, rNibErr)
}

func TestSaveE2TAddressesSdlFailure(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlMock := initSdlMock()

	e2tAddresses := []string{address}
	data, err := json.Marshal(e2tAddresses)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInfoListSdlFailure - Failed to marshal E2TInfoList. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, E2TAddressesKey, data)
	sdlMock.On("Set", namespace, []interface{}{setExpected}).Return(expectedErr)

	rNibErr := w.SaveE2TAddresses(e2tAddresses)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestRemoveE2TInstanceSuccess(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlMock := initSdlMock()

	e2tAddresses := []string{fmt.Sprintf("E2TInstance:%s", address)}
	var e error
	sdlMock.On("Remove", namespace, e2tAddresses).Return(e)

	rNibErr := w.RemoveE2TInstance(address)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestRemoveE2TInstanceSdlFailure(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlMock := initSdlMock()

	e2tAddresses := []string{fmt.Sprintf("E2TInstance:%s", address)}
	expectedErr := errors.New("expected error")
	sdlMock.On("Remove", namespace, e2tAddresses).Return(expectedErr)

	rNibErr := w.RemoveE2TInstance(address)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestRemoveE2TInstanceEmptyAddressFailure(t *testing.T) {
	w, sdlMock := initSdlMock()

	rNibErr := w.RemoveE2TInstance("")
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNodebInfoOnConnectionStatusInversionSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestUpdateNodebInfoOnConnectionStatusInversionSuccess - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("SetAndPublish", namespace, []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)
	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoOnConnectionStatusInversionMissingInventoryNameFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlMock := initSdlMock()
	nodebInfo := &entities.NodebInfo{}
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestUpdateNodebInfoOnConnectionStatusInversionMissingInventoryNameFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("SetAndPublish", namespace, []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)

	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateNodebInfoOnConnectionStatusInversionMissingGlobalNbId(t *testing.T) {
	inventoryName := "name"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlMock := initSdlMock()
	nodebInfo := &entities.NodebInfo{}
	nodebInfo.RanName = inventoryName
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestUpdateNodebInfoOnConnectionStatusInversionMissingInventoryNameFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	setExpected = append(setExpected, nodebNameKey, data)
	sdlMock.On("SetAndPublish", namespace, []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)

	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoOnConnectionStatusInversionSdlFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	data, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestUpdateNodebInfoOnConnectionStatusInversionSuccess - Failed to marshal NodeB entity. Error: %v", err)
	}
	e := errors.New("expected error")
	var setExpected []interface{}

	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	setExpected = append(setExpected, nodebNameKey, data)
	setExpected = append(setExpected, nodebIdKey, data)

	sdlMock.On("SetAndPublish", namespace, []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestSaveGeneralConfiguration(t *testing.T) {
	w, sdlMock := initSdlMock()

	key := common.BuildGeneralConfigurationKey()
	configurationData := "{\"enableRic\":true}"
	configuration := &entities.GeneralConfiguration{}
	configuration.EnableRic = true

	sdlMock.On("Set", namespace, []interface{}{[]interface{}{key, []byte(configurationData)}}).Return(nil)
	rNibErr := w.SaveGeneralConfiguration(configuration)

	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestSaveGeneralConfigurationDbError(t *testing.T) {
	w, sdlMock := initSdlMock()

	key := common.BuildGeneralConfigurationKey()
	configurationData := "{\"enableRic\":true}"
	configuration := &entities.GeneralConfiguration{}
	configuration.EnableRic = true

	expectedErr := errors.New("expected error")

	sdlMock.On("Set", namespace, []interface{}{[]interface{}{key, []byte(configurationData)}}).Return(expectedErr)
	rNibErr := w.SaveGeneralConfiguration(configuration)

	assert.NotNil(t, rNibErr)
}

func TestRemoveServedCellsFailure(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedCellsToRemove := generateServedCells("whatever1", "whatever2")
	expectedErr := errors.New("expected error")
	sdlMock.On("Remove", namespace, buildServedCellInfoKeysToRemove(RanName, servedCellsToRemove)).Return(expectedErr)

	rNibErr := w.RemoveServedCells(RanName, servedCellsToRemove)

	assert.NotNil(t, rNibErr)
}

func TestRemoveServedCellsSuccess(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedCellsToRemove := generateServedCells("whatever1", "whatever2")
	sdlMock.On("Remove", namespace, buildServedCellInfoKeysToRemove(RanName, servedCellsToRemove)).Return(nil)
	err := w.RemoveServedCells(RanName, servedCellsToRemove)
	assert.Nil(t, err)
}

func TestUpdateEnbInvalidNodebInfoFailure(t *testing.T) {
	w, sdlMock := initSdlMock()
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := &entities.NodebInfo{}
	sdlMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateEnbInvalidCellFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedCells := []*entities.ServedCellInfo{{CellId: ""}}
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	sdlMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateEnbRnibKeyValidationError(t *testing.T) {
	//Empty RAN name fails RNIB validation
	inventoryName := ""
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, _ := initSdlMock()
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells

	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateEnbSdlFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	setExpected := getUpdateEnbCellsSetExpected(t, nodebInfo, servedCells)
	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(errors.New("expected error"))
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestUpdateEnbSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlMock := initSdlMock()
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	setExpected := getUpdateEnbCellsSetExpected(t, nodebInfo, servedCells)

	var e error
	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(e)
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.Nil(t, rNibErr)
}

func getUpdateEnbSetExpected(t *testing.T, nodebInfo *entities.NodebInfo, servedCells []*entities.ServedCellInfo) []interface{} {

	nodebInfoData, err := proto.Marshal(nodebInfo)
	if err != nil {
		t.Fatalf("#rNibWriter_test.getUpdateEnbSetExpected - Failed to marshal NodeB entity. Error: %s", err)
	}

	nodebNameKey, _ := common.ValidateAndBuildNodeBNameKey(nodebInfo.RanName)
	nodebIdKey, _ := common.ValidateAndBuildNodeBIdKey(nodebInfo.NodeType.String(), nodebInfo.GlobalNbId.PlmnId, nodebInfo.GlobalNbId.NbId)
	setExpected := []interface{}{nodebNameKey, nodebInfoData, nodebIdKey, nodebInfoData}

	for _, v := range servedCells {

		cellEntity := entities.ServedCellInfo{CellId: "some cell id", EutraMode: entities.Eutra_FDD, CsgId: "some csg id"}
		cellData, err := proto.Marshal(&cellEntity)

		if err != nil {
			t.Fatalf("#rNibWriter_test.getUpdateEnbSetExpected - Failed to marshal cell entity. Error: %s", err)
		}

		nrCellIdKey, _ := common.ValidateAndBuildNrCellIdKey(v.GetCellId())
		cellNamePciKey, _ := common.ValidateAndBuildCellNamePciKey(nodebInfo.RanName, v.GetPci())
		setExpected = append(setExpected, nrCellIdKey, cellData, cellNamePciKey, cellData)
	}
	return setExpected
}

func TestRemoveEnbSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_MANIPULATION"
	eventName := inventoryName + "_" + "DELETED"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	var e error

	expectedKeys := []string{}
	cell1Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[0].CellId)
	cell1PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[0].Pci)
	cell2Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[1].CellId)
	cell2PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[1].Pci)
	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	expectedKeys = append(expectedKeys, cell1Key, cell1PciKey, cell2Key, cell2PciKey, nodebNameKey, nodebIdKey)
	sdlMock.On("RemoveAndPublish", namespace, []string{channelName, eventName}, expectedKeys).Return(e)

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestRemoveEnbRnibKeyValidationError(t *testing.T) {
	//Empty RAN name fails RNIB key validation
	inventoryName := ""
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, _ := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.NotNil(t, rNibErr)
}

func TestRemoveEnbRemoveAndPublishError(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_MANIPULATION"
	eventName := inventoryName + "_" + "DELETED"
	w, sdlMock := initSdlMock()
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	expectedKeys := []string{}
	cell1Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[0].CellId)
	cell1PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[0].Pci)
	cell2Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[1].CellId)
	cell2PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[1].Pci)
	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	expectedKeys = append(expectedKeys, cell1Key, cell1PciKey, cell2Key, cell2PciKey, nodebNameKey, nodebIdKey)
	sdlMock.On("RemoveAndPublish", namespace, []string{channelName, eventName}, expectedKeys).Return(errors.New("for test"))

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestRemoveNbIdentitySuccess(t *testing.T) {
	w, sdlMock := initSdlMock()
	nbIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#TestRemoveNbIdentitySuccess - failed to Marshal NbIdentity")
	}

	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), []interface{}{nbIdData}).Return(nil)

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nbIdentity)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestRemoveNbIdentityMarshalNilFailure(t *testing.T) {
	w, _ := initSdlMock()

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nil)
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	assert.Equal(t, expectedErr, rNibErr)
}

func TestRemoveNbIdentityError(t *testing.T) {
	w, sdlMock := initSdlMock()
	nbIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#TestRemoveNbIdentitySuccess - failed to Marshal NbIdentity")
	}

	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), []interface{}{nbIdData}).Return(fmt.Errorf("for test"))

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nbIdentity)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestAddEnb(t *testing.T) {
	ranName := "RAN:" + RanName
	w, sdlMock := initSdlMock()
	nb := entities.NodebInfo{
		RanName:          RanName,
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   "4a952a0a",
			PlmnId: "02f829",
		},
	}

	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaff", Pci: 3}
	cellEntity := entities.Cell{Type: entities.Cell_LTE_CELL, Cell: &entities.Cell_ServedCellInfo{ServedCellInfo: cell}}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	data, err := proto.Marshal(&nb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB entity. Error: %v", err)
	}
	var e error

	cellData, err := proto.Marshal(&cellEntity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal Cell entity. Error: %v", err)
	}
	var setExpected []interface{}
	setExpected = append(setExpected, ranName, data)
	setExpected = append(setExpected, "ENB:02f829:4a952a0a", data)
	setExpected = append(setExpected, fmt.Sprintf("CELL:%s", cell.GetCellId()), cellData)
	setExpected = append(setExpected, fmt.Sprintf("PCI:%s:%02x", RanName, cell.GetPci()), cellData)

	sdlMock.On("SetAndPublish", namespace, []string{"RAN_MANIPULATION", RanName + "_" + RanAddedEvent}, []interface{}{setExpected}).Return(e)

	rNibErr := w.AddEnb(&nb)
	assert.Nil(t, rNibErr)
}

func TestAddEnbMarshalNilFailure(t *testing.T) {
	w, _ := initSdlMock()

	rNibErr := w.AddEnb(nil)
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	assert.Equal(t, expectedErr, rNibErr)
}

func TestAddEnbCellIdValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{}
	nb.RanName = "name"
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.AddEnb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestAddEnbInventoryNameValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId:   "4a952a0a",
			PlmnId: "02f829",
		},
	}
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaa", Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.AddEnb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNodeBNameKey - an empty inventory name received", rNibErr.Error())
}

func TestAddEnbGlobalNbIdPlmnValidationFailure(t *testing.T) {
	w, _ := initSdlMock()
	nb := entities.NodebInfo{
		RanName:          "name",
		NodeType:         entities.Node_ENB,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Ip:               "localhost",
		Port:             5656,
		GlobalNbId: &entities.GlobalNbId{
			NbId: "4a952a0a",
			//Empty PLMNID fails RNIB validation
			PlmnId: "",
		},
	}
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaa", Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}
	rNibErr := w.AddEnb(&nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNodeBIdKey - an empty plmnId received", rNibErr.Error())
}

func TestUpdateNbIdentityOneMemberSuccess(t *testing.T) {
	w, sdlMock := initSdlMock()

	proto, nbIdentity := createNbIdentityProto(t, "ran1", "plmnId1", "nbId1", entities.ConnectionStatus_DISCONNECTED)
	val := []interface{}{proto}

	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), val).Return(nil)

	protoAdd, nbIdentityAdd := createNbIdentityProto(t, "ran1_add", "plmnId1_add", "nbId1_add", entities.ConnectionStatus_CONNECTED)
	sdlMock.On("AddMember", namespace, entities.Node_ENB.String(), []interface{}{protoAdd}).Return(nil)

	newNbIdIdentities := []*entities.NbIdentity{nbIdentityAdd}
	oldNbIdIdentities := []*entities.NbIdentity{nbIdentity}

	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentitySuccess(t *testing.T) {
	w, sdlMock := initSdlMock()

	var nbIdIdentitiesProtoToRemove []interface{}
	protoRan1, _ := createNbIdentityProto(t, "ran1", "plmnId1", "nbId1", entities.ConnectionStatus_DISCONNECTED)
	protoRan2, _ := createNbIdentityProto(t, "ran2", "plmnId2", "nbId2", entities.ConnectionStatus_DISCONNECTED)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan1)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan2)
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToRemove).Return(nil)

	var nbIdIdentitiesProtoToAdd []interface{}
	protoRan1Add, _ := createNbIdentityProto(t, "ran1_add", "plmnId1_add", "nbId1_add", entities.ConnectionStatus_CONNECTED)
	protoRan2Add, _ := createNbIdentityProto(t, "ran2_add", "plmnId2_add", "nbId2_add", entities.ConnectionStatus_CONNECTED)
	nbIdIdentitiesProtoToAdd = append(nbIdIdentitiesProtoToAdd, protoRan1Add)
	nbIdIdentitiesProtoToAdd = append(nbIdIdentitiesProtoToAdd, protoRan2Add)
	sdlMock.On("AddMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToAdd).Return(nil)

	var newNbIdIdentities []*entities.NbIdentity
	firstNewNbIdIdentity := &entities.NbIdentity{InventoryName: "ran1_add", ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1_add", NbId: "nbId1_add"}}
	secondNewNbIdIdentity := &entities.NbIdentity{InventoryName: "ran2_add", ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2_add", NbId: "nbId2_add"}}
	newNbIdIdentities = append(newNbIdIdentities, firstNewNbIdIdentity)
	newNbIdIdentities = append(newNbIdIdentities, secondNewNbIdIdentity)

	var oldNbIdIdentities []*entities.NbIdentity
	firstOldNbIdIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	secondOldNbIdIdentity := &entities.NbIdentity{InventoryName: "ran2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}}
	oldNbIdIdentities = append(oldNbIdIdentities, firstOldNbIdIdentity)
	oldNbIdIdentities = append(oldNbIdIdentities, secondOldNbIdIdentity)

	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.Nil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentityOldIdentityMarshalNilFailure(t *testing.T) {
	w, _ := initSdlMock()

	oldNbIdIdentities := []*entities.NbIdentity{nil}
	newNbIdIdentities := []*entities.NbIdentity{
		&entities.NbIdentity{
			InventoryName:    "ran1_add",
			ConnectionStatus: entities.ConnectionStatus_CONNECTED,
			GlobalNbId:       &entities.GlobalNbId{PlmnId: "plmnId1_add", NbId: "nbId1_add"},
		},
	}

	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.Equal(t, expectedErr, rNibErr)
}

func TestUpdateNbIdentityNewIdentityMarshalNilFailure(t *testing.T) {
	w, sdlMock := initSdlMock()

	var nbIdIdentitiesProtoToRemove []interface{}
	protoRan1, _ := createNbIdentityProto(t, "ran1", "plmnId1", "nbId1", entities.ConnectionStatus_DISCONNECTED)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan1)
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToRemove).Return(nil)

	oldNbIdIdentities := []*entities.NbIdentity{
		&entities.NbIdentity{
			InventoryName:    "ran1",
			ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,
			GlobalNbId:       &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"},
		},
	}
	newNbIdIdentities := []*entities.NbIdentity{nil}

	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.Equal(t, expectedErr, rNibErr)
}

func TestUpdateNbIdentityRemoveFailure(t *testing.T) {
	w, sdlMock := initSdlMock()

	var nbIdIdentitiesProtoToRemove []interface{}
	protoRan1, _ := createNbIdentityProto(t, "ran1", "plmnId1", "nbId1", entities.ConnectionStatus_DISCONNECTED)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan1)
	protoRan2, _ := createNbIdentityProto(t, "ran2", "plmnId2", "nbId2", entities.ConnectionStatus_DISCONNECTED)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan2)

	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToRemove).Return(fmt.Errorf("for test"))

	var oldNbIdIdentities []*entities.NbIdentity
	firstOldNbIdIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	secondOldNbIdIdentity := &entities.NbIdentity{InventoryName: "ran2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}}
	oldNbIdIdentities = append(oldNbIdIdentities, firstOldNbIdIdentity)
	oldNbIdIdentities = append(oldNbIdIdentities, secondOldNbIdIdentity)

	var newNbIdIdentities []*entities.NbIdentity

	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentitySdlAddMemberFailure(t *testing.T) {
	w, sdlMock := initSdlMock()

	var nbIdIdentitiesProtoToRemove []interface{}
	protoRan1, _ := createNbIdentityProto(t, "ran1", "plmnId1", "nbId1", entities.ConnectionStatus_DISCONNECTED)
	nbIdIdentitiesProtoToRemove = append(nbIdIdentitiesProtoToRemove, protoRan1)
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToRemove).Return(nil)

	var nbIdIdentitiesProtoToAdd []interface{}
	protoRan1Add, _ := createNbIdentityProto(t, "ran1_add", "plmnId1_add", "nbId1_add", entities.ConnectionStatus_CONNECTED)
	nbIdIdentitiesProtoToAdd = append(nbIdIdentitiesProtoToAdd, protoRan1Add)
	sdlMock.On("AddMember", namespace, entities.Node_ENB.String(), nbIdIdentitiesProtoToAdd).Return(fmt.Errorf("for test"))

	var oldNbIdIdentities []*entities.NbIdentity
	firstOldNbIdIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	oldNbIdIdentities = append(oldNbIdIdentities, firstOldNbIdIdentity)

	var newNbIdIdentities []*entities.NbIdentity
	firstNewNbIdIdentity := &entities.NbIdentity{InventoryName: "ran1_add", ConnectionStatus: entities.ConnectionStatus_CONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1_add", NbId: "nbId1_add"}}
	newNbIdIdentities = append(newNbIdIdentities, firstNewNbIdIdentity)

	rNibErr := w.UpdateNbIdentities(entities.Node_ENB, oldNbIdIdentities, newNbIdIdentities)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentityAddFailure(t *testing.T) {
	w, sdlMock := initSdlMock()
	nbIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#TestRemoveNbIdentitySuccess - failed to Marshal NbIdentity")
	}
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), []interface{}{nbIdData}).Return(fmt.Errorf("for test"))

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nbIdentity)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentityNoNbIdentityToRemove(t *testing.T) {
	w, sdlMock := initSdlMock()
	nbIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#TestRemoveNbIdentitySuccess - failed to Marshal NbIdentity")
	}
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), []interface{}{nbIdData}).Return(fmt.Errorf("for test"))

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nbIdentity)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func TestUpdateNbIdentityNoNbIdentityToAdd(t *testing.T) {
	w, sdlMock := initSdlMock()
	nbIdentity := &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#TestRemoveNbIdentitySuccess - failed to Marshal NbIdentity")
	}
	sdlMock.On("RemoveMember", namespace, entities.Node_ENB.String(), []interface{}{nbIdData}).Return(fmt.Errorf("for test"))

	rNibErr := w.RemoveNbIdentity(entities.Node_ENB, nbIdentity)
	assert.NotNil(t, rNibErr)
	sdlMock.AssertExpectations(t)
}

func createNbIdentityProto(t *testing.T, ranName string, plmnId string, nbId string, connectionStatus entities.ConnectionStatus) ([]byte, *entities.NbIdentity) {
	nbIdentity := &entities.NbIdentity{InventoryName: ranName, ConnectionStatus: connectionStatus, GlobalNbId: &entities.GlobalNbId{PlmnId: plmnId, NbId: nbId}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#createNbIdentityProto - failed to Marshal NbIdentity")
	}
	return nbIdData, nbIdentity
}
