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
)

type UpdateGnbManager struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	nodebValidator  *NodebValidator
}

func NewUpdateGnbManager(logger *logger.Logger, rnibDataService services.RNibDataService, nodebValidator *NodebValidator) *UpdateGnbManager {
	return &UpdateGnbManager{
		logger:          logger,
		rnibDataService: rnibDataService,
		nodebValidator:  nodebValidator,
	}
}

func (h *UpdateGnbManager) Validate(request models.Request) error {

	updateGnbRequest := request.(*models.UpdateGnbRequest)

	if err := h.nodebValidator.IsGnbValid(updateGnbRequest.Gnb); err != nil {
		h.logger.Errorf("#UpdateGnbManager.Validate - RAN name: %s - validation failure: %s is a mandatory field and cannot be empty", updateGnbRequest.RanName, err)
		return err
	}

	return nil
}

func (h *UpdateGnbManager) RemoveNodebCells(nodeb *entities.NodebInfo) error {

	if nodeb.NodeType != entities.Node_GNB {
		h.logger.Errorf("#UpdateGnbManager.RemoveNodebCells - RAN name: %s - node type isn't gNB", nodeb.GetRanName())
		return e2managererrors.NewRequestValidationError()
	}

	servedNrCells := nodeb.GetGnb().GetServedNrCells()

	if len(servedNrCells) == 0 {
		h.logger.Infof("#UpdateGnbManager.RemoveNodebCells - RAN name: %s - gNB cells are nil or empty - no cells to remove", nodeb.GetRanName())
		return nil
	}

	err := h.rnibDataService.RemoveServedNrCells(nodeb.GetRanName(), servedNrCells)
	if err != nil {
		h.logger.Errorf("#UpdateGnbManager.RemoveNodebCells - RAN name: %s - Failed removing gNB cells", nodeb.GetRanName())
		return e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#UpdateGnbManager.RemoveNodebCells - RAN name: %s - Successfully removed gNB cells", nodeb.GetRanName())
	return nil
}

func (h *UpdateGnbManager) SetNodeb(nodeb *entities.NodebInfo, request models.Request) {

	updateGnbRequest := request.(*models.UpdateGnbRequest)
	nodeb.GetGnb().ServedNrCells = updateGnbRequest.ServedNrCells
}

func (h *UpdateGnbManager) UpdateNodeb(nodeb *entities.NodebInfo) error {

	err := h.rnibDataService.UpdateGnbCells(nodeb, nodeb.GetGnb().GetServedNrCells())
	if err != nil {
		h.logger.Errorf("#UpdateGnbManager.UpdateNodeb - RAN name: %s - Failed updating gNB cells. Error: %s", nodeb.GetRanName(), err)
		return e2managererrors.NewRnibDbError()
	}
	h.logger.Infof("#UpdateGnbManager.UpdateNodeb - RAN name: %s - Successfully updated gNB cells", nodeb.GetRanName())

	return nil
}