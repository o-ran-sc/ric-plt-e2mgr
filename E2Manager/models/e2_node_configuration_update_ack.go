//
// Copyright (c) 2022 Samsung Electronics Co., Ltd. All Rights Reserved.
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

type E2nodeConfigurationUpdateAcknowledgeE2APPDU struct {
	XMLName xml.Name `xml:"E2AP-PDU"`
	Text    string   `xml:",chardata"`
	Outcome interface{}
}

type E2nodeConfigurationUpdateAcknowledgeSuccessfulOutcome struct {
	XMLName       xml.Name `xml:"successfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                                 string `xml:",chardata"`
		E2nodeConfigurationUpdateAcknowledge struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text                                    string                                    `xml:",chardata"`
				E2nodeConfigurationUpdateAcknowledgeIEs []E2nodeConfigurationUpdateAcknowledgeIEs `xml:"E2nodeConfigurationUpdateAcknowledge-IEs"`
			} `xml:"protocolIEs"`
		} `xml:"E2nodeConfigurationUpdateAcknowledge"`
	} `xml:"value"`
}

type E2nodeConfigurationUpdateAcknowledgeIEs struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value interface{} `xml:"value"`
}

type E2nodeConfigurationUpdateAcknowledgeTransID struct {
	Text          string `xml:",chardata"`
	TransactionID string `xml:"TransactionID"`
}

type E2nodeComponentConfigAdditionAckList struct {
	Text                                 string `xml:",chardata"`
	E2nodeComponentConfigAdditionAckList struct {
		Text                      string                                  `xml:",chardata"`
		ProtocolIESingleContainer []AdditionListProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"E2nodeComponentConfigAdditionAck-List"`
}

type AdditionListProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                                 string             `xml:",chardata"`
		E2nodeComponentConfigAdditionAckItem ComponentAckDetail `xml:"E2nodeComponentConfigAdditionAck-Item"`
	} `xml:"value"`
}

type E2nodeComponentConfigUpdateAckList struct {
	Text                               string `xml:",chardata"`
	E2nodeComponentConfigUpdateAckList struct {
		Text                      string                            `xml:",chardata"`
		ProtocolIESingleContainer []UpdateProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"E2nodeComponentConfigUpdateAck-List"`
}

type UpdateProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                               string             `xml:",chardata"`
		E2nodeComponentConfigUpdateAckItem ComponentAckDetail `xml:"E2nodeComponentConfigUpdateAck-Item"`
	} `xml:"value"`
}

type E2nodeComponentConfigRemovalAckList struct {
	Text                                string `xml:",chardata"`
	E2nodeComponentConfigRemovalAckList struct {
		Text                      string                             `xml:",chardata"`
		ProtocolIESingleContainer []RemovalProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"E2nodeComponentConfigRemovalAck-List"`
}

type RemovalProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                                string             `xml:",chardata"`
		E2nodeComponentConfigRemovalAckItem ComponentAckDetail `xml:"E2nodeComponentConfigRemovalAck-Item"`
	} `xml:"value"`
}

type ComponentAckDetail struct {
	Text                         string                `xml:",chardata"`
	E2nodeComponentInterfaceType E2NodeComponentType   `xml:"e2nodeComponentInterfaceType"`
	E2nodeComponentID            E2NodeComponentIDResp `xml:"e2nodeComponentID"`
	E2nodeConfigUpdateAck        E2nodeConfigUpdateAckResp
}
