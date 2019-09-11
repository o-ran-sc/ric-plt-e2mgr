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
//

package httpmsghandlers

import (
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

const (
	X2_RESET_ACTIVITY_NAME = "X2_RESET"
)

type X2ResetRequestHandler struct {
	rNibDataService services.RNibDataService
	rmrService     *services.RmrService
	logger         *logger.Logger
}

func NewX2ResetRequestHandler(logger *logger.Logger, rmrService *services.RmrService, rNibDataService services.RNibDataService) *X2ResetRequestHandler {
	return &X2ResetRequestHandler{
		rNibDataService: rNibDataService,
		rmrService:     rmrService,
		logger:         logger,
	}
}

func (handler *X2ResetRequestHandler) Handle(request models.Request) error {
	resetRequest := request.(models.ResetRequest)

	if len(resetRequest.Cause) == 0 {
		resetRequest.Cause = e2pdus.OmInterventionCause
	}
	payload, ok := e2pdus.KnownCausesToX2ResetPDU(resetRequest.Cause)
	if !ok {
		handler.logger.Errorf("#reset_request_handler.Handle - Unknown cause (%s)", resetRequest.Cause)
		return e2managererrors.NewRequestValidationError()
	}

	nodeb, err := handler.rNibDataService.GetNodeb(resetRequest.RanName)
	if err != nil {
		handler.logger.Errorf("#reset_request_handler.Handle - failed to get status of RAN: %s from RNIB. Error: %s", resetRequest.RanName, err.Error())
		_, ok := err.(*common.ResourceNotFoundError)
		if ok {
			return e2managererrors.NewResourceNotFoundError()
		}
		return e2managererrors.NewRnibDbError()
	}

	if nodeb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		handler.logger.Errorf("#reset_request_handler.Handle - RAN: %s in wrong state (%s)", resetRequest.RanName, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
		return e2managererrors.NewWrongStateError(X2_RESET_ACTIVITY_NAME, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
	}

	response := models.NotificationResponse{MgsType: rmrCgo.RIC_X2_RESET, RanName: resetRequest.RanName, Payload: payload}
	if err := handler.rmrService.SendRmrMessage(&response); err != nil {
		handler.logger.Errorf("#reset_request_handler.Handle - failed to send reset message to RMR: %s", err)
		return e2managererrors.NewRmrError()
	}

	handler.logger.Infof("#reset_request_handler.Handle - sent x2 reset to RAN: %s with cause: %s", resetRequest.RanName, resetRequest.Cause)
	return nil
}
