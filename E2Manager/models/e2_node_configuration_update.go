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

type E2nodeConfigurationUpdateMessage struct {
	XMLName xml.Name `xml:"E2nodeConfigurationUpdateMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU struct {
		Text              string `xml:",chardata"`
		InitiatingMessage struct {
			Text          string `xml:",chardata"`
			ProcedureCode string `xml:"procedureCode"`
			Criticality   struct {
				Text   string `xml:",chardata"`
				Reject string `xml:"reject"`
			} `xml:"criticality"`
			Value struct {
				Text                      string                    `xml:",chardata"`
				E2nodeConfigurationUpdate E2nodeConfigurationUpdate `xml:"E2nodeConfigurationUpdate"`
			} `xml:"value"`
		} `xml:"initiatingMessage"`
	} `xml:"E2AP-PDU"`
}

type E2nodeConfigurationUpdate struct {
	Text        string `xml:",chardata"`
	ProtocolIEs struct {
		Text                         string                        `xml:",chardata"`
		E2nodeConfigurationUpdateIEs []E2nodeConfigurationUpdateIE `xml:"E2nodeConfigurationUpdate-IEs"`
	} `xml:"protocolIEs"`
}

type E2nodeConfigurationUpdateIE struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                              string                            `xml:",chardata"`
		E2nodeComponentConfigAdditionList E2nodeComponentConfigAdditionList `xml:"E2nodeComponentConfigAddition-List"`
		E2nodeComponentConfigUpdateList   E2nodeComponentConfigUpdateList   `xml:"E2nodeComponentConfigUpdate-List"`
		E2nodeComponentConfigRemovalList  E2nodeComponentConfigRemovalList  `xml:"E2nodeComponentConfigRemoval-List"`
	} `xml:"value"`
}

type E2nodeComponentConfigAdditionList struct {
	Text                      string `xml:",chardata"`
	ProtocolIESingleContainer []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"id"`
		Criticality struct {
			Text   string `xml:",chardata"`
			Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
			Text                              string `xml:",chardata"`
			E2nodeComponentConfigAdditionItem struct {
				Text                         string                       `xml:",chardata"`
				E2nodeComponentInterfaceType E2nodeComponentInterfaceType `xml:"e2nodeComponentInterfaceType"`
				E2nodeComponentID            E2nodeComponentID            `xml:"e2nodeComponentID"`
				E2nodeComponentConfiguration E2nodeComponentConfiguration `xml:"e2nodeComponentConfiguration"`
			} `xml:"E2nodeComponentConfigAddition-Item"`
		} `xml:"value"`
	} `xml:"ProtocolIE-SingleContainer"`
}

type E2nodeComponentConfigUpdateList struct {
	Text                      string `xml:",chardata"`
	ProtocolIESingleContainer []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"id"`
		Criticality struct {
			Text   string `xml:",chardata"`
			Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
			Text                            string `xml:",chardata"`
			E2nodeComponentConfigUpdateItem struct {
				Text                         string                       `xml:",chardata"`
				E2nodeComponentInterfaceType E2nodeComponentInterfaceType `xml:"e2nodeComponentInterfaceType"`
				E2nodeComponentID            E2nodeComponentID            `xml:"e2nodeComponentID"`
				E2nodeComponentConfiguration E2nodeComponentConfiguration `xml:"e2nodeComponentConfiguration"`
			} `xml:"E2nodeComponentConfigUpdate-Item"`
		} `xml:"value"`
	} `xml:"ProtocolIE-SingleContainer"`
}

type E2nodeComponentConfigRemovalList struct {
	Text                      string `xml:",chardata"`
	ProtocolIESingleContainer []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"id"`
		Criticality struct {
			Text   string `xml:",chardata"`
			Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
			Text                             string `xml:",chardata"`
			E2nodeComponentConfigRemovalItem struct {
				Text                         string                       `xml:",chardata"`
				E2nodeComponentInterfaceType E2nodeComponentInterfaceType `xml:"e2nodeComponentInterfaceType"`
				E2nodeComponentID            E2nodeComponentID            `xml:"e2nodeComponentID"`
			} `xml:"E2nodeComponentConfigRemoval-Item"`
		} `xml:"value"`
	} `xml:"ProtocolIE-SingleContainer"`
}

type E2nodeComponentInterfaceType struct {
	Text string    `xml:",chardata"`
	Ng   *struct{} `xml:"ng"`
	Xn   *struct{} `xml:"xn"`
	E1   *struct{} `xml:"e1"`
	F1   *struct{} `xml:"f1"`
	W1   *struct{} `xml:"w1"`
	S1   *struct{} `xml:"s1"`
	X2   *struct{} `xml:"x2"`
}

type E2nodeComponentID struct {
	Text                           string `xml:",chardata"`
	E2nodeComponentInterfaceTypeNG struct {
		Text    string `xml:",chardata"`
		AmfName string `xml:"amf-name"`
	} `xml:"e2nodeComponentInterfaceTypeNG"`
	E2nodeComponentInterfaceTypeXn struct {
		Text              string `xml:",chardata"`
		GlobalNGRANNodeID string `xml:"global-NG-RAN-Node-ID"`
	} `xml:"e2nodeComponentInterfaceTypeXn"`
	E2nodeComponentInterfaceTypeE1 struct {
		Text      string `xml:",chardata"`
		GNBCUCPID string `xml:"gNB-CU-CP-ID"`
	} `xml:"e2nodeComponentInterfaceTypeE1"`
	E2nodeComponentInterfaceTypeF1 struct {
		Text    string `xml:",chardata"`
		GNBDUID string `xml:"gNB-DU-ID"`
	} `xml:"e2nodeComponentInterfaceTypeF1"`
	E2nodeComponentInterfaceTypeW1 struct {
		Text      string `xml:",chardata"`
		NgENBDUID string `xml:"ng-eNB-DU-ID"`
	} `xml:"e2nodeComponentInterfaceTypeW1"`
	E2nodeComponentInterfaceTypeS1 struct {
		Text    string `xml:",chardata"`
		MmeName string `xml:"mme-name"`
	} `xml:"e2nodeComponentInterfaceTypeS1"`
	E2nodeComponentInterfaceTypeX2 struct {
		Text          string `xml:",chardata"`
		GlobalENBID   string `xml:"global-eNB-ID"`
		GlobalEnGNBID string `xml:"global-en-gNB-ID"`
	} `xml:"e2nodeComponentInterfaceTypeX2"`
}

type E2nodeComponentConfiguration struct {
	Text                        string `xml:",chardata"`
	E2nodeComponentRequestPart  string `xml:"e2nodeComponentRequestPart"`
	E2nodeComponentResponsePart string `xml:"e2nodeComponentResponsePart"`
}
