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

package rNibWriter

import (
	"e2mgr/mocks"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitRNibWriter(t *testing.T) {
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	available, created := writerPool.Stats()
	assert.Equal(t, available, 0, "number of available objects in the writerPool should be 0")
	assert.Equal(t, created, 0, "number of created objects in the writerPool should be 0")
	w := GetRNibWriter()
	assert.NotNil(t, w)
}

func TestInitPool(t *testing.T) {
	writerPool = nil
	sdlInstanceMock := new(mocks.MockSdlInstance)
	initPool(1, func() interface{} {
		sdlI := common.ISdlInstance(sdlInstanceMock)
		return &rNibWriterInstance{sdl: &sdlI, namespace: namespace}
	},
		func(obj interface{}) {
		},
	)
	assert.NotNil(t, writerPool)
	assert.NotNil(t, writerPool.New)
	assert.NotNil(t, writerPool.Destroy)
	available, created := writerPool.Stats()
	assert.Equal(t, 0, available, "number of available objects in the writerPool should be 0")
	assert.Equal(t, 0, created, "number of created objects in the writerPool should be 0")
}

var namespace = "namespace"

func initSdlInstanceMock(namespace string, poolSize int) *mocks.MockSdlInstance {
	sdlInstanceMock := new(mocks.MockSdlInstance)
	initPool(poolSize, func() interface{} {
		sdlI := common.ISdlInstance(sdlInstanceMock)
		return &rNibWriterInstance{sdl: &sdlI, namespace: namespace}
	},
		func(obj interface{}) {
		},
	)
	return sdlInstanceMock
}

func TestSaveEnb(t *testing.T) {
	name := "name"
	ranName := "RAN:" + name
	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId:"aaff", Pci:3}
	cellEntity := entities.Cell{Type:entities.Cell_LTE_CELL, Cell:&entities.Cell_ServedCellInfo{ServedCellInfo:cell}}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb:&enb}
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

	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(e)

	nbIdData, err := proto.Marshal(&entities.NbIdentity{InventoryName: name})
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnb - Failed to marshal nbIdentity entity. Error: %v", err)
	}
	sdlInstanceMock.On("RemoveMember", entities.Node_UNKNOWN.String(), []interface{}{nbIdData}).Return(e)

	nbIdentity := &entities.NbIdentity{InventoryName:name, GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
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
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{Pci:3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb:&enb}

	nbIdentity := &entities.NbIdentity{InventoryName:name, GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.Equal(t, common.VALIDATION_ERROR, rNibErr.GetCode())
	assert.Equal(t, "3 VALIDATION_ERROR - #utils.ValidateAndBuildCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveEnbInventoryNameValidationFailure(t *testing.T) {
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	cell := &entities.ServedCellInfo{CellId:"aaa",Pci:3}
	enb.ServedCells = []*entities.ServedCellInfo{cell}
	nb.Configuration = &entities.NodebInfo_Enb{Enb:&enb}

	nbIdentity := &entities.NbIdentity{InventoryName:"", GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.Equal(t, common.VALIDATION_ERROR, rNibErr.GetCode())
	assert.Equal(t, "3 VALIDATION_ERROR - #utils.ValidateAndBuildNodeBNameKey - an empty inventory name received", rNibErr.Error())
}

func TestSaveEnbOnClosedPool(t *testing.T) {
	name := "name"
	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_ENB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	enb := entities.Enb{}
	nb.Configuration = &entities.NodebInfo_Enb{Enb:&enb}
	data, err := proto.Marshal(&nb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEnbOnClosedPool - Failed to marshal NodeB entity. Error: %v", err)
	}
	setExpected := []interface{}{name, data}
	var e error
	sdlInstanceMock.On("Set", setExpected).Return(e)
	writerPool.Close()
	nbIdentity := &entities.NbIdentity{}
	assert.Panics(t, func(){w.SaveNodeb(nbIdentity, &nb)})
}

func TestSaveGnbCellIdValidationFailure(t *testing.T) {
	name := "name"
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_GNB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	gnb := entities.Gnb{}
	cellInfo:= &entities.ServedNRCellInformation{NrPci:2}
	cell := &entities.ServedNRCell{ServedNrCellInformation:cellInfo}
	gnb.ServedNrCells = []*entities.ServedNRCell{cell}
	nb.Configuration = &entities.NodebInfo_Gnb{Gnb:&gnb}

	nbIdentity := &entities.NbIdentity{InventoryName:name, GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
	rNibErr := w.SaveNodeb(nbIdentity, &nb)
	assert.NotNil(t, rNibErr)
	assert.Equal(t, common.VALIDATION_ERROR, rNibErr.GetCode())
	assert.Equal(t, "3 VALIDATION_ERROR - #utils.ValidateAndBuildNrCellIdKey - an empty cell id received", rNibErr.Error())
}

func TestSaveGnb(t *testing.T) {
	name := "name"
	ranName := "RAN:" + name
	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	nb := entities.NodebInfo{}
	nb.NodeType = entities.Node_GNB
	nb.ConnectionStatus = 1
	nb.Ip = "localhost"
	nb.Port = 5656
	gnb := entities.Gnb{}
	cellInfo:= &entities.ServedNRCellInformation{NrPci:2,CellId:"ccdd"}
	cell := &entities.ServedNRCell{ServedNrCellInformation:cellInfo}
	cellEntity := entities.Cell{Type:entities.Cell_NR_CELL, Cell:&entities.Cell_ServedNrCell{ServedNrCell:cell}}
	gnb.ServedNrCells = []*entities.ServedNRCell{cell}
	nb.Configuration = &entities.NodebInfo_Gnb{Gnb:&gnb}
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
	nbIdentity := &entities.NbIdentity{InventoryName:name, GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
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

	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()


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
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()

	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	err:= w.SaveRanLoadInformation(inventoryName, nil)
	assert.Equal(t, expectedErr, err)
}

func TestSaveRanLoadInformationEmptyInventoryNameFailure(t *testing.T) {
	inventoryName := ""
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()

	err:= w.SaveRanLoadInformation(inventoryName, nil)
	assert.NotNil(t, err)
	assert.Equal(t, common.VALIDATION_ERROR, err.GetCode())
}

func TestSaveRanLoadInformationSdlFailure(t *testing.T) {
	inventoryName := "name2"

	loadKey, validationErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if validationErr != nil {
		t.Errorf("#rNibWriter_test.TestSaveRanLoadInformationSuccess - Failed to build ran load infromation key. Error: %v", validationErr)
	}


	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()

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
	assert.Equal(t, common.INTERNAL_ERROR, rNibErr.GetCode())
	assert.Equal(t, expectedErr, rNibErr.GetError())
}

func generateCellLoadInformation() *entities.CellLoadInformation {
	cellLoadInformation := entities.CellLoadInformation{}

	cellLoadInformation.CellId = "123"

	ulInterferenceOverloadIndication := entities.UlInterferenceOverloadIndication_HIGH_INTERFERENCE
	cellLoadInformation.UlInterferenceOverloadIndications = []entities.UlInterferenceOverloadIndication{ulInterferenceOverloadIndication}

	ulHighInterferenceInformation := entities.UlHighInterferenceInformation{
		TargetCellId:"456",
		UlHighInterferenceIndication:"xxx",
	}

	cellLoadInformation.UlHighInterferenceInfos = []*entities.UlHighInterferenceInformation{&ulHighInterferenceInformation }

	cellLoadInformation.RelativeNarrowbandTxPower = &entities.RelativeNarrowbandTxPower{
		RntpPerPrb:"xxx",
		RntpThreshold:entities.RntpThreshold_NEG_4,
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V1_ANT_PRT,
		PB: 1,
		PdcchInterferenceImpact:2,
		EnhancedRntp: &entities.EnhancedRntp{
			EnhancedRntpBitmap:"xxx",
			RntpHighPowerThreshold:entities.RntpThreshold_NEG_2,
			EnhancedRntpStartTime: &entities.StartTime{StartSfn:500,StartSubframeNumber:5},
		},
	}

	cellLoadInformation.AbsInformation = &entities.AbsInformation{
		Mode: entities.AbsInformationMode_ABS_INFO_FDD,
		AbsPatternInfo:"xxx",
		NumberOfCellSpecificAntennaPorts:entities.NumberOfCellSpecificAntennaPorts_V2_ANT_PRT,
		MeasurementSubset:"xxx",
	}

	cellLoadInformation.InvokeIndication = entities.InvokeIndication_ABS_INFORMATION

	cellLoadInformation.ExtendedUlInterferenceOverloadInfo = &entities.ExtendedUlInterferenceOverloadInfo{
		AssociatedSubframes:"xxx",
		ExtendedUlInterferenceOverloadIndications:cellLoadInformation.UlInterferenceOverloadIndications,
	}

	compInformationItem := &entities.CompInformationItem{
		CompHypothesisSets: []*entities.CompHypothesisSet{&entities.CompHypothesisSet{CellId: "789", CompHypothesis:"xxx"}},
		BenefitMetric:50,
	}

	cellLoadInformation.CompInformation = &entities.CompInformation{
		CompInformationItems:[]*entities.CompInformationItem{compInformationItem},
		CompInformationStartTime:&entities.StartTime{StartSfn:123,StartSubframeNumber:456},
	}

	cellLoadInformation.DynamicDlTransmissionInformation = &entities.DynamicDlTransmissionInformation{
		State: entities.NaicsState_NAICS_ACTIVE,
		TransmissionModes:"xxx",
		PB: 2,
		PAList:[]entities.PA{entities.PA_DB_NEG_3},
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
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	expectedErr := common.NewInternalError(errors.New("proto: Marshal called with nil"))
	nbIdentity := &entities.NbIdentity{}
	actualErr := w.SaveNodeb(nbIdentity, nil)
	assert.Equal(t, expectedErr, actualErr)
}

func TestSaveUnknownTypeEntityFailure(t *testing.T) {
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	expectedErr := common.NewValidationError(errors.New("#rNibWriter.saveNodeB - Unknown responding node type, entity: ip:\"localhost\" port:5656 "))
	nbIdentity := &entities.NbIdentity{InventoryName:"name", GlobalNbId:&entities.GlobalNbId{PlmnId:"02f829", NbId:"4a952a0a"}}
	nb := &entities.NodebInfo{}
	nb.Port = 5656
	nb.Ip = "localhost"
	actualErr := w.SaveNodeb(nbIdentity, nb)
	assert.Equal(t, expectedErr, actualErr)
}

func TestSaveEntityFailure(t *testing.T) {
	name := "name"
	plmnId := "02f829"
	nbId := "4a952a0a"

	writerPool = nil
	sdlInstanceMock := initSdlInstanceMock(namespace, 1)
	w := GetRNibWriter()
	gnb := entities.NodebInfo{}
	gnb.NodeType = entities.Node_GNB
	data, err := proto.Marshal(&gnb)
	if err != nil {
		t.Errorf("#rNibWriter_test.TestSaveEntityFailure - Failed to marshal NodeB entity. Error: %v", err)
	}
	nbIdentity := &entities.NbIdentity{InventoryName:name, GlobalNbId:&entities.GlobalNbId{PlmnId:plmnId, NbId:nbId}}
	setExpected := []interface{}{"RAN:" + name, data}
	setExpected = append(setExpected,"GNB:" + plmnId + ":" + nbId, data)
	expectedErr := errors.New("expected error")
	sdlInstanceMock.On("Set", []interface{}{setExpected}).Return(expectedErr)
	rNibErr := w.SaveNodeb(nbIdentity, &gnb)
	assert.NotEmpty(t, rNibErr)
}

func TestGetRNibWriterPoolNotInitializedFailure(t *testing.T) {
	writerPool = nil
	assert.Panics(t, func(){GetRNibWriter()})
}

func TestGetRNibWriter(t *testing.T) {
	writerPool = nil
	initSdlInstanceMock(namespace, 1)
	received := GetRNibWriter()
	assert.NotEmpty(t, received)
	available, created := writerPool.Stats()
	assert.Equal(t, 0, available, "number of available objects in the writerPool should be 0")
	assert.Equal(t, 1, created, "number of created objects in the writerPool should be 1")
	writerPool.Close()
}

func TestClose(t *testing.T) {
	writerPool = nil
	instanceMock := initSdlInstanceMock(namespace, 2)
	w1 := GetRNibWriter()
	w2 := GetRNibWriter()
	writerPool.Put(w1)
	writerPool.Put(w2)
	available, created := writerPool.Stats()
	assert.Equal(t, 2, available, "number of available objects in the writerPool should be 2")
	assert.Equal(t, 2, created, "number of created objects in the writerPool should be 2")
	var e error
	instanceMock.On("Close").Return(e)
	Close()
	available, created = writerPool.Stats()
	assert.Equal(t, 0, available, "number of available objects in the writerPool should be 0")
	assert.Equal(t, 0, created, "number of created objects in the writerPool should be 0")
}

func TestCloseOnClosedPoolFailure(t *testing.T) {
	writerPool = nil
	instanceMock := initSdlInstanceMock(namespace, 1)
	w1 := GetRNibWriter()
	writerPool.Put(w1)
	available, created := writerPool.Stats()
	assert.Equal(t, 1, available, "number of available objects in the writerPool should be 1")
	assert.Equal(t, 1, created, "number of created objects in the writerPool should be 1")
	var e error
	instanceMock.On("Close").Return(e)
	Close()
	assert.Panics(t, func(){Close()})
}

func TestCloseFailure(t *testing.T) {
	writerPool = nil
	instanceMock := initSdlInstanceMock(namespace, 2)
	w1 := GetRNibWriter()
	writerPool.Put(w1)
	available, created := writerPool.Stats()
	assert.Equal(t, 1, available, "number of available objects in the writerPool should be 1")
	assert.Equal(t, 1, created, "number of created objects in the writerPool should be 1")
	e := errors.New("expected error")
	instanceMock.On("Close").Return(e)
	Close()
	available, created = writerPool.Stats()
	assert.Equal(t, 0, available, "number of available objects in the writerPool should be 0")
	assert.Equal(t, 0, created, "number of created objects in the writerPool should be 0")
}

func TestInit(t *testing.T) {
	writerPool = nil
	Init("", 1)
	assert.NotNil(t, writerPool)
	assert.NotNil(t, writerPool.New)
	assert.NotNil(t, writerPool.Destroy)
	available, created := writerPool.Stats()
	assert.Equal(t, 0, available, "number of available objects in the writerPool should be 0")
	assert.Equal(t, 0, created, "number of created objects in the writerPool should be 0")
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
//		nbIdentity := &entities.NbIdentity{InventoryName: fmt.Sprintf("nameEnb%d" ,i), GlobalNbId:&entities.GlobalNbId{PlmnId:fmt.Sprintf("%02x", plmnId + i), NbId:fmt.Sprintf("%02x", nbId + i)}}
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
//		nbIdentity = &entities.NbIdentity{InventoryName: fmt.Sprintf("nameGnb%d" ,i), GlobalNbId:&entities.GlobalNbId{PlmnId:fmt.Sprintf("%02x", plmnId - i), NbId:fmt.Sprintf("%02x", nbId - i)}}
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