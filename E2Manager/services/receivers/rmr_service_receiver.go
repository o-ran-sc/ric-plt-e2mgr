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

package receivers

import (
	"e2mgr/managers/notificationmanager"
	"e2mgr/services"
)

// RmrService holds an instance of RMR messenger as well as its configuration
type RmrServiceReceiver struct {
	services.RmrService
	nManager *notificationmanager.NotificationManager
}

// NewRmrService instantiates a new Rmr service instance
func NewRmrServiceReceiver(rmrService services.RmrService, nManager *notificationmanager.NotificationManager) *RmrServiceReceiver {

	return &RmrServiceReceiver{
		RmrService: rmrService,
		nManager:   nManager,
	}
}

// ListenAndHandle waits for messages coming from rmr_rcv_msg and sends it to a designated message handler
func (r *RmrServiceReceiver) ListenAndHandle() {

	for {
		mbuf, err := (*r.Messenger).RecvMsg()
		r.Config.Logger.Debugf("#rmr_service_receiver.ListenAndHandle - Going to handle received message: %#v\n", mbuf)

		// TODO: one mbuf received immediately execute goroutine
		if err != nil {
			continue //TODO log error
		}

		r.nManager.HandleMessage(r.Config.Logger, r.E2sessions, mbuf, r.RmrResponse)
	}
}
