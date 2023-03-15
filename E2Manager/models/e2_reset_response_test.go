//
// Copyright 2022 Samsung Electronics Co.
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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package models_test

import (
	"e2mgr/models"
	"e2mgr/utils"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ResetResponseXMLPath = "../tests/resources/reset/reset-response.xml"
)

func getResetResponseMessage(t *testing.T, xmlPath string) *models.E2ResetResponseMessage {
	resetResponse := utils.ReadXmlFile(t, xmlPath)
	resetResponseMsg := &models.E2ResetResponseMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(resetResponse), &resetResponseMsg.E2ApPdu)
	assert.Nil(t, err)
	return resetResponseMsg
}

func TestParseResetResponse(t *testing.T) {
	rr := getResetResponseMessage(t, ResetResponseXMLPath)
	assert.NotEqual(t, nil, rr, "xml is not parsed correctly")
	assert.Equal(t, models.ProcedureCode_id_Reset, rr.E2ApPdu.SuccessfulOutcome.ProcedureCode)
	assert.Equal(t, 1, len(rr.E2ApPdu.SuccessfulOutcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs))

	txid := rr.E2ApPdu.SuccessfulOutcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs[0]

	assert.Equal(t, models.ProtocolIE_ID_id_TransactionID, txid.ID)
}
