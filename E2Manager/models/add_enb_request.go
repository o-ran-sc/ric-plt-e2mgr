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

type AddEnbRawRequest struct {
	RanName    string
	GlobalNbId json.RawMessage
	Ip         string
	Port       uint32
	Enb        json.RawMessage
}

type AddEnbRequest struct {
	RanName    string
	GlobalNbId *entities.GlobalNbId
	Ip         string
	Port       uint32
	Enb        *entities.Enb
}

func (r *AddEnbRequest) UnmarshalJSON(data []byte) error {
	addEnbRawRequest := AddEnbRawRequest{}
	err := json.Unmarshal(data, &addEnbRawRequest)

	if err != nil {
		return err
	}

	r.RanName = addEnbRawRequest.RanName
	r.Ip = addEnbRawRequest.Ip
	r.Port = addEnbRawRequest.Port

	globalNbId := entities.GlobalNbId{}
	err = jsonpb.UnmarshalString(string(addEnbRawRequest.GlobalNbId), &globalNbId)

	if err != nil {
		return err
	}

	r.GlobalNbId = &globalNbId

	enb := entities.Enb{}
	err = jsonpb.UnmarshalString(string(addEnbRawRequest.Enb), &enb)

	if err != nil {
		return err
	}

	r.Enb = &enb
	return nil
}