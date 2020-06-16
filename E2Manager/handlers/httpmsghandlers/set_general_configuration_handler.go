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

import "C"
import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
)

type SetGeneralConfigurationHandler struct {
	rnibDataService services.RNibDataService
	logger          *logger.Logger
}

func NewSetGeneralConfigurationHandler(logger *logger.Logger, rnibDataService services.RNibDataService) *SetGeneralConfigurationHandler {
	return &SetGeneralConfigurationHandler{
		logger:          logger,
		rnibDataService: rnibDataService,
	}
}

func (h *SetGeneralConfigurationHandler) Handle(request models.Request) (models.IResponse, error) {
	h.logger.Infof("#SetGeneralConfigurationHandler.Handle - handling set general parameters")

	configuration := request.(models.GeneralConfigurationRequest)

	existingConfig, err := h.rnibDataService.GetGeneralConfiguration()

	if err != nil {
		return nil, err
	}

	if existingConfig.EnableRic != configuration.EnableRic {

		existingConfig.EnableRic = configuration.EnableRic
		err := h.rnibDataService.SaveGeneralConfiguration(existingConfig)

		if err != nil {
			return nil, err
		}

	}
	response := &models.GeneralConfigurationResponse{EnableRic: configuration.EnableRic}

	return response, nil
}
