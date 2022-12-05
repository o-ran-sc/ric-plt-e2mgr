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
	XMLName xml.Name `xml:"E2ResetRequestMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU struct {
		XMLName           xml.Name `xml:"E2AP-PDU"`
		Text              string   `xml:",chardata"`
		SuccessfulOutcome struct {
			Text          string `xml:",chardata"`
			ProcedureCode string `xml:"procedureCode"`
			Criticality   struct {
				Text   string `xml:",chardata"`
				Ignore string `xml:"ignore"`
			} `xml:"criticality"`
			Value struct {
				Text          string `xml:",chardata"`
				ResetResponse struct {
					Text        string `xml:",chardata"`
					ProtocolIEs struct {
						Text             string `xml:",chardata"`
						ResetResponseIEs []struct {
							Text        string `xml:",chardata"`
							ID          string `xml:"id"`
							Criticality struct {
								Text   string `xml:",chardata"`
								Ignore string `xml:"ignore"`
							} `xml:"criticality"`
							Value struct {
								Text          string `xml:",chardata"`
								TransactionID string `xml:"TransactionID"`
							} `xml:"value"`
						} `xml:"ResetResponseIEs"`
					} `xml:"protocolIEs"`
				} `xml:"ResetResponse"`
			} `xml:"value"`
		} `xml:"successfulOutcome"`
	} `xml:"E2AP-PDU"`
}
