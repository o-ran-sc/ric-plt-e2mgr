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
	"encoding/xml"
)

type E2SetupSuccessResponseMessage struct {
	XMLName xml.Name `xml:"E2SetupSuccessResponseMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU struct {
		Text              string `xml:",chardata"`
		SuccessfulOutcome struct {
			Text          string `xml:",chardata"`
			ProcedureCode string `xml:"procedureCode"`
			Criticality   struct {
				Text   string `xml:",chardata"`
				Reject string `xml:"reject"`
			} `xml:"criticality"`
			Value struct {
				Text            string `xml:",chardata"`
				E2setupResponse struct {
					Text        string `xml:",chardata"`
					ProtocolIEs struct {
						Text               string `xml:",chardata"`
						E2setupResponseIEs struct {
							Text        string `xml:",chardata"`
							ID          string `xml:"id"`
							Criticality struct {
								Text   string `xml:",chardata"`
								Reject string `xml:"reject"`
							} `xml:"criticality"`
							Value struct {
								Text        string `xml:",chardata"`
								GlobalRICID struct {
									Text         string `xml:",chardata"`
									PLMNIdentity string `xml:"pLMN-Identity"`
									RicID        string `xml:"ric-ID"`
								} `xml:"GlobalRIC-ID"`
							} `xml:"value"`
						} `xml:"E2setupResponseIEs"`
					} `xml:"protocolIEs"`
				} `xml:"E2setupResponse"`
			} `xml:"value"`
		} `xml:"successfulOutcome"`
	} `xml:"E2AP-PDU"`
}


func (m *E2SetupSuccessResponseMessage) SetPlmnId(plmnId string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs.Value.GlobalRICID.PLMNIdentity = plmnId
}

func (m *E2SetupSuccessResponseMessage) SetNbId(ricID string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs.Value.GlobalRICID.RicID = ricID
}