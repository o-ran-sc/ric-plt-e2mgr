//
// Copyright 2022 Samsung Electronics Co.
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
	"encoding/xml"
)

type E2ResetResponseMessage struct {
	XMLName xml.Name `xml:"E2ResetSuccessResponseMessage"`
	Text    string   `xml:",chardata"`
	E2ApPdu E2ApPdu  `xml:"E2AP-PDU"`
}

type E2ApPdu struct {
	XMLName           xml.Name          `xml:"E2AP-PDU"`
	Text              string            `xml:",chardata"`
	SuccessfulOutcome successfulOutcome `xml:"successfulOutcome"`
}

type successfulOutcome struct {
	XMLName       xml.Name `xml:"successfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value SuccessfulOutcomeValue `xml:"value"`
}

type Criticality struct {
	Ignore string `xml:"ignore"`
}

type SuccessfulOutcomeValue struct {
	Text          string        `xml:",chardata"`
	ResetResponse ResetResponse `xml:"ResetResponse"`
}

type ResetResponse struct {
	Text        string      `xml:",chardata"`
	ProtocolIEs ProtocolIEs `xml:"protocolIEs"`
}

type ProtocolIEs struct {
	Text             string             `xml:",chardata"`
	ResetResponseIEs []ResetResponseIEs `xml:"ResetResponseIEs"`
}

type ResetResponseIEs struct {
	Text        string                `xml:",chardata"`
	ID          string                `xml:"id"`
	Criticality Criticality           `xml:"criticality"`
	Value       ResetResponseIEsValue `xml:"value"`
}

type ResetResponseIEsValue struct {
	TransactionID string `xml:"TransactionID"`
}

func NewE2ResetResponseMessage(request *E2ResetRequestMessage) E2ResetResponseMessage {
	outcome := successfulOutcome{}
	outcome.ProcedureCode = request.E2ApPDU.InitiatingMessage.ProcedureCode
	e2ResetRequestIes := request.E2ApPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs
	numOfIes := len(e2ResetRequestIes)

	outcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs = make([]ResetResponseIEs, numOfIes)
	for ieCount := 0; ieCount < numOfIes; ieCount++ {
		outcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs[ieCount].ID = request.E2ApPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs[ieCount].ID
		outcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs[ieCount].Criticality.Ignore = request.E2ApPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs[ieCount].Criticality.Ignore
		outcome.Value.ResetResponse.ProtocolIEs.ResetResponseIEs[ieCount].Value.TransactionID = request.E2ApPDU.InitiatingMessage.Value.ResetRequest.ProtocolIEs.ResetRequestIEs[ieCount].Value.TransactionID
	}
	return E2ResetResponseMessage{E2ApPDU: E2ApPDU{SuccessfulOutcome: outcome}}
}
