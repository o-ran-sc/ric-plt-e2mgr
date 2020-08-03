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

package managers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
)

type UpdateEnbManager struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	nodebValidator  *NodebValidator
}

func NewUpdateEnbManager(logger *logger.Logger, rnibDataService services.RNibDataService, nodebValidator *NodebValidator) *UpdateEnbManager {
	return &UpdateEnbManager{
		logger:          logger,
		rnibDataService: rnibDataService,
		nodebValidator:  nodebValidator,
	}
}

func (h *UpdateEnbManager) Validate(request models.Request) error {

	updateEnbRequest := request.(*models.UpdateEnbRequest)

	h.logger.Infof("#UpdateEnbManager.Validate - Validate incoming request, ran name: %s", updateEnbRequest.RanName)

	if err := h.validateRequestBody(updateEnbRequest); err != nil {
		h.logger.Errorf("#UpdateEnbManager.Validate - validation failure: %s is a mandatory field and cannot be empty", err)
		return err
	}

	return nil
}

func (h *UpdateEnbManager) RemoveNodebCells(nodeb *entities.NodebInfo) error {

	if nodeb.NodeType != entities.Node_ENB {
		h.logger.Errorf("#UpdateEnbManager.RemoveNodebCells - RAN name: %s - RAN is not eNB.", nodeb.RanName)
		return e2managererrors.NewRequestValidationError()
	}

	servedCells := nodeb.GetEnb().GetServedCells()

	if len(servedCells) == 0 {
		h.logger.Infof("#UpdateGnbManager.RemoveNodebCells - RAN name: %s - eNB cells are nil or empty - no cells to remove", nodeb.GetRanName())
		return nil
	}

	err := h.rnibDataService.RemoveServedCells(nodeb.GetRanName(), servedCells)
	if err != nil {
		h.logger.Errorf("#UpdateEnbManager.RemoveNodebCells - RAN name: %s - Failed removing eNB served cells", nodeb.GetRanName())
		return e2managererrors.NewRnibDbError()
	}
	h.logger.Infof("#UpdateEnbManager.RemoveNodebCells - RAN name: %s - Successfully removed eNB served cells", nodeb.GetRanName())

	return nil
}

func (h *UpdateEnbManager) SetNodeb(nodeb *entities.NodebInfo, request models.Request) {
	updateEnbRequest := request.(*models.UpdateEnbRequest)
	updateEnbRequest.Enb.EnbType = nodeb.GetEnb().GetEnbType()
	nodeb.Configuration = &entities.NodebInfo_Enb{Enb: updateEnbRequest.Enb}
}

func (h *UpdateEnbManager) UpdateNodeb(nodeb *entities.NodebInfo) error {

	err := h.rnibDataService.UpdateEnb(nodeb, nodeb.GetEnb().GetServedCells())
	if err != nil {
		h.logger.Errorf("#UpdateEnbManager.UpdateNodeb - RAN name: %s - Failed updating eNB. Error: %s", nodeb.GetRanName(), err)
		return e2managererrors.NewRnibDbError()
	}
	h.logger.Infof("#UpdateEnbManager.UpdateNodeb - RAN name: %s - Successfully updated eNB", nodeb.GetRanName())

	return nil
}

func (h *UpdateEnbManager) validateRequestBody(request *models.UpdateEnbRequest) error {

	if request.Enb == nil {
		return errors.New("enb")
	}

	if err := h.nodebValidator.IsEnbValid(request.Enb); err != nil {
		return err
	}

	if h.nodebValidator.IsNgEnbType(request.Enb.GetEnbType()){
		return errors.New("enb.enbType")
	}

	return nil
}
