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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"e2mgr/models"
)

func CreateInitialNodeInfo(requestDetails *models.SetupRequest, protocol entities.E2ApplicationProtocol) (*entities.NodebInfo, *entities.NbIdentity) {
	nodebInfo := &entities.NodebInfo{}
	nodebInfo.Ip = requestDetails.RanIp
	nodebInfo.Port = uint32(requestDetails.RanPort)
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTING
	nodebInfo.E2ApplicationProtocol = protocol
	nodebInfo.RanName = requestDetails.RanName
	nodebInfo.ConnectionAttempts = 0

	nodebIdentity := &entities.NbIdentity{}
	nodebIdentity.InventoryName = requestDetails.RanName
	return nodebInfo, nodebIdentity
}