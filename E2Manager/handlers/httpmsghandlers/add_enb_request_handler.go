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
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
)

type AddEnbRequestHandler struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
}

func NewAddEnbRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService) *AddEnbRequestHandler {
	return &AddEnbRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
	}
}

func (h *AddEnbRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	addEnbRequest := request.(*models.AddEnbRequest)

	h.logger.Infof("#AddEnbRequestHandler.Handle - Ran name: %s", addEnbRequest.RanName)

	err := h.validateRequestBody(addEnbRequest)

	if err != nil {
		h.logger.Errorf("#AddEnbRequestHandler.Handle - validation failure: %s is a mandatory field and cannot be empty", err)
		return nil, e2managererrors.NewRequestValidationError()
	}

	_, err = h.rNibDataService.GetNodeb(addEnbRequest.RanName)

	if err == nil {
		h.logger.Errorf("#AddEnbRequestHandler.Handle - RAN name: %s - RAN already exists. quit", addEnbRequest.RanName)
		return nil, e2managererrors.NewNodebExistsError()
	}

	_, ok := err.(*common.ResourceNotFoundError)
	if !ok {
		h.logger.Errorf("#AddEnbRequestHandler.Handle - RAN name: %s - failed to get nodeb entity from RNIB. Error: %s", addEnbRequest.RanName, err)
		return nil, e2managererrors.NewRnibDbError()
	}

	nbIdentity := h.createNbIdentity(addEnbRequest)
	nodebInfo := h.createNodebInfo(addEnbRequest)

	err = h.rNibDataService.SaveNodeb(nbIdentity, nodebInfo)

	if err != nil {
		h.logger.Errorf("#AddEnbRequestHandler.Handle - RAN name: %s - failed to save nodeb entity in RNIB. Error: %s", addEnbRequest.RanName, err)
		return nil, e2managererrors.NewRnibDbError()
	}

	return models.NewAddEnbResponse(nodebInfo), nil
}

func (h *AddEnbRequestHandler) createNodebInfo(addEnbRequest *models.AddEnbRequest) *entities.NodebInfo {
	nodebInfo := entities.NodebInfo{
		RanName:          addEnbRequest.RanName,
		Ip:               addEnbRequest.Ip,
		Port:             addEnbRequest.Port,
		GlobalNbId:       addEnbRequest.GlobalNbId,
		Configuration:    &entities.NodebInfo_Enb{Enb: addEnbRequest.Enb},
		ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,
	}

	return &nodebInfo
}

func (h *AddEnbRequestHandler) createNbIdentity(addEnbRequest *models.AddEnbRequest) *entities.NbIdentity {
	nbIdentity := entities.NbIdentity{
		GlobalNbId:    addEnbRequest.GlobalNbId,
		InventoryName: addEnbRequest.RanName,
	}

	return &nbIdentity
}

func (h *AddEnbRequestHandler) validateRequestBody(addEnbRequest *models.AddEnbRequest) error {

	if addEnbRequest.RanName == "" {
		return errors.New("ranName")
	}

	if addEnbRequest.GlobalNbId == nil {
		return errors.New("globalNbId")
	}

	if err := isGlobalNbIdValid(addEnbRequest.GlobalNbId); err != nil {
		return err
	}

	if addEnbRequest.Enb == nil {
		return errors.New("enb")
	}

	if err := isEnbValid(addEnbRequest.Enb); err != nil {
		return err
	}

	return nil
}

func isGlobalNbIdValid(globalNbId *entities.GlobalNbId) error {
	if globalNbId.PlmnId == "" {
		return errors.New("globalNbId.plmnId")
	}

	if globalNbId.NbId == "" {
		return errors.New("globalNbId.nbId")
	}

	return nil
}

func isEnbValid(enb *entities.Enb) error {
	if enb.EnbType == entities.EnbType_UNKNOWN_ENB_TYPE {
		return errors.New("enb.enbType")
	}

	if enb.ServedCells == nil || len(enb.ServedCells) == 0 {
		return errors.New("enb.servedCells")
	}

	for _, servedCell := range enb.ServedCells {
		err := isServedCellValid(servedCell)

		if err != nil {
			return err
		}
	}

	return nil
}

func isServedCellValid(servedCell *entities.ServedCellInfo) error {

	if servedCell.CellId == "" {
		return errors.New("servedCell.cellId")
	}

	if servedCell.EutraMode == entities.Eutra_UNKNOWN {
		return errors.New("servedCell.eutraMode")
	}

	if servedCell.Tac == "" {
		return errors.New("servedCell.tac")
	}

	if servedCell.BroadcastPlmns == nil || len(servedCell.BroadcastPlmns) == 0 {
		return errors.New("servedCell.broadcastPlmns")
	}

	if servedCell.ChoiceEutraMode == nil {
		return errors.New("servedCell.choiceEutraMode")
	}

	return isChoiceEutraModeValid(servedCell.ChoiceEutraMode)
}

func isChoiceEutraModeValid(choiceEutraMode *entities.ChoiceEUTRAMode) error {
	if choiceEutraMode.Fdd != nil {
		return isFddInfoValid(choiceEutraMode.Fdd)
	}

	if choiceEutraMode.Tdd != nil {
		return isTddInfoValid(choiceEutraMode.Tdd)
	}

	return errors.New("servedCell.fdd / servedCell.tdd")
}

func isTddInfoValid(tdd *entities.TddInfo) error {
	return nil
}

func isFddInfoValid(fdd *entities.FddInfo) error {
	return nil
}
