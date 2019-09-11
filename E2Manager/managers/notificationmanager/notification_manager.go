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

package notificationmanager

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"fmt"
	"time"
)

type NotificationManager struct {
	notificationHandlerProvider *rmrmsghandlerprovider.NotificationHandlerProvider
}

func NewNotificationManager(rnibDataService services.RNibDataService, ranReconnectionManager *managers.RanReconnectionManager) *NotificationManager {
	notificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider(rnibDataService, ranReconnectionManager)

	return &NotificationManager{
		notificationHandlerProvider: notificationHandlerProvider,
	}
}

//TODO add NEWHandler with log
func (m NotificationManager) HandleMessage(logger *logger.Logger, mbuf *rmrCgo.MBuf, responseChannel chan<- *models.NotificationResponse) {

	notificationHandler, err := m.notificationHandlerProvider.GetNotificationHandler(mbuf.MType)

	if err != nil {
		logger.Errorf(fmt.Sprintf("%s", err))
		return
	}

	notificationRequest := models.NewNotificationRequest(mbuf.Meid, *mbuf.Payload, time.Now(), string(*mbuf.XAction))
	go notificationHandler.Handle(logger, notificationRequest, responseChannel)
}
