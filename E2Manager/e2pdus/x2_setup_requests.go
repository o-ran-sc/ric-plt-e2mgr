package e2pdus

// #cgo CFLAGS: -I../asn1codec/inc/ -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2setup_request_wrapper.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)

const (
	ShortMacro_eNB_ID = 18
	Macro_eNB_ID      = 20
	LongMacro_eNB_ID  = 21
	Home_eNB_ID       = 28
)

var PackedEndcX2setupRequest []byte
var PackedX2setupRequest []byte
var PackedEndcX2setupRequestAsString string
var PackedX2setupRequestAsString string

func PreparePackedEndcX2SetupRequest(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int,pLMNId []byte, eNB_Id []byte /*18, 20, 21, 28 bits length*/, bitqty uint, ricFlag []byte) ([]byte, string, error) {
	packedBuf := make([]byte, maxAsn1PackedBufferSize)
	errBuf := make([]C.char, maxAsn1CodecMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))
	pduAsString := ""

	if !C.build_pack_endc_x2setup_request(
			(*C.uchar)(unsafe.Pointer(&pLMNId[0])) /*pLMN_Identity*/,
			(*C.uchar)(unsafe.Pointer(&eNB_Id[0])),
			C.uint(bitqty),
			(*C.uchar)(unsafe.Pointer(&ricFlag[0])) /*pLMN_Identity*/,
			&packedBufSize,
			(*C.uchar)(unsafe.Pointer(&packedBuf[0])),
			C.ulong(len(errBuf)),
			&errBuf[0]) {
		return nil, "", errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errBuf[0])))
	}

	pdu:= C.new_pdu(C.size_t(1)) //TODO: change signature
	defer C.delete_pdu(pdu)
	if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])),C.size_t(len(errBuf)), &errBuf[0]){
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		pduAsString = C.GoString(&errBuf[0])
	}

	return packedBuf[:packedBufSize], pduAsString, nil
}

func PreparePackedX2SetupRequest(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int,pLMNId []byte, eNB_Id []byte /*18, 20, 21, 28 bits length*/, bitqty uint, ricFlag []byte) ([]byte, string, error)  {
	packedBuf := make([]byte, maxAsn1PackedBufferSize)
	errBuf := make([]C.char, maxAsn1CodecMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))
	pduAsString := ""

	if !C.build_pack_x2setup_request(
		(*C.uchar)(unsafe.Pointer(&pLMNId[0])) /*pLMN_Identity*/,
		(*C.uchar)(unsafe.Pointer(&eNB_Id[0])),
		C.uint(bitqty),
		(*C.uchar)(unsafe.Pointer(&ricFlag[0])) /*pLMN_Identity*/,
		&packedBufSize,
		(*C.uchar)(unsafe.Pointer(&packedBuf[0])),
		C.ulong(len(errBuf)),
		&errBuf[0]) {
		return nil, "", errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errBuf[0])))
	}

	pdu:= C.new_pdu(C.size_t(1)) //TODO: change signature
	defer C.delete_pdu(pdu)
	if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])),C.size_t(len(errBuf)), &errBuf[0]){
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		pduAsString = C.GoString(&errBuf[0])
	}
	return packedBuf[:packedBufSize], pduAsString, nil
}

