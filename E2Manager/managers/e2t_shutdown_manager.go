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
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)


type IE2TShutdownManager interface {
	Shutdown(e2tInstance *entities.E2TInstance) error
}

type E2TShutdownManager struct {
	logger              *logger.Logger
	rnibDataService     services.RNibDataService
	e2TInstancesManager IE2TInstancesManager
}

func NewE2TShutdownManager(logger *logger.Logger, rnibDataService services.RNibDataService, e2TInstancesManager IE2TInstancesManager) E2TShutdownManager {
	return E2TShutdownManager{
		logger:              logger,
		rnibDataService:     rnibDataService,
		e2TInstancesManager: e2TInstancesManager,
	}
}

func (h E2TShutdownManager) Shutdown(e2tInstance *entities.E2TInstance) error{
	h.logger.Infof("#E2TShutdownManager.Shutdown - E2T %s is Dead, RIP", e2tInstance.Address)

	return nil
}
