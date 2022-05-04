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
	"strconv"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
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

func (m *E2nodeConfigurationUpdateMessage) ExtractConfigAdditionList() []entities.E2NodeComponentConfig {
	var result []entities.E2NodeComponentConfig
	e2nodeConfigUpdateIEs := m.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs

	var additionList *E2nodeComponentConfigAdditionList
	for _, v := range e2nodeConfigUpdateIEs {
		if v.ID == ProtocolIE_ID_id_E2nodeComponentConfigAddition {
			additionList = &(v.Value.E2nodeComponentConfigAdditionList)
			break
		}
	}

	// need not to check for empty addtionList
	// as list defined as SIZE(1..maxofE2nodeComponents)
	if additionList != nil {
		for _, val := range additionList.ProtocolIESingleContainer {
			componentItem := val.Value.E2nodeComponentConfigAdditionItem
			if componentItem.E2nodeComponentInterfaceType.Ng != nil { // NG Interface
				result = append(result, entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_ng,
					E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
					E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeNG{
						E2NodeComponentInterfaceTypeNG: &entities.E2NodeComponentInterfaceNG{
							AmfName: componentItem.E2nodeComponentID.E2nodeComponentInterfaceTypeNG.AmfName,
						},
					},
				})
			}

			// TODO - Not Supported Yet
			if componentItem.E2nodeComponentInterfaceType.Xn != nil { // xn inetrface
				// result = append(result, entities.E2NodeComponentConfig{
				// 	E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_xn,
				// 	E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
				// 	E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
				// 	E2NodeComponentID:            &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeXn{},
				// })
			}

			if componentItem.E2nodeComponentInterfaceType.E1 != nil { // e1 interface
				gnbCuCpId, err := strconv.ParseInt(componentItem.E2nodeComponentID.E2nodeComponentInterfaceTypeE1.GNBCUCPID, 10, 64)
				if err != nil {
					continue
				}

				result = append(result, entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_e1,
					E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
					E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeE1{
						E2NodeComponentInterfaceTypeE1: &entities.E2NodeComponentInterfaceE1{
							GNBCuCpId: gnbCuCpId,
						},
					},
				})
			}

			if componentItem.E2nodeComponentInterfaceType.F1 != nil { // f1 interface
				gnbDuId, err := strconv.ParseInt(componentItem.E2nodeComponentID.E2nodeComponentInterfaceTypeF1.GNBDUID, 10, 64)
				if err != nil {
					continue
				}
				result = append(result, entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_f1,
					E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
					E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeF1{
						E2NodeComponentInterfaceTypeF1: &entities.E2NodeComponentInterfaceF1{
							GNBDuId: gnbDuId,
						},
					},
				})
			}

			if componentItem.E2nodeComponentInterfaceType.W1 != nil { // w1 interface
				ngenbDuId, err := strconv.ParseInt(componentItem.E2nodeComponentID.E2nodeComponentInterfaceTypeW1.NgENBDUID, 10, 64)
				if err != nil {
					continue
				}

				result = append(result, entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_w1,
					E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
					E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeW1{
						E2NodeComponentInterfaceTypeW1: &entities.E2NodeComponentInterfaceW1{
							NgenbDuId: ngenbDuId,
						},
					},
				})
			}

			if componentItem.E2nodeComponentInterfaceType.S1 != nil { // s1 interface
				result = append(result, entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_s1,
					E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
					E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeS1{
						E2NodeComponentInterfaceTypeS1: &entities.E2NodeComponentInterfaceS1{
							MmeName: componentItem.E2nodeComponentID.E2nodeComponentInterfaceTypeS1.MmeName,
						},
					},
				})
			}

			// TODO - Not Supported Yet
			if componentItem.E2nodeComponentInterfaceType.X2 != nil { // x2 interface
				// result = append(result, entities.E2NodeComponentConfig{
				// 	E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
				// 	E2NodeComponentRequestPart:   []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentRequestPart),
				// 	E2NodeComponentResponsePart:  []byte(componentItem.E2nodeComponentConfiguration.E2nodeComponentResponsePart),
				// 	E2NodeComponentID:            &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{},
				// })

			}
		}
	}
	return result
}
