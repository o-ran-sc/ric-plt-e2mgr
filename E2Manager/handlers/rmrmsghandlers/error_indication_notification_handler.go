// Copyright 2023 Nokia
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
//  platform project (RICP)

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/utils"
	"encoding/xml"
	"fmt"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"sync"
)

var e2ErrorIndicationMessage = models.ErrorIndicationMessage{}

var E2SETUP_PROCEDURE string = "1"
var RICSERVICEUPDATE_PROCEDURE string = "7"

type ErrorIndicationHandler struct {
	logger                  *logger.Logger
	ranDisconnectionManager managers.IRanDisconnectionManager
	RicServiceUpdateManager managers.IRicServiceUpdateManager
	procedureMapMutex       sync.RWMutex
}


func ErrorIndicationNotificationHandler(logger *logger.Logger, ranDisconnectionManager managers.IRanDisconnectionManager, RicServiceUpdateManager managers.IRicServiceUpdateManager) *ErrorIndicationHandler {
	return &ErrorIndicationHandler{
		logger:                  logger,
		ranDisconnectionManager: ranDisconnectionManager,
		RicServiceUpdateManager: RicServiceUpdateManager,
	}
}
func (errorIndicationHandler *ErrorIndicationHandler) Handle (request *models.NotificationRequest) {
	ranName := request.RanName
	errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-Received Error Indication from E2Node - %s", ranName)

	errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-Received ErrorIndication payload at E2M is - %x", request.Payload)
	errorIndicationMessage, err := errorIndicationHandler.parseErrorIndication(request.Payload)
	if err != nil {
		errorIndicationHandler.logger.Errorf("#ErrorIndicationHandler.Handle- Parsing is not successful")
		return
	}
	errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle ERROR INDICATION from E2Node has been parsed successfully- %+v", errorIndicationMessage)
	errorIndicationIE := errorIndicationMessage.E2APPDU.InitiatingMessage.Value.ErrorIndication.ProtocolIEs.ErrorIndicationIEs
	fmt.Printf("errorIndicationIE value is %+v", errorIndicationIE)
	
	for i := 0 ; i < len(errorIndicationIE) ; i++ {
		if errorIndicationIE[i].ID == 2 {
			errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-CD is: %+v", errorIndicationIE[i].Value.CriticalityDiagnostics)
			if errorIndicationIE[i].Value.CriticalityDiagnostics.ProcedureCode != "" && errorIndicationIE[i].Value.CriticalityDiagnostics.TriggeringMessage.SuccessfulOutcome != nil {
				procedureCode := errorIndicationIE[i].Value.CriticalityDiagnostics.ProcedureCode
				errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-procedureCode present is: %+v", procedureCode)
				errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle- before triggeringMessage present is: %+v", errorIndicationIE[i].Value.CriticalityDiagnostics.TriggeringMessage)
				triggeringMessageValue := &errorIndicationIE[i].Value.CriticalityDiagnostics.TriggeringMessage.SuccessfulOutcome
				errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-triggeringMessage present is: %+v", *triggeringMessageValue)
				if procedureCode != "" && triggeringMessageValue != nil {
					errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.handleErrorIndicationBasedOnProcedureCode for all scenarios")
					switch procedureCode {
					case E2SETUP_PROCEDURE:
						if triggeringMessageValue != nil {
							errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication happened at E2Setup procedure")
							err = errorIndicationHandler.ranDisconnectionManager.DisconnectRan(ranName)
							errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-Cleanup Completed !!")
							return
						} else {
							errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication recieved for unsuccessful-outcome, no action taken")
						}
					case RICSERVICEUPDATE_PROCEDURE:
						if triggeringMessageValue != nil {
							errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication happened at Ric Service Update procedure")
							err = errorIndicationHandler.RicServiceUpdateManager.RevertRanFunctions(ranName)
							if err != nil {
								errorIndicationHandler.logger.Errorf("#ErrorIndicationHandler.Handle-reverting RanFunctions and updating the nodebInfo failed due to error %+v", err)
							}
							return
						} else {
							errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication recieved for unsuccessful-outcome, no action taken")
						}
					default:
						errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-problem in handling of error indication")
						return
					}
				}
			}
		}
	}
	errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-CriticalityDiagnostics IEs unsuccessful hence Retrieving based on procedureMap")
	errorIndicationHandler.HandleBasedOnProcedureType(ranName)
	errorIndicationHandler.logger.Debugf("#ErrorIndicationHandler.Handle-Cleanup Completed !!")
}



func (errorIndicationHandler *ErrorIndicationHandler) HandleBasedOnProcedureType(ranName string) error {
	errorIndicationHandler.procedureMapMutex.RLock()
	procedureType, ok := models.ProcedureMap[ranName]
	errorIndicationHandler.procedureMapMutex.RUnlock()
	if !ok {
		errorIndicationHandler.logger.Errorf("#ErrorIndicationHandler.Handle-Error ProcedureType not found for ranName %s", ranName)
	} else {
		switch procedureType {
		case models.E2SetupProcedureCompleted:
			errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication happened at E2Setup procedure")
			err := errorIndicationHandler.ranDisconnectionManager.DisconnectRan(ranName)
			if err != nil {
				errorIndicationHandler.logger.Errorf("#ErrorIndicationHandler.Handle-Disconnect RAN and updating the nodebInfo failed due to error %+v", err)
			}
		case models.RicServiceUpdateCompleted:
			errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication happened at Ric Service Update procedure")
			err := errorIndicationHandler.RicServiceUpdateManager.RevertRanFunctions(ranName)
			if err != nil {
				errorIndicationHandler.logger.Errorf("#ErrorIndicationHandler.Handle-reverting RanFunctions and updating the nodebInfo failed due to error %+v", err)
			}
		case models.E2SetupProcedureFailure, models.RicServiceUpdateFailure:
			errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-ErrorIndication occcured before successful outcome hence ignoring")
		default:
			errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.Handle-Error in handling the ErrorIndication based on enum")
		}
	}
	return nil
}

func (errorIndicationHandler *ErrorIndicationHandler) parseErrorIndication(payload []byte) (*models.ErrorIndicationMessage, error) {
	pipInd := bytes.IndexByte(payload, '|')
	if pipInd < 0 {
		return nil, common.NewInternalError(fmt.Errorf("#ErrorIndicationHandler.parseErrorIndication - Error parsing ERROR INDICATION failed extract Payload: no | separator found"))
	}
	errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.parseErrorIndication - payload: %s", payload)
	errorIndicationHandler.logger.Infof("#ErrorIndicationHandler.parseErrorIndication - payload: %s", payload[pipInd+1:])
	errorIndicationMessage := &models.ErrorIndicationMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload[pipInd+1:]), &errorIndicationMessage.E2APPDU)
	if err != nil {
		return nil, common.NewInternalError(fmt.Errorf("#ErrorIndicationHandler.parseErrorIndication - Error unmarshalling ERROR INDICATION payload: %x", payload))
	}
	return errorIndicationMessage, nil
}