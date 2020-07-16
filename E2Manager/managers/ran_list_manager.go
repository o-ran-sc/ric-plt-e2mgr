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

package managers

import (
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"sync"
)

type ranListManagerInstance struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	mux             sync.Mutex
	nbIdentityMap   map[string]*entities.NbIdentity
}

type RanListManager interface {
	InitNbIdentityMap() error
	AddNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error
	UpdateNbIdentityConnectionStatus(nodeType entities.Node_Type, ranName string, connectionStatus entities.ConnectionStatus) error
	RemoveNbIdentity(nodeType entities.Node_Type, ranName string) error
	GetNbIdentityList() []*entities.NbIdentity
	UpdateRanState(nodebInfo *entities.NodebInfo) error // TODO: replace with UpdateNbIdentityConnectionStatus
}

func NewRanListManager(logger *logger.Logger, rnibDataService services.RNibDataService) RanListManager {
	return &ranListManagerInstance{
		logger:          logger,
		rnibDataService: rnibDataService,
		nbIdentityMap:   make(map[string]*entities.NbIdentity),
	}
}

// TODO: replace with UpdateNbIdentityConnectionStatus
func (m *ranListManagerInstance) UpdateRanState(nodebInfo *entities.NodebInfo) error {
	m.logger.Infof("#ranListManagerInstance.UpdateRanState - RAN name: %s - Updating state...", nodebInfo.RanName)
	return nil
}

func (m *ranListManagerInstance) InitNbIdentityMap() error {
	nbIds, err := m.rnibDataService.GetListNodebIds()

	if err != nil {
		m.logger.Errorf("#ranListManagerInstance.InitNbIdentityMap - Failed fetching RAN list from DB. error: %s", err)
		return err
	}

	for _, v := range nbIds {
		m.nbIdentityMap[v.InventoryName] = v
	}

	m.logger.Infof("#ranListManagerInstance.InitNbIdentityMap - Successfully initiated nodeb identity map")
	m.logger.Debugf("#ranListManagerInstance.InitNbIdentityMap - nodeb Identity map: %s", m.nbIdentityMap)
	return nil
}

func (m *ranListManagerInstance) AddNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.nbIdentityMap[nbIdentity.InventoryName] = nbIdentity

	err := m.rnibDataService.AddNbIdentity(nodeType, nbIdentity)

	if err != nil {
		m.logger.Errorf("#ranListManagerInstance.AddNbIdentity - RAN name: %s - Failed adding nodeb identity to DB. error: %s", nbIdentity.InventoryName, err)
		return err
	}

	m.logger.Infof("#ranListManagerInstance.AddNbIdentity - RAN name: %s - Successfully added nodeb identity", nbIdentity.InventoryName)
	m.logger.Debugf("#ranListManagerInstance.AddNbIdentity - nodeb Identity map: %s", m.nbIdentityMap)
	return nil
}

func (m *ranListManagerInstance) UpdateNbIdentityConnectionStatus(nodeType entities.Node_Type, ranName string, connectionStatus entities.ConnectionStatus) error {
	//TODO: implement
	return nil
}

func (m *ranListManagerInstance) RemoveNbIdentity(nodeType entities.Node_Type, ranName string) error {
	//TODO: implement
	return nil
}

func (m *ranListManagerInstance) GetNbIdentityList() []*entities.NbIdentity {
	nbIds := make([]*entities.NbIdentity, len(m.nbIdentityMap))
	var index = 0
	for _, v := range m.nbIdentityMap {
		nbIds[index] = v
		index++
	}

	return nbIds
}
