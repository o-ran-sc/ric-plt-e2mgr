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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRicServiceUpdateAckMessageSuccess(t *testing.T) {
	item1 := models.RicServiceAckRANFunctionIDItem{RanFunctionID: 100, RanFunctionRevision: 200}
	item2 := models.RicServiceAckRANFunctionIDItem{RanFunctionID: 101, RanFunctionRevision: 201}
	serviceupdateAckFunctionIds := []models.RicServiceAckRANFunctionIDItem{
		item1,
		item2,
	}

	serviceUpdateAck := models.NewServiceUpdateAck(serviceupdateAckFunctionIds)
	ies := serviceUpdateAck.InitiatingMessage.(models.RicServiceUpdateAckSuccessfulOutcome).Value.RICserviceUpdateAcknowledge.ProtocolIEs.RICserviceUpdateAcknowledgeIEs
	assert.Equal(t, "9", ies[0].ID)
	assert.Equal(t, "6", ies[0].Value.RANfunctionsIDList.ProtocolIESingleContainer[0].Id)
	assert.Equal(t, "6", ies[0].Value.RANfunctionsIDList.ProtocolIESingleContainer[1].Id)
	assert.Equal(t, item1, ies[0].Value.RANfunctionsIDList.ProtocolIESingleContainer[0].Value.RANfunctionIDItem)
	assert.Equal(t, item2, ies[0].Value.RANfunctionsIDList.ProtocolIESingleContainer[1].Value.RANfunctionIDItem)
}

func TestRicServiceUpdateAckMessageNoRanFunctionIdItemsSuccess(t *testing.T) {
	successfulOutcome := models.RicServiceUpdateAckSuccessfulOutcome{}
	successfulOutcome.ProcedureCode = "7"

	expectedAck := models.RicServiceUpdateAckE2APPDU{InitiatingMessage: successfulOutcome}
	serviceUpdateAck := models.NewServiceUpdateAck(nil)
	assert.Equal(t, expectedAck, serviceUpdateAck)
}
