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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
)

type NodebValidator struct {
}

func NewNodebValidator() *NodebValidator {
	return &NodebValidator{}
}

func (h *NodebValidator) IsGnbValid(gnb *entities.Gnb) error {
	if len(gnb.ServedNrCells) == 0 {
		return errors.New("gnb.servedNrCells")
	}

	for _, servedNrCell := range gnb.ServedNrCells {
		if servedNrCell.ServedNrCellInformation == nil {
			return errors.New("gnb.servedNrCellInformation")
		}

		err := isServedNrCellInformationValid(servedNrCell.ServedNrCellInformation)

		if err != nil {
			return err
		}

		if len(servedNrCell.NrNeighbourInfos) == 0 {
			continue
		}

		for _, nrNeighbourInformation := range servedNrCell.NrNeighbourInfos {

			err := isNrNeighbourInformationValid(nrNeighbourInformation)

			if err != nil {
				return err
			}

		}
	}

	return nil
}

func isServedNrCellInformationValid(servedNrCellInformation *entities.ServedNRCellInformation) error {
	if servedNrCellInformation.CellId == "" {
		return errors.New("servedNrCellInformation.cellId")
	}

	if servedNrCellInformation.ChoiceNrMode == nil {
		return errors.New("servedNrCellInformation.choiceNrMode")
	}

	if servedNrCellInformation.NrMode == entities.Nr_UNKNOWN {
		return errors.New("servedNrCellInformation.nrMode")
	}

	if len(servedNrCellInformation.ServedPlmns) == 0 {
		return errors.New("servedNrCellInformation.servedPlmns")
	}

	return isServedNrCellInfoChoiceNrModeValid(servedNrCellInformation.ChoiceNrMode)
}

func isServedNrCellInfoChoiceNrModeValid(choiceNrMode *entities.ServedNRCellInformation_ChoiceNRMode) error {
	if choiceNrMode.Fdd != nil {
		return isServedNrCellInfoFddValid(choiceNrMode.Fdd)
	}

	if choiceNrMode.Tdd != nil {
		return isServedNrCellInfoTddValid(choiceNrMode.Tdd)
	}

	return errors.New("servedNrCellInformation.choiceNrMode.fdd / servedNrCellInformation.choiceNrMode.tdd")
}

func isServedNrCellInfoTddValid(tdd *entities.ServedNRCellInformation_ChoiceNRMode_TddInfo) error {
	return nil
}

func isServedNrCellInfoFddValid(fdd *entities.ServedNRCellInformation_ChoiceNRMode_FddInfo) error {
	return nil
}

func isNrNeighbourInformationValid(nrNeighbourInformation *entities.NrNeighbourInformation) error {
	if nrNeighbourInformation.NrCgi == "" {
		return errors.New("nrNeighbourInformation.nrCgi")
	}

	if nrNeighbourInformation.ChoiceNrMode == nil {
		return errors.New("nrNeighbourInformation.choiceNrMode")
	}

	if nrNeighbourInformation.NrMode == entities.Nr_UNKNOWN {
		return errors.New("nrNeighbourInformation.nrMode")
	}

	return isNrNeighbourInfoChoiceNrModeValid(nrNeighbourInformation.ChoiceNrMode)
}

func isNrNeighbourInfoChoiceNrModeValid(choiceNrMode *entities.NrNeighbourInformation_ChoiceNRMode) error {
	if choiceNrMode.Fdd != nil {
		return isNrNeighbourInfoFddValid(choiceNrMode.Fdd)
	}

	if choiceNrMode.Tdd != nil {
		return isNrNeighbourInfoTddValid(choiceNrMode.Tdd)
	}

	return errors.New("nrNeighbourInformation.choiceNrMode.fdd / nrNeighbourInformation.choiceNrMode.tdd")
}

func isNrNeighbourInfoTddValid(tdd *entities.NrNeighbourInformation_ChoiceNRMode_TddInfo) error {
	return nil
}

func isNrNeighbourInfoFddValid(fdd *entities.NrNeighbourInformation_ChoiceNRMode_FddInfo) error {
	return nil
}

func (h *NodebValidator) IsEnbValid(enb *entities.Enb) error {
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
