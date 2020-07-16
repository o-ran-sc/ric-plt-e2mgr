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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type DeleteEnbRequestHandler struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
	ranListManager  managers.RanListManager
}

func NewDeleteEnbRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService, ranListManager managers.RanListManager) *DeleteEnbRequestHandler {
	return &DeleteEnbRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
		ranListManager:  ranListManager,
	}
}

func (h *DeleteEnbRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	deleteEnbRequest := request.(*models.DeleteEnbRequest)

	h.logger.Infof("#DeleteEnbRequestHandler.Handle - RAN name: %s", deleteEnbRequest.RanName)

	nodebInfo, err := h.rNibDataService.GetNodeb(deleteEnbRequest.RanName)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			h.logger.Errorf("#DeleteEnbRequestHandler.Handle - RAN name: %s - failed to get nodeb entity from RNIB. Error: %s", deleteEnbRequest.RanName, err)
			return nil, e2managererrors.NewRnibDbError()
		}

		h.logger.Errorf("#DeleteEnbRequestHandler.Handle - RAN name: %s - RAN not found on RNIB. Error: %s", deleteEnbRequest.RanName, err)
		return nil, e2managererrors.NewResourceNotFoundError()
	}

	if nodebInfo.NodeType != entities.Node_ENB {
		h.logger.Errorf("#DeleteEnbRequestHandler.Handle - RAN name: %s - RAN is not eNB.", deleteEnbRequest.RanName)
		return nil, e2managererrors.NewRequestValidationError()
	}

	err = h.rNibDataService.RemoveEnb(nodebInfo)
	if err != nil {
		h.logger.Errorf("#DeleteEnbRequestHandler.Handle - RAN name: %s - failed to delete nodeb entity in RNIB. Error: %s", deleteEnbRequest.RanName, err)
		return nil, e2managererrors.NewRnibDbError()
	}

	err = h.ranListManager.RemoveNbIdentity(entities.Node_ENB, deleteEnbRequest.RanName)
	if err != nil {
		h.logger.Errorf("#DeleteEnbRequestHandler.Handle - RAN name: %s - failed to delete nbIdentity in RNIB. Error: %s", deleteEnbRequest.RanName, err)
		return nil, e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#DeleteEnbRequestHandler.Handle - RAN name: %s - deleted successfully.", deleteEnbRequest.RanName)
	return models.NewNodebResponse(nodebInfo), nil
}
