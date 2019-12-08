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

package managers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"math"
	"sync"
)

type E2TInstancesManager struct {
	rnibDataService services.RNibDataService
	logger          *logger.Logger
	mux             sync.Mutex
}

type IE2TInstancesManager interface {
	GetE2TInstance(e2tAddress string) (*entities.E2TInstance, error)
	GetE2TInstances() ([]*entities.E2TInstance, error)
	AddE2TInstance(e2tAddress string) error
	RemoveE2TInstance(e2tInstance *entities.E2TInstance) error
	SelectE2TInstance() (string, error)
	AssociateRan(ranName string, e2tAddress string) error
	DissociateRan(ranName string, e2tAddress string) error
}

func NewE2TInstancesManager(rnibDataService services.RNibDataService, logger *logger.Logger) *E2TInstancesManager {
	return &E2TInstancesManager{
		rnibDataService: rnibDataService,
		logger:          logger,
	}
}

func (m *E2TInstancesManager) GetE2TInstance(e2tAddress string) (*entities.E2TInstance, error) {
	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#GetE2TInstance - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
	}

	return e2tInstance, err
}

func (m *E2TInstancesManager) GetE2TInstances() ([]*entities.E2TInstance, error) {
	e2tAddresses, err := m.rnibDataService.GetE2TAddresses()

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.GetE2TInstances - Failed retrieving E2T addresses. error: %s", err)
		return nil, err
	}

	if len(e2tAddresses) == 0 {
		m.logger.Warnf("#E2TInstancesManager.GetE2TInstances - Empty E2T addresses list")
		return []*entities.E2TInstance{}, nil
	}

	e2tInstances, err := m.rnibDataService.GetE2TInstances(e2tAddresses)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.GetE2TInstances - Failed retrieving E2T instances list. error: %s", err)
		return e2tInstances, err
	}

	if len(e2tInstances) == 0 {
		m.logger.Warnf("#E2TInstancesManager.GetE2TInstances - Empty E2T instances list")
		return e2tInstances, nil
	}

	return e2tInstances, nil
}

func findActiveE2TInstanceWithMinimumAssociatedRans(e2tInstances []*entities.E2TInstance) *entities.E2TInstance {
	var minInstance *entities.E2TInstance
	minAssociatedRanCount := math.MaxInt32

	for _, v := range e2tInstances {
		if v.State == entities.Active && len(v.AssociatedRanList) < minAssociatedRanCount {
			minAssociatedRanCount = len(v.AssociatedRanList)
			minInstance = v
		}
	}

	return minInstance
}

func (m *E2TInstancesManager) AddE2TInstance(e2tAddress string) error {

	e2tInstance := entities.NewE2TInstance(e2tAddress)
	err := m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed saving E2T instance. error: %s", e2tInstance.Address, err)
		return err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tAddresses, err := m.rnibDataService.GetE2TAddresses()

	if err != nil {

		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed retrieving E2T addresses list. error: %s", e2tInstance.Address, err)
			return err
		}
	}

	e2tAddresses = append(e2tAddresses, e2tInstance.Address)

	err = m.rnibDataService.SaveE2TAddresses(e2tAddresses)

	if err != nil {
		m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed saving E2T addresses list. error: %s", e2tInstance.Address, err)
		return err
	}

	m.logger.Infof("#AddE2TInstance - E2T Instance address: %s - successfully completed", e2tInstance.Address)
	return nil
}

func (m *E2TInstancesManager) DissociateRan(ranName string, e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#DissociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	i := 0 // output index
	for _, v := range e2tInstance.AssociatedRanList {
		if v != ranName {
			// copy and increment index
			e2tInstance.AssociatedRanList[i] = v
			i++
		}
	}

	e2tInstance.AssociatedRanList = e2tInstance.AssociatedRanList[:i]

	err = m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#DissociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	return nil
}

func (m *E2TInstancesManager) RemoveE2TInstance(e2tInstance *entities.E2TInstance) error {
	return nil
}
func (m *E2TInstancesManager) SelectE2TInstance() (string, error) {

	e2tInstances, err := m.GetE2TInstances()

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.SelectE2TInstance - failed retrieving E2T instances. error: %s", err)
		return "", e2managererrors.NewRnibDbError()
	}

	if len(e2tInstances) == 0 {
		m.logger.Errorf("#E2TInstancesManager.SelectE2TInstance - No E2T instance found")
		return "", e2managererrors.NewE2TInstanceAbsenceError()
	}

	min := findActiveE2TInstanceWithMinimumAssociatedRans(e2tInstances)

	if min == nil {
		m.logger.Errorf("#SelectE2TInstance - No active E2T instance found")
		return "", e2managererrors.NewE2TInstanceAbsenceError()
	}

	return min.Address, nil
}

func (m *E2TInstancesManager) AssociateRan(ranName string, e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		return e2managererrors.NewRnibDbError()
	}

	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, ranName)

	err = m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return e2managererrors.NewRnibDbError()
	}

	return nil
}
