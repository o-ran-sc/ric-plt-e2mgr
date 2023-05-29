//
// Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.
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

package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/services"
	"e2mgr/tests"
	"testing"
)

const (
	E2ResetXmlPath            = "../../tests/resources/reset/reset-request.xml"
	StateChangeMessageChannel = "RAN_CONNECTION_STATUS_CHANGE"
)

func initE2ResetMocks(t *testing.T) (*E2ResetRequestHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{
		RnibRetryIntervalMs:       10,
		MaxRnibConnectionAttempts: 3,
		E2ResetTimeOutSec:         10,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel: StateChangeMessageChannel,
		}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := &mocks.RanAlarmServiceMock{}
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	ranResetManager := managers.NewRanResetManager(logger, rnibDataService, ranConnectStatusChangeManager)
	handler := NewE2ResetRequestHandler(logger, rmrSender, rnibDataService, ranResetManager)
	return handler, readerMock, writerMock, rmrMessengerMock
}
