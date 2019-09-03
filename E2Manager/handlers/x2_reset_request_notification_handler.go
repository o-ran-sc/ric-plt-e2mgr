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
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type X2ResetRequestNotificationHandler struct {
	rnibReaderProvider func() reader.RNibReader
}

func NewX2ResetRequestNotificationHandler(rnibReaderProvider func() reader.RNibReader) X2ResetRequestNotificationHandler {
	return X2ResetRequestNotificationHandler{
		rnibReaderProvider: rnibReaderProvider,
	}
}

func (src X2ResetRequestNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	logger.Debugf("#X2ResetRequestNotificationHandler.Handle - Ran name: %s", request.RanName)

	nb, rNibErr := src.rnibReaderProvider().GetNodeb(request.RanName)
	if rNibErr != nil {
		logger.Errorf("#X2ResetRequestNotificationHandler.Handle - failed to retrieve nodeB entity. RanName: %s. Error: %s", request.RanName, rNibErr.Error())
		printHandlingSetupResponseElapsedTimeInMs(logger, "#X2ResetRequestNotificationHandler.Handle - Summary: Elapsed time for receiving and handling reset request message from E2 terminator", request.StartTime)

		return
	}
	logger.Debugf("#X2ResetRequestNotificationHandler.Handle - nodeB entity retrieved. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)

	if nb.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		logger.Warnf("#X2ResetRequestNotificationHandler.Handle - nodeB entity in incorrect state. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)
		printHandlingSetupResponseElapsedTimeInMs(logger, "#X2ResetRequestNotificationHandler.Handle - Summary: Elapsed time for receiving and handling reset request message from E2 terminator", request.StartTime)

		return
	}

	if nb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		logger.Errorf("#X2ResetRequestNotificationHandler.Handle - nodeB entity in incorrect state. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)
		printHandlingSetupResponseElapsedTimeInMs(logger, "#X2ResetRequestNotificationHandler.Handle - Summary: Elapsed time for receiving and handling reset request message from E2 terminator", request.StartTime)

		return
	}
	response := models.NotificationResponse{RanName: request.RanName, Payload: e2pdus.PackedX2ResetResponse, MgsType: rmrCgo.RIC_X2_RESET_RESP}
	messageChannel <- &response

	//TODO change name of printHandlingSetupResponseElapsedTimeInMs (remove setup response) and move to utils?
	printHandlingSetupResponseElapsedTimeInMs(logger, "#X2ResetRequestNotificationHandler.Handle - Summary: Elapsed time for receiving and handling reset request message from E2 terminator", request.StartTime)
}