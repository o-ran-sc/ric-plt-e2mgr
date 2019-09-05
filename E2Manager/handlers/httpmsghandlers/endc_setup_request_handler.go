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

package httpmsghandlers

import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"e2mgr/rnibBuilders"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"sync"
	"time"

	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
)

type EndcSetupRequestHandler struct {
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewEndcSetupRequestHandler(rnibWriterProvider func() rNibWriter.RNibWriter) *EndcSetupRequestHandler {
	return &EndcSetupRequestHandler{
		rnibWriterProvider: rnibWriterProvider,
	}
}

func (handler EndcSetupRequestHandler) PreHandle(logger *logger.Logger, details *models.RequestDetails) error {
	nodebInfo, nodebIdentity := rnibBuilders.CreateInitialNodeInfo(details,entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)

	rNibErr := handler.rnibWriterProvider().SaveNodeb(nodebIdentity, nodebInfo)
	if rNibErr != nil {
		logger.Errorf("#endc_setup_request_handler.PreHandle - failed to initial nodeb entity for ran name: %v in RNIB. Error: %s", details.RanName, rNibErr.Error())
	} else {
		logger.Infof("#endc_setup_request_handler.PreHandle - initial nodeb entity for ran name: %v was saved to RNIB ", details.RanName)
	}

	return rNibErr
}

func (EndcSetupRequestHandler) CreateMessage(logger *logger.Logger, requestDetails *models.RequestDetails, messageChannel chan *models.E2RequestMessage, e2sessions sessions.E2Sessions, startTime time.Time, wg sync.WaitGroup) {

	wg.Add(1)

	transactionId := requestDetails.RanName
	e2sessions[transactionId] = sessions.E2SessionDetails{SessionStart: startTime, Request: requestDetails}
	setupRequestMessage := models.NewE2RequestMessage(transactionId, requestDetails.RanIp, requestDetails.RanPort, requestDetails.RanName, e2pdus.PackedEndcX2setupRequest)

	logger.Debugf("#endc_setup_request_handler.CreateMessage - PDU: %s", e2pdus.PackedEndcX2setupRequestAsString)
	logger.Debugf("#endc_setup_request_handler.CreateMessage - setupRequestMessage was created successfuly. setup request details(transactionId = [%s]): %+v", transactionId, setupRequestMessage)
	messageChannel <- setupRequestMessage

	wg.Done()
}

func (EndcSetupRequestHandler) GetMessageType() int {
	return rmrCgo.RIC_ENDC_X2_SETUP_REQ
}
