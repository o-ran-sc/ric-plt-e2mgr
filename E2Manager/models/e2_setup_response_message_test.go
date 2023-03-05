//
// Copyright 2020 AT&T Intellectual Property
// Copyright 2020 Nokia
// Copyright (c) 2022 Samsung Electronics Co., Ltd. All Rights Reserved.
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
	e2SetupRespGnbSetupRequestXmlPath             = "../tests/resources/setupRequest/setupRequest_gnb.xml"
	e2SetupRespEnGnbSetupRequestXmlPath           = "../tests/resources/setupRequest/setupRequest_gnb_without_functions.xml"
	e2SetupRespGnbSetupRequestIntTypeF1XmlPath    = "../tests/resources/setupRequest/setupRequest_gnb_inttype_f1.xml"
	e2SetupRespGnbSetupRequestIntTypeE1XmlPath    = "../tests/resources/setupRequest/setupRequest_gnb_inttype_e1.xml"
	e2SetupRespGnbSetupRequestIntTypeW1XmlPath    = "../tests/resources/setupRequest/setupRequest_gnb_inttype_w1.xml"
	e2SetupRespGnbSetupRequestIntTypeS1XmlPath    = "../tests/resources/setupRequest/setupRequest_gnb_inttype_s1.xml"
	e2SetupRespGnbSetupRequestIntTypeXngnbXmlPath = "../tests/resources/setupRequest/setupRequest_gnb_inttype_xngnb.xml"
	e2SetupRespGnbSetupRequestIntTypeXnenbXmlPath = "../tests/resources/setupRequest/setupRequest_gnb_inttype_xnenb.xml"
	e2SetupRespGnbSetupRequestIntTypeX2gnbXmlPath = "../tests/resources/setupRequest/setupRequest_gnb_inttype_x2gnb.xml"
	e2SetupRespGnbSetupRequestIntTypeX2enbXmlPath = "../tests/resources/setupRequest/setupRequest_gnb_inttype_x2enb.xml"
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
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageWithoutRanFunctionsSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespEnGnbSetupRequestXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, 3, len(respIEs))
}

func TestNewE2SetupFailureResponseMessageSuccess(t *testing.T) {
	waitTime := models.TimeToWaitEnum.V60s
	cause := models.Cause{Misc: &models.CauseMisc{OmIntervention: &struct{}{}}}
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestXmlPath)

	resp := models.NewE2SetupFailureResponseMessage(waitTime, cause, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.UnsuccessfulOutcome).Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs
	assert.Equal(t, models.CauseID, respIEs[1].ID)
	assert.Equal(t, cause, respIEs[1].Value.Value.(models.Cause))
}

func TestNewE2SetupSuccessResponseMessageIntTypeF1Success(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeF1XmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeS1Success(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeS1XmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeE1Success(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeE1XmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeW1Success(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeW1XmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeXngnbSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeXngnbXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeXnenbSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeXnenbXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeX2gnbSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeX2gnbXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}

func TestNewE2SetupSuccessResponseMessageIntTypeX2enbSuccess(t *testing.T) {
	plmn := "23F749"
	ricNearRtId := "10101010110011001110"
	setupRequest := getE2SetupRespTestE2SetupRequest(t, e2SetupRespGnbSetupRequestIntTypeX2enbXmlPath)

	resp := models.NewE2SetupSuccessResponseMessage(plmn, ricNearRtId, setupRequest)
	respIEs := resp.E2APPDU.Outcome.(models.SuccessfulOutcome).Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs
	assert.Equal(t, models.GlobalRicID, respIEs[1].ID)
	assert.Equal(t, plmn, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.PLMNIdentity)
	assert.Equal(t, ricNearRtId, respIEs[1].Value.(models.GlobalRICID).GlobalRICID.RicID)
	assert.Equal(t, models.RanFunctionsAcceptedID, respIEs[2].ID)
}
