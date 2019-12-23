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


package utils

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"strings"
)

func ConvertNodebIdListToProtoMessageList(l []*entities.NbIdentity) []proto.Message {
	protoMessageList := make([]proto.Message, len(l))

	for i, d := range l {
		protoMessageList[i] = d
	}

	return protoMessageList
}

func MarshalProtoMessageListToJsonArray(msgList []proto.Message) (string, error){
	m := jsonpb.Marshaler{}
	ms := "["

	for _, msg := range msgList {
		s, err :=m.MarshalToString(msg)

		if (err != nil) {
			return s, err;
		}


		ms+=s+","
	}

	return strings.TrimSuffix(ms,",") +"]", nil
}
