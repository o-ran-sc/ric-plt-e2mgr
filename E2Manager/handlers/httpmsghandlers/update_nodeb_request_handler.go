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
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
)

type UpdateNodebRequestHandler struct {
	logger             *logger.Logger
	rNibDataService    services.RNibDataService
	updateNodebManager managers.IUpdateNodebManager
}

func NewUpdateNodebRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService, updateNodebManager managers.IUpdateNodebManager) *UpdateNodebRequestHandler {
	return &UpdateNodebRequestHandler{
		logger:             logger,
		rNibDataService:    rNibDataService,
		updateNodebManager: updateNodebManager,
	}
}

func (h *UpdateNodebRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	ranName := h.getRanName(request)

	h.logger.Infof("#UpdateNodebRequestHandler.Handle - Ran name: %s", ranName)

	err := h.updateNodebManager.Validate(request)
	if err != nil {
		return nil, e2managererrors.NewRequestValidationError()
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			h.logger.Errorf("#UpdateNodebRequestHandler.Handle - RAN name: %s - failed to get nodeb entity from RNIB. Error: %s", ranName, err)
			return nil, e2managererrors.NewRnibDbError()
		}

		h.logger.Errorf("#UpdateNodebRequestHandler.Handle - RAN name: %s - RAN not found on RNIB. Error: %s", ranName, err)
		return nil, e2managererrors.NewResourceNotFoundError()
	}

	err = h.updateNodebManager.RemoveNodebCells(nodebInfo)
	if err != nil {
		return nil, err
	}

	err = h.updateNodebManager.SetNodeb(nodebInfo, request)
	if err != nil {
		return nil, err
	}

	err = h.updateNodebManager.UpdateNodeb(nodebInfo)
	if err != nil {
		return nil, err
	}

	return models.NewNodebResponse(nodebInfo), nil
}

func (h *UpdateNodebRequestHandler) getRanName(request models.Request) string {

	var ranName string
	updateEnbRequest, ok := request.(*models.UpdateEnbRequest)
	if !ok {
		//updateGnbRequest := request.(*models.UpdateGnbRequest)
		//ranName = updateGnbRequest.RanName
	} else {
		ranName = updateEnbRequest.RanName
	}
	return ranName
}