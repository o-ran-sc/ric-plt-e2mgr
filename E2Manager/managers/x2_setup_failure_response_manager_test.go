//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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
//
package managers

import (
	"e2mgr/converters"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPopulateX2NodebByPduFailure(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	handler := NewX2SetupFailureResponseManager(converters.NewX2SetupFailureResponseConverter(logger))
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createRandomPayload())
	assert.NotNil(t, err)
}

func TestPopulateX2NodebByPduSuccess(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	handler := NewX2SetupFailureResponseManager(converters.NewX2SetupFailureResponseConverter(logger))
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createX2SetupFailureResponsePayload(t))
	assert.Nil(t, err)
	assert.Equal(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.Equal(t, entities.Failure_X2_SETUP_FAILURE, nodebInfo.FailureType)

}

func createX2SetupFailureResponsePayload(t *testing.T) []byte {
	packedPdu := "4006001a0000030005400200000016400100001140087821a00000008040"
	var payload []byte
	_, err := fmt.Sscanf(packedPdu, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}
	return payload
}
