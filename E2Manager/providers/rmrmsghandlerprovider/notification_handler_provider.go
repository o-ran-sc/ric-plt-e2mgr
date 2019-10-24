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
	"e2mgr/converters"
	"e2mgr/handlers/rmrmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"fmt"
)

type NotificationHandlerProvider struct {
	notificationHandlers map[int]rmrmsghandlers.NotificationHandler
}

func NewNotificationHandlerProvider(logger *logger.Logger, rnibDataService services.RNibDataService, ranReconnectionManager *managers.RanReconnectionManager, ranStatusChangeManager *managers.RanStatusChangeManager, rmrSender *rmrsender.RmrSender, x2SetupResponseManager *managers.X2SetupResponseManager, x2SetupFailureResponseManager *managers.X2SetupFailureResponseManager) *NotificationHandlerProvider {
	return &NotificationHandlerProvider{
		notificationHandlers: initNotificationHandlersMap(logger, rnibDataService, ranReconnectionManager, ranStatusChangeManager, rmrSender, x2SetupResponseManager, x2SetupFailureResponseManager),
	}
}

func initNotificationHandlersMap(logger *logger.Logger, rnibDataService services.RNibDataService, ranReconnectionManager *managers.RanReconnectionManager, ranStatusChangeManager *managers.RanStatusChangeManager, rmrSender *rmrsender.RmrSender, x2SetupResponseManager *managers.X2SetupResponseManager, x2SetupFailureResponseManager *managers.X2SetupFailureResponseManager) map[int]rmrmsghandlers.NotificationHandler {
	return map[int]rmrmsghandlers.NotificationHandler{
		rmrCgo.RIC_X2_SETUP_RESP:           rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupResponseManager, ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP),
		rmrCgo.RIC_X2_SETUP_FAILURE:        rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupFailureResponseManager, nil, rmrCgo.RIC_X2_SETUP_FAILURE),
		rmrCgo.RIC_ENDC_X2_SETUP_RESP:      rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, managers.NewEndcSetupResponseManager(), ranStatusChangeManager, rmrCgo.RIC_ENDC_X2_SETUP_RESP),
		rmrCgo.RIC_ENDC_X2_SETUP_FAILURE:   rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, managers.NewEndcSetupFailureResponseManager(), nil, rmrCgo.RIC_ENDC_X2_SETUP_FAILURE),
		rmrCgo.RIC_SCTP_CONNECTION_FAILURE: rmrmsghandlers.NewRanLostConnectionHandler(logger, ranReconnectionManager),
		rmrCgo.RIC_ENB_LOAD_INFORMATION:    rmrmsghandlers.NewEnbLoadInformationNotificationHandler(logger, rnibDataService, converters.NewEnbLoadInformationExtractor(logger)),
		rmrCgo.RIC_ENB_CONF_UPDATE:         rmrmsghandlers.NewX2EnbConfigurationUpdateHandler(logger, rmrSender),
		rmrCgo.RIC_ENDC_CONF_UPDATE:        rmrmsghandlers.NewEndcConfigurationUpdateHandler(logger, rmrSender),
		rmrCgo.RIC_X2_RESET_RESP:           rmrmsghandlers.NewX2ResetResponseHandler(logger, rnibDataService, ranStatusChangeManager, converters.NewX2ResetResponseExtractor(logger)),
		rmrCgo.RIC_X2_RESET:                rmrmsghandlers.NewX2ResetRequestNotificationHandler(logger, rnibDataService, ranStatusChangeManager, rmrSender),
		rmrCgo.RIC_E2_TERM_INIT:            rmrmsghandlers.NewE2TermInitNotificationHandler(logger, ranReconnectionManager, rnibDataService),
	}
}

func (provider NotificationHandlerProvider) GetNotificationHandler(messageType int) (rmrmsghandlers.NotificationHandler, error) {
	handler, ok := provider.notificationHandlers[messageType]

	if !ok {
		return nil, fmt.Errorf("notification handler not found for message %d", messageType)
	}

	return handler, nil
}
