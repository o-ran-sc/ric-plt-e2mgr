package e2pdus

import (
	"e2mgr/logger"
	"fmt"
	"strings"
	"testing"
)

func TestKnownCausesToX2ResetPDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	var testCases = []struct {
		cause string
		packedPdu        string
	}{
		{
			cause:     OmInterventionCause,
			packedPdu: "000700080000010005400164",
		},
		{
			cause:     "PROTOCOL:transfer-syntax-error",
			packedPdu: "000700080000010005400140",
		},
		{
			cause:     "transport:transport-RESOURCE-unavailable",
			packedPdu: "000700080000010005400120",
		},

		{
			cause:     "radioNetwork:invalid-MME-groupid",
			packedPdu: "00070009000001000540020680",
		},

	}

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			payload, ok := KnownCausesToX2ResetPDU(tc.cause)
			if !ok {
				t.Errorf("want: success, got: not found.\n")
			} else {
				tmp := fmt.Sprintf("%x", payload)
				if len(tmp) != len(tc.packedPdu) {
					t.Errorf("want packed len:%d, got: %d\n", len(tc.packedPdu)/2, len(payload)/2)
				}

				if strings.Compare(tmp, tc.packedPdu) != 0 {
					t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", tc.packedPdu, tmp)
				}
			}
		})
	}
}


func TestKnownCausesToX2ResetPDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	_, ok  := KnownCausesToX2ResetPDU("xxxx")
	if ok {
		t.Errorf("want: not found, got: success.\n")
	}
}


func TestPrepareX2ResetPDUsFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareX2ResetPDUs(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#reset_request_handler.Handle - failed to build and pack the reset message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big:"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}