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

// #cgo CFLAGS: -I../../asn1codec/inc/ -I../../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../../asn1codec/lib/ -L../../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
import "C"
import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type X2ResetResponseHandler struct {
	rnibReaderProvider func() reader.RNibReader
}

func NewX2ResetResponseHandler(rnibReaderProvider func() reader.RNibReader) X2ResetResponseHandler {
	return X2ResetResponseHandler{
		rnibReaderProvider: rnibReaderProvider,
	}
}

func (src X2ResetResponseHandler) Handle(logger *logger.Logger, request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	logger.Infof("#x2ResetResponseHandler.Handle - received reset response. Payload: %x", request.Payload)

	if nb, rNibErr := src.rnibReaderProvider().GetNodeb(request.RanName); rNibErr != nil {
		logger.Errorf("#x2ResetResponseHandler.Handle - failed to retrieve nb entity. RanName: %s. Error: %s", request.RanName, rNibErr.Error())
	} else {
		logger.Debugf("#x2ResetResponseHandler.Handle - nb entity retrieved. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)
		refinedMessage, err := converters.UnpackX2apPduAndRefine(logger, e2pdus.MaxAsn1CodecAllocationBufferSize , request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)
		if err != nil {
			logger.Errorf("#x2ResetResponseHandler.Handle - failed to unpack reset response message. RanName %s, Payload: %s", request.RanName, request.Payload)
		} else {
			logger.Debugf("#x2ResetResponseHandler.Handle - reset response message payload unpacked. RanName %s, Message: %s", request.RanName, refinedMessage.PduPrint)
		}
	}
}
