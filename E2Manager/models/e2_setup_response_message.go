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
)

const (
	CauseID                       = "1"
	GlobalE2nodeID                = "3"
	GlobalRicID                   = "4"
	RanFunctionIDItemID           = "6"
	RanFunctionsAcceptedID        = "9"
	RanFunctionsAddedID           = "10"
	TimeToWaitID                  = "31"
	TransactionID                 = "49"
	E2nodeConfigAdditionID        = "50"
	E2nodeConfigAdditionItemID    = "51"
	E2nodeConfigAdditionAckID     = "52"
	E2nodeConfigAdditionAckItemID = "53"
)

type TimeToWait = int

var TimeToWaitEnum = struct {
	V60s TimeToWait
	V20s TimeToWait
	V10s TimeToWait
	V5s  TimeToWait
	V2s  TimeToWait
	V1s  TimeToWait
}{60, 20, 10, 5, 2, 1}

var timeToWaitMap = map[TimeToWait]interface{}{
	TimeToWaitEnum.V60s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V60s    string   `xml:"v60s"`
	}{},
	TimeToWaitEnum.V20s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V20s    string   `xml:"v20s"`
	}{},
	TimeToWaitEnum.V10s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V10s    string   `xml:"v10s"`
	}{},
	TimeToWaitEnum.V5s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V5s     string   `xml:"v5s"`
	}{},
	TimeToWaitEnum.V2s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V2s     string   `xml:"v2s"`
	}{},
	TimeToWaitEnum.V1s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V1s     string   `xml:"v1s"`
	}{},
}

type ConfigStatus = int

var ConfigStatusEnum = struct {
	Success ConfigStatus
	Failure ConfigStatus
}{0, 1}

var configStatusMap = map[ConfigStatus]interface{}{
	ConfigStatusEnum.Success: struct {
		XMLName xml.Name `xml:"updateOutcome"`
		Text    string   `xml:",chardata"`
		Success string   `xml:"success"`
	}{},
	ConfigStatusEnum.Failure: struct {
		XMLName xml.Name `xml:"updateOutcome"`
		Text    string   `xml:",chardata"`
		Failure string   `xml:"failure"`
	}{},
}

func NewE2SetupSuccessResponseMessage(plmnId string, ricId string, request *E2SetupRequestMessage) E2SetupResponseMessage {
	outcome := SuccessfulOutcome{}
	outcome.ProcedureCode = "1"

	e2SetupRequestIes := request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs
	numOfIes := len(e2SetupRequestIes)

	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs = make([]E2setupResponseIEs, numOfIes)

	for ieCount := 0; ieCount < numOfIes; ieCount++ {
		switch e2SetupRequestIes[ieCount].ID {
		case TransactionID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = TransactionID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = TransID{
				TransactionID: request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[ieCount].Value.TransactionID,
			}

		case GlobalE2nodeID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = GlobalRicID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = GlobalRICID{GlobalRICID: struct {
				Text         string `xml:",chardata"`
				PLMNIdentity string `xml:"pLMN-Identity"`
				RicID        string `xml:"ric-ID"`
			}{PLMNIdentity: plmnId, RicID: ricId}}

		case RanFunctionsAddedID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = RanFunctionsAcceptedID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = RANfunctionsIDList{RANfunctionsIDList: struct {
				Text                      string                      `xml:",chardata"`
				ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
			}{ProtocolIESingleContainer: extractRanFunctionsIDList(request, ieCount)}}

		case E2nodeConfigAdditionID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = E2nodeConfigAdditionAckID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = E2NodeConfigUpdateAckList{E2NodeConfigUpdateAckList: struct {
				Text                        string                        `xml:",chardata"`
				E2NodeConfigSingleContainer []E2NodeConfigSingleContainer `xml:"ProtocolIE-SingleContainer"`
			}{E2NodeConfigSingleContainer: extractE2NodeConfigUpdateList(request, ieCount, 0)}}
		}
	}

	return E2SetupResponseMessage{E2APPDU: E2APPDU{Outcome: outcome}}
}

func NewE2SetupFailureResponseMessage(timeToWait TimeToWait, cause Cause, request *E2SetupRequestMessage) E2SetupResponseMessage {
	outcome := UnsuccessfulOutcome{}

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs = make([]E2setupFailureIEs, 3)
	outcome.ProcedureCode = "1"

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].ID = TransactionID
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].Value.Value = TransFailID{
		ID: request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[0].Value.TransactionID,
	}

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].ID = CauseID
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].Value.Value = cause

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[2].ID = TimeToWaitID
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[2].Value.Value = timeToWaitMap[timeToWait]

	return E2SetupResponseMessage{E2APPDU: E2APPDU{Outcome: outcome}}
}

type E2SetupResponseMessage struct {
	XMLName xml.Name `xml:"E2SetupSuccessResponseMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU E2APPDU
}

type E2APPDU struct {
	XMLName xml.Name `xml:"E2AP-PDU"`
	Text    string   `xml:",chardata"`
	Outcome interface{}
}

type SuccessfulOutcome struct {
	XMLName       xml.Name `xml:"successfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text            string `xml:",chardata"`
		E2setupResponse struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text               string               `xml:",chardata"`
				E2setupResponseIEs []E2setupResponseIEs `xml:"E2setupResponseIEs"`
			} `xml:"protocolIEs"`
		} `xml:"E2setupResponse"`
	} `xml:"value"`
}

type E2setupResponseIEs struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value interface{} `xml:"value"`
}

type GlobalRICID struct {
	Text        string `xml:",chardata"`
	GlobalRICID struct {
		Text         string `xml:",chardata"`
		PLMNIdentity string `xml:"pLMN-Identity"`
		RicID        string `xml:"ric-ID"`
	} `xml:"GlobalRIC-ID"`
}

type RANfunctionsIDList struct {
	Text               string `xml:",chardata"`
	RANfunctionsIDList struct {
		Text                      string                      `xml:",chardata"`
		ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"RANfunctionsID-List"`
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
			RanFunctionID       uint32 `xml:"ranFunctionID"`
			RanFunctionRevision uint32 `xml:"ranFunctionRevision"`
		} `xml:"RANfunctionID-Item"`
	} `xml:"value"`
}

type TransID struct {
	Text          string `xml:",chardata"`
	TransactionID string `xml:"TransactionID"`
}

type TransFailID struct {
	XMLName xml.Name `xml:"TransactionID"`
	ID      string   `xml:",chardata"`
}

type E2NodeConfigUpdateAckList struct {
	Text                      string `xml:",chardata"`
	E2NodeConfigUpdateAckList struct {
		Text                        string                        `xml:",chardata"`
		E2NodeConfigSingleContainer []E2NodeConfigSingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"E2nodeComponentConfigAdditionAck-List"`
}

type E2NodeConfigSingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text                      string `xml:",chardata"`
		E2NodeConfigUpdateAckItem struct {
			Text                  string              `xml:",chardata"`
			E2nodeComponentType   E2NodeComponentType `xml:"e2nodeComponentInterfaceType"`
			E2nodeComponentID     E2NodeComponentIDResp
			E2nodeConfigUpdateAck E2nodeConfigUpdateAckResp
		} `xml:"E2nodeComponentConfigAdditionAck-Item"`
	} `xml:"value"`
}

type E2NodeComponentIDResp struct {
	XMLName xml.Name `xml:"e2nodeComponentID"`
	Value   interface{}
}

type E2nodeConfigUpdateAckResp struct {
	XMLName xml.Name `xml:"e2nodeComponentConfigurationAck"`
	Value   interface{}
}

type UnsuccessfulOutcome struct {
	XMLName       xml.Name `xml:"unsuccessfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text           string `xml:",chardata"`
		E2setupFailure struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text              string              `xml:",chardata"`
				E2setupFailureIEs []E2setupFailureIEs `xml:"E2setupFailureIEs"`
			} `xml:"protocolIEs"`
		} `xml:"E2setupFailure"`
	} `xml:"value"`
}

type E2setupFailureIEs struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text  string `xml:",chardata"`
		Value interface{}
	} `xml:"value"`
}

func extractRanFunctionsIDList(request *E2SetupRequestMessage, index int) []ProtocolIESingleContainer {
	list := &request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[index].Value.RANfunctionsList
	ids := make([]ProtocolIESingleContainer, len(list.ProtocolIESingleContainer))
	for i := 0; i < len(ids); i++ {
		ids[i] = convertToRANfunctionID(list, i)
	}
	return ids
}

func convertToRANfunctionID(list *RANfunctionsList, i int) ProtocolIESingleContainer {
	id := ProtocolIESingleContainer{}

	id.ID = RanFunctionIDItemID
	id.Value.RANfunctionIDItem.RanFunctionID = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionID
	id.Value.RANfunctionIDItem.RanFunctionRevision = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionRevision

	return id
}

func extractE2NodeConfigUpdateList(request *E2SetupRequestMessage, index int, outcome int) []E2NodeConfigSingleContainer {
	list := &request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[index].Value.E2NodeConfigList
	ids := make([]E2NodeConfigSingleContainer, len(list.ProtocolIESingleContainer))
	for i := 0; i < len(ids); i++ {
		ids[i] = convertToE2NodeConfig(list, i, outcome)
	}
	return ids
}

func convertToE2NodeConfig(list *E2NodeConfigList, i int, outcome int) E2NodeConfigSingleContainer {
	id := E2NodeConfigSingleContainer{}
	id.ID = E2nodeConfigAdditionAckItemID

	id.Value.E2NodeConfigUpdateAckItem.E2nodeConfigUpdateAck.Value = configStatusMap[outcome]

	id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentType = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType

	if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.NG != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeNG{
			AMFName: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeNG.AMFName,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.XN != nil {
		ifXn := E2NodeIFTypeXN{}
		if gnbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID; gnbid != nil {
			ifXn.GlobalNGRANNodeID.GlobalgNBID = &GNB{PLMNID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID.PLMNID, GnbID: GnbID{GnbID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalgNBID.GnbID.GnbID}}
		} else if ngenb := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID; ngenb != nil {
			if ngenbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdMacro; ngenbid != "" && len(ngenbid) == 20 {
				ifXn.GlobalNGRANNodeID.GlobalngeNBID = &NgeNBID{PLMNID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID, EnbID: &EnbID_Xn{EnbIdMacro: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdMacro}}
			} else if ngenbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdShortMacro; ngenbid != "" && len(ngenbid) == 18 {
				ifXn.GlobalNGRANNodeID.GlobalngeNBID = &NgeNBID{PLMNID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID, EnbID: &EnbID_Xn{EnbIdShortMacro: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdShortMacro}}
			} else if ngenbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdLongMacro; ngenbid != "" && len(ngenbid) == 21 {
				ifXn.GlobalNGRANNodeID.GlobalngeNBID = &NgeNBID{PLMNID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.PLMNID, EnbID: &EnbID_Xn{EnbIdLongMacro: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNGRANNodeID.GlobalngeNBID.EnbID.EnbIdLongMacro}}
			}
		} else {
			//not valid
		}
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = ifXn
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.E1 != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeE1{
			GNBCUCPID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeE1.GNBCUCPID,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.F1 != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeF1{
			GNBDUID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeF1.GNBDUID,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.W1 != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeW1{
			NGENBDUID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeW1.NGENBDUID,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.S1 != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeS1{
			MMENAME: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeS1.MMENAME,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.X2 != nil {
		ifX2 := E2NodeIFTypeX2{}
		if gnbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID; gnbid != nil {
			ifX2.GlobalenGNBID = &GlobalenGNBID{PLMNIdentity: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID.PLMNIdentity, GNBID: GNBID{GNBID: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalenGNBID.GNBID.GNBID}}
		} else if enb := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID; enb != nil {
			if enbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.MacroEnbId; enbid != "" && len(enbid) == 20 {
				ifX2.GlobalENBID = &GlobalENBID{PLMNIdentity: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity, ENBID: &ENBID_X2{MacroEnbId: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.MacroEnbId}}
			} else if enbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.HomeEnbId; enbid != "" && len(enbid) == 28 {
				ifX2.GlobalENBID = &GlobalENBID{PLMNIdentity: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity, ENBID: &ENBID_X2{HomeEnbId: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.HomeEnbId}}
			} else if enbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.ShortMacroEnbId; enbid != "" && len(enbid) == 18 {
				ifX2.GlobalENBID = &GlobalENBID{PLMNIdentity: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity, ENBID: &ENBID_X2{ShortMacroEnbId: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.ShortMacroEnbId}}
			} else if enbid := list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.LongMacroEnbId; enbid != "" && len(enbid) == 21 {
				ifX2.GlobalENBID = &GlobalENBID{PLMNIdentity: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity, ENBID: &ENBID_X2{LongMacroEnbId: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.LongMacroEnbId}}
			}
		} else {
			//not valid
		}
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = ifX2
	}

	return id
}
