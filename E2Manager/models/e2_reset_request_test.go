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
	ResetRequstXMLPath = "../tests/resources/reset/reset-request.xml"
)

func getResetRequestMessage(t *testing.T, reqXmlPath string) *models.E2ResetRequestMessage {
	resetRequest := utils.ReadXmlFile(t, reqXmlPath)
	resetRequestMsg := &models.E2ResetRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(resetRequest), &resetRequestMsg.E2APPDU)
	assert.Nil(t, err)
	return resetRequestMsg
}

func TestParseResetRequest(t *testing.T) {
	rr := getResetRequestMessage(t, ResetRequstXMLPath)
	assert.NotEqual(t, nil, rr, "xml is not parsed correctly")
	assert.Equal(t, models.ProcedureCode_id_Reset, rr.E2APPDU.InitiatingMessage.ProcedureCode)
	assert.Equal(t, 2, len(rr.E2APPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs))

	txid := rr.E2APPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs[0]
	cause := rr.E2APPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs[1]

	assert.Equal(t, models.ProtocolIE_ID_id_TransactionID, txid.ID)
	assert.Equal(t, models.ProtocolIE_ID_id_Cause, cause.ID)

	assert.Equal(t, false, cause.Value.Cause.E2Node.E2nodeComponentUnknown == nil)
	assert.Equal(t, true, cause.Value.Cause.Misc.ControlProcessingOverload == nil)
}
