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

import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2SetupFailureResponseNotificationHandler struct{}

func (src X2SetupFailureResponseNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	e2session, ok := e2Sessions[request.TransactionId]

	failureResponse, err := unpackX2SetupFailureResponseAndExtract(logger, MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		logger.Errorf("#x2Setup_failure_response_notification_handler.Handle - unpack failed. Error: %v", err)
	}

	printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#x2Setup_failure_response_notification_handler.handle - transactionId %s: Summary: Elapsed time for receiving and handling setup response from E2 terminator", request.TransactionId), request.StartTime)
	if ok {
		if failureResponse != nil {
			nb := &entities.NodebInfo{}
			nbIdentity := &entities.NbIdentity{}

			nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED_SETUP_FAILED
			nb.Ip = e2session.Request.RanIp
			nb.Port = uint32(e2session.Request.RanPort)
			nb.SetupFailure = failureResponse
			nb.FailureType = entities.Failure_X2_SETUP_FAILURE
			nbIdentity.InventoryName = e2session.Request.RanName

			//insert/update database
			if rNibErr := rNibWriter.GetRNibWriter().SaveNodeb(nbIdentity, nb); rNibErr != nil {
				logger.Errorf("#x2Setup_failure_response_notification_handler.Handle - transactionId %s: rNibWriter failed to save failure response data. Error: %s", request.TransactionId, rNibErr.Error())
			} else {
				logger.Infof("#x2Setup_failure_response_notification_handler.Handle - transactionId %s: saved to rNib", request.TransactionId)
				if logger.DebugEnabled() {
					logger.Debugf("x2Setup_failure_response_notification_handler.Handle - transactionId %s: saved to rNib , value:[%s]", request.TransactionId, fmt.Sprintf("%s %s", nb.ConnectionStatus, failureResponse))
				}
			}
		}
		printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#x2Setup_failure_response_notification_handler.handle - transactionId %s: Summary: Total roundtrip elapsed time", request.TransactionId), e2session.SessionStart)
		delete(e2Sessions, request.TransactionId) // Avoid pinning memory (help GC)
	}
}
