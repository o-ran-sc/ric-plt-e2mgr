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

// #cgo CFLAGS: -I../asn1codec/inc/ -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
import "C"
import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/sessions"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type X2ResetResponseHandler struct{
	rnibReaderProvider func() reader.RNibReader
}

func NewX2ResetResponseHandler(rnibReaderProvider func() reader.RNibReader) X2ResetResponseHandler {
	return X2ResetResponseHandler{
		rnibReaderProvider: rnibReaderProvider,
	}
}

func (src X2ResetResponseHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {
	request.TransactionId = request.RanName
		logger.Debugf("#x2ResetResponseHandler.Handle - transactionId %s: received reset response. Payload: %s", request.TransactionId, request.Payload)

		if nb, rNibErr := src.rnibReaderProvider().GetNodeb(request.RanName); rNibErr != nil {
			logger.Errorf("#x2ResetResponseHandler.Handle - transactionId %s: failed to retrieve nb entity. RanName: %s. Error: %s", request.TransactionId, request.RanName, rNibErr.Error())
		} else {
			logger.Debugf("#x2ResetResponseHandler.Handle - transactionId %s: nb entity retrieved. RanName %s, ConnectionStatus %s", request.TransactionId, nb.RanName, nb.ConnectionStatus)
			//TODO: only returned in debug mode
			refinedMessage, err := unpackX2apPduAndRefine(logger, MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize /*message buffer*/)
			if err != nil {
				logger.Errorf("#x2ResetResponseHandler.Handle - transactionId %s: failed to unpack reset response message. RanName %s, Payload: %s", request.TransactionId , request.RanName, request.Payload)
			} else {
				logger.Infof("#x2ResetResponseHandler.Handle - transactionId %s: reset response message payload unpacked. RanName %s, Message: %s", request.TransactionId , request.RanName, refinedMessage.pduPrint)
			}
		}
		e2session, ok := e2Sessions[request.TransactionId]
		if ok {
			printHandlingSetupResponseElapsedTimeInMs(logger, fmt.Sprintf("#x2ResetResponseHandler.Handle- Summary: Total roundtrip elapsed time for transactionId %s", request.TransactionId), e2session.SessionStart)
			delete(e2Sessions, request.TransactionId)
		}

}