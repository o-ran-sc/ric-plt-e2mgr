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

type E2ResetRequestMessage struct {
	E2ApPDU E2ApPDU `xml:"E2AP-PDU"`
}

type E2ApPDU struct {
	InitiatingMessage InitiatingMessageY `xml:"initiatingMessage"`
}

type InitiatingMessageY struct {
	ProcedureCode int64                        `xml:"procedureCode"`
	Criticality   InitiatingMessageCriticality `xml:"criticality"`
	Value         InitiatingMessageValue       `xml:"value"`
}

type InitiatingMessageCriticality struct {
	Reject string `xml:"reject"`
}

type InitiatingMessageValue struct {
	E2ResetRequest E2ResetRequest `xml:"ResetRequest"`
}

type E2ResetRequest struct {
	ProtocolIes ProtocolIes `xml:"protocolIEs"`
}

type ProtocolIes struct {
	ResetRequestIEs []ResetRequestIEs `xml:"ResetRequestIEs"`
}

type ResetRequestIEs struct {
	ID          int64                   `xml:"id"`
	Criticality ResetRequestCriticality `xml:"criticality"`
	Value       ResetRequestValue       `xml:"value"`
}

type ResetRequestCriticality struct {
	Ignore string `xml:"ignore"`
}

type ResetRequestValue struct {
	TransactionID *int64             `xml:"TransactionID"`
	Cause         *CauseResetRequest `xml:"Cause"`
}

type CauseResetRequest struct {
	Text       string     `xml:",chardata"`
	E2Node     E2Node     `xml:"e2Node"`
	RicRequest RicRequest `xml:"ricRequest"`
	Misc       Misc       `xml:"misc"`
	Protocol   Protocol   `xml:"protocol"`
	Transport  Transport  `xml:"transport"`
	RicService RicService `xml:"ricService"`
}

type E2Node struct {
	Text                   string    `xml:",chardata"`
	E2nodeComponentUnknown *struct{} `xml:"e2node-component-unknown"`
}

type RicRequest struct {
	Text                                       string    `xml:",chardata"`
	RanFunctionIDInvalid                       *struct{} `xml:"ran-function-id-invalid"`
	ActionNotSupported                         *struct{} `xml:"action-not-supported"`
	ExcessiveActions                           *struct{} `xml:"excessive-actions"`
	DuplicateAction                            *struct{} `xml:"duplicate-action"`
	DuplicateEventTrigger                      *struct{} `xml:"duplicate-event-trigger"`
	FunctionResourceLimit                      *struct{} `xml:"function-resource-limit"`
	RequestIDUnknown                           *struct{} `xml:"request-id-unknown"`
	InconsistentActionSubsequentActionSequence *struct{} `xml:"inconsistent-action-subsequent-action-sequence"`
	ControlMessageInvalid                      *struct{} `xml:"control-message-invalid"`
	RicCallProcessIDInvalid                    *struct{} `xml:"ric-call-process-id-invalid"`
	ControlTimerExpired                        *struct{} `xml:"control-timer-expired"`
	ControlFailedToExecute                     *struct{} `xml:"control-failed-to-execute"`
	SystemNotReady                             *struct{} `xml:"system-not-ready"`
	Unspecified                                *struct{} `xml:"unspecified"`
}

type Misc struct {
	Text                      string    `xml:",chardata"`
	ControlProcessingOverload *struct{} `xml:"control-processing-overload"`
	HardwareFailure           *struct{} `xml:"hardware-failure"`
	OmIntervention            *struct{} `xml:"om-intervention"`
	Unspecified               *struct{} `xml:"unspecified"`
}

type Protocol struct {
	Text                                         string    `xml:",chardata"`
	TransferSyntaxError                          *struct{} `xml:"transfer-syntax-error"`
	AbstractSyntaxErrorReject                    *struct{} `xml:"abstract-syntax-error-reject"`
	AbstractSyntaxErrorIgnoreAndNotify           *struct{} `xml:"abstract-syntax-error-ignore-and-notify"`
	MessageNotCompatibleWithReceiverState        *struct{} `xml:"message-not-compatible-with-receiver-state"`
	SemanticError                                *struct{} `xml:"semantic-error"`
	AbstractSyntaxErrorFalselyConstructedMessage *struct{} `xml:"abstract-syntax-error-falsely-constructed-message"`
	Unspecified                                  *struct{} `xml:"unspecified"`
}

type Transport struct {
	Text                         string    `xml:",chardata"`
	Unspecified                  *struct{} `xml:"unspecified"`
	TransportResourceUnavailable *struct{} `xml:"transport-resource-unavailable"`
}

type RicService struct {
	Text                    string    `xml:",chardata"`
	RanFunctionNotSupported *struct{} `xml:"ran-function-not-supported"`
	ExcessiveFunctions      *struct{} `xml:"excessive-functions"`
	RicResourceLimit        *struct{} `xml:"ric-resource-limit"`
}
