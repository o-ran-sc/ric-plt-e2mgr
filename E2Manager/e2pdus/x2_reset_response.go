/*******************************************************************************
 *
 *   Copyright (c) 2019 AT&T Intellectual Property.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 *
 *******************************************************************************/
 package e2pdus

// #cgo CFLAGS: -I../asn1codec/inc/  -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <x2reset_response_wrapper.h>
import "C"
import (
	"fmt"
	"unsafe"
)

var PackedX2ResetResponse []byte

func prepareX2ResetResponsePDU(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {

	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)
	var payloadSize = C.ulong(maxAsn1PackedBufferSize)

	if status := C.build_pack_x2reset_response(&payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
		return fmt.Errorf("#x2_reset_response.prepareX2ResetResponsePDU - failed to build and pack the reset response message %s ", C.GoString(&errorBuffer[0]))

	}
	PackedX2ResetResponse = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))

	return nil
}

func init() {
	if err := prepareX2ResetResponsePDU(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
}
