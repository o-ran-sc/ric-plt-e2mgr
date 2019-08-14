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

package providers

import (
	"e2mgr/handlers"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type NotificationHandlerProvider struct{
	notificationHandlers map[int]handlers.NotificationHandler
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewNotificationHandlerProvider(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) *NotificationHandlerProvider {
	return &NotificationHandlerProvider{
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
		notificationHandlers: initNotificationHandlersMap(rnibReaderProvider, rnibWriterProvider),
	}
}

func initNotificationHandlersMap(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) map[int]handlers.NotificationHandler{
	return  map[int]handlers.NotificationHandler{
		//TODO change handlers.NotificationHandler to *handlers.NotificationHandler
		rmrCgo.RIC_X2_SETUP_RESP:           handlers.X2SetupResponseNotificationHandler{},
		rmrCgo.RIC_X2_SETUP_FAILURE:        handlers.X2SetupFailureResponseNotificationHandler{},
		rmrCgo.RIC_ENDC_X2_SETUP_RESP:      handlers.EndcX2SetupResponseNotificationHandler{},
		rmrCgo.RIC_ENDC_X2_SETUP_FAILURE:   handlers.EndcX2SetupFailureResponseNotificationHandler{},
		rmrCgo.RIC_SCTP_CONNECTION_FAILURE: handlers.NewRanLostConnectionHandler(rnibReaderProvider, rnibWriterProvider),
		rmrCgo.RIC_ENB_LOAD_INFORMATION:    handlers.RicEnbLoadInformationNotificationHandler{},
		rmrCgo.RIC_ENB_CONF_UPDATE:    		handlers.X2EnbConfigurationUpdateHandler{},
		rmrCgo.RIC_ENDC_CONF_UPDATE:    	handlers.EndcConfigurationUpdateHandler{},
		rmrCgo.RIC_X2_RESET_RESP:			handlers.NewX2ResetResponseHandler(rnibReaderProvider),
	}
}

func (provider NotificationHandlerProvider) GetNotificationHandler(messageType int) (handlers.NotificationHandler, error) {
	handler, ok := provider.notificationHandlers[messageType]

	if !ok {
		return nil, fmt.Errorf("notification handler not found for message %d",messageType)
	}

	return handler, nil

}
