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
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
)

type E2TermInitNotificationHandler struct {
	logger                 *logger.Logger
	rnibDataService        services.RNibDataService
	ranReconnectionManager *managers.RanReconnectionManager
}

func NewE2TermInitNotificationHandler(logger *logger.Logger, ranReconnectionManager *managers.RanReconnectionManager, rnibDataService services.RNibDataService) E2TermInitNotificationHandler {
	return E2TermInitNotificationHandler{
		logger:                 logger,
		rnibDataService:        rnibDataService,
		ranReconnectionManager: ranReconnectionManager,
	}
}

func (h E2TermInitNotificationHandler) Handle(request *models.NotificationRequest) {

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - Handling E2_TERM_INIT")

	nbIdentityList, err := h.rnibDataService.GetListNodebIds()
	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed to get nodes list from RNIB. Error: %s", err.Error())
		return
	}

	if len(nbIdentityList) == 0 {
		h.logger.Warnf("#E2TermInitNotificationHandler.Handle - The Nodes list in RNIB is empty")
		return
	}

	for _, nbIdentity := range nbIdentityList {

		if err := h.ranReconnectionManager.ReconnectRan(nbIdentity.InventoryName); err != nil {
			h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Ran name: %s - connection attempt failure, error: %s", (*nbIdentity).GetInventoryName(), err.Error())
			_, ok := err.(*common.ResourceNotFoundError)
			if !ok {
				break
			}
		}
	}

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - Completed handling of E2_TERM_INIT")
}
