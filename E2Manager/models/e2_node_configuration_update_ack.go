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

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
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

func prepareAdditionAckList(e2nodeConfigs []entities.E2NodeComponentConfig) []AdditionListProtocolIESingleContainer {
	additionListAckSingle := []AdditionListProtocolIESingleContainer{}
	for _, v := range e2nodeConfigs {
		c := convertEntitiyToModelComponent(v)

		t := AdditionListProtocolIESingleContainer{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigAdditionAck_Item,
			Value: struct {
				Text                                 string             `xml:",chardata"`
				E2nodeComponentConfigAdditionAckItem ComponentAckDetail `xml:"E2nodeComponentConfigAdditionAck-Item"`
			}{
				E2nodeComponentConfigAdditionAckItem: *c,
			},
		}
		additionListAckSingle = append(additionListAckSingle, t)
	}
	return additionListAckSingle
}

func prepareUpdateAckList(e2nodeConfigs []entities.E2NodeComponentConfig) []UpdateProtocolIESingleContainer {
	updateListAckSingle := []UpdateProtocolIESingleContainer{}
	for _, v := range e2nodeConfigs {
		c := convertEntitiyToModelComponent(v)

		t := UpdateProtocolIESingleContainer{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigUpdateAck_Item,
			Value: struct {
				Text                               string             `xml:",chardata"`
				E2nodeComponentConfigUpdateAckItem ComponentAckDetail `xml:"E2nodeComponentConfigUpdateAck-Item"`
			}{
				E2nodeComponentConfigUpdateAckItem: *c,
			},
		}
		updateListAckSingle = append(updateListAckSingle, t)
	}
	return updateListAckSingle
}

func prepareRemovalAckList(e2nodeConfigs []entities.E2NodeComponentConfig) []RemovalProtocolIESingleContainer {
	removalListAckSingle := []RemovalProtocolIESingleContainer{}
	for _, v := range e2nodeConfigs {
		c := convertEntitiyToModelComponent(v)

		t := RemovalProtocolIESingleContainer{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigRemovalAck_Item,
			Value: struct {
				Text                                string             `xml:",chardata"`
				E2nodeComponentConfigRemovalAckItem ComponentAckDetail `xml:"E2nodeComponentConfigRemovalAck-Item"`
			}{
				E2nodeComponentConfigRemovalAckItem: *c,
			},
		}
		removalListAckSingle = append(removalListAckSingle, t)
	}
	return removalListAckSingle
}

func updateIDAndStatus(t *ComponentAckDetail, c entities.E2NodeComponentConfig, succss bool) {
	switch c.E2NodeComponentInterfaceType {
	case entities.E2NodeComponentInterfaceType_ng:
		t.E2nodeComponentID.Value = E2NodeIFTypeNG{
			AMFName: c.GetE2NodeComponentInterfaceTypeNG().AmfName,
		}
	case entities.E2NodeComponentInterfaceType_e1:
		t.E2nodeComponentID.Value = E2NodeIFTypeE1{
			GNBCUCPID: c.GetE2NodeComponentInterfaceTypeE1().GetGNBCuCpId(),
		}
	case entities.E2NodeComponentInterfaceType_s1:
		t.E2nodeComponentID.Value = E2NodeIFTypeS1{
			MMENAME: c.GetE2NodeComponentInterfaceTypeS1().GetMmeName(),
		}
	case entities.E2NodeComponentInterfaceType_f1:
		t.E2nodeComponentID.Value = E2NodeIFTypeF1{
			GNBDUID: c.GetE2NodeComponentInterfaceTypeF1().GetGNBDuId(),
		}
	case entities.E2NodeComponentInterfaceType_w1:
		t.E2nodeComponentID.Value = E2NodeIFTypeW1{
			NGENBDUID: c.GetE2NodeComponentInterfaceTypeW1().GetNgenbDuId(),
		}
	}

	if succss {
		t.E2nodeConfigUpdateAck = E2nodeConfigUpdateAckResp{
			Value: struct {
				XMLName xml.Name `xml:"updateOutcome"`
				Text    string   `xml:",chardata"`
				Success string   `xml:"success"`
			}{},
		}
	} else {
		t.E2nodeConfigUpdateAck = E2nodeConfigUpdateAckResp{
			Value: struct {
				XMLName xml.Name `xml:"updateOutcome"`
				Text    string   `xml:",chardata"`
				Success string   `xml:"failure"`
			}{},
		}
	}
}

func updateInterfaceType(t *ComponentAckDetail, c entities.E2NodeComponentConfig) {
	switch c.E2NodeComponentInterfaceType {
	case entities.E2NodeComponentInterfaceType_ng:
		t.E2nodeComponentInterfaceType = E2NodeComponentType{
			NG: &struct{}{},
		}
	case entities.E2NodeComponentInterfaceType_e1:
		t.E2nodeComponentInterfaceType = E2NodeComponentType{
			E1: &struct{}{},
		}
	case entities.E2NodeComponentInterfaceType_f1:
		t.E2nodeComponentInterfaceType = E2NodeComponentType{
			F1: &struct{}{},
		}
	case entities.E2NodeComponentInterfaceType_w1:
		t.E2nodeComponentInterfaceType = E2NodeComponentType{
			W1: &struct{}{},
		}
	case entities.E2NodeComponentInterfaceType_s1:
		t.E2nodeComponentInterfaceType = E2NodeComponentType{
			S1: &struct{}{},
		}
	}
}

func convertEntitiyToModelComponent(component entities.E2NodeComponentConfig) *ComponentAckDetail {
	componentAckDetail := &ComponentAckDetail{}
	updateInterfaceType(componentAckDetail, component)
	updateIDAndStatus(componentAckDetail, component, true)
	return componentAckDetail
}

func NewE2nodeConfigurationUpdateSuccessResponseMessage(e2nodeConfigupdateMessage *E2nodeConfigurationUpdateMessage) *E2nodeConfigurationUpdateAcknowledgeE2APPDU {
	successfulOutcome := E2nodeConfigurationUpdateAcknowledgeSuccessfulOutcome{
		ProcedureCode: ProcedureCode_id_E2nodeConfigurationUpdate,
	}

	e2nodeConfigurationUpdateAckIEs := []E2nodeConfigurationUpdateAcknowledgeIEs{}
	txIEs := E2nodeConfigurationUpdateAcknowledgeIEs{
		ID: ProtocolIE_ID_id_TransactionID,
		Value: E2nodeConfigurationUpdateAcknowledgeTransID{
			TransactionID: e2nodeConfigupdateMessage.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.TransactionID,
		},
	}

	e2nodeConfigurationUpdateAckIEs = append(e2nodeConfigurationUpdateAckIEs, txIEs)

	items := e2nodeConfigupdateMessage.ExtractConfigAdditionList()
	if len(items) > 0 {
		addtionListAckIEs := E2nodeConfigurationUpdateAcknowledgeIEs{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigAdditionAck,
			Value: E2nodeComponentConfigAdditionAckList{
				E2nodeComponentConfigAdditionAckList: struct {
					Text                      string                                  `xml:",chardata"`
					ProtocolIESingleContainer []AdditionListProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
				}{
					ProtocolIESingleContainer: prepareAdditionAckList(items),
				},
			},
		}
		e2nodeConfigurationUpdateAckIEs = append(e2nodeConfigurationUpdateAckIEs, addtionListAckIEs)
	}

	items = e2nodeConfigupdateMessage.ExtractConfigUpdateList()
	if len(items) > 0 {
		updateListAckIEs := E2nodeConfigurationUpdateAcknowledgeIEs{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigUpdateAck,
			Value: E2nodeComponentConfigUpdateAckList{
				E2nodeComponentConfigUpdateAckList: struct {
					Text                      string                            `xml:",chardata"`
					ProtocolIESingleContainer []UpdateProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
				}{
					ProtocolIESingleContainer: prepareUpdateAckList(items),
				},
			},
		}
		e2nodeConfigurationUpdateAckIEs = append(e2nodeConfigurationUpdateAckIEs, updateListAckIEs)
	}

	items = e2nodeConfigupdateMessage.ExtractConfigDeletionList()
	if len(items) > 0 {
		removalListAckIEs := E2nodeConfigurationUpdateAcknowledgeIEs{
			ID: ProtocolIE_ID_id_E2nodeComponentConfigRemovalAck,
			Value: E2nodeComponentConfigRemovalAckList{
				E2nodeComponentConfigRemovalAckList: struct {
					Text                      string                             `xml:",chardata"`
					ProtocolIESingleContainer []RemovalProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
				}{
					ProtocolIESingleContainer: prepareRemovalAckList(items),
				},
			},
		}
		e2nodeConfigurationUpdateAckIEs = append(e2nodeConfigurationUpdateAckIEs, removalListAckIEs)
	}

	successfulOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs = e2nodeConfigurationUpdateAckIEs
	response := &E2nodeConfigurationUpdateAcknowledgeE2APPDU{
		Outcome: successfulOutcome,
	}
	return response
}
