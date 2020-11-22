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
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAddEnb = map[string]interface{}{
	"ranName": "test",
	"globalNbId": map[string]interface{}{
		"plmnId": "whatever",
		"nbId":   "whatever2",
	},
	"enb": map[string]interface{}{
		"enbType": 1,
		"servedCells": []interface{}{
			map[string]interface{}{
				"cellId": "whatever",
				"choiceEutraMode": map[string]interface{}{
					"fdd": map[string]interface{}{},
				},
				"eutraMode": 1,
				"pci":       1,
				"tac":       "whatever3",
				"broadcastPlmns": []interface{}{
					"whatever",
				},
			},
		},
	},
}

func TestAddEnbRequestSuccess(t *testing.T) {
	addEnbRequest := models.AddEnbRequest{}
	buf, err := json.Marshal(testAddEnb)
	assert.Nil(t, err)

	err = addEnbRequest.UnmarshalJSON(buf)
	assert.Nil(t, err)
	assert.Equal(t, "test", addEnbRequest.RanName)
	assert.Equal(t, "whatever", addEnbRequest.GlobalNbId.PlmnId)
	assert.Equal(t, "whatever2", addEnbRequest.GlobalNbId.NbId)
	assert.Equal(t, entities.EnbType_MACRO_ENB, addEnbRequest.Enb.EnbType)
	assert.Equal(t, uint32(1), addEnbRequest.Enb.ServedCells[0].Pci)
	assert.Equal(t, "whatever", addEnbRequest.Enb.ServedCells[0].CellId)
	assert.Equal(t, "whatever3", addEnbRequest.Enb.ServedCells[0].Tac)
	assert.Equal(t, "whatever", addEnbRequest.Enb.ServedCells[0].BroadcastPlmns[0])
	assert.Equal(t, entities.Eutra_FDD, addEnbRequest.Enb.ServedCells[0].EutraMode)
}

func TestAddEnbRequestJsonUnmarshalError(t *testing.T) {
	addEnbRequest := models.AddEnbRequest{}

	// Invalid json: attribute name without quotes (should be "cause":).
	err := addEnbRequest.UnmarshalJSON([]byte("{ranName:\"test\""))
	assert.NotNil(t, err)
}
