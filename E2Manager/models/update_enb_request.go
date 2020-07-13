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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package models

import (
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
)

type UpdateEnbRawRequest struct {
	RanName    string
	Enb        json.RawMessage
}

type UpdateEnbRequest struct {
	RanName    string
	Enb        *entities.Enb
}

func (r *UpdateEnbRequest) UnmarshalJSON(data []byte) error {
	updateEnbRawRequest := UpdateEnbRawRequest{}
	err := json.Unmarshal(data, &updateEnbRawRequest)

	if err != nil {
		return err
	}

	if updateEnbRawRequest.Enb != nil {
		enb := entities.Enb{}
		err = jsonpb.UnmarshalString(string(updateEnbRawRequest.Enb), &enb)

		if err != nil {
			return err
		}

		r.Enb = &enb
	}

	return nil
}