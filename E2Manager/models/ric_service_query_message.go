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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type RanFunctionIdItem struct {
	Text                  string `xml:",chardata"`
	RanFunctionId         uint32 `xml:"ranFunctionID"`
	RanFunctionRevision   uint32 `xml:"ranFunctionRevision"`
}

type RicServiceQueryProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	Id          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text               string          `xml:",chardata"`
		RanFunctionIdItem  RanFunctionIdItem `xml:"RANfunctionID-Item"`
	} `xml:"value"`
}

type RICServiceQueryIEs struct {
	Text        string `xml:",chardata"`
	Id          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text   string `xml:",chardata"`
		RANFunctionIdList struct {
			Text string `xml:",chardata"`
			ProtocolIESingleContainer []RicServiceQueryProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
		} `xml:"RANfunctionsID-List"`
	} `xml:"value"`
}

type RICServiceQuery struct {
	Text    string   `xml:",chardata"`
	ProtocolIEs struct {
		Text              string `xml:",chardata"`
		RICServiceQueryIEs []RICServiceQueryIEs `xml:"RICserviceQuery-IEs"`
	} `xml:"protocolIEs"`
}

type InitiatingMessage struct {
	Text          string `xml:",chardata"`
	ProcedureCode string `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text           string         `xml:",chardata"`
		RICServiceQuery RICServiceQuery `xml:"RICserviceQuery"`
	} `xml:"value"`
}

type RicServiceQueryE2APPDU struct {
	XMLName xml.Name `xml:"E2AP-PDU"`
	Text              string `xml:",chardata"`
	InitiatingMessage InitiatingMessage `xml:"initiatingMessage"`
}


type RICServiceQueryMessage struct{
	XMLName xml.Name `xml:"RICserviceQueryMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU RicServiceQueryE2APPDU  `xml:"E2AP-PDU"`
}

func NewRicServiceQueryMessage(ranFunctions []*entities.RanFunction) RICServiceQueryMessage {
	initiatingMessage := InitiatingMessage{}
	initiatingMessage.ProcedureCode = "6"
	initiatingMessage.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs = make([]RICServiceQueryIEs,1)
	initiatingMessage.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[0].Id = "9"
	protocolIESingleContainer := make([]RicServiceQueryProtocolIESingleContainer,len(ranFunctions))
	for i := 0; i < len(ranFunctions); i++ {
		protocolIESingleContainer[i].Id = "6"
		protocolIESingleContainer[i].Value.RanFunctionIdItem.RanFunctionId = ranFunctions[i].RanFunctionId
		protocolIESingleContainer[i].Value.RanFunctionIdItem.RanFunctionRevision = ranFunctions[i].RanFunctionRevision
	}
	initiatingMessage.Value.RICServiceQuery.ProtocolIEs.RICServiceQueryIEs[0].Value.RANFunctionIdList.ProtocolIESingleContainer = protocolIESingleContainer

	return RICServiceQueryMessage{E2APPDU:RicServiceQueryE2APPDU{InitiatingMessage:initiatingMessage}}
}


