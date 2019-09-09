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
	"e2mgr/configuration"
	"e2mgr/handlers/rmrmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/services"
	"e2mgr/sessions"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"strings"
	"testing"

	"e2mgr/rmrCgo"
)

/*
 * Verify support for known providers.
 */

func TestGetNotificationHandlerSuccess(t *testing.T) {

	logger := initLog(t)
	rmrService := getRmrService(&mocks.RmrMessengerMock{}, logger)

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibReaderProvider, rnibWriterProvider, rmrService)

	var testCases = []struct {
		msgType int
		handler rmrmsghandlers.NotificationHandler
	}{
		{rmrCgo.RIC_X2_SETUP_RESP, rmrmsghandlers.NewSetupResponseNotificationHandler(rnibReaderProvider, rnibWriterProvider, managers.NewX2SetupResponseManager(), "X2 Setup Response")},
		{rmrCgo.RIC_X2_SETUP_FAILURE, rmrmsghandlers.NewSetupResponseNotificationHandler(rnibReaderProvider, rnibWriterProvider, managers.NewX2SetupFailureResponseManager(),"X2 Setup Failure Response")},
		{rmrCgo.RIC_ENDC_X2_SETUP_RESP, rmrmsghandlers.NewSetupResponseNotificationHandler(rnibReaderProvider, rnibWriterProvider, managers.NewEndcSetupResponseManager(),"ENDC Setup Response")},
		{rmrCgo.RIC_ENDC_X2_SETUP_FAILURE, rmrmsghandlers.NewSetupResponseNotificationHandler(rnibReaderProvider, rnibWriterProvider, managers.NewEndcSetupFailureResponseManager(),"ENDC Setup Failure Response"),},
		{rmrCgo.RIC_SCTP_CONNECTION_FAILURE, rmrmsghandlers.NewRanLostConnectionHandler(ranReconnectionManager)},
		{rmrCgo.RIC_ENB_LOAD_INFORMATION, rmrmsghandlers.NewEnbLoadInformationNotificationHandler(rnibWriterProvider)},
		{rmrCgo.RIC_ENB_CONF_UPDATE, rmrmsghandlers.X2EnbConfigurationUpdateHandler{}},
		{rmrCgo.RIC_ENDC_CONF_UPDATE, rmrmsghandlers.EndcConfigurationUpdateHandler{}},
		{rmrCgo.RIC_E2_TERM_INIT, rmrmsghandlers.NewE2TermInitNotificationHandler(ranReconnectionManager, rnibReaderProvider)},
	}

	for _, tc := range testCases {

		provider := NewNotificationHandlerProvider(rnibReaderProvider, rnibWriterProvider, ranReconnectionManager)
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

	logger := initLog(t)
	rmrService := getRmrService(&mocks.RmrMessengerMock{}, logger)

	var testCases = []struct {
		msgType   int
		errorText string
	}{
		{9999 /*unknown*/, "notification handler not found"},
	}
	for _, tc := range testCases {
		readerMock := &mocks.RnibReaderMock{}
		rnibReaderProvider := func() reader.RNibReader {
			return readerMock
		}
		writerMock := &mocks.RnibWriterMock{}
		rnibWriterProvider := func() rNibWriter.RNibWriter {
			return writerMock
		}

		ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibReaderProvider, rnibWriterProvider, rmrService)

		provider := NewNotificationHandlerProvider(rnibReaderProvider, rnibWriterProvider, ranReconnectionManager)
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
func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *services.RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return services.NewRmrService(services.NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger, make(sessions.E2Sessions), messageChannel)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
