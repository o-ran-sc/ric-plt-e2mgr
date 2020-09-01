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
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"sync"
	"time"
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
	UpdateHealthcheckTimeStampReceived(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity)
	UpdateHealthcheckTimeStampSent(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity)
	UpdateNbIdentities(nodeType entities.Node_Type, oldNbIdentities []*entities.NbIdentity, newNbIdentities []*entities.NbIdentity) error
}

func NewRanListManager(logger *logger.Logger, rnibDataService services.RNibDataService) RanListManager {
	return &ranListManagerInstance{
		logger:          logger,
		rnibDataService: rnibDataService,
		nbIdentityMap:   make(map[string]*entities.NbIdentity),
	}
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
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logger.Infof("#ranListManagerInstance.UpdateNbIdentityConnectionStatus - RAN name: %s - updating nodeb identity connection status", ranName)

	oldNbIdentity, ok := m.nbIdentityMap[ranName]
	if !ok {
		m.logger.Errorf("#ranListManagerInstance.UpdateNbIdentityConnectionStatus - RAN name: %s - nodeb identity not found in nbIdentityMap", ranName)
		return e2managererrors.NewInternalError()
	}

	newNbIdentity := &entities.NbIdentity{
		GlobalNbId:       oldNbIdentity.GlobalNbId,
		InventoryName:    ranName,
		ConnectionStatus: connectionStatus,
		HealthCheckTimestampSent: oldNbIdentity.HealthCheckTimestampSent,
		HealthCheckTimestampReceived: oldNbIdentity.HealthCheckTimestampReceived,
	}
	m.nbIdentityMap[ranName] = newNbIdentity

	err := m.rnibDataService.UpdateNbIdentity(nodeType, oldNbIdentity, newNbIdentity)
	if err != nil {
		m.logger.Errorf("#ranListManagerInstance.UpdateNbIdentityConnectionStatus - RAN name: %s - Failed updating nodeb identity in DB. error: %s", ranName, err)
		return err
	}
	m.logger.Infof("#ranListManagerInstance.UpdateNbIdentityConnectionStatus - RAN name: %s - Successfully updated nodeb identity", ranName)
	return nil
}

func (m *ranListManagerInstance) RemoveNbIdentity(nodeType entities.Node_Type, ranName string) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.logger.Infof("#ranListManagerInstance.RemoveNbIdentity - RAN name: %s - deleting nodeb identity from memory and db...", ranName)

	nbIdentity, ok := m.nbIdentityMap[ranName]
	if !ok {
		m.logger.Infof("#ranListManagerInstance.RemoveNbIdentity - RAN name: %s - nodeb identity not found", ranName)
		return nil
	}

	delete(m.nbIdentityMap, ranName)

	err := m.rnibDataService.RemoveNbIdentity(nodeType, nbIdentity)
	if err != nil {
		m.logger.Errorf("#ranListManagerInstance.RemoveNbIdentity - RAN name: %s - Failed removing nodeb identity from DB. error: %s", ranName, err)
		return err
	}

	m.logger.Infof("#ranListManagerInstance.RemoveNbIdentity - RAN name: %s - Successfully deleted nodeb identity", ranName)
	return nil
}

func (m *ranListManagerInstance) GetNbIdentityList() []*entities.NbIdentity {
	nbIds := make([]*entities.NbIdentity, 0, len(m.nbIdentityMap))
	for _, v := range m.nbIdentityMap {
		nbIds = append(nbIds, v)
	}

	m.logger.Infof("#ranListManagerInstance.GetNbIdentityList - %d identity returned", len(nbIds))

	return nbIds
}

func (m *ranListManagerInstance) UpdateHealthcheckTimeStampSent(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity){
	currentTimeStamp := time.Now().UnixNano()
	oldNbIdentity := m.nbIdentityMap[oldRRanName]

	newNbIdentity := &entities.NbIdentity{
		GlobalNbId:       oldNbIdentity.GlobalNbId,
		InventoryName:    oldNbIdentity.InventoryName,
		ConnectionStatus: oldNbIdentity.ConnectionStatus,
		HealthCheckTimestampSent: currentTimeStamp,
		HealthCheckTimestampReceived: oldNbIdentity.HealthCheckTimestampReceived,
	}

	m.nbIdentityMap[oldNbIdentity.InventoryName] = newNbIdentity
	return oldNbIdentity, newNbIdentity
}

func (m *ranListManagerInstance) UpdateHealthcheckTimeStampReceived(oldRRanName string) (*entities.NbIdentity, *entities.NbIdentity){
	currentTimeStamp := time.Now().UnixNano()
	oldNbIdentity := m.nbIdentityMap[oldRRanName]

	newNbIdentity := &entities.NbIdentity{
		GlobalNbId:       oldNbIdentity.GlobalNbId,
		InventoryName:    oldNbIdentity.InventoryName,
		ConnectionStatus: oldNbIdentity.ConnectionStatus,
		HealthCheckTimestampSent: oldNbIdentity.HealthCheckTimestampSent,
		HealthCheckTimestampReceived: currentTimeStamp,
	}

	m.nbIdentityMap[oldNbIdentity.InventoryName] = newNbIdentity
	return oldNbIdentity, newNbIdentity
}

func (m *ranListManagerInstance) UpdateNbIdentities(nodeType entities.Node_Type, oldNbIdentities []*entities.NbIdentity, newNbIdentities []*entities.NbIdentity) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	err:= m.rnibDataService.UpdateNbIdentities(nodeType, oldNbIdentities, newNbIdentities)

	return err
}
