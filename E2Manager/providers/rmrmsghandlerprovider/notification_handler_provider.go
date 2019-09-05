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

package rmrmsghandlerprovider

import (
	"e2mgr/handlers/rmrmsghandlers"
	"e2mgr/managers"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type NotificationHandlerProvider struct {
	notificationHandlers map[int]rmrmsghandlers.NotificationHandler
}

func NewNotificationHandlerProvider(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter, ranReconnectionManager *managers.RanReconnectionManager) *NotificationHandlerProvider {
	return &NotificationHandlerProvider{
		notificationHandlers: initNotificationHandlersMap(rnibReaderProvider, rnibWriterProvider, ranReconnectionManager),
	}
}

func initNotificationHandlersMap(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter, ranReconnectionManager *managers.RanReconnectionManager) map[int]rmrmsghandlers.NotificationHandler {
	return map[int]rmrmsghandlers.NotificationHandler{
		//TODO change handlers.NotificationHandler to *handlers.NotificationHandler
		rmrCgo.RIC_X2_SETUP_RESP:           rmrmsghandlers.X2SetupResponseNotificationHandler{},
		rmrCgo.RIC_X2_SETUP_FAILURE:        rmrmsghandlers.X2SetupFailureResponseNotificationHandler{},
		rmrCgo.RIC_ENDC_X2_SETUP_RESP:      rmrmsghandlers.EndcX2SetupResponseNotificationHandler{},
		rmrCgo.RIC_ENDC_X2_SETUP_FAILURE:   rmrmsghandlers.EndcX2SetupFailureResponseNotificationHandler{},
		rmrCgo.RIC_SCTP_CONNECTION_FAILURE: rmrmsghandlers.NewRanLostConnectionHandler(ranReconnectionManager),
		rmrCgo.RIC_ENB_LOAD_INFORMATION:    rmrmsghandlers.NewEnbLoadInformationNotificationHandler(rnibWriterProvider),
		rmrCgo.RIC_ENB_CONF_UPDATE:         rmrmsghandlers.X2EnbConfigurationUpdateHandler{},
		rmrCgo.RIC_ENDC_CONF_UPDATE:        rmrmsghandlers.EndcConfigurationUpdateHandler{},
		rmrCgo.RIC_X2_RESET_RESP:           rmrmsghandlers.NewX2ResetResponseHandler(rnibReaderProvider),
		rmrCgo.RIC_X2_RESET:                rmrmsghandlers.NewX2ResetRequestNotificationHandler(rnibReaderProvider),
		rmrCgo.RIC_E2_TERM_INIT:            rmrmsghandlers.NewE2TermInitNotificationHandler(ranReconnectionManager, rnibReaderProvider ),
	}
}

func (provider NotificationHandlerProvider) GetNotificationHandler(messageType int) (rmrmsghandlers.NotificationHandler, error) {
	handler, ok := provider.notificationHandlers[messageType]

	if !ok {
		return nil, fmt.Errorf("notification handler not found for message %d", messageType)
	}

	return handler, nil

}
