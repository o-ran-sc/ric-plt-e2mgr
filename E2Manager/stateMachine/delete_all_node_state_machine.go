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

package stateMachine

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)


var nodeStateMap = map[entities.ConnectionStatus]entities.ConnectionStatus{
	entities.ConnectionStatus_CONNECTING:	entities.ConnectionStatus_SHUTTING_DOWN,
	entities.ConnectionStatus_CONNECTED:	entities.ConnectionStatus_SHUTTING_DOWN,
	entities.ConnectionStatus_CONNECTED_SETUP_FAILED:	entities.ConnectionStatus_SHUTTING_DOWN,
	entities.ConnectionStatus_DISCONNECTED:	entities.ConnectionStatus_SHUT_DOWN,
}

func NodeNextStateDeleteAll(state entities.ConnectionStatus) (entities.ConnectionStatus, bool) {
	nextState, error := nodeStateMap[state]

	if !error {
		return state, false
	}

	return nextState, true
}
