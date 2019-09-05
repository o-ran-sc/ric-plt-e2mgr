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
	"bytes"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewSetupRequestHandler(t *testing.T) {

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	h := NewSetupRequestHandler(rnibWriterProvider)
	assert.NotNil(t, h)
}

func TestCreateMessageSuccess(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#setup_request_handler_test.TestCreateMessageSuccess - failed to initialize logger, error: %s", err)
	}
	messageChannel := make(chan *models.E2RequestMessage)
	assert.NotPanics(t, func() { createMsg(log, messageChannel) })
	assert.NotEmpty(t, <-messageChannel)
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
			eNBIdBitqty: e2pdus.ShortMacro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcd0e/20",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0xe},
			eNBIdBitqty: e2pdus.Macro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcd07/21",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0x7}, /*00000111 -> 00111000*/
			eNBIdBitqty: e2pdus.LongMacro_eNB_ID,
		},
		{
			ricId:       "bbbccc-abcdef08/28",
			pLMNId:      []byte{0xbb, 0xbc, 0xcc},
			eNBId:       []byte{0xab, 0xcd, 0xef, 0x8},
			eNBIdBitqty: e2pdus.Home_eNB_ID,
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
func createMsg(log *logger.Logger, messageChannel chan *models.E2RequestMessage) {
	h := SetupRequestHandler{}
	E2Sessions := make(sessions.E2Sessions)
	var wg sync.WaitGroup
	var rd models.RequestDetails
	go h.CreateMessage(log, &rd, messageChannel, E2Sessions, time.Now(), wg)
	wg.Wait()
}
