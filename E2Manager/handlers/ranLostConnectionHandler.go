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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type RanLostConnectionHandler struct{
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewRanLostConnectionHandler(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) RanLostConnectionHandler {
	return RanLostConnectionHandler{
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
	}
}
func (src RanLostConnectionHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	logger.Warnf("#ranLostConnectionHandler.Handle - Received lost connection (transaction id = %s): %s", request.TransactionId, request.Payload)

	var nb *entities.NodebInfo
	var rNibErr common.IRNibError
	if nb, rNibErr = src.rnibReaderProvider().GetNodeb(request.RanName); rNibErr != nil {
		logger.Errorf("#ranLostConnectionHandler.Handle - transactionId %s: rNib reader failed to retrieve nb entity with RanName: %s. Error: %s", request.TransactionId, request.RanName, rNibErr.Error())
	} else {
		logger.Debugf("#ranLostConnectionHandler.Handle - transactionId %s: nb entity has been retrieved. RanName %s, ConnectionStatus %s", request.TransactionId, nb.RanName, nb.ConnectionStatus)
		changeNodebState(logger, nb)
		nbIdentity := &entities.NbIdentity{InventoryName:nb.RanName, GlobalNbId:nb.GlobalNbId}
		if rNibErr = src.rnibWriterProvider().SaveNodeb(nbIdentity, nb); rNibErr != nil {
			logger.Errorf("#ranLostConnectionHandler.Handle - transactionId %s: rNibWriter failed to save nb entity %s. Error: %s", request.TransactionId, nb.RanName, rNibErr.Error())
		} else {
			logger.Infof("#ranLostConnectionHandler.Handle - transactionId %s: saved to rNib", request.TransactionId)
			logger.Debugf("#ranLostConnectionHandler.Handle - transactionId %s: saved to rNib. RanName %s, ConnectionStatus %v", request.TransactionId, nb.RanName, nb.ConnectionStatus)

		}
	}
	e2session, ok := e2Sessions[request.TransactionId]
	printHandlingSetupResponseElapsedTimeInMs(logger, "#ranLostConnectionHandler.Handle - Summary: Elapsed time for receiving and handling sctp error response from E2 terminator", request.StartTime)
	if ok {
		printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#ranLostConnectionHandler.Handle- Summary: Total roundtrip elapsed time for transactionId %s", request.TransactionId), e2session.SessionStart)
		delete(e2Sessions, request.TransactionId) // Avoid pinning memory (help GC)
	}

}

func changeNodebState(logger *logger.Logger, nb *entities.NodebInfo) {
	switch nb.ConnectionStatus{
	case entities.ConnectionStatus_CONNECTED, entities.ConnectionStatus_CONNECTING, entities.ConnectionStatus_CONNECTED_SETUP_FAILED:
		nb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	case entities.ConnectionStatus_DISCONNECTED:
		logger.Infof("#ranLostConnectionHandler.changeNodebState - nb entity with ConnectionStatus %v occurred. RanName: %s", nb.ConnectionStatus, nb.RanName)
	default:
		nb.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	}
}
