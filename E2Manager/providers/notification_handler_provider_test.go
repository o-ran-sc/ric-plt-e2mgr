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
	"e2mgr/mocks"
	"e2mgr/rNibWriter"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"strings"
	"testing"

	"e2mgr/handlers"
	"e2mgr/rmrCgo"
)

/*
 * Verify support for known providers.
 */

func TestGetNotificationHandlerSuccess(t *testing.T) {
	readerMock :=&mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	var testCases = []struct {
		msgType int
		handler handlers.NotificationHandler
	}{
		{rmrCgo.RIC_X2_SETUP_RESP /*successful x2 setup response*/, handlers.X2SetupResponseNotificationHandler{}},
		{rmrCgo.RIC_X2_SETUP_FAILURE /*unsuccessful x2 setup response*/, handlers.X2SetupFailureResponseNotificationHandler{}},
		{rmrCgo.RIC_ENDC_X2_SETUP_RESP /*successful en-dc x2 setup response*/,handlers.EndcX2SetupResponseNotificationHandler{}},
		{rmrCgo.RIC_ENDC_X2_SETUP_FAILURE /*unsuccessful en-dc x2 setup response*/,handlers.EndcX2SetupFailureResponseNotificationHandler{}},
		{rmrCgo.RIC_SCTP_CONNECTION_FAILURE /*sctp errors*/, handlers.NewRanLostConnectionHandler(rnibReaderProvider, rnibWriterProvider)},
		{rmrCgo.RIC_ENB_LOAD_INFORMATION, handlers.RicEnbLoadInformationNotificationHandler{}},
		{rmrCgo.RIC_ENB_CONF_UPDATE, handlers.X2EnbConfigurationUpdateHandler{}},
		{rmrCgo.RIC_ENDC_CONF_UPDATE, handlers.EndcConfigurationUpdateHandler{}},
	}
	for _, tc := range testCases {

		provider := NewNotificationHandlerProvider(rnibReaderProvider, rnibWriterProvider)
		t.Run(fmt.Sprintf("%d", tc.msgType), func(t *testing.T) {
			handler, err := provider.GetNotificationHandler(tc.msgType)
			if err != nil {
				t.Errorf("want: handler for message type %d, got: error %s", tc.msgType, err)
			}
			//Note struct is empty, so it will match any other empty struct.
			// https://golang.org/ref/spec#Comparison_operators: Struct values are comparable if all their fields are comparable. Two struct values are equal if their corresponding non-blank fields are equal.
			if /*handler != tc.handler &&*/ strings.Compare(fmt.Sprintf("%T", handler), fmt.Sprintf("%T", tc.handler)) != 0 {
				t.Errorf("want: handler %T for message type %d, got: %T", tc.handler,tc.msgType, handler)
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
		readerMock :=&mocks.RnibReaderMock{}
		rnibReaderProvider := func() reader.RNibReader {
			return readerMock
		}
		writerMock := &mocks.RnibWriterMock{}
		rnibWriterProvider := func() rNibWriter.RNibWriter {
			return writerMock
		}
		provider := NewNotificationHandlerProvider(rnibReaderProvider, rnibWriterProvider)
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
