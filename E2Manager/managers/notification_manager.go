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
	"e2mgr/models"
	"e2mgr/providers"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"time"
)

type NotificationManager struct {
	notificationHandlerProvider *providers.NotificationHandlerProvider
}

func NewNotificationManager(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) *NotificationManager {
	notificationHandlerProvider := providers.NewNotificationHandlerProvider(rnibReaderProvider, rnibWriterProvider)

	return &NotificationManager{
		notificationHandlerProvider: notificationHandlerProvider,
	}
}

//TODO add NEWHandler with log
func (m NotificationManager) HandleMessage(logger *logger.Logger, e2Sessions sessions.E2Sessions, mbuf *rmrCgo.MBuf, responseChannel chan<- *models.NotificationResponse) {

	notificationHandler, err := m.notificationHandlerProvider.GetNotificationHandler(mbuf.MType)

	if err != nil {
		logger.Errorf(fmt.Sprintf("%s", err))
		return
	}

	notificationRequest := models.NewNotificationRequest(mbuf.Meid, *mbuf.Payload, time.Now(), string(*mbuf.XAction))
	go notificationHandler.Handle(logger, e2Sessions, notificationRequest, responseChannel)
}
