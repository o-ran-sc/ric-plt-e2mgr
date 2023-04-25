//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package controllers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/services"
	"encoding/json"
	"net/http"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type ISymptomdataController interface {
	GetSymptomData(writer http.ResponseWriter, r *http.Request)
}

type SymptomdataController struct {
	logger          *logger.Logger
	handlerProvider *httpmsghandlerprovider.IncomingRequestHandlerProvider
	rNibDataService services.RNibDataService
	ranListManager  managers.RanListManager
}

func NewSymptomdataController(l *logger.Logger, hp *httpmsghandlerprovider.IncomingRequestHandlerProvider, rnib services.RNibDataService, rl managers.RanListManager) *SymptomdataController {
	return &SymptomdataController{
		logger:          l,
		handlerProvider: hp,
		rNibDataService: rnib,
		ranListManager:  rl,
	}
}

func (s *SymptomdataController) GetSymptomData(w http.ResponseWriter, r *http.Request) {
	e2TList := s.handleRequest(httpmsghandlerprovider.GetE2TInstancesRequest)
	nodeBList := s.ranListManager.GetNbIdentityList()

	s.logger.Infof("nodeBList=%+v, e2TList=%+v", nodeBList, e2TList)

	var NodebInfoList []entities.NodebInfo
	for _, v := range nodeBList {
		if n, err := s.rNibDataService.GetNodeb(v.InventoryName); err == nil {
			NodebInfoList = append(NodebInfoList, *n)
		}
	}
	s.logger.Infof("NodebInfoList=%+v", NodebInfoList)

	symptomdata := struct {
		E2TList       models.IResponse       `json:"e2TList"`
		NodeBList     []*entities.NbIdentity `json:"nodeBList"`
		NodebInfoList []entities.NodebInfo   `json:"nodeBInfo"`
	}{
		e2TList,
		nodeBList,
		NodebInfoList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=platform/e2_info.json")
	w.WriteHeader(http.StatusOK)
	resp, _ := json.MarshalIndent(symptomdata, "", "    ")
	w.Write(resp)
}

func (s *SymptomdataController) handleRequest(requestName httpmsghandlerprovider.IncomingRequest) models.IResponse {
	handler, err := s.handlerProvider.GetHandler(requestName)
	if err != nil {
		return nil
	}

	resp, err := handler.Handle(nil)
	return resp
}
