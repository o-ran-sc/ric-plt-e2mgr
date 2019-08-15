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
// #include <x2setup_request_wrapper.h>
import "C"
import (
	"e2mgr/logger"
	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)

type X2PduRefinedResponse struct {
	pduPrint string
}

//func unpackX2apPduUPer(logger *logger.Logger, allocationBufferSize int, packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (*C.E2AP_PDU_t, error) {
//	pdu := C.new_pdu(C.ulong(allocationBufferSize))
//
//	if pdu == nil {
//		return nil, errors.New("allocation failure (pdu)")
//	}
//
//	logger.Debugf("#x2apPdu_asn1_unpacker.unpackX2apPduUPer - Packed pdu(%d):%x", packedBufferSize, packedBuf)
//
//	errBuf := make([]C.char, maxMessageBufferSize)
//	if !C.unpack_pdu_aux(pdu, C.ulong(packedBufferSize), (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.ulong(len(errBuf)), &errBuf[0], C.ATS_UNALIGNED_BASIC_PER) {
//		return nil, errors.New(fmt.Sprintf("unpacking error: %s", C.GoString(&errBuf[0])))
//	}
//
//	if logger.DebugEnabled() {
//		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
//		logger.Debugf("#x2apPdu_asn1_unpacker.unpackX2apPduUPer - PDU: %v  packed size:%d", C.GoString(&errBuf[0]), packedBufferSize)
//	}
//
//	return pdu, nil
//}

func unpackX2apPdu(logger *logger.Logger, allocationBufferSize int, packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (*C.E2AP_PDU_t, error) {
	pdu := C.new_pdu(C.ulong(allocationBufferSize))

	if pdu == nil {
		return nil, errors.New("allocation failure (pdu)")
	}

	logger.Infof("#x2apPdu_asn1_unpacker.unpackX2apPdu - Packed pdu(%d):%x", packedBufferSize, packedBuf)

	errBuf := make([]C.char, maxMessageBufferSize)
	if !C.per_unpack_pdu(pdu, C.ulong(packedBufferSize), (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.ulong(len(errBuf)), &errBuf[0]) {
		return nil, errors.New(fmt.Sprintf("unpacking error: %s", C.GoString(&errBuf[0])))
	}

	if logger.DebugEnabled() {
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		logger.Debugf("#x2apPdu_asn1_unpacker.unpackX2apPdu - PDU: %v  packed size:%d", C.GoString(&errBuf[0]), packedBufferSize)
	}

	return pdu, nil
}

func unpackX2apPduAndRefine(logger *logger.Logger, allocationBufferSize int, packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (*X2PduRefinedResponse, error) {
	pdu, err := unpackX2apPdu(logger, allocationBufferSize, packedBufferSize, packedBuf, maxMessageBufferSize)
	if err != nil {
		return nil, err
	}

	defer C.delete_pdu(pdu)

	var refinedResponse = ""
	if logger.DebugEnabled() {
		buf := make([]C.char, 16*maxMessageBufferSize)
		C.asn1_pdu_printer(pdu, C.size_t(len(buf)), &buf[0])
		refinedResponse  = C.GoString(&buf[0])
	}

	return &X2PduRefinedResponse{pduPrint: refinedResponse}, nil
}



