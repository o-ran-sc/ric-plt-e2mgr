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
	"time"
)

type E2TInstancesManager struct {
	rnibDataService services.RNibDataService
	logger          *logger.Logger
	mux             sync.Mutex
}

type IE2TInstancesManager interface {
	GetE2TInstance(e2tAddress string) (*entities.E2TInstance, error)
	GetE2TInstances() ([]*entities.E2TInstance, error)
	GetE2TInstancesNoLogs() ([]*entities.E2TInstance, error)
	AddE2TInstance(e2tAddress string) error
	RemoveE2TInstance(e2tInstance *entities.E2TInstance) error
	SelectE2TInstance() (string, error)
	AssociateRan(ranName string, e2tAddress string) error
	DissociateRan(ranName string, e2tAddress string) error
	ActivateE2TInstance(e2tInstance *entities.E2TInstance) error
	ResetKeepAliveTimestamp(e2tAddress string) error
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

		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#E2TInstancesManager.GetE2TInstance - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		} else {
			m.logger.Infof("#E2TInstancesManager.GetE2TInstance - E2T Instance address: %s not found on DB", e2tAddress)
		}
	}

	return e2tInstance, err
}

func (m *E2TInstancesManager) GetE2TInstancesNoLogs() ([]*entities.E2TInstance, error) {
	e2tAddresses, err := m.rnibDataService.GetE2TAddressesNoLogs()

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#E2TInstancesManager.GetE2TInstancesNoLogs - Failed retrieving E2T addresses. error: %s", err)
		}
		return nil, err
	}

	if len(e2tAddresses) == 0 {
		return []*entities.E2TInstance{}, nil
	}

	e2tInstances, err := m.rnibDataService.GetE2TInstancesNoLogs(e2tAddresses)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#E2TInstancesManager.GetE2TInstancesNoLogs - Failed retrieving E2T instances list. error: %s", err)
		}
		return e2tInstances, err
	}

	return e2tInstances, nil
}

func (m *E2TInstancesManager) GetE2TInstances() ([]*entities.E2TInstance, error) {
	e2tAddresses, err := m.rnibDataService.GetE2TAddresses()

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.GetE2TInstances - Failed retrieving E2T addresses. error: %s", err)
		return nil, err
	}

	if len(e2tAddresses) == 0 {
		m.logger.Infof("#E2TInstancesManager.GetE2TInstances - Empty E2T addresses list")
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

func (m *E2TInstancesManager) ResetKeepAliveTimestampsForAllE2TInstances() {

	e2tInstances, err := m.GetE2TInstances()

	if err != nil {
		m.logger.Errorf("E2TInstancesManager.ResetKeepAliveTimestampForAllE2TInstances - Couldn't reset timestamps due to a DB error")
		return
	}

	if len(e2tInstances) == 0 {
		m.logger.Infof("E2TInstancesManager.ResetKeepAliveTimestampForAllE2TInstances - No instances, ignoring reset")
		return
	}

	for _, v := range e2tInstances {

		if v.State != entities.Active {
			continue
		}

		v.KeepAliveTimestamp = time.Now().UnixNano()

		err := m.rnibDataService.SaveE2TInstance(v)

		if err != nil {
			m.logger.Errorf("E2TInstancesManager.ResetKeepAliveTimestampForAllE2TInstances - E2T address: %s - failed resetting e2t instance keep alive timestamp. error: %s", v.Address, err)
		}
	}

	m.logger.Infof("E2TInstancesManager.ResetKeepAliveTimestampForAllE2TInstances - Done with reset")
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

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance := entities.NewE2TInstance(e2tAddress)
	err := m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.AddE2TInstance - E2T Instance address: %s - Failed saving E2T instance. error: %s", e2tInstance.Address, err)
		return err
	}

	e2tAddresses, err := m.rnibDataService.GetE2TAddresses()

	if err != nil {

		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#E2TInstancesManager.AddE2TInstance - E2T Instance address: %s - Failed retrieving E2T addresses list. error: %s", e2tInstance.Address, err)
			return err
		}
	}

	e2tAddresses = append(e2tAddresses, e2tInstance.Address)

	err = m.rnibDataService.SaveE2TAddresses(e2tAddresses)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.AddE2TInstance - E2T Instance address: %s - Failed saving E2T addresses list. error: %s", e2tInstance.Address, err)
		return err
	}

	m.logger.Infof("#E2TInstancesManager.AddE2TInstance - E2T Instance address: %s - successfully added E2T instance", e2tInstance.Address)
	return nil
}

func (m *E2TInstancesManager) DissociateRan(ranName string, e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.DissociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
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
		m.logger.Errorf("#E2TInstancesManager.DissociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	m.logger.Infof("#E2TInstancesManager.DissociateRan - successfully dissociated RAN %s from E2T %s", ranName, e2tInstance.Address)
	return nil
}

func (m *E2TInstancesManager) RemoveE2TInstance(e2tInstance *entities.E2TInstance) error {
	return nil
}
func (m *E2TInstancesManager) SelectE2TInstance() (string, error) {

	e2tInstances, err := m.GetE2TInstances()

	if err != nil {
		return "", e2managererrors.NewRnibDbError()
	}

	if len(e2tInstances) == 0 {
		m.logger.Errorf("#E2TInstancesManager.SelectE2TInstance - No E2T instance found")
		return "", e2managererrors.NewE2TInstanceAbsenceError()
	}

	min := findActiveE2TInstanceWithMinimumAssociatedRans(e2tInstances)

	if min == nil {
		m.logger.Errorf("#E2TInstancesManager.SelectE2TInstance - No active E2T instance found")
		return "", e2managererrors.NewE2TInstanceAbsenceError()
	}

	m.logger.Infof("#E2TInstancesManager.SelectE2TInstance - successfully selected E2T instance. address: %s", min.Address)
	return min.Address, nil
}

func (m *E2TInstancesManager) AssociateRan(ranName string, e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.AssociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		return e2managererrors.NewRnibDbError()
	}

	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, ranName)

	err = m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.AssociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return e2managererrors.NewRnibDbError()
	}

	m.logger.Infof("#E2TInstancesManager.AssociateRan - successfully associated RAN %s with E2T %s", ranName, e2tInstance.Address)
	return nil
}

func (h E2TInstancesManager) ActivateE2TInstance(e2tInstance *entities.E2TInstance) error{

	if e2tInstance == nil {
		h.logger.Errorf("#E2TInstancesManager.ActivateE2TInstance - e2tInstance empty")
		return e2managererrors.NewInternalError()
	}

	h.logger.Infof("#E2TInstancesManager.ActivateE2TInstance - E2T Address: %s - activate E2T instance", e2tInstance.Address)

	e2tInstance.State = entities.Active
	e2tInstance.KeepAliveTimestamp = time.Now().UnixNano()

	err := h.rnibDataService.SaveE2TInstance(e2tInstance)
	if err != nil {
		h.logger.Errorf("#E2TInstancesManager.ActivateE2TInstance - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tInstance.Address, err)
		return err
	}
	return nil
}

func (m *E2TInstancesManager) ResetKeepAliveTimestamp(e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInstance, err := m.rnibDataService.GetE2TInstanceNoLogs(e2tAddress)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.ResetKeepAliveTimestamp - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	if e2tInstance.State == entities.ToBeDeleted || e2tInstance.State == entities.RoutingManagerFailure {
		m.logger.Warnf("#E2TInstancesManager.ResetKeepAliveTimestamp - Ignore. This Instance is about to deleted")
		return nil

	}

	e2tInstance.KeepAliveTimestamp = time.Now().UnixNano()
	err = m.rnibDataService.SaveE2TInstanceNoLogs(e2tInstance)

	if err != nil {
		m.logger.Errorf("#E2TInstancesManager.ResetKeepAliveTimestamp - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	return nil
}
