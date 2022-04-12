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
	causeID                       = "1"
	globalE2nodeID                = "3"
	globalRICID                   = "4"
	ranFunctionIDItemID           = "6"
	ranFunctionsAcceptedID        = "9"
	ranFunctionsAddedID           = "10"
	timeToWaitID                  = "31"
	transactionID                 = "49"
	e2nodeConfigAdditionID        = "50"
	e2nodeConfigAdditionItemID    = "51"
	e2nodeConfigAdditionAckID     = "52"
	e2nodeConfigAdditionAckItemID = "53"
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
		case transactionID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = transactionID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = TransID{
				TransactionID: request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[ieCount].Value.TransactionID,
			}

		case globalE2nodeID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = globalRICID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = GlobalRICID{GlobalRICID: struct {
				Text         string `xml:",chardata"`
				PLMNIdentity string `xml:"pLMN-Identity"`
				RicID        string `xml:"ric-ID"`
			}{PLMNIdentity: plmnId, RicID: ricId}}

		case ranFunctionsAddedID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = ranFunctionsAcceptedID
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].Value = RANfunctionsIDList{RANfunctionsIDList: struct {
				Text                      string                      `xml:",chardata"`
				ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
			}{ProtocolIESingleContainer: extractRanFunctionsIDList(request, ieCount)}}

		case e2nodeConfigAdditionID:
			outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[ieCount].ID = e2nodeConfigAdditionAckID
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

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].ID = transactionID
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].Value.Value = TransFailID{
		ID: request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[0].Value.TransactionID,
	}

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].ID = causeID
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].Value.Value = cause

	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[2].ID = timeToWaitID
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

	id.ID = ranFunctionIDItemID
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
	id.ID = e2nodeConfigAdditionAckItemID

	id.Value.E2NodeConfigUpdateAckItem.E2nodeConfigUpdateAck.Value = configStatusMap[outcome]

	id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentType = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType

	if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.NG != nil {
		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = E2NodeIFTypeNG{
			AMFName: list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeNG.AMFName,
		}
	} else if list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentType.XN != nil {
		ifXn := E2NodeIFTypeXN{}
		ifXn.GlobalNgENBID.GNB.PLMNID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.GNB.PLMNID
		ifXn.GlobalNgENBID.GNB.GnbID.GnbID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.GNB.GnbID.GnbID
		ifXn.GlobalNgENBID.NGENB.PLMNID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.NGENB.PLMNID
		ifXn.GlobalNgENBID.NGENB.GnbID.ENBIDMacro = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.NGENB.GnbID.ENBIDMacro
		ifXn.GlobalNgENBID.NGENB.GnbID.ENBIDShortMacro = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.NGENB.GnbID.ENBIDShortMacro
		ifXn.GlobalNgENBID.NGENB.GnbID.ENBIDLongMacro = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeXN.GlobalNgENBID.NGENB.GnbID.ENBIDLongMacro

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
		ifX2.GlobalENBID.PLMNIdentity = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.PLMNIdentity
		ifX2.GlobalENBID.ENBID.MacroENBID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.MacroENBID
		ifX2.GlobalENBID.ENBID.HomeENBID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.HomeENBID
		ifX2.GlobalENBID.ENBID.ShortMacroENBID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.ShortMacroENBID
		ifX2.GlobalENBID.ENBID.LongMacroENBID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalENBID.ENBID.LongMacroENBID

		ifX2.GlobalEnGNBID.PLMNIdentity = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalEnGNBID.PLMNIdentity
		ifX2.GlobalEnGNBID.GNBID.GNBID = list.ProtocolIESingleContainer[i].Value.E2nodeConfigAdditionItem.E2nodeComponentID.E2NodeIFTypeX2.GlobalEnGNBID.GNBID.GNBID

		id.Value.E2NodeConfigUpdateAckItem.E2nodeComponentID.Value = ifX2
	}

	return id
}
