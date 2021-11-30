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

type RANfunctionsItemProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	Id          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text            string `xml:",chardata"`
		RANfunctionItem struct {
			Text                  string `xml:",chardata"`
			RanFunctionID         uint32 `xml:"ranFunctionID"`
			RanFunctionDefinition string `xml:"ranFunctionDefinition"`
			RanFunctionRevision   uint32 `xml:"ranFunctionRevision"`
			RanFunctionOID        string `xml:"ranFunctionOID"`
		} `xml:"RANfunction-Item"`
	} `xml:"value"`
}

type RANfunctionsItemIDProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	Id          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text              string `xml:",chardata"`
		RANfunctionIDItem struct {
			Text                string `xml:",chardata"`
			RanFunctionID       uint32 `xml:"ranFunctionID"`
			RanFunctionRevision uint32 `xml:"ranFunctionRevision"`
		} `xml:"RANfunctionID-Item"`
	} `xml:"value"`
}

type RICServiceUpdateIEs struct {
	Text        string `xml:",chardata"`
	ID          int    `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text             string `xml:",chardata"`
		RANfunctionsList struct {
			Text                                      string                                      `xml:",chardata"`
			RANfunctionsItemProtocolIESingleContainer []RANfunctionsItemProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
		} `xml:"RANfunctions-List"`
		RANfunctionsIDList struct {
			Text                                        string                                        `xml:",chardata"`
			RANfunctionsItemIDProtocolIESingleContainer []RANfunctionsItemIDProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
		} `xml:"RANfunctionsID-List"`
	} `xml:"value"`
}

type RICServiceUpdateInitiatingMessage struct {
	Text          string `xml:",chardata"`
	ProcedureCode string `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text             string `xml:",chardata"`
		RICServiceUpdate struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text                string                `xml:",chardata"`
				RICServiceUpdateIEs []RICServiceUpdateIEs `xml:"RICserviceUpdate-IEs"`
			} `xml:"protocolIEs"`
		} `xml:"RICserviceUpdate"`
	} `xml:"value"`
}

type RICServiceUpdateE2APPDU struct {
	XMLName           xml.Name                          `xml:"E2AP-PDU"`
	Text              string                            `xml:",chardata"`
	InitiatingMessage RICServiceUpdateInitiatingMessage `xml:"initiatingMessage"`
}

type RICServiceUpdateMessage struct {
	XMLName xml.Name                `xml:"RICserviceUpdateMessage"`
	Text    string                  `xml:",chardata"`
	E2APPDU RICServiceUpdateE2APPDU `xml:"E2AP-PDU"`
}

func (m *RICServiceUpdateE2APPDU) ExtractRanFunctionsList() []*entities.RanFunction {
	serviceUpdateRequestIes := m.InitiatingMessage.Value.RICServiceUpdate.ProtocolIEs.RICServiceUpdateIEs
	if len(serviceUpdateRequestIes) < 2 {
		return nil
	}

	ranFunctionsListContainer := serviceUpdateRequestIes[1].Value.RANfunctionsList.RANfunctionsItemProtocolIESingleContainer
	funcs := make([]*entities.RanFunction, len(ranFunctionsListContainer))
	for i := 0; i < len(funcs); i++ {
		ranFunctionItem := ranFunctionsListContainer[i].Value.RANfunctionItem

		funcs[i] = &entities.RanFunction{
			RanFunctionId:         ranFunctionItem.RanFunctionID,
			RanFunctionDefinition: ranFunctionItem.RanFunctionDefinition,
			RanFunctionRevision:   ranFunctionItem.RanFunctionRevision,
			RanFunctionOid:        ranFunctionItem.RanFunctionOID,
		}
	}
	return funcs
}
