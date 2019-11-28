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


package rmrreceiver

import (
	"e2mgr/configuration"
	"e2mgr/converters"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/managers/notificationmanager"
	"e2mgr/mocks"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"testing"
	"time"
)

func TestListenAndHandle(t *testing.T) {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#rmr_service_test.TestListenAndHandle - failed to initialize logger, error: %s", err)
	}
	rmrReceiver := initRmrReceiver(log)
	go rmrReceiver.ListenAndHandle()
	time.Sleep(time.Microsecond * 10)
}

func initRmrMessenger(log *logger.Logger) *rmrCgo.RmrMessenger {
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)

	// TODO: that's not good since we don't actually test anything. if the error is populated then the loop will just continue and it's sort of a "workaround" for that method to be called
	var buf *rmrCgo.MBuf
	e := fmt.Errorf("test error")
	rmrMessengerMock.On("RecvMsg").Return(buf, e)
	return &rmrMessenger
}

func initRmrReceiver(logger *logger.Logger) *RmrReceiver {
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	rnibDataService := services.NewRnibDataService(logger, config, rnibReaderProvider, rnibWriterProvider)
	rmrMessenger := initRmrMessenger(logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, ranSetupManager)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(logger, rmrSender)
	x2SetupResponseConverter := converters.NewX2SetupResponseConverter(logger)
	x2SetupResponseManager := managers.NewX2SetupResponseManager(x2SetupResponseConverter)
	x2SetupFailureResponseConverter := converters.NewX2SetupFailureResponseConverter(logger)
	x2SetupFailureResponseManager := managers.NewX2SetupFailureResponseManager(x2SetupFailureResponseConverter)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider(logger, rnibDataService, ranReconnectionManager, ranStatusChangeManager, rmrSender, x2SetupResponseManager, x2SetupFailureResponseManager )
	notificationManager := notificationmanager.NewNotificationManager(logger, rmrNotificationHandlerProvider)
	return NewRmrReceiver(logger, rmrMessenger, notificationManager)
}
