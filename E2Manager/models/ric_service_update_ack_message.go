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
		Text string `xml:",chardata"`
		//Ignore string `xml:"ignore"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text              string                         `xml:",chardata"`
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
	Value interface{} `xml:"value"`
}

type RICserviceUpdateAcknowledgeTransactionID struct {
	Text          string `xml:",chardata"`
	TransactionID string `xml:"TransactionID"`
}

type RICserviceUpdateAcknowledgeRANfunctionsList struct {
	Text               string `xml:",chardata"`
	RANfunctionsIDList struct {
		Text                      string                                                 `xml:",chardata"`
		ProtocolIESingleContainer []RICserviceUpdateAcknowledgeProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"RANfunctionsID-List"`
}
type RicServiceUpdateAckSuccessfulOutcome struct {
	XMLName       xml.Name `xml:"successfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                        string                      `xml:",chardata"`
		RICserviceUpdateAcknowledge RICserviceUpdateAcknowledge `xml:"RICserviceUpdateAcknowledge"`
	} `xml:"value"`
}

type RICserviceUpdateAcknowledge struct {
	Text        string `xml:",chardata"`
	ProtocolIEs struct {
		Text                           string                           `xml:",chardata"`
		RICserviceUpdateAcknowledgeIEs []RICserviceUpdateAcknowledgeIEs `xml:"RICserviceUpdateAcknowledge-IEs"`
	} `xml:"protocolIEs"`
}

type RicServiceUpdateAckE2APPDU struct {
	XMLName           xml.Name `xml:"E2AP-PDU"`
	Text              string   `xml:",chardata"`
	SuccessfulOutcome interface{}
}

func NewServiceUpdateAck(ricServiceUpdate []RicServiceAckRANFunctionIDItem, txId string) RicServiceUpdateAckE2APPDU {

	txIE := RICserviceUpdateAcknowledgeIEs{
		ID: ProtocolIE_ID_id_TransactionID,
		Value: RICserviceUpdateAcknowledgeTransactionID{
			TransactionID: txId,
		},
	}

	protocolIESingleContainer := make([]RICserviceUpdateAcknowledgeProtocolIESingleContainer, len(ricServiceUpdate))
	for i := 0; i < len(ricServiceUpdate); i++ {
		protocolIESingleContainer[i].Value.RANfunctionIDItem = ricServiceUpdate[i]
		protocolIESingleContainer[i].Id = ProtocolIE_ID_id_RANfunctionID_Item
	}

	// mandatory RANFunctionsAccepted
	functionListIE := RICserviceUpdateAcknowledgeIEs{
		ID: ProtocolIE_ID_id_RANfunctionsAccepted,
		Value: RICserviceUpdateAcknowledgeRANfunctionsList{
			RANfunctionsIDList: struct {
				Text                      string                                                 `xml:",chardata"`
				ProtocolIESingleContainer []RICserviceUpdateAcknowledgeProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
			}{
				ProtocolIESingleContainer: protocolIESingleContainer,
			},
		},
	}

	successfulOutcome := RicServiceUpdateAckSuccessfulOutcome{
		ProcedureCode: SuccessfulOutcome_value_PR_RICserviceUpdateAcknowledge,
		Value: struct {
			Text                        string                      `xml:",chardata"`
			RICserviceUpdateAcknowledge RICserviceUpdateAcknowledge `xml:"RICserviceUpdateAcknowledge"`
		}{
			RICserviceUpdateAcknowledge: RICserviceUpdateAcknowledge{
				ProtocolIEs: struct {
					Text                           string                           `xml:",chardata"`
					RICserviceUpdateAcknowledgeIEs []RICserviceUpdateAcknowledgeIEs `xml:"RICserviceUpdateAcknowledge-IEs"`
				}{
					RICserviceUpdateAcknowledgeIEs: []RICserviceUpdateAcknowledgeIEs{txIE, functionListIE},
				},
			},
		},
	}

	return RicServiceUpdateAckE2APPDU{SuccessfulOutcome: successfulOutcome}
}
