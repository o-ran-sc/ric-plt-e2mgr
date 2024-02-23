package models

import (
	"encoding/xml"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"sync"

)


type ProcedureType int

const (
	E2SetupProcedureNotInitiated ProcedureType = iota
	E2SetupProcedureOngoing
	E2SetupProcedureCompleted
	E2SetupProcedureFailure
	RicServiceUpdateCompleted
	RicServiceUpdateFailure
)

var(
	ProcedureMap = make(map[string]ProcedureType)
	procedureMapMutex sync.RWMutex	
) 

func UpdateProcedureType(ranName string, newProcedureType ProcedureType) {
	procedureMapMutex.Lock()
	defer procedureMapMutex.Unlock()
	ProcedureMap[ranName] = newProcedureType
}

var ExistingRanFunctiuonsMap = make(map[string][]*entities.RanFunction)

type ErrorIndicationMessage struct {
	XMLName xml.Name                `xml:"ErrorIndicationMessage"`
	Text    string                  `xml:",chardata"`
	E2APPDU ErrorIndicationE2APPDU `xml:"E2AP-PDU"`
}
type ErrorIndicationE2APPDU struct {
	XMLName           xml.Name                          `xml:"E2AP-PDU"`
	Text              string                            `xml:",chardata"`
	InitiatingMessage ErrorIndicationInitiatingMessage `xml:"initiatingMessage"`
}
type ErrorIndicationInitiatingMessage struct {
	Text          string `xml:",chardata"`
	ProcedureCode string `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text             string `xml:",chardata"`
		ErrorIndication struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text                string                `xml:",chardata"`
				ErrorIndicationIEs []ErrorIndicationIEs `xml:"ErrorIndication-IEs"`
			} `xml:"protocolIEs"`
		} `xml:"ErrorIndication"`
	} `xml:"value"`
}
type ErrorIndicationIEs struct {
	Text        string `xml:",chardata"`
	ID          int    `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text             string `xml:",chardata"`
		TransactionID    string `xml:"TransactionID"`
		RICrequestID  struct {
			Text           string `xml:",chardata"`
			RicRequestorID int32 `xml:"ricRequestorID"`
			RicInstanceID  int32 `xml:"ricInstanceID"`
		} `xml:"RICrequestID"`
		RANfunctionID     int32 `xml:"RANfunctionID"`
		CriticalityDiagnostics struct {
			Text        string `xml:",chardata"`
			ProcedureCode string `xml:"procedureCode"`
			TriggeringMessage TriggeringMessage `xml:"triggeringMessage"`
		} `xml:"CriticalityDiagnostics"`
	}`xml:"value"`
}

type TriggeringMessage struct {
	Text string    `xml:",chardata"`
	InitiatingMessage  *struct{} `xml:"initiatingMessage"`
	SuccessfulOutcome *struct{} `xml:"successful-outcome"`
	UnsuccessfulOutcome  *struct{} `xml:"unsuccessful-outcome"`
}

