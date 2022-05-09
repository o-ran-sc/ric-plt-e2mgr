//
// Copyright (c) 2022 Samsung Electronics Co., Ltd. All Rights Reserved.
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

package rmrmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"encoding/xml"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

var (
	toReplaceTags = []string{"reject", "ignore", "procedureCode", "id", "RANfunctionID-Item", "RANfunctionsID-List", "success", "s1", "ng", "e1", "f1", "w1", "x1", "xn"}
)

type E2nodeConfigUpdateNotificationHandler struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
	rmrSender       *rmrsender.RmrSender
}

func NewE2nodeConfigUpdateNotificationHandler(logger *logger.Logger, rNibDataService services.RNibDataService, rmrSender *rmrsender.RmrSender) *E2nodeConfigUpdateNotificationHandler {
	return &E2nodeConfigUpdateNotificationHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
		rmrSender:       rmrSender,
	}
}

func (e *E2nodeConfigUpdateNotificationHandler) Handle(request *models.NotificationRequest) {
	e.logger.Infof("#E2nodeConfigUpdateNotificationHandler.Handle - RAN name: %s - received E2_Config_Update. Payload: %x", request.RanName, request.Payload)
	e2NodeConfig, err := e.parseE2NodeConfigurationUpdate(request.Payload)
	if err != nil {
		e.logger.Errorf(err.Error())
		return
	}
	e.logger.Debugf("#E2nodeConfigUpdateNotificationHandler.Handle - RIC_E2_Node_Config_Update parsed successfully")

	nodebInfo, err := e.rNibDataService.GetNodeb(request.RanName)

	if err != nil {
		switch v := err.(type) {
		case *common.ResourceNotFoundError:
			e.logger.Errorf("#E2nodeConfigUpdateNotificationHandler.Handle - failed to get nodeB, E2nodeConfigUpdate will not be processed further.")
		default:
			e.logger.Errorf("#E2nodeConfigUpdateNotificationHandler.Handle - nobeB entity of RanName:%s absent in RNIB. Error: %s", request.RanName, v)
		}
		return
	}
	e.updateE2nodeConfig(e2NodeConfig, nodebInfo)
}

func (e *E2nodeConfigUpdateNotificationHandler) updateE2nodeConfig(e2nodeConfig *models.E2nodeConfigurationUpdateMessage, nodebInfo *entities.NodebInfo) {
	e.handleAddConfig(e2nodeConfig, nodebInfo)
	e.handleUpdateConfig(e2nodeConfig, nodebInfo)
	e.rNibDataService.UpdateNodebInfoAndPublish(nodebInfo)
}

func (e *E2nodeConfigUpdateNotificationHandler) compareConfigIDs(n1, n2 entities.E2NodeComponentConfig) bool {
	if n1.E2NodeComponentInterfaceType != n2.E2NodeComponentInterfaceType {
		return false
	}

	switch n1.E2NodeComponentInterfaceType {
	case entities.E2NodeComponentInterfaceType_ng:
		return n1.GetE2NodeComponentInterfaceTypeNG().GetAmfName() == n2.GetE2NodeComponentInterfaceTypeNG().GetAmfName()
	case entities.E2NodeComponentInterfaceType_xn:
		// TODO -- Not supported yet
		e.logger.Infof("#E2nodeConfigUpdateNotificationHandler.Handle - Interface type Xn is not supported")
	case entities.E2NodeComponentInterfaceType_e1:
		return n1.GetE2NodeComponentInterfaceTypeE1().GetGNBCuCpId() == n2.GetE2NodeComponentInterfaceTypeE1().GetGNBCuCpId()
	case entities.E2NodeComponentInterfaceType_f1:
		return n1.GetE2NodeComponentInterfaceTypeF1().GetGNBDuId() == n2.GetE2NodeComponentInterfaceTypeF1().GetGNBDuId()
	case entities.E2NodeComponentInterfaceType_w1:
		return n1.GetE2NodeComponentInterfaceTypeW1().GetNgenbDuId() == n2.GetE2NodeComponentInterfaceTypeW1().GetNgenbDuId()
	case entities.E2NodeComponentInterfaceType_s1:
		return n1.GetE2NodeComponentInterfaceTypeS1().GetMmeName() == n2.GetE2NodeComponentInterfaceTypeS1().GetMmeName()

	case entities.E2NodeComponentInterfaceType_x2:
		// TODO -- Not supported yet
		e.logger.Infof("#E2nodeConfigUpdateNotificationHandler.Handle - Interface type X2 is not supported")
	default:
		e.logger.Errorf("#E2nodeConfigUpdateNotificationHandler.Handle - Interface type not supported")
	}
	return false
}

func (e *E2nodeConfigUpdateNotificationHandler) handleAddConfig(e2nodeConfig *models.E2nodeConfigurationUpdateMessage, nodebInfo *entities.NodebInfo) {
	var result []*entities.E2NodeComponentConfig

	additionList := e2nodeConfig.ExtractConfigAdditionList()
	for i, _ := range additionList {
		result = append(result, &additionList[i])
	}

	if nodebInfo.NodeType == entities.Node_ENB {
		nodebInfo.GetEnb().NodeConfigs = append(result, nodebInfo.GetEnb().NodeConfigs...)

	} else {
		nodebInfo.GetGnb().NodeConfigs = append(result, nodebInfo.GetGnb().NodeConfigs...)
	}
}

func (e *E2nodeConfigUpdateNotificationHandler) handleUpdateConfig(e2nodeConfig *models.E2nodeConfigurationUpdateMessage, nodebInfo *entities.NodebInfo) {
	updateList := e2nodeConfig.ExtractConfigUpdateList()
	if nodebInfo.GetNodeType() == entities.Node_GNB {
		for i := 0; i < len(updateList); i++ {
			u := updateList[i]
			if nodebInfo.GetNodeType() == entities.Node_GNB {
				for j := 0; j < len(nodebInfo.GetGnb().NodeConfigs); j++ {
					if e.compareConfigIDs(u, *nodebInfo.GetGnb().NodeConfigs[j]) {
						e.logger.Debugf("#E2nodeConfigUpdateNotificationHandler.Handle - item at position [%d] should be updated", i)
						nodebInfo.GetGnb().NodeConfigs[i] = &u
						break
					} else {
						e.logger.Debugf("#E2nodeConfigUpdateNotificationHandler.Handle - dint match")
					}
				}
			}
		}
	} else {
		for i := 0; i < len(updateList); i++ {
			u := updateList[i]
			if nodebInfo.GetNodeType() == entities.Node_ENB {
				for j := 0; j < len(nodebInfo.GetEnb().NodeConfigs); j++ {
					v := nodebInfo.GetEnb().NodeConfigs[j]
					if e.compareConfigIDs(u, *v) {
						e.logger.Debugf("#E2nodeConfigUpdateNotificationHandler.Handle - item at position [%d] should be updated", i)
						nodebInfo.GetEnb().NodeConfigs[i] = &u
						break
					}
				}
			}

		}
	}
}

func (e *E2nodeConfigUpdateNotificationHandler) parseE2NodeConfigurationUpdate(payload []byte) (*models.E2nodeConfigurationUpdateMessage, error) {
	e2nodeConfig := models.E2nodeConfigurationUpdateMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload), &(e2nodeConfig.E2APPDU))

	if err != nil {
		e.logger.Errorf("#E2nodeConfigUpdateNotificationHandler.Handle - error in parsing request message: %+v", err)
		return nil, err
	}
	e.logger.Debugf("#E2nodeConfigUpdateNotificationHandler.Handle - Unmarshalling is successful %v", e2nodeConfig.E2APPDU.InitiatingMessage.ProcedureCode)
	return &e2nodeConfig, nil
}
