//
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
//  platform project (RICP).

package managers


import(
	"e2mgr/logger"
	"e2mgr/services"
	"errors"
	"e2mgr/models"
)


type IRicServiceUpdateManager interface {
	RevertRanFunctions(ranName string) error
	StoreExistingRanFunctions(ranName string) error
}


type RicServiceUpdateManager struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
}

func NewRicServiceUpdateManager(logger *logger.Logger, rNibDataService services.RNibDataService) *RicServiceUpdateManager {
	return &RicServiceUpdateManager{
		logger:          logger,
		rNibDataService: rNibDataService,
	}
}


func (h *RicServiceUpdateManager) StoreExistingRanFunctions(ranName string) error {
	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil {
		h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions - failed to get nodeB entity for ran name: %v due to RNIB Error: %s", ranName, err)
	}
	if nodebInfo.GetGnb() == nil {
        h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions - GNB is nil for RAN name: %s", ranName)
        return errors.New("There is empty gnb nodebInfo")
    }
	models.ExistingRanFunctiuonsMap[ranName] =  nodebInfo.GetGnb().RanFunctions
	h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions - Updated ranFunctions for reverting the changes are %v:", models.ExistingRanFunctiuonsMap[ranName])
	return nil
}

func (h *RicServiceUpdateManager) RevertRanFunctions(ranName string) error {
	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)
	if err != nil {
		h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions - failed to get nodeB entity for ran name: %v due to RNIB Error: %s", ranName, err)
	}

	if nodebInfo.GetGnb() != nil && nodebInfo.GetGnb().RanFunctions != nil {
		nodebInfo.GetGnb().RanFunctions = models.ExistingRanFunctiuonsMap[ranName]
	} else {
		h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions returned nil")
	}
	err = h.rNibDataService.UpdateNodebInfoAndPublish(nodebInfo)
	if err != nil {
		h.logger.Errorf("#RicServiceUpdateManager.revertRanFunctions - RAN name: %s - Failed at UpdateNodebInfoAndPublish. error: %s", nodebInfo.RanName, err)
		return err
	}

	h.logger.Infof("#RicServiceUpdateManager.revertRanFunctions - Revert ranFunctions for RAN name: %s", ranName)
	return nil
}