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

var namespace = "namespace"

const (
	RanName = "test"
)

func initSdlInstanceMock(namespace string) (w RNibWriter, sdlInstanceMock *mocks.MockSdlInstance) {
	sdlInstanceMock = new(mocks.MockSdlInstance)
	w = GetRNibWriter(sdlInstanceMock, configuration.RnibWriterConfig{StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE", RanManipulationMessageChannel: "RAN_MANIPULATION"})
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
				Fdd: &entities.ServedNRCellInformation_ChoiceNRMode_FddInfo{

				},
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
			Pci:       uint32(i + 1),
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
			Pci:       uint32(i + 1),
		})
	}

	return servedCells
}

func TestRemoveServedNrCellsSuccess(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCellsToRemove := generateServedNrCells("whatever1", "whatever2")
	sdlInstanceMock.On("Remove", buildServedNRCellKeysToRemove(RanName, servedNrCellsToRemove)).Return(nil)
	err := w.RemoveServedNrCells(RanName, servedNrCellsToRemove)
	assert.Nil(t, err)
}

func TestRemoveServedNrCellsFailure(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCellsToRemove := generateServedNrCells("whatever1", "whatever2")
	sdlInstanceMock.On("Remove", buildServedNRCellKeysToRemove(RanName, servedNrCellsToRemove)).Return(errors.New("expected error"))
	err := w.RemoveServedNrCells(RanName, servedNrCellsToRemove)
	assert.IsType(t, &common.InternalError{}, err)
}

func TestUpdateGnbCellsInvalidNodebInfoFailure(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := &entities.NodebInfo{}
	sdlInstanceMock.AssertNotCalled(t, "Set")
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateGnbCellsInvalidCellFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCells := []*entities.ServedNRCell{{ServedNrCellInformation: &entities.ServedNRCellInformation{}}}
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	sdlInstanceMock.AssertNotCalled(t, "Set")
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
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	setExpected := getUpdateGnbCellsSetExpected(t, nodebInfo, servedNrCells)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(errors.New("expected error"))
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestUpdateGnbCellsSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedNrCells := generateServedNrCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_GNB, plmnId, nbId)
	nodebInfo.GetGnb().ServedNrCells = servedNrCells
	setExpected := getUpdateGnbCellsSetExpected(t, nodebInfo, servedNrCells)
	var e error
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)
	rNibErr := w.UpdateGnbCells(nodebInfo, servedNrCells)
	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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

	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)
	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoMissingInventoryNameFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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

	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)

	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateNodebInfoMissingGlobalNbId(t *testing.T) {
	inventoryName := "name"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfo(nodebInfo)

	assert.Nil(t, rNibErr)
}

func TestSaveEnb(t *testing.T) {
	name := "name"
	ranName := "RAN:" + name
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	nb := entities.NodebInfo{}
	nb.RanName = name
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
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
	setExpected = append(setExpected, fmt.Sprintf("PCI:%s:%02x", name, cell.GetPci()), cellData)

	sdlInstanceMock.On("SetAndPublish", []string{"RAN_MANIPULATION", name + "_" + RanAddedEvent}, []interface{}{setExpected}).Return(e)

	nbIdData, err := proto.Marshal(&entities.NbIdentity{InventoryName: name})
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal nbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_UNKNOWN.String(), []interface{}{nbIdData}).Return(e)

	nbIdentity := &entities.NbIdentity{InventoryName: name, GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	nbIdData, err = proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal NodeB Identity entity. Error: %v", err)
	}
	sdlInstanceMock.On("AddMember", "ENB", []interface{}{nbIdData}).Return(e)

	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.Nil(t, rNibErr)
}

func TestSaveEnbCellIdValidationFailure(t *testing.T) {
	name := "name"
	w, _ := initSdlInstanceMock(namespace)
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}

	nbIdentity := &entities.NbIdentity{InventoryName: name, GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveEnbInventoryNameValidationFailure(t *testing.T) {
	w, _ := initSdlInstanceMock(namespace)
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId: "aaa", Pci: 3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb: &enb}

	nbIdentity := &entities.NbIdentity{InventoryName: "", GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNodeBNameKey - an empty inventory name received", rNibErr.Error())
}

func TestSaveGnbCellIdValidationFailure(t *testing.T) {
	name := "name"
	w, _ := initSdlInstanceMock(namespace)
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_GNB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	gnb := entities.Gnb{}
	cellInfo := &entities.ServedNRCellInformation{NrPci: 2}
	cell := &entities.ServedNRCell{ServedNrCellInformation: cellInfo}
	gnb.ServedNrCells = []*entities.ServedNRCell{cell}
	nb.Configuration = &entities.NodebInfo_Gnb{Gnb: &gnb}

	nbIdentity := &entities.NbIdentity{InventoryName: name, GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	assert.Equal(t, "#utils.ValidateAndBuildNrCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveGnb(t *testing.T) {
	name := "name"
	ranName := "RAN:" + name
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_GNB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
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
	setExpected = append(setExpected, fmt.Sprintf("PCI:%s:%02x", name, cell.GetServedNrCellInformation().GetNrPci()), cellData)

	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)
	nbIdentity := &entities.NbIdentity{InventoryName: name, GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveGnb - Failed to marshal NodeB Identity entity. Error: %v", err)
	}
	sdlInstanceMock.On("AddMember", "GNB", []interface{}{nbIdData}).Return(e)

	nbIdData, err = proto.Marshal(&entities.NbIdentity{InventoryName: name})
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal nbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_UNKNOWN.String(), []interface{}{nbIdData}).Return(e)

	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.Nil(t, rNibErr)
}

func TestSaveRanLoadInformationSuccess(t *testing.T) {
	inventoryName := "name"
	loadKey, validationErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformationSuccess - Failed to build ran load infromation key. Error: %v", validationErr)
	}

	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	ranLoadInformation := generateRanLoadInformation()
	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformation - Failed to marshal RanLoadInformation entity. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveRanLoadInformation(inventoryName, ranLoadInformation)
	assert.Nil(t, rNibErr)
}

func TestSaveRanLoadInformationMarshalNilFailure(t *testing.T) {
	inventoryName := "name2"
	w, _ := initSdlInstanceMock(namespace)

	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	err := w.SaveRanLoadInformation(inventoryName, nil)
	assert.Equal(t, expectedErr, err)
}

func TestSaveRanLoadInformationEmptyInventoryNameFailure(t *testing.T) {
	inventoryName := ""
	w, _ := initSdlInstanceMock(namespace)

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

	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	ranLoadInformation := generateRanLoadInformation()
	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformation - Failed to marshal RanLoadInformation entity. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(expectedErr)

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
	w, _ := initSdlInstanceMock(namespace)
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	nbIdentity := &entities.NbIdentity{}
	actualErr := w.SaveNodeb(nbIdentity, nil)
	assert.Equal(t, expectedErr, actualErr)
}

func TestSaveUnknownTypeEntityFailure(t *testing.T) {
	w, _ := initSdlInstanceMock(namespace)
	nbIdentity := &entities.NbIdentity{InventoryName: "name", GlobalNbId: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"}}
	nb := &entities.NodebInfo{}
	nb.Port = 5656
	nb.Ip = "localhost"
	actualErr := w.SaveNodeb(nbIdentity, nb)
	assert.IsType(t, &common.ValidationError{}, actualErr)
}

func TestSaveEntityFailure(t *testing.T) {
	name := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"

	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	gnb := entities.NodebInfo{}
	gnb.NodeType = entities.Node_GNB
	data, err := proto.Marshal(&gnb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEntityFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	nbIdentity := &entities.NbIdentity{InventoryName: name, GlobalNbId: &entities.GlobalNbId{PlmnId: plmnId, NbId: nbId}}
	setExpected := []interface{}{"RAN:" + name, data}
	setExpected = append(setExpected, "GNB:"+plmnId+":"+nbId, data)
	expectedErr := errors.New("expected error")
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(expectedErr)
	rNibErr := w.SaveNodeb(nbIdentity, &gnb)
	assert.NotEmpty(t, rNibErr)
}

func TestGetRNibWriter(t *testing.T) {
	received, _ := initSdlInstanceMock(namespace)
	assert.NotEmpty(t, received)
}

func TestSaveE2TInstanceSuccess(t *testing.T) {
	address := "10.10.2.15:9800"
	loadKey, validationErr := common.ValidateAndBuildE2TInstanceKey(address)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSuccess - Failed to build E2T Instance key. Error: %v", validationErr)
	}

	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tInstance := generateE2tInstance(address)
	data, err := json.Marshal(e2tInstance)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSuccess - Failed to marshal E2tInstance entity. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveE2TInstance(e2tInstance)
	assert.Nil(t, rNibErr)
}

func TestSaveE2TInstanceNullE2tInstanceFailure(t *testing.T) {
	w, _ := initSdlInstanceMock(namespace)
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

	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tInstance := generateE2tInstance(address)
	data, err := json.Marshal(e2tInstance)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInstanceSdlFailure - Failed to marshal E2tInstance entity. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, loadKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(expectedErr)

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
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tAddresses := []string{address}
	data, err := json.Marshal(e2tAddresses)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInfoListSuccess - Failed to marshal E2TInfoList. Error: %v", err)
	}

	var e error
	var setExpected []interface{}
	setExpected = append(setExpected, E2TAddressesKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	rNibErr := w.SaveE2TAddresses(e2tAddresses)
	assert.Nil(t, rNibErr)
}

func TestSaveE2TAddressesSdlFailure(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tAddresses := []string{address}
	data, err := json.Marshal(e2tAddresses)

	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveE2TInfoListSdlFailure - Failed to marshal E2TInfoList. Error: %v", err)
	}

	expectedErr := errors.New("expected error")
	var setExpected []interface{}
	setExpected = append(setExpected, E2TAddressesKey, data)
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(expectedErr)

	rNibErr := w.SaveE2TAddresses(e2tAddresses)
	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestRemoveE2TInstanceSuccess(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tAddresses := []string{fmt.Sprintf("E2TInstance:%s", address)}
	var e error
	sdlInstanceMock.On("Remove", e2tAddresses).Return(e)

	rNibErr := w.RemoveE2TInstance(address)
	assert.Nil(t, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

func TestRemoveE2TInstanceSdlFailure(t *testing.T) {
	address := "10.10.2.15:9800"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	e2tAddresses := []string{fmt.Sprintf("E2TInstance:%s", address)}
	expectedErr := errors.New("expected error")
	sdlInstanceMock.On("Remove", e2tAddresses).Return(expectedErr)

	rNibErr := w.RemoveE2TInstance(address)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestRemoveE2TInstanceEmptyAddressFailure(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	rNibErr := w.RemoveE2TInstance("")
	assert.IsType(t, &common.ValidationError{}, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

func TestUpdateNodebInfoOnConnectionStatusInversionSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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

	sdlInstanceMock.On("SetAndPublish", []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)
	assert.Nil(t, rNibErr)
}

func TestUpdateNodebInfoOnConnectionStatusInversionMissingInventoryNameFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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

	sdlInstanceMock.On("SetAndPublish", []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)

	assert.NotNil(t, rNibErr)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateNodebInfoOnConnectionStatusInversionMissingGlobalNbId(t *testing.T) {
	inventoryName := "name"
	channelName := "RAN_CONNECTION_STATUS_CHANGE"
	eventName := inventoryName + "_" + "CONNECTED"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
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
	sdlInstanceMock.On("SetAndPublish", []string{channelName, eventName}, []interface{}{setExpected}).Return(e)

	rNibErr := w.UpdateNodebInfoOnConnectionStatusInversion(nodebInfo, eventName)

	assert.Nil(t, rNibErr)
}

func TestSaveGeneralConfiguration(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	key := common.BuildGeneralConfigurationKey()
	configurationData := "{\"enableRic\":true}"
	configuration := &entities.GeneralConfiguration{}
	configuration.EnableRic = true

	sdlInstanceMock.On("Set", []interface{}{[]interface{}{key, []byte(configurationData)}}).Return(nil)
	rNibErr := w.SaveGeneralConfiguration(configuration)

	assert.Nil(t, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

func TestSaveGeneralConfigurationDbError(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)

	key := common.BuildGeneralConfigurationKey()
	configurationData := "{\"enableRic\":true}"
	configuration := &entities.GeneralConfiguration{}
	configuration.EnableRic = true

	expectedErr := errors.New("expected error")

	sdlInstanceMock.On("Set", []interface{}{[]interface{}{key, []byte(configurationData)}}).Return(expectedErr)
	rNibErr := w.SaveGeneralConfiguration(configuration)

	assert.NotNil(t, rNibErr)
}

func TestRemoveServedCellsFailure(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCellsToRemove := generateServedCells("whatever1", "whatever2")
	expectedErr := errors.New("expected error")
	sdlInstanceMock.On("Remove", buildServedCellInfoKeysToRemove(RanName, servedCellsToRemove)).Return(expectedErr)

	rNibErr := w.RemoveServedCells(RanName, servedCellsToRemove)

	assert.NotNil(t, rNibErr)
}

func TestRemoveServedCellsSuccess(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCellsToRemove := generateServedCells("whatever1", "whatever2")
	sdlInstanceMock.On("Remove", buildServedCellInfoKeysToRemove(RanName, servedCellsToRemove)).Return(nil)
	err := w.RemoveServedCells(RanName, servedCellsToRemove)
	assert.Nil(t, err)
}

func TestUpdateEnbInvalidNodebInfoFailure(t *testing.T) {
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := &entities.NodebInfo{}
	sdlInstanceMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateEnbInvalidCellFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCells := []*entities.ServedCellInfo{{CellId: ""}}
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	sdlInstanceMock.AssertNotCalled(t, "SetAndPublish")
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.ValidationError{}, rNibErr)
}

func TestUpdateEnbSdlFailure(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	setExpected := getUpdateEnbCellsSetExpected(t, nodebInfo, servedCells)
	sdlInstanceMock.On("SetAndPublish", []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(errors.New("expected error"))
	rNibErr := w.UpdateEnb(nodebInfo, servedCells)
	assert.IsType(t, &common.InternalError{}, rNibErr)
}

func TestUpdateEnbSuccess(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	servedCells := generateServedCells("test1", "test2")
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = servedCells
	setExpected := getUpdateEnbCellsSetExpected(t, nodebInfo, servedCells)

	var e error
	sdlInstanceMock.On("SetAndPublish", []string{"RAN_MANIPULATION", inventoryName + "_" + RanUpdatedEvent}, []interface{}{setExpected}).Return(e)
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
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	var e error

	nbIdentity := &entities.NbIdentity{InventoryName: nodebInfo.RanName, GlobalNbId: nodebInfo.GetGlobalNbId()}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestRemoveEnbSuccess - Failed to marshal NbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_ENB.String(), []interface{}{nbIdData}).Return(e)

	expectedKeys := []string{}
	cell1Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[0].CellId)
	cell1PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[0].Pci)
	cell2Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[1].CellId)
	cell2PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[1].Pci)
	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	expectedKeys = append(expectedKeys, cell1Key, cell1PciKey, cell2Key, cell2PciKey, nodebNameKey, nodebIdKey)
	sdlInstanceMock.On("RemoveAndPublish", []string{channelName, eventName}, expectedKeys).Return(e)

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.Nil(t, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

func TestRemoveEnbRemoveNbIdentityError(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	nbIdentity := &entities.NbIdentity{InventoryName: nodebInfo.RanName, GlobalNbId: nodebInfo.GetGlobalNbId()}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestRemoveEnbSuccess - Failed to marshal NbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_ENB.String(), []interface{}{nbIdData}).Return(errors.New("for test"))

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.NotNil(t, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

func TestRemoveEnbRemoveAndPublishError(t *testing.T) {
	inventoryName := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"
	channelName := "RAN_MANIPULATION"
	eventName := inventoryName + "_" + "DELETED"
	w, sdlInstanceMock := initSdlInstanceMock(namespace)
	nodebInfo := generateNodebInfo(inventoryName, entities.Node_ENB, plmnId, nbId)
	nodebInfo.GetEnb().ServedCells = generateServedCellInfos("cell1", "cell2")

	var e error

	nbIdentity := &entities.NbIdentity{InventoryName: nodebInfo.RanName, GlobalNbId: nodebInfo.GetGlobalNbId()}
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestRemoveEnbSuccess - Failed to marshal NbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_ENB.String(), []interface{}{nbIdData}).Return(e)

	expectedKeys := []string{}
	cell1Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[0].CellId)
	cell1PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[0].Pci)
	cell2Key := fmt.Sprintf("CELL:%s", nodebInfo.GetEnb().ServedCells[1].CellId)
	cell2PciKey := fmt.Sprintf("PCI:%s:%02x", inventoryName, nodebInfo.GetEnb().ServedCells[1].Pci)
	nodebNameKey := fmt.Sprintf("RAN:%s", inventoryName)
	nodebIdKey := fmt.Sprintf("ENB:%s:%s", plmnId, nbId)
	expectedKeys = append(expectedKeys, cell1Key, cell1PciKey, cell2Key, cell2PciKey, nodebNameKey, nodebIdKey)
	sdlInstanceMock.On("RemoveAndPublish", []string{channelName, eventName}, expectedKeys).Return(errors.New("for test"))

	rNibErr := w.RemoveEnb(nodebInfo)
	assert.NotNil(t, rNibErr)
	sdlInstanceMock.AssertExpectations(t)
}

//Integration tests
//
//func TestSaveEnbGnbInteg(t *testing.T){
//	for i := 0; i<10; i++{
//		Init("e2Manager", 1)
//		w := GetRNibWriter()
//		nb := entities.NodebInfo{}
//		nb.NodeType = entities.Node_ENB
//		nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
//		nb.Ip = "localhost"
//		nb.Port = uint32(5656 + i)
//		enb := entities.Enb{}
//		cell1 := &entities.ServedCellInfo{CellId:fmt.Sprintf("%02x",111 + i), Pci:uint32(11 + i)}
//		cell2 := &entities.ServedCellInfo{CellId:fmt.Sprintf("%02x",222 + i), Pci:uint32(22 + i)}
//		cell3 := &entities.ServedCellInfo{CellId:fmt.Sprintf("%02x",333 + i), Pci:uint32(33 + i)}
//		enb.ServedCells = []*entities.ServedCellInfo{cell1, cell2, cell3}
//		nb.Configuration = &entities.NodebInfo_Enb{Enb:&enb}
//		plmnId := 0x02f828
//		nbId := 0x4a952a0a
//		nbIdentity := &entities.NbIdentity{InventoryName: fmt.Sprintf("nameEnb%d" ,i), GlobalNbId:&entities.GlobalNbId{RicId:fmt.Sprintf("%02x", plmnId + i), NbId:fmt.Sprintf("%02x", nbId + i)}}
//		err := w.SaveNodeb(nbIdentity, &nb)
//		if err != nil{
//			t.Errorf("#rNibWriter_test.TestSaveEnbInteg - Failed to save NodeB entity. Error: %v", err)
//		}
//
//		nb1 := entities.NodebInfo{}
//		nb1.NodeType = entities.Node_GNB
//		nb1.ConnectionStatus = entities.ConnectionStatus_CONNECTED
//		nb1.Ip = "localhost"
//		nb1.Port =  uint32(6565 + i)
//		gnb := entities.Gnb{}
//		gCell1 := &entities.ServedNRCell{ServedNrCellInformation:&entities.ServedNRCellInformation{CellId:fmt.Sprintf("%02x",1111 + i), NrPci:uint32(1 + i)}}
//		gCell2 := &entities.ServedNRCell{ServedNrCellInformation:&entities.ServedNRCellInformation{CellId:fmt.Sprintf("%02x",2222 + i), NrPci:uint32(2 + i)}}
//		gCell3 := &entities.ServedNRCell{ServedNrCellInformation:&entities.ServedNRCellInformation{CellId:fmt.Sprintf("%02x",3333 + i), NrPci:uint32(3 + i)}}
//		gnb.ServedNrCells = []*entities.ServedNRCell{gCell1, gCell2, gCell3,}
//		nb1.Configuration = &entities.NodebInfo_Gnb{Gnb:&gnb}
//		nbIdentity = &entities.NbIdentity{InventoryName: fmt.Sprintf("nameGnb%d" ,i), GlobalNbId:&entities.GlobalNbId{RicId:fmt.Sprintf("%02x", plmnId - i), NbId:fmt.Sprintf("%02x", nbId - i)}}
//		err = w.SaveNodeb(nbIdentity, &nb1)
//		if err != nil{
//			t.Errorf("#rNibWriter_test.TestSaveEnbInteg - Failed to save NodeB entity. Error: %v", err)
//		}
//	}
//}
//
//func TestSaveNbRanNamesInteg(t *testing.T){
//	for i := 0; i<10; i++{
//		Init("e2Manager", 1)
//		w := GetRNibWriter()
//		nb := entities.NodebInfo{}
//		nb.ConnectionStatus = entities.ConnectionStatus_CONNECTING
//		nb.Ip = "localhost"
//		nb.Port = uint32(5656 + i)
//		nbIdentity := &entities.NbIdentity{InventoryName: fmt.Sprintf("nameOnly%d" ,i)}
//		err := w.SaveNodeb(nbIdentity, &nb)
//		if err != nil{
//			t.Errorf("#rNibWriter_test.TestSaveEnbInteg - Failed to save NodeB entity. Error: %v", err)
//		}
//	}
//}
//
//func TestSaveRanLoadInformationInteg(t *testing.T){
//		Init("e2Manager", 1)
//		w := GetRNibWriter()
//		ranLoadInformation := generateRanLoadInformation()
//		err := w.SaveRanLoadInformation("ran_integ", ranLoadInformation)
//		if err != nil{
//			t.Errorf("#rNibWriter_test.TestSaveRanLoadInformationInteg - Failed to save RanLoadInformation entity. Error: %v", err)
//		}
//}
