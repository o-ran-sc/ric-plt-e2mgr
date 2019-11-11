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

package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// TODO: message command id / source / dest

type MessageInfo struct {
	MessageTimestamp int64  `json:"messageTimestamp"`
	MessageType      int    `json:"messageType"`
	Meid             string `json:"meid"`
	Payload          string `json:"payload"`
	TransactionId    string `json:"transactionId"`
}

func GetMessageInfoAsJson(messageType int, meid string, payload []byte, transactionId []byte) string {
	messageInfo := MessageInfo{
		MessageTimestamp: time.Now().Unix(),
		MessageType:      messageType,
		Meid:             meid,
		Payload:          fmt.Sprintf("%x", payload),
		TransactionId:    string(transactionId),
	}

	jsonData, _ := json.Marshal(messageInfo)

	return string(jsonData)
}
