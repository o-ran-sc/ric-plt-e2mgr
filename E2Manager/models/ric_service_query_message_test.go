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

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
)

func getTestRicServiceQueryRanFunctions(t *testing.T) []*entities.RanFunction {
	gnbSetupRequestXmlPath := "../tests/resources/setupRequest/setupRequest_gnb.xml"
	xmlgnb := utils.ReadXmlFile(t, gnbSetupRequestXmlPath)

	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(xmlgnb), &setupRequest.E2APPDU)
	if err != nil {
		t.Fatal(err)
	}

	ranFunctionList := setupRequest.ExtractRanFunctionsList()
	return ranFunctionList
}

func TestRicServiceQueryMessageSuccess(t *testing.T) {
	ranFunctionList := getTestRicServiceQueryRanFunctions(t)

	serviceQuery := models.NewRicServiceQueryMessage(ranFunctionList)
	initMsg := serviceQuery.E2APPDU.InitiatingMessage
	assert.Equal(t, models.ProcedureCode_id_RICserviceQuery, initMsg.ProcedureCode)
	assert.Equal(t, models.ProtocolIE_ID_id_TransactionID, initMsg.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[0].Id)
	assert.Equal(t, models.ProtocolIE_ID_id_RANfunctionsAccepted, initMsg.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[1].Id)
	assert.Equal(t, 3, len(initMsg.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[1].Value.(models.RICServiceQueryRANFunctionIdList).RANFunctionIdList.ProtocolIESingleContainer))
}

func TestTransactionIdServiceQuery(t *testing.T) {
	ranFunctionList := getTestRicServiceQueryRanFunctions(t)
	serviceQuery := models.NewRicServiceQueryMessage(ranFunctionList)
	txIE := serviceQuery.E2APPDU.InitiatingMessage.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[0].Value.(models.RICServiceQueryTransactionID)
	assert.NotEmpty(t, txIE.TransactionID, "transaction ID should not be empty")
}
