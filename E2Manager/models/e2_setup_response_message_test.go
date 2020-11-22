//
// Copyright 2020 AT&T Intellectual Property
// Copyright 2020 Nokia
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
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	e2SetupRespGnbSetupRequestXmlPath   = "../tests/resources/setupRequest/setupRequest_gnb.xml"
	e2SetupRespEnGnbSetupRequestXmlPath = "../tests/resources/setupRequest/setupRequest_gnb_without_functions.xml"
)

func getE2SetupRespTestE2SetupRequest(t *testing.T, reqXmlPath string) *models.E2SetupRequestMessage {
	xmlGnb := utils.ReadXmlFile(t, reqXmlPath)
	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(xmlGnb), &setupRequest.E2APPDU)
	assert.Nil(t, err)
	return setupRequest
}

func TestNewE2SetupSuccessResponseMessageSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, "4", respIEs[0].ID)
	assert.Equal(t, plmn, respIEs[0].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[0].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, "9", respIEs[1].ID)
}

func TestNewE2SetupSuccessResponseMessageWithoutRanFunctionsSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespEnGnbSetupRequestXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, "4", respIEs[0].ID)
	assert.Equal(t, plmn, respIEs[0].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[0].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, 1, len(respIEs))
}

func TestNewE2SetupFailureResponseMessageSuccess(t *testing.T) {
	waitTime := models.TimeToWaitEnum.V60s
	cause := models.Cause{Misc: &models.CauseMisc{OmIntervention: &struct{}{}}}

	resp := models.NewE2SetupFailureResponseMessage(waitTime, cause)
	respIEs := resp.E2APPDU.Outcome.(models.UnsuccessfulOutcome).Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs
	assert.Equal(t, "1", respIEs[0].ID)
	assert.Equal(t, cause, respIEs[0].Value.Value.(models.Cause))
}
