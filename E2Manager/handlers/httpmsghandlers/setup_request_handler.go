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

package httpmsghandlers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rnibBuilders"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	X2SetupActivityName   = "X2_SETUP"
	EndcSetupActivityName = "ENDC_SETUP"
)

type SetupRequestHandler struct {
	rNibDataService services.RNibDataService
	logger          *logger.Logger
	ranSetupManager *managers.RanSetupManager
	protocol        entities.E2ApplicationProtocol
}

func NewSetupRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService,
	ranSetupManager *managers.RanSetupManager, protocol entities.E2ApplicationProtocol) *SetupRequestHandler {
	return &SetupRequestHandler{
		logger:          logger,
		rNibDataService:  rNibDataService,
		ranSetupManager: ranSetupManager,
		protocol:        protocol,
	}
}

func (handler *SetupRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	setupRequest := request.(models.SetupRequest)

	err := handler.validateRequestDetails(setupRequest)
	if err != nil {
		return nil, err
	}

	nodebInfo, err := handler.rNibDataService.GetNodeb(setupRequest.RanName)
	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			handler.logger.Errorf("#SetupRequestHandler.Handle - failed to get nodeB entity for ran name: %v from RNIB. Error: %s",
				setupRequest.RanName, err.Error())
			return nil, e2managererrors.NewRnibDbError()
		}

		result := handler.connectNewRan(&setupRequest, handler.protocol)
		return nil, result
	}

	result := handler.connectExistingRan(nodebInfo)
	return nil, result
}

func (handler *SetupRequestHandler) connectExistingRan(nodebInfo *entities.NodebInfo) error {

	if nodebInfo.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		handler.logger.Errorf("#SetupRequestHandler.connectExistingRan - RAN: %s in wrong state (%s)", nodebInfo.RanName, entities.ConnectionStatus_name[int32(nodebInfo.ConnectionStatus)])
		return e2managererrors.NewWrongStateError(handler.getActivityName(handler.protocol), entities.ConnectionStatus_name[int32(nodebInfo.ConnectionStatus)])
	}

	status := entities.ConnectionStatus_CONNECTING
	if nodebInfo.ConnectionStatus == entities.ConnectionStatus_CONNECTED{
		status = nodebInfo.ConnectionStatus
	}
	nodebInfo.ConnectionAttempts = 0

	result := handler.ranSetupManager.ExecuteSetup(nodebInfo, status)
	return result
}

func (handler *SetupRequestHandler) connectNewRan(request *models.SetupRequest, protocol entities.E2ApplicationProtocol) error {

	nodebInfo, nodebIdentity := rnibBuilders.CreateInitialNodeInfo(request, protocol)

	rNibErr := handler.rNibDataService.SaveNodeb(nodebIdentity, nodebInfo)
	if rNibErr != nil {
		handler.logger.Errorf("#SetupRequestHandler.connectNewRan - failed to initial nodeb entity for ran name: %v in RNIB. Error: %s", request.RanName, rNibErr.Error())
		return e2managererrors.NewRnibDbError()
	}
	handler.logger.Infof("#SetupRequestHandler.connectNewRan - initial nodeb entity for ran name: %v was saved to RNIB ", request.RanName)

	result := handler.ranSetupManager.ExecuteSetup(nodebInfo, entities.ConnectionStatus_CONNECTING)
	return result
}

func (handler *SetupRequestHandler) validateRequestDetails(request models.SetupRequest) error {

	if request.RanPort == 0 {
		handler.logger.Errorf("#SetupRequestHandler.validateRequestDetails - validation failure: port cannot be zero")
		return e2managererrors.NewRequestValidationError()
	}
	err := validation.ValidateStruct(&request,
		validation.Field(&request.RanIp, validation.Required, is.IP),
		validation.Field(&request.RanName, validation.Required),
	)

	if err != nil {
		handler.logger.Errorf("#SetupRequestHandler.validateRequestDetails - validation failure, error: %v", err)
		return e2managererrors.NewRequestValidationError()
	}

	return nil
}

func (handler *SetupRequestHandler) getActivityName(protocol entities.E2ApplicationProtocol) string {
	if protocol == entities.E2ApplicationProtocol_X2_SETUP_REQUEST {
		return X2SetupActivityName
	}
	return EndcSetupActivityName
}
