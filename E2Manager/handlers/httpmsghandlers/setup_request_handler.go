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
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"os"
	"sync"
	"time"

	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
)

const (
	ENV_RIC_ID                       = "RIC_ID"
	MaxAsn1CodecAllocationBufferSize = 64 * 1024
	MaxAsn1PackedBufferSize          = 4096
	MaxAsn1CodecMessageBufferSize    = 4096
)


/*The Ric Id is the combination of pLMNId and ENBId*/
var pLMNId []byte
var eNBId []byte
var eNBIdBitqty uint
var ricFlag = []byte{0xbb, 0xbc, 0xcc} /*pLMNId [3]bytes*/

type SetupRequestHandler struct {
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewSetupRequestHandler(rnibWriterProvider func() rNibWriter.RNibWriter) *SetupRequestHandler {
	return &SetupRequestHandler{
		rnibWriterProvider: rnibWriterProvider,
	}
}

func (handler SetupRequestHandler) PreHandle(logger *logger.Logger, details *models.RequestDetails) error {
	nodebInfo, nodebIdentity := rnibBuilders.CreateInitialNodeInfo(details, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	rNibErr := handler.rnibWriterProvider().SaveNodeb(nodebIdentity, nodebInfo)
	if rNibErr != nil {
		logger.Errorf("#setup_request_handler.PreHandle - failed to save initial nodeb entity for ran name: %v in RNIB. Error: %s", details.RanName, rNibErr.Error())
	} else {
		logger.Infof("#setup_request_handler.PreHandle - initial nodeb entity for ran name: %v was saved to RNIB ", details.RanName)
	}

	return rNibErr
}

func (SetupRequestHandler) CreateMessage(logger *logger.Logger, requestDetails *models.RequestDetails, messageChannel chan *models.E2RequestMessage, e2sessions sessions.E2Sessions, startTime time.Time, wg sync.WaitGroup) {

	wg.Add(1)

	transactionId := requestDetails.RanName
	e2sessions[transactionId] = sessions.E2SessionDetails{SessionStart: startTime, Request: requestDetails}
	setupRequestMessage := models.NewE2RequestMessage(transactionId, requestDetails.RanIp, requestDetails.RanPort, requestDetails.RanName, e2pdus.PackedX2setupRequest)

	logger.Debugf("#setup_request_handler.CreateMessage - PDU: %s", e2pdus.PackedX2setupRequestAsString)
	logger.Debugf("#setup_request_handler.CreateMessage - setupRequestMessage was created successfully. setup request details(transactionId = [%s]): %+v", transactionId, setupRequestMessage)
	messageChannel <- setupRequestMessage

	wg.Done()
}

//Expected value in RIC_ID = pLMN_Identity-eNB_ID/<eNB_ID size in bits>
//<6 hex digits>-<6 or 8 hex digits>/<18|20|21|28>
//Each byte is represented by two hex digits, the value in the lowest byte of the eNB_ID must be assigned to the lowest bits
//For example, to get the value of ffffeab/28  the last byte must be 0x0b, not 0xb0 (-ffffea0b/28).
func parseRicID(ricId string) error {
	if _, err := fmt.Sscanf(ricId, "%6x-%8x/%2d", &pLMNId, &eNBId, &eNBIdBitqty); err != nil {
		return fmt.Errorf("unable to extract the value of %s: %s", ENV_RIC_ID, err)
	}

	if len(pLMNId) < 3 {
		return fmt.Errorf("invalid value for %s, len(pLMNId:%v) != 3", ENV_RIC_ID, pLMNId)
	}

	if len(eNBId) < 3 {
		return fmt.Errorf("invalid value for %s, len(eNBId:%v) != 3 or 4", ENV_RIC_ID, eNBId)
	}

	if eNBIdBitqty != e2pdus.ShortMacro_eNB_ID && eNBIdBitqty != e2pdus.Macro_eNB_ID && eNBIdBitqty != e2pdus.LongMacro_eNB_ID && eNBIdBitqty != e2pdus.Home_eNB_ID {
		return fmt.Errorf("invalid value for %s, eNBIdBitqty: %d", ENV_RIC_ID, eNBIdBitqty)
	}

	return nil
}

//TODO: remove Get
func (SetupRequestHandler) GetMessageType() int {
	return rmrCgo.RIC_X2_SETUP_REQ
}

func init() {
	var err error
	ricId := os.Getenv(ENV_RIC_ID)
	//ricId="bbbccc-ffff0e/20"
	//ricId="bbbccc-abcd0e/20"
	if err = parseRicID(ricId); err != nil {
		panic(err)
	}

	e2pdus.PackedEndcX2setupRequest,e2pdus.PackedEndcX2setupRequestAsString, err = e2pdus.PreparePackedEndcX2SetupRequest(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize,pLMNId, eNBId, eNBIdBitqty, ricFlag )
	if err != nil{
		panic(err)
	}
	e2pdus.PackedX2setupRequest,e2pdus.PackedX2setupRequestAsString, err = e2pdus.PreparePackedX2SetupRequest(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize,pLMNId, eNBId, eNBIdBitqty, ricFlag )
	if err != nil{
		panic(err)
	}
}
