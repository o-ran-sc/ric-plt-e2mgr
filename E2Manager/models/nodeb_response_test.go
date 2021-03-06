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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodebIdResponseMarshalSuccess(t *testing.T) {
	nodebInfo := &entities.NodebInfo{
		RanName:          "test",
		Ip:               "10.20.30.40",
		Port:             1234,
		NodeType:         entities.Node_GNB,
		GlobalNbId:       &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"},
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
	}
	response := models.NewNodebResponse(nodebInfo)
	resp, err := response.Marshal()

	m := jsonpb.Marshaler{}
	expectedData, _ := m.MarshalToString(nodebInfo)
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedData), resp)
}
