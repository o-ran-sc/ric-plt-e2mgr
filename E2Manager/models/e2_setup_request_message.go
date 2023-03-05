//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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
	"strings"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type Gnb struct {
	Text        string `xml:",chardata"`
	GlobalGNBID struct {
		Text   string `xml:",chardata"`
		PlmnID string `xml:"plmn-id"`
		GnbID  struct {
			Text  string `xml:",chardata"`
			GnbID string `xml:"gnb-ID"`
		} `xml:"gnb-id"`
	} `xml:"global-gNB-ID"`
}

type EnGnb struct {
	Text        string `xml:",chardata"`
	GlobalGNBID struct {
		Text   string `xml:",chardata"`
		PlmnID string `xml:"pLMN-Identity"`
		GnbID  struct {
			Text  string `xml:",chardata"`
			GnbID string `xml:"gNB-ID"`
		} `xml:"gNB-ID"`
	} `xml:"global-gNB-ID"`
}

type NgEnbId struct {
	Text            string `xml:",chardata"`
	EnbIdMacro      string `xml:"enb-ID-macro"`
	EnbIdShortMacro string `xml:"enb-ID-shortmacro"`
	EnbIdLongMacro  string `xml:"enb-ID-longmacro"`
}

type NgEnb struct {
	Text          string `xml:",chardata"`
	GlobalNgENBID struct {
		Text   string  `xml:",chardata"`
		PlmnID string  `xml:"plmn-id"`
		EnbID  NgEnbId `xml:"enb-id"`
	} `xml:"global-ng-eNB-ID"`
}

type EnbId struct {
	Text            string `xml:",chardata"`
	MacroEnbId      string `xml:"macro-eNB-ID"`
	HomeEnbId       string `xml:"home-eNB-ID"`
	ShortMacroEnbId string `xml:"short-Macro-eNB-ID"`
	LongMacroEnbId  string `xml:"long-Macro-eNB-ID"`
}

type Enb struct {
	Text        string `xml:",chardata"`
	GlobalENBID struct {
		Text   string `xml:",chardata"`
		PlmnID string `xml:"pLMN-Identity"`
		EnbID  EnbId  `xml:"eNB-ID"`
	} `xml:"global-eNB-ID"`
}

type GlobalE2NodeId struct {
	Text  string `xml:",chardata"`
	GNB   Gnb    `xml:"gNB"`
	EnGNB EnGnb  `xml:"en-gNB"`
	NgENB NgEnb  `xml:"ng-eNB"`
	ENB   Enb    `xml:"eNB"`
}

type E2SetupRequest struct {
	Text        string `xml:",chardata"`
	ProtocolIEs struct {
		Text              string `xml:",chardata"`
		E2setupRequestIEs []struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id"`
			Criticality struct {
				Text   string `xml:",chardata"`
				Reject string `xml:"reject"`
			} `xml:"criticality"`
			Value struct {
				Text             string           `xml:",chardata"`
				TransactionID    string           `xml:"TransactionID"`
				GlobalE2nodeID   GlobalE2NodeId   `xml:"GlobalE2node-ID"`
				RANfunctionsList RANfunctionsList `xml:"RANfunctions-List"`
				E2NodeConfigList E2NodeConfigList `xml:"E2nodeComponentConfigAddition-List"`
			} `xml:"value"`
		} `xml:"E2setupRequestIEs"`
	} `xml:"protocolIEs"`
}

type E2SetupRequestMessage struct {
	XMLName xml.Name `xml:"E2SetupRequestMessage"`
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
				Text           string         `xml:",chardata"`
				E2setupRequest E2SetupRequest `xml:"E2setupRequest"`
			} `xml:"value"`
		} `xml:"initiatingMessage"`
	} `xml:"E2AP-PDU"`
}

type RanFunctionItem struct {
	Text                  string `xml:",chardata"`
	RanFunctionID         uint32 `xml:"ranFunctionID"`
	RanFunctionDefinition string `xml:"ranFunctionDefinition"`
	RanFunctionRevision   uint32 `xml:"ranFunctionRevision"`
	RanFunctionOID        string `xml:"ranFunctionOID"`
}

type RANfunctionsList struct {
	Text                      string `xml:",chardata"`
	ProtocolIESingleContainer []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"id"`
		Criticality struct {
			Text   string `xml:",chardata"`
			Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
			Text            string          `xml:",chardata"`
			RANfunctionItem RanFunctionItem `xml:"RANfunction-Item"`
		} `xml:"value"`
	} `xml:"ProtocolIE-SingleContainer"`
}

type E2NodeConfigList struct {
	Text                      string `xml:",chardata"`
	ProtocolIESingleContainer []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"id"`
		Criticality struct {
			Text   string `xml:",chardata"`
			Reject string `xml:"reject"`
		} `xml:"criticality"`
		Value struct {
			Text                     string                   `xml:",chardata"`
			E2nodeConfigAdditionItem E2NodeConfigAdditionItem `xml:"E2nodeComponentConfigAddition-Item"`
		} `xml:"value"`
	} `xml:"ProtocolIE-SingleContainer"`
}

type E2NodeComponentType struct {
	Text string    `xml:",chardata"`
	NG   *struct{} `xml:"ng"`
	XN   *struct{} `xml:"xn"`
	E1   *struct{} `xml:"e1"`
	F1   *struct{} `xml:"f1"`
	W1   *struct{} `xml:"w1"`
	S1   *struct{} `xml:"s1"`
	X2   *struct{} `xml:"x2"`
}

type E2NodeConfigAdditionItem struct {
	Text                string              `xml:",chardata"`
	E2nodeComponentType E2NodeComponentType `xml:"e2nodeComponentInterfaceType"`
	E2nodeComponentID   E2NodeComponentId   `xml:"e2nodeComponentID"`
	E2nodeConfiguration E2NodeConfigValue   `xml:"e2nodeComponentConfiguration"`
}

type E2NodeConfigValue struct {
	Text               string `xml:",chardata"`
	E2NodeRequestPart  []byte `xml:"e2nodeComponentRequestPart"`
	E2NodeResponsePart []byte `xml:"e2nodeComponentResponsePart"`
}

type E2NodeComponentId struct {
	Text           string `xml:",chardata"`
	E2NodeIFTypeNG E2NodeIFTypeNG
	E2NodeIFTypeXN E2NodeIFTypeXN
	E2NodeIFTypeE1 E2NodeIFTypeE1
	E2NodeIFTypeF1 E2NodeIFTypeF1
	E2NodeIFTypeW1 E2NodeIFTypeW1
	E2NodeIFTypeS1 E2NodeIFTypeS1
	E2NodeIFTypeX2 E2NodeIFTypeX2
}

type E2NodeIFTypeNG struct {
	XMLName xml.Name `xml:"e2nodeComponentInterfaceTypeNG"`
	Text    string   `xml:",chardata"`
	AMFName string   `xml:"amf-name"`
}

type E2NodeIFTypeXN struct {
	XMLName           xml.Name `xml:"e2nodeComponentInterfaceTypeXn"`
	Text              string   `xml:",chardata"`
	GlobalNGRANNodeID struct {
		Text          string   `xml:",chardata"`
		GlobalgNBID   *GNB     `xml:"gNB,omitempty"`
		GlobalngeNBID *NgeNBID `xml:"ng-eNB,omitempty"`
	} `xml:"global-NG-RAN-Node-ID"`
}

type GNB struct {
	Text   string `xml:",chardata"`
	PLMNID string `xml:"plmn-id"`
	GnbID  GnbID  `xml:"gnb-id"`
}
type GnbID struct {
	Text  string `xml:",chardata"`
	GnbID string `xml:"gnb-ID"`
}

type NgeNBID struct {
	Text   string    `xml:",chardata"`
	PLMNID string    `xml:"plmn-id"`
	EnbID  *EnbID_Xn `xml:"enb-id"`
}

type EnbID_Xn struct {
	Text            string `xml:",chardata"`
	EnbIdMacro      string `xml:"enb-ID-macro,omitempty"`
	EnbIdShortMacro string `xml:"enb-ID-shortmacro,omitempty"`
	EnbIdLongMacro  string `xml:"enb-ID-longmacro,omitempty"`
}

type E2NodeIFTypeE1 struct {
	XMLName   xml.Name `xml:"e2nodeComponentInterfaceTypeE1"`
	Text      string   `xml:",chardata"`
	GNBCUCPID int64    `xml:"gNB-CU-CP-ID"`
}

type E2NodeIFTypeF1 struct {
	XMLName xml.Name `xml:"e2nodeComponentInterfaceTypeF1"`
	Text    string   `xml:",chardata"`
	GNBDUID int64    `xml:"gNB-DU-ID"`
}

type E2NodeIFTypeW1 struct {
	XMLName   xml.Name `xml:"e2nodeComponentInterfaceTypeW1"`
	Text      string   `xml:",chardata"`
	NGENBDUID int64    `xml:"ng-eNB-DU-ID"`
}

type E2NodeIFTypeS1 struct {
	XMLName xml.Name `xml:"e2nodeComponentInterfaceTypeS1"`
	Text    string   `xml:",chardata"`
	MMENAME string   `xml:"mme-name"`
}
type E2NodeIFTypeX2 struct {
	XMLName       xml.Name       `xml:"e2nodeComponentInterfaceTypeX2"`
	Text          string         `xml:",chardata"`
	GlobalENBID   *GlobalENBID   `xml:"global-eNB-ID,omitempty"`
	GlobalenGNBID *GlobalenGNBID `xml:"global-en-gNB-ID,omitempty"`
}
type GlobalENBID struct {
	Text         string    `xml:",chardata"`
	PLMNIdentity string    `xml:"pLMN-Identity"`
	ENBID        *ENBID_X2 `xml:"eNB-ID"`
}
type ENBID_X2 struct {
	Text            string `xml:",chardata"`
	MacroEnbId      string `xml:"macro-eNB-ID,omitempty"`
	HomeEnbId       string `xml:"home-eNB-ID,omitempty"`
	ShortMacroEnbId string `xml:"short-Macro-eNB-ID,omitempty"`
	LongMacroEnbId  string `xml:"long-Macro-eNB-ID,omitempty"`
}
type GlobalenGNBID struct {
	Text         string `xml:",chardata"`
	PLMNIdentity string `xml:"pLMN-Identity"`
	GNBID        GNBID  `xml:"gNB-ID"`
}

type GNBID struct {
	Text  string `xml:",chardata"`
	GNBID string `xml:"gNB-ID"`
}

func (m *E2SetupRequestMessage) ExtractRanFunctionsList() []*entities.RanFunction {
	// TODO: verify e2SetupRequestIEs structure with Adi
	e2SetupRequestIes := m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs

	var ranFuntionsList RANfunctionsList
	var isPopulated bool

	for _, v := range e2SetupRequestIes {
		if v.ID == RanFunctionsAddedID {
			ranFuntionsList = v.Value.RANfunctionsList
			isPopulated = true
			break
		}
	}

	if !isPopulated {
		return nil
	}

	ranFunctionsListContainer := ranFuntionsList.ProtocolIESingleContainer
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

func (m *E2SetupRequestMessage) ExtractE2NodeConfigList() []*entities.E2NodeComponentConfig {
	e2SetupRequestIes := m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs
	numOfIes := len(e2SetupRequestIes)
	var e2NodeConfigListContainer E2NodeConfigList
	var isPopulated bool

	for ieCount := 0; ieCount < numOfIes; ieCount++ {

		if e2SetupRequestIes[ieCount].ID == E2nodeConfigAdditionID {
			e2NodeConfigListContainer = e2SetupRequestIes[ieCount].Value.E2NodeConfigList
			isPopulated = true
			break
		}
	}

	if !isPopulated {
		return nil
	}

	e2nodeComponentConfigs := make([]*entities.E2NodeComponentConfig, len(e2NodeConfigListContainer.ProtocolIESingleContainer))
	for i := 0; i < len(e2nodeComponentConfigs); i++ {
		e2NodeConfigItem := e2NodeConfigListContainer.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem

		if e2NodeConfigItem.E2nodeComponentType.NG != nil {
			e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
				E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_ng,
				E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
				E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
				E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeNG{
					E2NodeComponentInterfaceTypeNG: &entities.E2NodeComponentInterfaceNG{
						AmfName: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeNG.AMFName,
					},
				},
			}
		} else if e2NodeConfigItem.E2nodeComponentType.E1 != nil {
			e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
				E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_e1,
				E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
				E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
				E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeE1{
					E2NodeComponentInterfaceTypeE1: &entities.E2NodeComponentInterfaceE1{
						GNBCuCpId: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeE1.GNBCUCPID,
					},
				},
			}
		} else if e2NodeConfigItem.E2nodeComponentType.F1 != nil {
			e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
				E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_f1,
				E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
				E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
				E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeF1{
					E2NodeComponentInterfaceTypeF1: &entities.E2NodeComponentInterfaceF1{
						GNBDuId: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeF1.GNBDUID,
					},
				},
			}
		} else if e2NodeConfigItem.E2nodeComponentType.W1 != nil {
			e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
				E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_w1,
				E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
				E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
				E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeW1{
					E2NodeComponentInterfaceTypeW1: &entities.E2NodeComponentInterfaceW1{
						NgenbDuId: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeW1.NGENBDUID,
					},
				},
			}
		} else if e2NodeConfigItem.E2nodeComponentType.S1 != nil {
			e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
				E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_s1,
				E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
				E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
				E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeS1{
					E2NodeComponentInterfaceTypeS1: &entities.E2NodeComponentInterfaceS1{
						MmeName: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeS1.MMENAME,
					},
				},
			}
		} else if e2NodeConfigItem.E2nodeComponentType.XN != nil {
			if gnbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID; gnbid != nil {
				e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_xn,
					E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
					E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeXn{
						E2NodeComponentInterfaceTypeXn: &entities.E2NodeComponentInterfaceXn{
							GlobalNgRanNodeId: &entities.E2NodeComponentInterfaceXn_GlobalGnbId{
								GlobalGnbId: &entities.GlobalGNBID{
									PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID.PLMNID,
									GnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID.GnbID.GnbID,
									GnbType:      1,
								},
							},
						},
					},
				}

			} else if ngenb := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID; ngenb != nil {
				if ngenbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdMacro; ngenbid != "" && len(ngenbid) == 20 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_xn,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeXn{
							E2NodeComponentInterfaceTypeXn: &entities.E2NodeComponentInterfaceXn{
								GlobalNgRanNodeId: &entities.E2NodeComponentInterfaceXn_GlobalNgenbId{
									GlobalNgenbId: &entities.GlobalNGENBID{
										PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID,
										EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdMacro,
										EnbType:      1,
									},
								},
							},
						},
					}
				} else if ngenbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdShortMacro; ngenbid != "" && len(ngenbid) == 18 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_xn,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeXn{
							E2NodeComponentInterfaceTypeXn: &entities.E2NodeComponentInterfaceXn{
								GlobalNgRanNodeId: &entities.E2NodeComponentInterfaceXn_GlobalNgenbId{
									GlobalNgenbId: &entities.GlobalNGENBID{
										PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID,
										EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdShortMacro,
										EnbType:      3,
									},
								},
							},
						},
					}
				} else if ngenbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdLongMacro; ngenbid != "" && len(ngenbid) == 21 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_xn,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeXn{
							E2NodeComponentInterfaceTypeXn: &entities.E2NodeComponentInterfaceXn{
								GlobalNgRanNodeId: &entities.E2NodeComponentInterfaceXn_GlobalNgenbId{
									GlobalNgenbId: &entities.GlobalNGENBID{
										PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID,
										EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdLongMacro,
										EnbType:      4,
									},
								},
							},
						},
					}
				}
			} else {
				//not valid
			}
		} else if e2NodeConfigItem.E2nodeComponentType.X2 != nil {
			if gnbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID; gnbid != nil {
				e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
					E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
					E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
					E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
					E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{
						E2NodeComponentInterfaceTypeX2: &entities.E2NodeComponentInterfaceX2{
							GlobalEngnbId: &entities.GlobalENGNBID{
								PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID.PLMNIdentity,
								GnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID.GNBID.GNBID,
								GnbType:      1,
							},
						},
					},
				}
			} else if enb := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID; enb != nil {
				if enbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.MacroEnbId; enbid != "" && len(enbid) == 20 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{
							E2NodeComponentInterfaceTypeX2: &entities.E2NodeComponentInterfaceX2{
								GlobalEnbId: &entities.GlobalENBID{
									PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity,
									EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.MacroEnbId,
									EnbType:      1,
								},
							},
						},
					}
				} else if enbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.HomeEnbId; enbid != "" && len(enbid) == 28 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{
							E2NodeComponentInterfaceTypeX2: &entities.E2NodeComponentInterfaceX2{
								GlobalEnbId: &entities.GlobalENBID{
									PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity,
									EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.HomeEnbId,
									EnbType:      2,
								},
							},
						},
					}
				} else if enbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.ShortMacroEnbId; enbid != "" && len(enbid) == 18 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{
							E2NodeComponentInterfaceTypeX2: &entities.E2NodeComponentInterfaceX2{
								GlobalEnbId: &entities.GlobalENBID{
									PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity,
									EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.ShortMacroEnbId,
									EnbType:      3,
								},
							},
						},
					}
				} else if enbid := e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.LongMacroEnbId; enbid != "" && len(enbid) == 21 {
					e2nodeComponentConfigs[i] = &entities.E2NodeComponentConfig{
						E2NodeComponentInterfaceType: entities.E2NodeComponentInterfaceType_x2,
						E2NodeComponentRequestPart:   e2NodeConfigItem.E2nodeConfiguration.E2NodeRequestPart,
						E2NodeComponentResponsePart:  e2NodeConfigItem.E2nodeConfiguration.E2NodeResponsePart,
						E2NodeComponentID: &entities.E2NodeComponentConfig_E2NodeComponentInterfaceTypeX2{
							E2NodeComponentInterfaceTypeX2: &entities.E2NodeComponentInterfaceX2{
								GlobalEnbId: &entities.GlobalENBID{
									PlmnIdentity: e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity,
									EnbId:        e2NodeConfigItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.LongMacroEnbId,
									EnbType:      4,
								},
							},
						},
					}
				}
			} else {
				//not valid
			}
		} //end of x2
	}
	return e2nodeComponentConfigs
}

func (m *E2SetupRequestMessage) getGlobalE2NodeId() GlobalE2NodeId {

	// TODO: Handle error case if GlobalE2NodeId not available
	e2SetupRequestIes := m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs
	numOfIes := len(e2SetupRequestIes)
	index := 1

	for ieCount := 0; ieCount < numOfIes; ieCount++ {
		if e2SetupRequestIes[ieCount].ID == GlobalE2nodeID {
			index = ieCount
			break
		}
	}

	return m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[index].Value.GlobalE2nodeID
}

func (m *E2SetupRequestMessage) GetPlmnId() string {
	globalE2NodeId := m.getGlobalE2NodeId()
	if id := globalE2NodeId.GNB.GlobalGNBID.PlmnID; id != "" {
		return m.trimSpaces(id)
	}
	if id := globalE2NodeId.EnGNB.GlobalGNBID.PlmnID; id != "" {
		return m.trimSpaces(id)
	}
	if id := globalE2NodeId.ENB.GlobalENBID.PlmnID; id != "" {
		return m.trimSpaces(id)
	}
	if id := globalE2NodeId.NgENB.GlobalNgENBID.PlmnID; id != "" {
		return m.trimSpaces(id)
	}
	return ""
}

func (m *E2SetupRequestMessage) getInnerEnbId(enbId EnbId) string {

	if id := enbId.HomeEnbId; id != "" {
		return id
	}

	if id := enbId.LongMacroEnbId; id != "" {
		return id
	}

	if id := enbId.MacroEnbId; id != "" {
		return id
	}

	if id := enbId.ShortMacroEnbId; id != "" {
		return id
	}

	return ""
}

func (m *E2SetupRequestMessage) getInnerNgEnbId(enbId NgEnbId) string {
	if id := enbId.EnbIdLongMacro; id != "" {
		return id
	}

	if id := enbId.EnbIdMacro; id != "" {
		return id
	}

	if id := enbId.EnbIdShortMacro; id != "" {
		return id
	}

	return ""
}

func (m *E2SetupRequestMessage) GetNbId() string {
	globalE2NodeId := m.getGlobalE2NodeId()

	if id := globalE2NodeId.GNB.GlobalGNBID.GnbID.GnbID; id != "" {
		return m.trimSpaces(id)
	}

	if id := globalE2NodeId.EnGNB.GlobalGNBID.GnbID.GnbID; id != "" {
		return m.trimSpaces(id)
	}

	if id := m.getInnerEnbId(globalE2NodeId.ENB.GlobalENBID.EnbID); id != "" {
		return m.trimSpaces(id)
	}

	if id := m.getInnerNgEnbId(globalE2NodeId.NgENB.GlobalNgENBID.EnbID); id != "" {
		return m.trimSpaces(id)
	}

	return ""
}

func (m *E2SetupRequestMessage) trimSpaces(str string) string {
	return strings.NewReplacer(" ", "", "\n", "").Replace(str)
}
