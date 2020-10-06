//
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package httpmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
)

type GetNodebIdRequestHandler struct {
	logger          *logger.Logger
	ranListManager  managers.RanListManager
}

func NewGetNodebIdRequestHandler(logger *logger.Logger, ranListManager managers.RanListManager) *GetNodebIdRequestHandler {
	return &GetNodebIdRequestHandler{
		logger:          logger,
		ranListManager:  ranListManager,
	}
}

func (h *GetNodebIdRequestHandler) Handle(request models.Request) (models.IResponse, error) {
	getNodebIdRequest := request.(models.GetNodebIdRequest)
	ranName := getNodebIdRequest.RanName

	nodebId, err := h.ranListManager.GetNbIdentity(ranName)
	if err != nil {
		return nil, err
	}

	return models.NewNodebIdResponse(nodebId), nil
}
