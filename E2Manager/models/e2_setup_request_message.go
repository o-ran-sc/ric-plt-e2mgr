//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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
	"github.com/golang/protobuf/ptypes/wrappers"
	"strings"
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
				GlobalE2nodeID   GlobalE2NodeId   `xml:"GlobalE2node-ID"`
				RANfunctionsList RANfunctionsList `xml:"RANfunctions-List"`
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
	Text                  string                `xml:",chardata"`
	RanFunctionID         uint32                `xml:"ranFunctionID"`
	RanFunctionDefinition RanFunctionDefinition `xml:"ranFunctionDefinition"`
	RanFunctionRevision   uint32                `xml:"ranFunctionRevision"`
}

type RanFunctionDefinition struct {
	Text                            string                          `xml:",chardata"`
	E2smGnbNrtRanFunctionDefinition E2smGnbNrtRanFunctionDefinition `xml:"E2SM-gNB-NRT-RANfunction-Definition"`
}

type RanFunctionName struct {
	Text                   string  `xml:",chardata"`
	RanFunctionShortName   string  `xml:"ranFunction-ShortName"`
	RanFunctionE2smOid     string  `xml:"ranFunction-E2SM-OID"`
	RanFunctionDescription string  `xml:"ranFunction-Description"`
	RanFunctionInstance    *uint32 `xml:"ranFunction-Instance"`
}

type RicEventTriggerStyleList struct {
	Text                      string `xml:",chardata"`
	RicEventTriggerStyleType  uint32 `xml:"ric-EventTriggerStyle-Type"`
	RicEventTriggerStyleName  string `xml:"ric-EventTriggerStyle-Name"`
	RicEventTriggerFormatType uint32 `xml:"ric-EventTriggerFormat-Type"`
}

type RanParameterDefItem struct {
	Text             string           `xml:",chardata"`
	RanParameterID   uint32           `xml:"ranParameter-ID"`
	RanParameterName string           `xml:"ranParameter-Name"`
	RanParameterType RanParameterType `xml:"ranParameter-Type"`
}

type RanParameterType struct {
	Text            string    `xml:",chardata"`
	Boolean         *struct{} `xml:"boolean,omitempty"`
	Integer         *struct{} `xml:"integer,omitempty"`
	Enumerated      *struct{} `xml:"enumerated,omitempty"`
	BitString       *struct{} `xml:"bit-string,omitempty"`
	OctetString     *struct{} `xml:"octet-string,omitempty"`
	PrintableString *struct{} `xml:"printable-string,omitempty"`
}

type RicReportStyleList struct {
	Text                         string `xml:",chardata"`
	RicReportStyleType           uint32 `xml:"ric-ReportStyle-Type"`
	RicReportStyleName           string `xml:"ric-ReportStyle-Name"`
	RicReportActionFormatType    uint32 `xml:"ric-ReportActionFormat-Type"`
	RicReportRanParameterDefList struct {
		Text                string                `xml:",chardata"`
		RanParameterDefItem []RanParameterDefItem `xml:"RANparameterDef-Item"`
	} `xml:"ric-ReportRanParameterDef-List"`
	RicIndicationHeaderFormatType  uint32 `xml:"ric-IndicationHeaderFormat-Type"`
	RicIndicationMessageFormatType uint32 `xml:"ric-IndicationMessageFormat-Type"`
}

type RicInsertStyleList struct {
	Text                         string `xml:",chardata"`
	RicInsertStyleType           uint32 `xml:"ric-InsertStyle-Type"`
	RicInsertStyleName           string `xml:"ric-InsertStyle-Name"`
	RicInsertActionFormatType    uint32 `xml:"ric-InsertActionFormat-Type"`
	RicInsertRanParameterDefList struct {
		Text                string                `xml:",chardata"`
		RanParameterDefItem []RanParameterDefItem `xml:"RANparameterDef-Item"`
	} `xml:"ric-InsertRanParameterDef-List"`
	RicIndicationHeaderFormatType  uint32 `xml:"ric-IndicationHeaderFormat-Type"`
	RicIndicationMessageFormatType uint32 `xml:"ric-IndicationMessageFormat-Type"`
	RicCallProcessIdFormatType     uint32 `xml:"ric-CallProcessIDFormat-Type"`
}

type RicControlStyleList struct {
	Text                        string `xml:",chardata"`
	RicControlStyleType         uint32 `xml:"ric-ControlStyle-Type"`
	RicControlStyleName         string `xml:"ric-ControlStyle-Name"`
	RicControlHeaderFormatType  uint32 `xml:"ric-ControlHeaderFormat-Type"`
	RicControlMessageFormatType uint32 `xml:"ric-ControlMessageFormat-Type"`
	RicCallProcessIdFormatType  uint32 `xml:"ric-CallProcessIDFormat-Type"`
}

type RicPolicyStyleList struct {
	Text                         string `xml:",chardata"`
	RicPolicyStyleType           uint32 `xml:"ric-PolicyStyle-Type"`
	RicPolicyStyleName           string `xml:"ric-PolicyStyle-Name"`
	RicPolicyActionFormatType    uint32 `xml:"ric-PolicyActionFormat-Type"`
	RicPolicyRanParameterDefList struct {
		Text                string                `xml:",chardata"`
		RanParameterDefItem []RanParameterDefItem `xml:"RANparameterDef-Item"`
	} `xml:"ric-PolicyRanParameterDef-List"`
}

type E2smGnbNrtRanFunctionDefinition struct {
	Text                     string          `xml:",chardata"`
	RanFunctionName          RanFunctionName `xml:"ranFunction-Name"`
	RicEventTriggerStyleList struct {
		Text                     string                     `xml:",chardata"`
		RicEventTriggerStyleList []RicEventTriggerStyleList `xml:"RIC-EventTriggerStyle-List"`
	} `xml:"ric-EventTriggerStyle-List"`
	RicReportStyleList struct {
		Text               string               `xml:",chardata"`
		RicReportStyleList []RicReportStyleList `xml:"RIC-ReportStyle-List"`
	} `xml:"ric-ReportStyle-List"`
	RicInsertStyleList struct {
		Text               string               `xml:",chardata"`
		RicInsertStyleList []RicInsertStyleList `xml:"RIC-InsertStyle-List"`
	} `xml:"ric-InsertStyle-List"`
	RicControlStyleList struct {
		Text                string                `xml:",chardata"`
		RicControlStyleList []RicControlStyleList `xml:"RIC-ControlStyle-List"`
	} `xml:"ric-ControlStyle-List"`
	RicPolicyStyleList struct {
		Text               string               `xml:",chardata"`
		RicPolicyStyleList []RicPolicyStyleList `xml:"RIC-PolicyStyle-List"`
	} `xml:"ric-PolicyStyle-List"`
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

func (m *E2SetupRequestMessage) ExtractRanFunctionsList() []*entities.RanFunction {
	// TODO: verify e2SetupRequestIEs structure with Adi
	e2SetupRequestIes := m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs
	if len(e2SetupRequestIes) < 2  {
		return nil
	}

	ranFunctionsListContainer := e2SetupRequestIes[1].Value.RANfunctionsList.ProtocolIESingleContainer
	funcs := make([]*entities.RanFunction, len(ranFunctionsListContainer))
	for i := 0; i < len(funcs); i++ {
		ranFunctionItem := ranFunctionsListContainer[i].Value.RANfunctionItem

		funcs[i] = &entities.RanFunction{
			RanFunctionId:         &wrappers.UInt32Value{Value: ranFunctionItem.RanFunctionID},
			RanFunctionDefinition: m.buildRanFunctionDefinitionProto(&ranFunctionItem.RanFunctionDefinition),
			RanFunctionRevision:   &wrappers.UInt32Value{Value: ranFunctionItem.RanFunctionRevision},
		}
	}
	return funcs
}

func (m *E2SetupRequestMessage) buildRanFunctionDefinitionProto(def *RanFunctionDefinition) *entities.RanFunctionDefinition {
	return &entities.RanFunctionDefinition{
		E2SmGnbNrtRanFunctionDefinition: &entities.E2SmGnbNrtRanFunctionDefinition{
			RanFunctionName:       buildRanFunctionNameProto(def),
			RicEventTriggerStyles: buildRicEventTriggerStylesProto(def),
			RicReportStyles:       buildRicReportStylesProto(def),
			RicInsertStyles:       buildRicInsertStylesProto(def),
			RicControlStyles:      buildRicControlStylesProto(def),
			RicPolicyStyles:       buildRicPolicyStylesProto(def),
		},
	}
}

func buildRanFunctionNameProto(def *RanFunctionDefinition) *entities.RanFunctionName {
	defRanFunctionName := def.E2smGnbNrtRanFunctionDefinition.RanFunctionName
	ranFunctionName := &entities.RanFunctionName{
		RanFunctionShortName:   &wrappers.StringValue{Value: defRanFunctionName.RanFunctionShortName},
		RanFunctionE2SmOid:     &wrappers.StringValue{Value: defRanFunctionName.RanFunctionE2smOid},
		RanFunctionDescription: &wrappers.StringValue{Value: defRanFunctionName.RanFunctionDescription},
	}

	if defRanFunctionName.RanFunctionInstance != nil {
		ranFunctionName.OptionalRanFunctionInstance = &entities.RanFunctionName_RanFunctionInstance{
			RanFunctionInstance: *defRanFunctionName.RanFunctionInstance,
		}
	}

	return ranFunctionName
}

func buildRicEventTriggerStylesProto(def *RanFunctionDefinition) []*entities.RicEventTriggerStyle {
	defRicEventTriggerStyleList := def.E2smGnbNrtRanFunctionDefinition.RicEventTriggerStyleList.RicEventTriggerStyleList
	ricEventTriggerStyles := make([]*entities.RicEventTriggerStyle, len(defRicEventTriggerStyleList))

	for i, v := range defRicEventTriggerStyleList {
		ricEventTriggerStyles[i] = &entities.RicEventTriggerStyle{
			RicEventTriggerStyleType:  &wrappers.UInt32Value{Value: v.RicEventTriggerStyleType},
			RicEventTriggerStyleName:  &wrappers.StringValue{Value: v.RicEventTriggerStyleName},
			RicEventTriggerFormatType: &wrappers.UInt32Value{Value: v.RicEventTriggerFormatType},
		}
	}

	return ricEventTriggerStyles
}

func buildRicReportStylesProto(def *RanFunctionDefinition) []*entities.RicReportStyle {
	defRicReportStyleList := def.E2smGnbNrtRanFunctionDefinition.RicReportStyleList.RicReportStyleList
	ricReportStyles := make([]*entities.RicReportStyle, len(defRicReportStyleList))

	for i, v := range defRicReportStyleList {
		ricReportStyles[i] = &entities.RicReportStyle{
			RicReportStyleType:             &wrappers.UInt32Value{Value: v.RicReportStyleType},
			RicReportStyleName:             &wrappers.StringValue{Value: v.RicReportStyleName},
			RicReportActionFormatType:      &wrappers.UInt32Value{Value: v.RicReportActionFormatType},
			RicReportRanParameterDefs:      buildRicReportRanParameterDefsProto(v),
			RicIndicationHeaderFormatType:  &wrappers.UInt32Value{Value: v.RicIndicationHeaderFormatType},
			RicIndicationMessageFormatType: &wrappers.UInt32Value{Value: v.RicIndicationMessageFormatType},
		}
	}

	return ricReportStyles
}

func buildRicReportRanParameterDefsProto(ricReportStyleList RicReportStyleList) []*entities.RanParameterDef {
	ricReportRanParameterDefList := ricReportStyleList.RicReportRanParameterDefList.RanParameterDefItem
	ranParameterDefs := make([]*entities.RanParameterDef, len(ricReportRanParameterDefList))

	for i, v := range ricReportRanParameterDefList {
		ranParameterDefs[i] = &entities.RanParameterDef{
			RanParameterId:   &wrappers.UInt32Value{Value: v.RanParameterID},
			RanParameterName: &wrappers.StringValue{Value: v.RanParameterName},
			RanParameterType: getRanParameterTypeEnumValue(v.RanParameterType),
		}
	}

	return ranParameterDefs
}

func getRanParameterTypeEnumValue(ranParameterType RanParameterType) entities.RanParameterType {
	if ranParameterType.Boolean != nil {
		return entities.RanParameterType_BOOLEAN
	}

	if ranParameterType.BitString != nil {
		return entities.RanParameterType_BIT_STRING
	}

	if ranParameterType.Enumerated != nil {
		return entities.RanParameterType_ENUMERATED
	}

	if ranParameterType.Integer != nil {
		return entities.RanParameterType_INTEGER
	}

	if ranParameterType.OctetString != nil {
		return entities.RanParameterType_OCTET_STRING
	}

	if ranParameterType.PrintableString != nil {
		return entities.RanParameterType_PRINTABLE_STRING
	}

	return entities.RanParameterType_UNKNOWN_RAN_PARAMETER_TYPE
}

func buildRicInsertStylesProto(def *RanFunctionDefinition) []*entities.RicInsertStyle {
	defRicInsertStyleList := def.E2smGnbNrtRanFunctionDefinition.RicInsertStyleList.RicInsertStyleList
	ricInsertStyles := make([]*entities.RicInsertStyle, len(defRicInsertStyleList))

	for i, v := range defRicInsertStyleList {
		ricInsertStyles[i] = &entities.RicInsertStyle{
			RicInsertStyleType:             &wrappers.UInt32Value{Value: v.RicInsertStyleType},
			RicInsertStyleName:             &wrappers.StringValue{Value: v.RicInsertStyleName},
			RicInsertActionFormatType:      &wrappers.UInt32Value{Value: v.RicInsertActionFormatType},
			RicInsertRanParameterDefs:      buildRicInsertRanParameterDefsProto(v),
			RicIndicationHeaderFormatType:  &wrappers.UInt32Value{Value: v.RicIndicationHeaderFormatType},
			RicIndicationMessageFormatType: &wrappers.UInt32Value{Value: v.RicIndicationMessageFormatType},
			RicCallProcessIdFormatType:     &wrappers.UInt32Value{Value: v.RicCallProcessIdFormatType},
		}
	}

	return ricInsertStyles
}

func buildRicInsertRanParameterDefsProto(ricInsertStyleList RicInsertStyleList) []*entities.RanParameterDef {
	ricInsertRanParameterDefList := ricInsertStyleList.RicInsertRanParameterDefList.RanParameterDefItem
	ranParameterDefs := make([]*entities.RanParameterDef, len(ricInsertRanParameterDefList))

	for i, v := range ricInsertRanParameterDefList {
		ranParameterDefs[i] = &entities.RanParameterDef{
			RanParameterId:   &wrappers.UInt32Value{Value: v.RanParameterID},
			RanParameterName: &wrappers.StringValue{Value: v.RanParameterName},
			RanParameterType: getRanParameterTypeEnumValue(v.RanParameterType),
		}
	}

	return ranParameterDefs
}

func buildRicControlStylesProto(def *RanFunctionDefinition) []*entities.RicControlStyle {
	defRicControlStyleList := def.E2smGnbNrtRanFunctionDefinition.RicControlStyleList.RicControlStyleList
	ricControlStyles := make([]*entities.RicControlStyle, len(defRicControlStyleList))

	for i, v := range defRicControlStyleList {
		ricControlStyles[i] = &entities.RicControlStyle{
			RicControlStyleType:         &wrappers.UInt32Value{Value: v.RicControlStyleType},
			RicControlStyleName:         &wrappers.StringValue{Value: v.RicControlStyleName},
			RicControlHeaderFormatType:  &wrappers.UInt32Value{Value: v.RicControlHeaderFormatType},
			RicControlMessageFormatType: &wrappers.UInt32Value{Value: v.RicControlMessageFormatType},
			RicCallProcessIdFormatType:  &wrappers.UInt32Value{Value: v.RicCallProcessIdFormatType},
		}
	}

	return ricControlStyles
}

func buildRicPolicyRanParameterDefsProto(ricPolicyStyleList RicPolicyStyleList) []*entities.RanParameterDef {
	ricPolicyRanParameterDefList := ricPolicyStyleList.RicPolicyRanParameterDefList.RanParameterDefItem
	ranParameterDefs := make([]*entities.RanParameterDef, len(ricPolicyRanParameterDefList))

	for i, v := range ricPolicyRanParameterDefList {
		ranParameterDefs[i] = &entities.RanParameterDef{
			RanParameterId:   &wrappers.UInt32Value{Value: v.RanParameterID},
			RanParameterName: &wrappers.StringValue{Value: v.RanParameterName},
			RanParameterType: getRanParameterTypeEnumValue(v.RanParameterType),
		}
	}

	return ranParameterDefs
}

func buildRicPolicyStylesProto(def *RanFunctionDefinition) []*entities.RicPolicyStyle {
	defRicPolicyStyleList := def.E2smGnbNrtRanFunctionDefinition.RicPolicyStyleList.RicPolicyStyleList
	ricPolicyStyles := make([]*entities.RicPolicyStyle, len(defRicPolicyStyleList))

	for i, v := range defRicPolicyStyleList {
		ricPolicyStyles[i] = &entities.RicPolicyStyle{
			RicPolicyStyleType:        &wrappers.UInt32Value{Value: v.RicPolicyStyleType},
			RicPolicyStyleName:        &wrappers.StringValue{Value: v.RicPolicyStyleName},
			RicPolicyActionFormatType: &wrappers.UInt32Value{Value: v.RicPolicyActionFormatType},
			RicPolicyRanParameterDefs: buildRicPolicyRanParameterDefsProto(v),
		}
	}

	return ricPolicyStyles
}

func (m *E2SetupRequestMessage) getGlobalE2NodeId() GlobalE2NodeId {
	return m.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[0].Value.GlobalE2nodeID
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
