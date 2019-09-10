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

package receivers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/managers/notificationmanager"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
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
	rmrMessengerMock := &mocks.RmrMessengerMock{}

	var buf *rmrCgo.MBuf
	e := fmt.Errorf("test error")
	rmrMessengerMock.On("RecvMsg").Return(buf, e)

	go getRmrServiceReceiver(rmrMessengerMock, log).ListenAndHandle()

	time.Sleep(time.Microsecond * 10)
}

func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, messageChannel)
}

func getRmrServiceReceiver(rmrMessengerMock *mocks.RmrMessengerMock, logger *logger.Logger) *RmrServiceReceiver {
	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	rmrService := getRmrService(rmrMessengerMock, logger)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrService, rNibWriter.GetRNibWriter)
	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibReaderProvider, rnibWriterProvider, ranSetupManager)
	nManager := notificationmanager.NewNotificationManager(rnibReaderProvider, rnibWriterProvider, ranReconnectionManager)

	return NewRmrServiceReceiver(*rmrService, nManager)
}
