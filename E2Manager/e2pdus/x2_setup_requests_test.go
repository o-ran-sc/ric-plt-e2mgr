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

import (
	"bytes"
	"e2mgr/logger"
	"fmt"
	"strings"
	"testing"
)

func TestPreparePackedEndcX2SetupRequest(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "0024003100000100f4002a0000020015000800bbbccc00abcde000fa0017000001f700bbbcccabcde0000000bbbccc000000000001"
	packedEndcX2setupRequest := PackedEndcX2setupRequest

	tmp := fmt.Sprintf("%x", packedEndcX2setupRequest)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedEndcX2setupRequest)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPreparePackedX2SetupRequest(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "0006002a0000020015000800bbbccc00abcde000140017000001f700bbbcccabcde0000000bbbccc000000000001"
	packedX2setupRequest := PackedX2setupRequest

	tmp := fmt.Sprintf("%x", packedX2setupRequest)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedX2setupRequest)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPreparePackedX2SetupRequestFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	_, _, err  = preparePackedX2SetupRequest(1, 4096, pLMNId, eNBId, eNBIdBitqty, ricFlag)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "packing error: #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}

func TestPreparePackedEndcSetupRequestFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	_, _, err  = preparePackedEndcX2SetupRequest(1, 4096, pLMNId, eNBId, eNBIdBitqty, ricFlag)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "packing error: #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}

func TestParseRicId(t *testing.T) {
	var testCases = []struct {
		ricId       string
		pLMNId      []byte
		eNBId       []byte
		eNBIdBitqty uint
		failure     error
	}{
		{
			ricId:       "bbbccc-abcd02/18",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0x2}, /*00000010 -> 10000000*/
			eNBIdBitqty: ShortMacro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcd0e/20",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0xe},
			eNBIdBitqty: Macro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcd07/21",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0x7}, /*00000111 -> 00111000*/
			eNBIdBitqty: LongMacro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcdef08/28",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0xef, 0x8},
			eNBIdBitqty: Home_eNB_ID,
		},
		{
			ricId:   "",
			failure: fmt.Errorf("unable to extract the value of RIC_ID: EOF"),
		},

		{
			ricId:   "bbbccc",
			failure: fmt.Errorf("unable to extract the value of RIC_ID: unexpected EOF"),
		},
		{
			ricId:   "bbbccc-",
			failure: fmt.Errorf("unable to extract the value of RIC_ID: EOF"),
		},
		{
			ricId:   "-bbbccc",
			failure: fmt.Errorf("%s", "unable to extract the value of RIC_ID: no hex data for %x string"),
		},
		{
			ricId:   "/20",
			failure: fmt.Errorf("%s", "unable to extract the value of RIC_ID: no hex data for %x string"),
		},
		{
			ricId:   "bbbcccdd-abcdef08/28", // pLMNId too long
			failure: fmt.Errorf("unable to extract the value of RIC_ID: input does not match format"),
		},
		{
			ricId:   "bbbccc-abcdef0809/28", // eNBId too long
			failure: fmt.Errorf("unable to extract the value of RIC_ID: input does not match format"),
		},

		{
			ricId:   "bbbc-abcdef08/28", // pLMNId too short
			failure: fmt.Errorf("invalid value for RIC_ID, len(pLMNId:[187 188]) != 3"),
		},
		{
			ricId:   "bbbccc-abcd/28", // eNBId too short
			failure: fmt.Errorf("invalid value for RIC_ID, len(eNBId:[171 205]) != 3 or 4"),
		},
		{
			ricId:   "bbbccc-abcdef08/239", // bit quantity too long - no error, will return 23 (which is invalid)
			failure: fmt.Errorf("invalid value for RIC_ID, eNBIdBitqty: 23"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.ricId, func(t *testing.T) {

			err := parseRicID(tc.ricId)
			if err != nil {
				if tc.failure == nil {
					t.Errorf("want: success, got: parse failed. Error: %v\n", err)
				} else {
					if strings.Compare(err.Error(), tc.failure.Error()) != 0 {
						t.Errorf("want: %s, got: %s\n", err, tc.failure)
					}
				}
			} else {
				if bytes.Compare(tc.pLMNId, pLMNId) != 0 {
					t.Errorf("want: pLMNId = %v, got: pLMNId = %v", tc.pLMNId, pLMNId)
				}

				if bytes.Compare(tc.eNBId, eNBId) != 0 {
					t.Errorf("want: eNBId = %v, got: eNBId = %v", tc.eNBId, eNBId)
				}

				if tc.eNBIdBitqty != eNBIdBitqty {
					t.Errorf("want: eNBIdBitqty = %d, got: eNBIdBitqty = %d", tc.eNBIdBitqty, eNBIdBitqty)
				}
			}
		})
	}
}