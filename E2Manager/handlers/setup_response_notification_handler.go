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

package handlers

import "C"
import (
	"e2mgr/models"
	"fmt"
	"time"

	"e2mgr/logger"
	"e2mgr/sessions"
)

type SetupResponseNotificationHandler struct{}

func (src SetupResponseNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	refinedResponse, err := unpackX2apPduAndRefine(logger, MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		logger.Errorf("#setup_response_notification_handler.Handle - unpack failed. Error: %v", err)
	}

	e2session, ok := e2Sessions[request.TransactionId]
	printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#setupResponseNotificationHandler.handle - transactionId %s: Summary: Elapsed time for receiving and handling setup response from E2 terminator", request.TransactionId), request.StartTime)
	if ok {
		printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#setupResponseNotificationHandler.handle - transactionId %s: Summary: Total roundtrip elapsed time", request.TransactionId), e2session.SessionStart)
		delete(e2Sessions, request.TransactionId) // Avoid pinning memory (help GC)
	}
	logger.Debugf("#setupResponseNotificationHandler.handle - transactionId %s: PDU: %v", request.TransactionId, refinedResponse.pduPrint)
}

func printHandlingSetupResponseElapsedTimeInMs(logger *logger.Logger, msg string, startTime time.Time) {
	logger.Infof("%s: %f ms", msg, float64(time.Since(startTime))/float64(time.Millisecond))
}
