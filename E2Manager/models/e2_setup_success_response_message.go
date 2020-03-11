package models
import (
	"encoding/xml"
)

type E2SetupSuccessResponseMessage struct {
	XMLName xml.Name `xml:"E2SetupSuccessResponseMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU struct {
		Text              string `xml:",chardata"`
		SuccessfulOutcome struct {
			Text          string `xml:",chardata"`
			ProcedureCode string `xml:"procedureCode"`
			Criticality   struct {
				Text   string `xml:",chardata"`
				Reject string `xml:"reject"`
			} `xml:"criticality"`
			Value struct {
				Text            string `xml:",chardata"`
				E2setupResponse struct {
					Text        string `xml:",chardata"`
					ProtocolIEs struct {
						Text               string `xml:",chardata"`
						E2setupResponseIEs struct {
							Text        string `xml:",chardata"`
							ID          string `xml:"id"`
							Criticality struct {
								Text   string `xml:",chardata"`
								Reject string `xml:"reject"`
							} `xml:"criticality"`
							Value struct {
								Text        string `xml:",chardata"`
								GlobalRICID struct {
									Text         string `xml:",chardata"`
									PLMNIdentity string `xml:"pLMN-Identity"`
									RicID        string `xml:"ric-ID"`
								} `xml:"GlobalRIC-ID"`
							} `xml:"value"`
						} `xml:"E2setupResponseIEs"`
					} `xml:"protocolIEs"`
				} `xml:"E2setupResponse"`
			} `xml:"value"`
		} `xml:"successfulOutcome"`
	} `xml:"E2AP-PDU"`
}


func (m *E2SetupSuccessResponseMessage) SetPlmnId(plmnId string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs.Value.GlobalRICID.PLMNIdentity = plmnId
}

func (m *E2SetupSuccessResponseMessage) SetNbId(ricID string){
	m.E2APPDU.SuccessfulOutcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs.Value.GlobalRICID.RicID = ricID
}