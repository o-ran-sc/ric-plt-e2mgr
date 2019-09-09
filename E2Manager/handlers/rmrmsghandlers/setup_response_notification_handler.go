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
package rmrmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type SetupResponseNotificationHandler struct {
	rnibReaderProvider   func() reader.RNibReader
	rnibWriterProvider   func() rNibWriter.RNibWriter
	setupResponseManager managers.ISetupResponseManager
	notificationType     string
}

func NewSetupResponseNotificationHandler(rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter, setupResponseManager managers.ISetupResponseManager, notificationType string) SetupResponseNotificationHandler {
	return SetupResponseNotificationHandler{
		rnibReaderProvider:   rnibReaderProvider,
		rnibWriterProvider:   rnibWriterProvider,
		setupResponseManager: setupResponseManager,
		notificationType:     notificationType,
	}
}

func (h SetupResponseNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions, request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {
	logger.Infof("#SetupResponseNotificationHandler - RAN name: %s - Received %s notification", request.RanName, h.notificationType)

	inventoryName := request.RanName

	nodebInfo, rnibErr := h.rnibReaderProvider().GetNodeb(inventoryName)

	if rnibErr != nil {
		logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Error fetching RAN from rNib: %v", request.RanName, rnibErr)
		return
	}

	if !isConnectionStatusValid(nodebInfo.ConnectionStatus) {
		logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Invalid RAN connection status: %s", request.RanName, nodebInfo.ConnectionStatus)
		return
	}

	nodebInfo.ConnectionAttempts = 0
	nbIdentity := &entities.NbIdentity{InventoryName: inventoryName}
	err := h.setupResponseManager.SetNodeb(logger, nbIdentity, nodebInfo, request.Payload)

	if err != nil {
		return
	}

	rnibErr = h.rnibWriterProvider().SaveNodeb(nbIdentity, nodebInfo)

	if rnibErr != nil {
		logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Error saving RAN to rNib: %v", request.RanName, rnibErr)
		return
	}

	logger.Infof("#SetupResponseNotificationHandler - RAN name: %s - Successfully saved RAN to rNib", request.RanName)
}

func isConnectionStatusValid(connectionStatus entities.ConnectionStatus) bool {
	return connectionStatus == entities.ConnectionStatus_CONNECTING || connectionStatus == entities.ConnectionStatus_CONNECTED
}
