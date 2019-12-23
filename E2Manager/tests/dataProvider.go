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


package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

const (
	MaxMsgSize  int    = 4096
	Port        int    = 3801
	Flags       int    = 0
	MessageType int    = 1001
	RanPort     uint16 = 879
	RanName     string = "test"
	RanIp       string = "10.0.0.3"
)

var (
	DummyPayload = []byte{1, 2, 3, 4}
	DummyXAction = []byte{5, 6, 7, 8}
)

func GetPort() string {
	return "tcp:" + strconv.Itoa(Port)
}

func GetHttpRequest() *http.Request {
	data := map[string]interface{}{"ranIp": RanIp,
		"ranPort": RanPort, "ranName": RanName}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data)
	req, _ := http.NewRequest("POST", "https://localhost:3800/request", b)
	return req
}

func GetInvalidRequestDetails() *http.Request {
	data := map[string]interface{}{"ranIp": "256.0.0.0",
		"ranPort": RanPort, "ranName": RanName}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data)
	req, _ := http.NewRequest("POST", "https://localhost:3800/request", b)
	return req
}

func GetInvalidMessageType() *http.Request {
	data := map[string]interface{}{"ranIp": "1.2.3.4",
		"ranPort": RanPort, "ranName": RanName}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data)
	req, _ := http.NewRequest("POST", "https://localhost:3800/request", b)
	return req
}

func GetPackedPayload(t *testing.T) []byte {
	inputPayloadAsStr := "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	payload := make([]byte, len(inputPayloadAsStr)/2)
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}
	return payload
}
