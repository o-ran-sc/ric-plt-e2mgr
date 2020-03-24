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
						E2setupResponseIEs []E2setupResponseIEs`xml:"E2setupResponseIEs"`
					} `xml:"protocolIEs"`
				} `xml:"E2setupResponse"`
			} `xml:"value"`
		} `xml:"successfulOutcome"`
	} `xml:"E2AP-PDU"`
}

type E2setupResponseIEs struct {
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
		RANfunctionsIDList struct {
			Text                      string `xml:",chardata"`
			ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
		} `xml:"RANfunctionsID-List"`
	} `xml:"value"`
}

type ProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text              string `xml:",chardata"`
		RANfunctionIDItem struct {
			Text                string `xml:",chardata"`
			RanFunctionID       string `xml:"ranFunctionID"`
			RanFunctionRevision string `xml:"ranFunctionRevision"`
		} `xml:"RANfunctionID-Item"`
	} `xml:"value"`
}

func NewE2SetupSuccessResponseMessage() *E2SetupSuccessResponseMessage{
	msg := &E2SetupSuccessResponseMessage{}
	msg.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs = make([]E2setupResponseIEs, 2)
	return msg
}

func (m *E2SetupSuccessResponseMessage) SetExtractRanFunctionsIDList(request *E2SetupRequestMessage) {
	list := &request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[1].Value.RANfunctionsList
	ids := make([]ProtocolIESingleContainer,len(list.ProtocolIESingleContainer))
	for i := 0; i< len(ids); i++{
		ids[i] = m.convertToRANfunctionID(list, i)
	}
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[1].Value.RANfunctionsIDList.ProtocolIESingleContainer = ids
}

func (m *E2SetupSuccessResponseMessage) convertToRANfunctionID(list *RANfunctionsList, i int) ProtocolIESingleContainer{
	id := ProtocolIESingleContainer{}
	id.Value.RANfunctionIDItem.RanFunctionID = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionID
	id.Value.RANfunctionIDItem.RanFunctionRevision = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionRevision
	return id
}

func (m *E2SetupSuccessResponseMessage) SetPlmnId(plmnId string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[0].Value.GlobalRICID.PLMNIdentity = plmnId
}

func (m *E2SetupSuccessResponseMessage) SetRicId(ricId string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[0].Value.GlobalRICID.RicID = ricId
}