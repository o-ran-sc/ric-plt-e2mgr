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

// #cgo CFLAGS: -I../asn1codec/inc/  -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2setup_request_wrapper.h>
import "C"
import (
	"e2mgr/logger"
	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)

func packX2apSetupRequest(logger *logger.Logger, allocationBufferSize int, maxPackedBufferSize int, maxMessageBufferSize int, pLMNId []byte, eNB_Id []byte /*18, 20, 21, 28 bits length*/, bitqty uint) ([]byte, error) {
	packedBuf := make([]byte, maxPackedBufferSize)
	errBuf := make([]C.char, maxMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))

	if !C.build_pack_x2setup_request((*C.uchar)(unsafe.Pointer(&pLMNId[0])) /*pLMN_Identity*/,
		(*C.uchar)(unsafe.Pointer(&eNB_Id[0])), C.uint(bitqty),(*C.uchar)(unsafe.Pointer(&ricFlag[0])) /*pLMN_Identity*/,
		&packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.ulong(len(errBuf)), &errBuf[0]) {
		return nil, errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errBuf[0])))
	}

	if logger.DebugEnabled() {
		pdu:= C.new_pdu(C.size_t(allocationBufferSize))
		defer C.delete_pdu(pdu)
		if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])),C.size_t(len(errBuf)), &errBuf[0]){
			C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
			logger.Debugf("#x2apSetupRequest_asn1_packer.packX2apSetupRequest - PDU:%s\n\npacked (%d):%x", C.GoString(&errBuf[0]), packedBufSize, packedBuf[:packedBufSize])
		}
	}
	return packedBuf[:packedBufSize], nil

}

