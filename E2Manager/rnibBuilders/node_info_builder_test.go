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

package rnibBuilders

import (
	"e2mgr/models"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

const ranName = "name"
const ranIP = "ip"
const ranPort = uint16(30000)

func TestCreateInitialNodeInfo(t *testing.T) {
	requestDetails :=  &models.RequestDetails{
		RanName: ranName,
		RanPort:ranPort,
		RanIp:ranIP,
	}
	nodebInfo, identity := CreateInitialNodeInfo(requestDetails, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	assert.Equal(t, identity.InventoryName, ranName)
	assert.Equal(t, nodebInfo.Ip, ranIP)
	assert.Equal(t, nodebInfo.ConnectionStatus, entities.ConnectionStatus_CONNECTING)
	assert.Equal(t, nodebInfo.E2ApplicationProtocol, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	assert.Equal(t, nodebInfo.Port, uint32(ranPort))
}