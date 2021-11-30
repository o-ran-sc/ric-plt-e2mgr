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

func getTestRicServiceUpdate(t *testing.T, xmlPath string) *models.RICServiceUpdateMessage {
	xmlServiceUpdate := utils.ReadXmlFile(t, xmlPath)
	ricServiceUpdate := &models.RICServiceUpdateMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(xmlServiceUpdate), &ricServiceUpdate.E2APPDU)
	assert.Nil(t, err)
	return ricServiceUpdate
}

func TestRicServiceUpdateMessageSuccess(t *testing.T) {
	serviceUpdate := getTestRicServiceUpdate(t, "../tests/resources/serviceUpdate/RicServiceUpdate_AddedFunction.xml")

	ranFunctions := serviceUpdate.E2APPDU.ExtractRanFunctionsList()
	assert.Equal(t, uint32(20), ranFunctions[0].RanFunctionId)
	assert.Equal(t, uint32(2), ranFunctions[0].RanFunctionRevision)
}

func TestRicServiceUpdateMessageNoRanFunctions(t *testing.T) {
	serviceUpdate := getTestRicServiceUpdate(t, "../tests/resources/serviceUpdate/RicServiceUpdate_Empty.xml")
	assert.Nil(t, serviceUpdate.E2APPDU.ExtractRanFunctionsList())
}

func TestRicServiceUpdateMessageWithOID(t *testing.T) {
	serviceUpdate := getTestRicServiceUpdate(t, "../tests/resources/serviceUpdate/RicServiceUpdate_AddedFunction_With_OID.xml")

	ranFunctions := serviceUpdate.E2APPDU.ExtractRanFunctionsList()

	assert.Equal(t, uint32(20), ranFunctions[0].RanFunctionId)
	assert.Equal(t, uint32(2), ranFunctions[0].RanFunctionRevision)
	assert.Equal(t, "OID20", ranFunctions[0].RanFunctionOid)
}
