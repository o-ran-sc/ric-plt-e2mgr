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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package httpmsghandlerprovider

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/handlers/httpmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func getRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

func setupTest(t *testing.T) *IncomingRequestHandlerProvider {
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	config.RoutingManager.BaseUrl = "http://10.10.2.15:12020/routingmanager"
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	rmrSender := getRmrSender(rmrMessengerMock, log)
	ranSetupManager := managers.NewRanSetupManager(log, rmrSender, rnibDataService)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	ranListManager := managers.NewRanListManager(log)
	ranAlarmService := services.NewRanAlarmService(log, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(log, rnibDataService,ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager)
	return NewIncomingRequestHandlerProvider(log, rmrSender, configuration.ParseConfiguration(), rnibDataService, ranSetupManager, e2tInstancesManager, e2tAssociationManager, rmClient)
}

func TestNewIncomingRequestHandlerProvider(t *testing.T) {
	provider := setupTest(t)

	assert.NotNil(t, provider)
}

func TestShutdownRequestHandler(t *testing.T) {
	provider := setupTest(t)
	handler, err := provider.GetHandler(ShutdownRequest)

	assert.NotNil(t, provider)
	assert.Nil(t, err)

	_, ok := handler.(*httpmsghandlers.DeleteAllRequestHandler)

	assert.True(t, ok)
}

func TestSetGeneralConfigurationHandler(t *testing.T) {
	provider := setupTest(t)
	handler, err := provider.GetHandler(SetGeneralConfigurationRequest)

	assert.NotNil(t, provider)
	assert.Nil(t, err)

	_, ok := handler.(*httpmsghandlers.SetGeneralConfigurationHandler)

	assert.True(t, ok)
}

func TestGetShutdownHandlerFailure(t *testing.T) {
	provider := setupTest(t)
	_, actual := provider.GetHandler("test")
	expected := &e2managererrors.InternalError{}

	assert.NotNil(t, actual)
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
