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

package rmrmsghandlers

import (
	"e2mgr/clients"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/services"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type E2TermInitNotificationHandler struct {
	logger                 *logger.Logger
	rnibDataService        services.RNibDataService
	ranReconnectionManager *managers.RanReconnectionManager
	e2tInstancesManager    managers.IE2TInstancesManager
	routingManagerClient   clients.IRoutingManagerClient
}

func NewE2TermInitNotificationHandler(logger *logger.Logger, ranReconnectionManager *managers.RanReconnectionManager, rnibDataService services.RNibDataService, e2tInstancesManager managers.IE2TInstancesManager, routingManagerClient clients.IRoutingManagerClient) E2TermInitNotificationHandler {
	return E2TermInitNotificationHandler{
		logger:                 logger,
		rnibDataService:        rnibDataService,
		ranReconnectionManager: ranReconnectionManager,
		e2tInstancesManager:    e2tInstancesManager,
		routingManagerClient:   routingManagerClient,
	}
}

func (h E2TermInitNotificationHandler) Handle(request *models.NotificationRequest) {

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - Handling E2_TERM_INIT")

	unmarshalledPayload := models.E2TermInitPayload{}
	err := json.Unmarshal(request.Payload, &unmarshalledPayload)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler - Error unmarshaling E2 Term Init payload: %s", err)
		return
	}

	e2tAddress := unmarshalledPayload.Address

	if len(e2tAddress) == 0 {
		h.logger.Errorf("#E2TermInitNotificationHandler - Empty E2T address received")
		return
	}

	e2tInstance, err := h.e2tInstancesManager.GetE2TInstance(e2tAddress)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
			return
		}

		h.HandleNewE2TInstance(e2tAddress)
		return
	}

	if len(e2tInstance.AssociatedRanList) == 0 {
		h.logger.Infof("#E2TermInitNotificationHandler.Handle - E2T Address: %s - E2T instance has no associated RANs", e2tInstance.Address)
		return
	}

	if e2tInstance.State == entities.ToBeDeleted{
		h.logger.Infof("#E2TermInitNotificationHandler.Handle - E2T Address: %s - E2T instance status is: %s, ignore", e2tInstance.Address, e2tInstance.State)
		return
	}

	if e2tInstance.State == entities.RoutingManagerFailure {
		err := h.e2tInstancesManager.ActivateE2TInstance(e2tInstance)
		if err != nil {
			return
		}
	}

	h.HandleExistingE2TInstance(e2tInstance)

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - Completed handling of E2_TERM_INIT")
}

func (h E2TermInitNotificationHandler) HandleExistingE2TInstance(e2tInstance *entities.E2TInstance) {

	for _, ranName := range e2tInstance.AssociatedRanList {

		if err := h.ranReconnectionManager.ReconnectRan(ranName); err != nil {
			h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Ran name: %s - connection attempt failure, error: %s", ranName, err)
			_, ok := err.(*common.ResourceNotFoundError)
			if !ok {
				break
			}
		}
	}
}

func (h E2TermInitNotificationHandler) HandleNewE2TInstance(e2tAddress string) {

	err := h.routingManagerClient.AddE2TInstance(e2tAddress)

	if err != nil{
		h.logger.Errorf("#E2TermInitNotificationHandler.HandleNewE2TInstance - e2t address: %s - routing manager call failure, error: %s", e2tAddress, err)
		return
	}

	_ = h.e2tInstancesManager.AddE2TInstance(e2tAddress)
}