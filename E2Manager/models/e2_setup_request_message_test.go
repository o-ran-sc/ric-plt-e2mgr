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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	e2SetupReqGnbSetupRequestXmlPath                      = "../tests/resources/setupRequest/setupRequest_gnb.xml"
	e2SetupReqEnGnbSetupRequestXmlPath                    = "../tests/resources/setupRequest/setupRequest_en-gNB.xml"
	e2SetupReqEnbSetupRequestXmlPath                      = "../tests/resources/setupRequest/setupRequest_enb.xml"
	e2SetupReqNgEnbSetupRequestXmlPath                    = "../tests/resources/setupRequest/setupRequest_ng-eNB.xml"
	e2SetupReqGnbSetupRequestWithOIDXmlPath               = "../tests/resources/setupRequest/setupRequest_with_oid_gnb.xml"
	e2SetupReqE1InterfaceGnbSetupRequestWithOIDXmlPath    = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_e1.xml"
	e2SetupReqF1InterfaceGnbSetupRequestWithOIDXmlPath    = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_f1.xml"
	e2SetupReqW1InterfaceGnbSetupRequestWithOIDXmlPath    = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_w1.xml"
	e2SetupReqS1InterfaceGnbSetupRequestWithOIDXmlPath    = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_s1.xml"
	e2SetupReqXngnbInterfaceGnbSetupRequestWithOIDXmlPath = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_xngnb.xml"
	e2SetupReqXnenbInterfaceGnbSetupRequestWithOIDXmlPath = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_xnenb.xml"
	e2SetupReqX2gnbInterfaceGnbSetupRequestWithOIDXmlPath = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_x2gnb.xml"
	e2SetupReqX2enbInterfaceGnbSetupRequestWithOIDXmlPath = "../tests/resources/setupRequest/setupRequest_with_oid_gnb_inttype_x2enb.xml"
)

func getTestE2SetupRequest(t *testing.T, reqXmlPath string) *models.E2SetupRequestMessage {
	xmlGnb := utils.ReadXmlFile(t, reqXmlPath)
	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(xmlGnb), &setupRequest.E2APPDU)
	assert.Nil(t, err)
	return setupRequest
}

func TestExtractRanFunctionsListFromGnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqGnbSetupRequestXmlPath)

	ranFuncList := setupRequest.ExtractRanFunctionsList()
	assert.Equal(t, uint32(1), ranFuncList[0].RanFunctionId)
	assert.Equal(t, uint32(2), ranFuncList[1].RanFunctionId)
	assert.Equal(t, uint32(3), ranFuncList[2].RanFunctionId)
	assert.Equal(t, uint32(1), ranFuncList[0].RanFunctionRevision)
	assert.Equal(t, uint32(1), ranFuncList[1].RanFunctionRevision)
	assert.Equal(t, uint32(1), ranFuncList[2].RanFunctionRevision)
}

func TestExtractRanFunctionsListFromGnbRequestwithOidSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqGnbSetupRequestWithOIDXmlPath)

	ranFuncList := setupRequest.ExtractRanFunctionsList()

	assert.Equal(t, uint32(1), ranFuncList[0].RanFunctionId)
	assert.Equal(t, uint32(2), ranFuncList[1].RanFunctionId)
	assert.Equal(t, uint32(3), ranFuncList[2].RanFunctionId)

	assert.Equal(t, uint32(1), ranFuncList[0].RanFunctionRevision)
	assert.Equal(t, uint32(1), ranFuncList[1].RanFunctionRevision)
	assert.Equal(t, uint32(1), ranFuncList[2].RanFunctionRevision)

	assert.Equal(t, "OID123", ranFuncList[0].RanFunctionOid)
	assert.Equal(t, "OID124", ranFuncList[1].RanFunctionOid)
	assert.Equal(t, "OID125", ranFuncList[2].RanFunctionOid)
}

func TestExtractE2nodeConfigSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 2, len(e2nodeConfigs))

	assert.Equal(t, "nginterf1", e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeNG().GetAmfName())
	assert.Equal(t, "nginterf2", e2nodeConfigs[1].GetE2NodeComponentInterfaceTypeNG().GetAmfName())

}

func TestGetPlmnIdFromGnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqGnbSetupRequestXmlPath)

	plmnID := setupRequest.GetPlmnId()
	assert.Equal(t, "02F829", plmnID)
}

func TestGetPlmnIdFromEnGnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqEnGnbSetupRequestXmlPath)

	plmnID := setupRequest.GetPlmnId()
	assert.Equal(t, "131014", plmnID)
}

func TestGetPlmnIdFromEnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqEnbSetupRequestXmlPath)

	plmnID := setupRequest.GetPlmnId()
	assert.Equal(t, "6359AB", plmnID)
}

func TestGetPlmnIdFromNgEnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqNgEnbSetupRequestXmlPath)

	plmnID := setupRequest.GetPlmnId()
	assert.Equal(t, "131014", plmnID)
}

func TestGetNbIdFromGnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqGnbSetupRequestXmlPath)

	nbID := setupRequest.GetNbId()
	assert.Equal(t, "001100000011000000110000", nbID)
}

func TestGetNbIdFromEnGnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqEnGnbSetupRequestXmlPath)

	nbID := setupRequest.GetNbId()
	assert.Equal(t, "11000101110001101100011111111000", nbID)
}

func TestGetNbIdFromEnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqEnbSetupRequestXmlPath)

	nbID := setupRequest.GetNbId()
	assert.Equal(t, "101010101010101010", nbID)
}

func TestGetNbIdFromNgEnbRequestSuccess(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqNgEnbSetupRequestXmlPath)

	nbID := setupRequest.GetNbId()
	assert.Equal(t, "101010101010101010", nbID)
}

func TestExtractE2nodeIntTypeE1ConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqE1InterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 2, len(e2nodeConfigs))

	assert.NotEqual(t, 1, int(e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeE1().GetGNBCuCpId()))
	assert.NotEqual(t, 2, int(e2nodeConfigs[1].GetE2NodeComponentInterfaceTypeE1().GetGNBCuCpId()))

}

func TestExtractE2nodeIntTypeF1ConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqF1InterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 2, len(e2nodeConfigs))

	assert.NotEqual(t, 1, int(e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeF1().GetGNBDuId()))
	assert.NotEqual(t, 2, int(e2nodeConfigs[1].GetE2NodeComponentInterfaceTypeF1().GetGNBDuId()))

}

func TestExtractE2nodeIntTypeW1ConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqW1InterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 2, len(e2nodeConfigs))

	assert.NotEqual(t, 1, int(e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeW1().GetNgenbDuId()))
	assert.NotEqual(t, 2, int(e2nodeConfigs[1].GetE2NodeComponentInterfaceTypeW1().GetNgenbDuId()))

}

func TestExtractE2nodeIntTypeS1ConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqS1InterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 2, len(e2nodeConfigs))

	assert.NotEqual(t, "e1interf1", e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeS1().GetMmeName())
	assert.NotEqual(t, "e1interf2", e2nodeConfigs[1].GetE2NodeComponentInterfaceTypeS1().GetMmeName())

}

func TestExtractE2nodeIntTypeXngnbConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqXngnbInterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 1, len(e2nodeConfigs))

	assert.NotEqual(t, "101010101010101010", (e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeXn().GetGlobalGnbId().GetGnbId()))
}

func TestExtractE2nodeIntTypeXnenbConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqXnenbInterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 1, len(e2nodeConfigs))

	assert.NotEqual(t, "101010101010101010", (e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeXn().GetGlobalNgenbId().GetEnbId()))
}

func TestExtractE2nodeIntTypeX2gnbConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqX2gnbInterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 1, len(e2nodeConfigs))

	assert.NotEqual(t, "101010101010101010", (e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeX2().GetGlobalEngnbId().GetGnbId()))
}

func TestExtractE2nodeIntTypeX2enbConfigFail(t *testing.T) {
	setupRequest := getTestE2SetupRequest(t, e2SetupReqX2enbInterfaceGnbSetupRequestWithOIDXmlPath)
	e2nodeConfigs := setupRequest.ExtractE2NodeConfigList()

	assert.Equal(t, 1, len(e2nodeConfigs))

	assert.NotEqual(t, "101010101010101010", (e2nodeConfigs[0].GetE2NodeComponentInterfaceTypeX2().GetGlobalEnbId().GetEnbId()))
}
