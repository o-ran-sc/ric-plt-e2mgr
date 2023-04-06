//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package rmrmsghandlerprovider

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/converters"
	"e2mgr/handlers/rmrmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"strings"
	"testing"

	"e2mgr/rmrCgo"
)

/*
 * Verify support for known providers.
 */

func initTestCase(t *testing.T) (*logger.Logger, *configuration.Configuration, services.RNibDataService, *rmrsender.RmrSender, managers.IE2TInstancesManager, clients.IRoutingManagerClient, *managers.E2TAssociationManager, managers.IRanConnectStatusChangeManager, managers.RanListManager) {
	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, RnibWriter: configuration.RnibWriterConfig{StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE", RanManipulationMessageChannel: "RAN_MANIPULATION"}}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	httpClient := &mocks.HttpClientMock{}

	rmrSender := initRmrSender(&mocks.RmrMessengerMock{}, logger)
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClient)
	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager)
	return logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager
}

func TestGetNotificationHandlerSuccess(t *testing.T) {

	logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager := initTestCase(t)

	ranDisconnectionManager := managers.NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(logger, rmrSender)
	ranResetManager := managers.NewRanResetManager(logger, rnibDataService, ranConnectStatusChangeManager)

	x2SetupResponseConverter := converters.NewX2SetupResponseConverter(logger)
	x2SetupResponseManager := managers.NewX2SetupResponseManager(x2SetupResponseConverter)

	x2SetupFailureResponseConverter := converters.NewX2SetupFailureResponseConverter(logger)
	x2SetupFailureResponseManager := managers.NewX2SetupFailureResponseManager(x2SetupFailureResponseConverter)

	endcSetupResponseConverter := converters.NewEndcSetupResponseConverter(logger)
	endcSetupResponseManager := managers.NewEndcSetupResponseManager(endcSetupResponseConverter)

	endcSetupFailureResponseConverter := converters.NewEndcSetupFailureResponseConverter(logger)
	endcSetupFailureResponseManager := managers.NewEndcSetupFailureResponseManager(endcSetupFailureResponseConverter)

	var testCases = []struct {
		msgType int
		handler rmrmsghandlers.NotificationHandler
	}{
		{rmrCgo.RIC_X2_SETUP_RESP, rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupResponseManager, ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_RESP)},
		{rmrCgo.RIC_X2_SETUP_FAILURE, rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, x2SetupFailureResponseManager, ranStatusChangeManager, rmrCgo.RIC_X2_SETUP_FAILURE)},
		{rmrCgo.RIC_ENDC_X2_SETUP_RESP, rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, endcSetupResponseManager, ranStatusChangeManager, rmrCgo.RIC_ENDC_X2_SETUP_RESP)},
		{rmrCgo.RIC_ENDC_X2_SETUP_FAILURE, rmrmsghandlers.NewSetupResponseNotificationHandler(logger, rnibDataService, endcSetupFailureResponseManager, ranStatusChangeManager, rmrCgo.RIC_ENDC_X2_SETUP_FAILURE)},
		{rmrCgo.RIC_SCTP_CONNECTION_FAILURE, rmrmsghandlers.NewRanLostConnectionHandler(logger, ranDisconnectionManager)},
		//{rmrCgo.RIC_ENB_LOAD_INFORMATION, rmrmsghandlers.NewEnbLoadInformationNotificationHandler(logger, rnibDataService, converters.NewEnbLoadInformationExtractor(logger))},
		{rmrCgo.RIC_ENB_CONF_UPDATE, rmrmsghandlers.NewX2EnbConfigurationUpdateHandler(logger, rmrSender)},
		{rmrCgo.RIC_ENDC_CONF_UPDATE, rmrmsghandlers.NewEndcConfigurationUpdateHandler(logger, rmrSender)},
		{rmrCgo.RIC_E2_TERM_INIT, rmrmsghandlers.NewE2TermInitNotificationHandler(logger, ranDisconnectionManager, e2tInstancesManager, routingManagerClient)},
		{rmrCgo.E2_TERM_KEEP_ALIVE_RESP, rmrmsghandlers.NewE2TKeepAliveResponseHandler(logger, rnibDataService, e2tInstancesManager)},
		{rmrCgo.RIC_X2_RESET_RESP, rmrmsghandlers.NewX2ResetResponseHandler(logger, rnibDataService, ranStatusChangeManager, converters.NewX2ResetResponseExtractor(logger))},
		{rmrCgo.RIC_X2_RESET, rmrmsghandlers.NewX2ResetRequestNotificationHandler(logger, rnibDataService, ranStatusChangeManager, rmrSender)},
		{rmrCgo.RIC_SERVICE_UPDATE, rmrmsghandlers.NewRicServiceUpdateHandler(logger, rmrSender, rnibDataService, ranListManager)},
		{rmrCgo.RIC_E2NODE_CONFIG_UPDATE, rmrmsghandlers.NewE2nodeConfigUpdateNotificationHandler(logger, rnibDataService, rmrSender)},
		{rmrCgo.RIC_E2_RESET_REQ, rmrmsghandlers.NewE2ResetRequestNotificationHandler(logger, rnibDataService, config, rmrSender, ranResetManager)},
	}

	for _, tc := range testCases {

		provider := NewNotificationHandlerProvider()
		provider.Init(logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager)
		t.Run(fmt.Sprintf("%d", tc.msgType), func(t *testing.T) {
			handler, err := provider.GetNotificationHandler(tc.msgType)
			if err != nil {
				t.Errorf("want: handler for message type %d, got: error %s", tc.msgType, err)
			}
			//Note struct is empty, so it will match any other empty struct.
			// https://golang.org/ref/spec#Comparison_operators: Struct values are comparable if all their fields are comparable. Two struct values are equal if their corresponding non-blank fields are equal.
			if /*handler != tc.handler &&*/ strings.Compare(fmt.Sprintf("%T", handler), fmt.Sprintf("%T", tc.handler)) != 0 {
				t.Errorf("want: handler %T for message type %d, got: %T", tc.handler, tc.msgType, handler)
			}
		})
	}
}

/*
 * Verify handling of a request for an unsupported message.
 */

func TestGetNotificationHandlerFailure(t *testing.T) {

	var testCases = []struct {
		msgType   int
		errorText string
	}{
		{9999 /*unknown*/, "notification handler not found"},
	}
	for _, tc := range testCases {

		logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager := initTestCase(t)
		provider := NewNotificationHandlerProvider()
		provider.Init(logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager)
		t.Run(fmt.Sprintf("%d", tc.msgType), func(t *testing.T) {
			_, err := provider.GetNotificationHandler(tc.msgType)
			if err == nil {
				t.Errorf("want: no handler for message type %d, got: success", tc.msgType)
			}
			if !strings.Contains(fmt.Sprintf("%s", err), tc.errorText) {
				t.Errorf("want: error [%s] for message type %d, got: %s", tc.errorText, tc.msgType, err)
			}
		})
	}
}

// TODO: extract to test_utils
func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
