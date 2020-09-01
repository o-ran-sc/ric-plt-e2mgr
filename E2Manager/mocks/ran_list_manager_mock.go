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

package mocks

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/mock"
)

type RanListManagerMock struct {
	mock.Mock
}

// TODO: remove after replaced with UpdateNbIdentityConnectionStatus
func (m *RanListManagerMock) UpdateRanState(nodebInfo *entities.NodebInfo) error {

	args := m.Called(nodebInfo)
	return args.Error(0)
}

func (m *RanListManagerMock) InitNbIdentityMap() error {
	args := m.Called()
	return args.Error(0)
}

func (m *RanListManagerMock) AddNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error {
	args := m.Called(nodeType, nbIdentity)
	return args.Error(0)
}

func (m *RanListManagerMock) UpdateNbIdentityConnectionStatus(nodeType entities.Node_Type, ranName string, connectionStatus entities.ConnectionStatus) error {
	args := m.Called(nodeType, ranName, connectionStatus)
	return args.Error(0)
}

func (m *RanListManagerMock) RemoveNbIdentity(nodeType entities.Node_Type, ranName string) error {
	args := m.Called(nodeType, ranName)
	return args.Error(0)
}

func (m *RanListManagerMock) GetNbIdentityList() []*entities.NbIdentity {
	args := m.Called()
	return args.Get(0).([]*entities.NbIdentity)
}

func (m *RanListManagerMock) UpdateHealthcheckTimeStampSent(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity){
	args := m.Called(oldRRanName)
	return args.Get(0).(*entities.NbIdentity), args.Get(1).(*entities.NbIdentity)
}

func (m *RanListManagerMock) UpdateHealthcheckTimeStampReceived(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity){
	args := m.Called(oldRRanName)
	return args.Get(0).(*entities.NbIdentity), args.Get(1).(*entities.NbIdentity)
}

func (m *RanListManagerMock) UpdateNbIdentities(nodeType entities.Node_Type, oldNbIdentities []*entities.NbIdentity, newNbIdentities []*entities.NbIdentity) error{
	args:= m.Called(nodeType, oldNbIdentities, newNbIdentities)
	return args.Error(0)
}

