package managers

import (
	"e2mgr/logger"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"sync"
)

type E2TInstancesManager struct {
	rnibDataService services.RNibDataService
	logger          *logger.Logger
	mux             sync.Mutex
}

type IE2TInstancesManager interface {
	GetE2TInstance(e2tAddress string) (*entities.E2TInstance, error)
	AddE2TInstance(e2tAddress string) error
	RemoveE2TInstance(e2tInstance *entities.E2TInstance) error
	SelectE2TInstance(e2tInstance *entities.E2TInstance) (string, error)
	AssociateRan(ranName string, e2tAddress string) error
	DeassociateRan(ranName string, e2tAddress string) error
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

func (m *E2TInstancesManager) AddE2TInstance(e2tAddress string) error {

	if len(e2tAddress) == 0 {
		m.logger.Errorf("#AddE2TInstance - Empty E2T address received")
		return fmt.Errorf("empty E2T address")
	}

	e2tInstance := entities.NewE2TInstance(e2tAddress)
	err := m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed saving E2T instance. error: %s", e2tInstance.Address, err)
		return err
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInfoList, err := m.rnibDataService.GetE2TInfoList()

	if err != nil {

		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed retrieving E2TInfoList. error: %s", e2tInstance.Address, err)
			return err
		}
	}

	e2tInstanceInfo := entities.NewE2TInstanceInfo(e2tInstance.Address)
	e2tInfoList = append(e2tInfoList, e2tInstanceInfo)

	err = m.rnibDataService.SaveE2TInfoList(e2tInfoList)

	if err != nil {
		m.logger.Errorf("#AddE2TInstance - E2T Instance address: %s - Failed saving E2TInfoList. error: %s", e2tInstance.Address, err)
		return err
	}

	m.logger.Infof("#AddE2TInstance - E2T Instance address: %s - successfully completed", e2tInstance.Address)
	return nil
}

func (m *E2TInstancesManager) DeassociateRan(ranName string, e2tAddress string) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	e2tInfoList, err := m.rnibDataService.GetE2TInfoList()

	if err != nil {
		m.logger.Errorf("#DeassociateRan - E2T Instance address: %s - Failed retrieving E2TInfoList. error: %s", e2tAddress, err)
		return err
	}

	isE2TInstanceFound := false

	for _, e2tInfoInstance := range e2tInfoList {
		if e2tInfoInstance.Address == e2tAddress {
			e2tInfoInstance.AssociatedRanCount--
			isE2TInstanceFound = true
			break
		}
	}

	if !isE2TInstanceFound {
		m.logger.Warnf("#DeassociateRan - E2T Instance address: %s - E2TInstance not found in E2TInfoList.", e2tAddress)
		return nil
	}

	err = m.rnibDataService.SaveE2TInfoList(e2tInfoList)

	if err != nil {
		m.logger.Errorf("#DeassociateRan - E2T Instance address: %s - Failed saving E2TInfoList. error: %s", e2tAddress, err)
		return err
	}

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#DeassociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
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
		m.logger.Errorf("#DeassociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	return nil
}

func (m *E2TInstancesManager) RemoveE2TInstance(e2tInstance *entities.E2TInstance) error {
	return nil
}
func (m *E2TInstancesManager) SelectE2TInstance(e2tInstance *entities.E2TInstance) (string, error) {
	return "", nil
}

func (m *E2TInstancesManager) AssociateRan(ranName string, e2tAddress string) error {

	e2tInfoList, err := m.rnibDataService.GetE2TInfoList()

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed retrieving E2TInfoList. error: %s", e2tAddress, err)
		return err
	}

	for _, e2tInfoInstance := range e2tInfoList {
		if e2tInfoInstance.Address == e2tAddress {
			e2tInfoInstance.AssociatedRanCount++
			break;
		}
	}

	err = m.rnibDataService.SaveE2TInfoList(e2tInfoList)

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed saving E2TInfoList. error: %s", e2tAddress, err)
		return err
	}

	e2tInstance, err := m.rnibDataService.GetE2TInstance(e2tAddress)

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed retrieving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, ranName)

	err = m.rnibDataService.SaveE2TInstance(e2tInstance)

	if err != nil {
		m.logger.Errorf("#AssociateRan - E2T Instance address: %s - Failed saving E2TInstance. error: %s", e2tAddress, err)
		return err
	}

	return nil
}
