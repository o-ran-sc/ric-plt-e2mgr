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
// #include <configuration_update_wrapper.h>
import "C"
import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"unsafe"
)

type X2EnbConfigurationUpdateHandler struct{}

func (src X2EnbConfigurationUpdateHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions, request *models.NotificationRequest,
	messageChannel chan<- *models.NotificationResponse) {

	var payloadSize C.ulong
	payloadSize = MaxAsn1PackedBufferSize
	packedBuffer := [MaxAsn1PackedBufferSize]C.uchar{}
	errorBuffer := [MaxAsn1PackedBufferSize]C.char{}
	refinedMessage, err := unpackX2apPduAndRefine(logger, MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		status :=  C.build_pack_x2enb_configuration_update_failure( &payloadSize, &packedBuffer[0], MaxAsn1PackedBufferSize, &errorBuffer[0])
		if status{
			payload := (*[1<<30]byte)(unsafe.Pointer(&packedBuffer))[:payloadSize:payloadSize]
			logger.Debugf("#x2enb_configuration_update_handler.Handle - Enb configuration update negative ack message payload: (%d) %02x", len(payload), payload)
			response := models.NotificationResponse{RanName:request.RanName, Payload:payload, MgsType: rmrCgo.RIC_ENB_CONFIGURATION_UPDATE_FAILURE}
			messageChannel<- &response
		} else {
			logger.Errorf("#x2enb_configuration_update_handler.Handle - failed to build and pack Enb configuration update unsuccessful outcome message. Error: %v", errorBuffer)
		}
		logger.Errorf("#x2enb_configuration_update_handler.Handle - unpack failed. Error: %v", err)
	} else {
		logger.Infof("#x2enb_configuration_update_handler.Handle - Enb configuration update initiating message received")
		logger.Debugf("#x2enb_configuration_update_handler.Handle - Enb configuration update initiating message payload: %s", refinedMessage.pduPrint)

		status := C.build_pack_x2enb_configuration_update_ack( &payloadSize, &packedBuffer[0], MaxAsn1PackedBufferSize, &errorBuffer[0])
		if status {
			payload := (*[1 << 30]byte)(unsafe.Pointer(&packedBuffer))[:payloadSize:payloadSize]
			logger.Debugf("#x2enb_configuration_update_handler.Handle - Enb configuration update positive ack message payload: (%d) %02x", len(payload), payload)
			response := models.NotificationResponse{RanName: request.RanName, Payload: payload, MgsType: rmrCgo.RIC_ENB_CONFIGURATION_UPDATE_ACK}
			messageChannel<- &response
		} else {
			logger.Errorf("#x2enb_configuration_update_handler.Handle - failed to build and pack enb configuration update successful outcome message. Error: %v", errorBuffer)
		}
	}
	printHandlingSetupResponseElapsedTimeInMs(logger, "#x2enb_configuration_update_handler.Handle - Summary: Elapsed time for receiving and handling enb configuration update initiating message from E2 terminator", request.StartTime)
}

