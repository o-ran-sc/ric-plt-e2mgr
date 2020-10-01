//
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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

type RicServiceAckRANFunctionIDItem struct {
	Text                string `xml:",chardata"`
        RanFunctionID       uint32 `xml:"ranFunctionID"`
        RanFunctionRevision uint32 `xml:"ranFunctionRevision"`
}

type RICserviceUpdateAcknowledgeProtocolIESingleContainer struct {
		Text        string `xml:",chardata"`
		Id          string `xml:"id"`
		Criticality struct {
				Text   string `xml:",chardata"`
				//Ignore string `xml:"ignore"`
				Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
				Text              string `xml:",chardata"`
				RANfunctionIDItem RicServiceAckRANFunctionIDItem `xml:"RANfunctionID-Item"`
		} `xml:"value"`
}

type RICserviceUpdateAcknowledgeIEs struct {
	Text        string `xml:",chardata"`
        ID          string `xml:"id"`
        Criticality struct {
                Text   string `xml:",chardata"`
                Reject string `xml:"reject"`
        } `xml:"criticality"`
        Value struct {
		Text   string `xml:",chardata"`
		RANfunctionsIDList struct {
			Text   string `xml:",chardata"`
			ProtocolIESingleContainer []RICserviceUpdateAcknowledgeProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
		} `xml:"RANfunctionsID-List"`
	}`xml:"value"`
}

type RicServiceUpdateAckSuccessfulOutcome struct {
	XMLName xml.Name        `xml:"successfulOutcome"`
        Text          string `xml:",chardata"`
        ProcedureCode string `xml:"procedureCode"`
        Criticality   struct {
                Text   string `xml:",chardata"`
                Reject string `xml:"reject"`
        } `xml:"criticality"`
	Value struct {
		Text            string `xml:",chardata"`
		RICserviceUpdateAcknowledge  RICserviceUpdateAcknowledge `xml:"RICserviceUpdateAcknowledge"`
	} `xml:"value"`
}

type RICserviceUpdateAcknowledge struct {
	Text        string `xml:",chardata"`
	ProtocolIEs struct {
		Text               string `xml:",chardata"`
		RICserviceUpdateAcknowledgeIEs []RICserviceUpdateAcknowledgeIEs `xml:"RICserviceUpdateAcknowledge-IEs"`
	} `xml:"protocolIEs"`
}

type RicServiceUpdateAckE2APPDU struct {
        XMLName xml.Name `xml:"E2AP-PDU"`
        Text    string `xml:",chardata"`
        InitiatingMessage interface{}
}

func NewServiceUpdateAck(ricServiceUpdate []RicServiceAckRANFunctionIDItem) RicServiceUpdateAckE2APPDU {
	successfulOutcome := RicServiceUpdateAckSuccessfulOutcome{}
	successfulOutcome.ProcedureCode = "7"
	if len(ricServiceUpdate) == 0 {
		return RicServiceUpdateAckE2APPDU{InitiatingMessage:successfulOutcome}
	}

	ricServiceUpdateAcknowledgeIEs := make([]RICserviceUpdateAcknowledgeIEs, 1)
	ricServiceUpdateAcknowledgeIEs[0].ID = "9"
	protocolIESingleContainer := make([]RICserviceUpdateAcknowledgeProtocolIESingleContainer, len(ricServiceUpdate))
	for i := 0; i < len(ricServiceUpdate); i++ {
		protocolIESingleContainer[i].Value.RANfunctionIDItem = ricServiceUpdate[i]
		protocolIESingleContainer[i].Id = "6"
	}
	ricServiceUpdateAcknowledgeIEs[0].Value.RANfunctionsIDList.ProtocolIESingleContainer = protocolIESingleContainer
	successfulOutcome.Value.RICserviceUpdateAcknowledge.ProtocolIEs.RICserviceUpdateAcknowledgeIEs = ricServiceUpdateAcknowledgeIEs
	return RicServiceUpdateAckE2APPDU{InitiatingMessage:successfulOutcome}
}

