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

package handlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type E2TermInitNotificationHandler struct {
	rnibReaderProvider func() reader.RNibReader
	ranReconnectionManager *managers.RanReconnectionManager
}

func NewE2TermInitNotificationHandler(ranReconnectionManager  *managers.RanReconnectionManager, rnibReaderProvider func() reader.RNibReader) E2TermInitNotificationHandler {
	return E2TermInitNotificationHandler{
		rnibReaderProvider: rnibReaderProvider,
		ranReconnectionManager: ranReconnectionManager,
	}
}


func (handler E2TermInitNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {
	logger.Infof("#e2_term_init_notification_handler.Handle - Received E2_TERM_INIT")
	return	//TODO: enable

	nbIdentityList, err := handler.rnibReaderProvider().GetListNodebIds()

	if err != nil {
		logger.Errorf("#e2_term_init_notification_handler.Handle - Failed to get nodes list from RNIB. Error: %s", err.Error())
		return
	}

	if len(nbIdentityList) == 0 {
		logger.Warnf("#e2_term_init_notification_handler.Handle - The Nodes list in RNIB is empty")
		return
	}

	for _,nbIdentity := range nbIdentityList{

		if err := handler.ranReconnectionManager.ReconnectRan(nbIdentity.InventoryName); err != nil {
			logger.Errorf("#e2_term_init_notification_handler.Handle - connection attempt failure, ran name: %s. Error: %s",
				(*nbIdentity).GetInventoryName(), err.Error())
			break
		}
	}
}
