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
package mocks

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"github.com/stretchr/testify/mock"
)

type RnibReaderMock struct {
	mock.Mock
}

func (m *RnibReaderMock) GetNodeb(inventoryName string) (*entities.NodebInfo, common.IRNibError) {
	args := m.Called(inventoryName)


	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*entities.NodebInfo), errArg.(common.IRNibError);
	}

	return args.Get(0).(*entities.NodebInfo), nil
}

func (m *RnibReaderMock) GetNodebByGlobalNbId(nodeType entities.Node_Type, globalNbId *entities.GlobalNbId) (*entities.NodebInfo, common.IRNibError) {
	args := m.Called(nodeType, globalNbId)

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*entities.NodebInfo), errArg.(common.IRNibError);
	}

	return args.Get(0).(*entities.NodebInfo), nil
}

func (m *RnibReaderMock)  GetCellList(inventoryName string) (*entities.Cells, common.IRNibError) {
	args := m.Called(inventoryName)

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*entities.Cells), errArg.(common.IRNibError);
	}

	return args.Get(0).(*entities.Cells), nil
}

func (m *RnibReaderMock) GetListGnbIds()(*[]*entities.NbIdentity, common.IRNibError) {
	args := m.Called()

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*[]*entities.NbIdentity), errArg.(common.IRNibError);
	}

	return args.Get(0).(*[]*entities.NbIdentity), nil
}

func (m *RnibReaderMock) GetListEnbIds()(*[]*entities.NbIdentity, common.IRNibError) {
	args := m.Called()

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*[]*entities.NbIdentity), errArg.(common.IRNibError);
	}

	return args.Get(0).(*[]*entities.NbIdentity), nil

}

func (m *RnibReaderMock) GetCountGnbList()(int, common.IRNibError) {
	args := m.Called()

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Int(0), errArg.(common.IRNibError);
	}

	return args.Int(0), nil

}

func (m *RnibReaderMock) GetCell(inventoryName string, pci uint32) (*entities.Cell, common.IRNibError) {
	args := m.Called(inventoryName, pci)

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*entities.Cell), errArg.(common.IRNibError);
	}

	return args.Get(0).(*entities.Cell), nil
}

func (m *RnibReaderMock) GetCellById(cellType entities.Cell_Type, cellId string) (*entities.Cell, common.IRNibError) {
	args := m.Called(cellType, cellId)

	errArg := args.Get(1);

	if (errArg != nil) {
		return args.Get(0).(*entities.Cell), errArg.(common.IRNibError);
	}

	return args.Get(0).(*entities.Cell), nil
}

func (m *RnibReaderMock) GetListNodebIds()([]*entities.NbIdentity, common.IRNibError){
	args := m.Called()

	errArg := args.Get(1)

	if errArg != nil {
		return args.Get(0).([]*entities.NbIdentity), errArg.(common.IRNibError)
	}

	return args.Get(0).([]*entities.NbIdentity), nil
}

func (m *RnibReaderMock) GetRanLoadInformation(inventoryName string) (*entities.RanLoadInformation, common.IRNibError){
	args := m.Called()

	errArg := args.Get(1)

	if errArg != nil {
		return args.Get(0).(*entities.RanLoadInformation), errArg.(common.IRNibError)
	}

	return args.Get(0).(*entities.RanLoadInformation), nil
}