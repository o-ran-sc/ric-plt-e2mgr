//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
// Copyright 2023 Capgemini
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
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

const cleanUpDurationNanoSec uint64 = 10000000000 // cleanUpDuration = 10sec (value in nanoSecond=10000000000)

var (
	emptyTagsToReplaceToSelfClosingTags = []string{"reject", "ignore", "transport-resource-unavailable", "om-intervention", "request-id-unknown",
		"unspecified", "message-not-compatible-with-receiver-state", "control-processing-overload",
		"v60s", "v20s", "v10s", "v5s", "v2s", "v1s", "ng", "xn", "e1", "f1", "w1", "s1", "x2", "success", "failure"}
	gnbTypesMap = map[string]entities.GnbType{
		"gnb":    entities.GnbType_GNB,
		"en_gnb": entities.GnbType_EN_GNB,
	}
	enbTypesMap = map[string]entities.EnbType{
		"enB_macro":         entities.EnbType_MACRO_ENB,
		"enB_home":          entities.EnbType_HOME_ENB,
		"enB_shortmacro":    entities.EnbType_SHORT_MACRO_ENB,
		"enB_longmacro":     entities.EnbType_LONG_MACRO_ENB,
		"ng_enB_macro":      entities.EnbType_MACRO_NG_ENB,
		"ng_enB_shortmacro": entities.EnbType_SHORT_MACRO_NG_ENB,
		"ng_enB_longmacro":  entities.EnbType_LONG_MACRO_NG_ENB,
	}
)

type E2SetupRequestNotificationHandler struct {
	logger                        *logger.Logger
	config                        *configuration.Configuration
	e2tInstancesManager           managers.IE2TInstancesManager
	rmrSender                     *rmrsender.RmrSender
	rNibDataService               services.RNibDataService
	e2tAssociationManager         *managers.E2TAssociationManager
	ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager
	ranListManager                managers.RanListManager
}

func NewE2SetupRequestNotificationHandler(logger *logger.Logger, config *configuration.Configuration, e2tInstancesManager managers.IE2TInstancesManager, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService, e2tAssociationManager *managers.E2TAssociationManager, ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager, ranListManager managers.RanListManager) *E2SetupRequestNotificationHandler {
	return &E2SetupRequestNotificationHandler{
		logger:                        logger,
		config:                        config,
		e2tInstancesManager:           e2tInstancesManager,
		rmrSender:                     rmrSender,
		rNibDataService:               rNibDataService,
		e2tAssociationManager:         e2tAssociationManager,
		ranConnectStatusChangeManager: ranConnectStatusChangeManager,
		ranListManager:                ranListManager,
	}
}

func (h *E2SetupRequestNotificationHandler) Handle(request *models.NotificationRequest) {
	ranName := request.RanName
	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - received E2_SETUP_REQUEST. Payload: %x", ranName, request.Payload)

	generalConfiguration, err := h.rNibDataService.GetGeneralConfiguration()

	if err != nil {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - Failed retrieving e2m general configuration - quitting e2 setup flow. error: %s", err)
		return
	}

	setupRequest, e2tIpAddress, err := h.parseSetupRequest(request.Payload)
	if err != nil {
		h.logger.Errorf(err.Error())
		return
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - E2T Address: %s - handling E2_SETUP_REQUEST", e2tIpAddress)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.Handle - E2_SETUP_REQUEST has been parsed successfully %+v", setupRequest)

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - got general configuration from rnib - enableRic: %t", generalConfiguration.EnableRic)

	if !generalConfiguration.EnableRic {
		cause := models.Cause{Misc: &models.CauseMisc{OmIntervention: &struct{}{}}}
		h.handleUnsuccessfulResponse(ranName, request, cause, setupRequest)
		return
	}

	_, err = h.e2tInstancesManager.GetE2TInstance(e2tIpAddress)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
		return
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)

	var functionsModified bool

	if err != nil {

		if _, ok := err.(*common.ResourceNotFoundError); !ok {
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to retrieve nodebInfo entity. Error: %s", ranName, err)
			return
		}

		if nodebInfo, err = h.handleNewRan(ranName, e2tIpAddress, setupRequest); err != nil {
			if _, ok := err.(*e2managererrors.UnknownSetupRequestRanNameError); ok {
				cause := models.Cause{RicRequest: &models.CauseRic{RequestIdUnknown: &struct{}{}}}
				h.handleUnsuccessfulResponse(ranName, request, cause, setupRequest)
			}
			return
		}

	} else {

		functionsModified, err = h.handleExistingRan(ranName, nodebInfo, setupRequest)

		if err != nil {
			h.fillCauseAndSendUnsuccessfulResponse(nodebInfo, request, setupRequest)
			return
		}
	}

	ranStatusChangePublished, err := h.e2tAssociationManager.AssociateRan(e2tIpAddress, nodebInfo)

	if err != nil {

		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to associate E2T to nodeB entity. Error: %s", ranName, err)
		if _, ok := err.(*e2managererrors.RoutingManagerError); ok {

			if err = h.handleUpdateAndPublishNodebInfo(functionsModified, ranStatusChangePublished, nodebInfo); err != nil {
				return
			}

			cause := models.Cause{Transport: &models.CauseTransport{TransportResourceUnavailable: &struct{}{}}}
			h.handleUnsuccessfulResponse(nodebInfo.RanName, request, cause, setupRequest)
		}
		return
	}

	if err = h.handleUpdateAndPublishNodebInfo(functionsModified, ranStatusChangePublished, nodebInfo); err != nil {
		return
	}

	h.handleSuccessfulResponse(ranName, request, setupRequest)
}

func (h *E2SetupRequestNotificationHandler) handleUpdateAndPublishNodebInfo(functionsModified bool, ranStatusChangePublished bool, nodebInfo *entities.NodebInfo) error {

	if ranStatusChangePublished || !functionsModified {
		return nil
	}

	err := h.rNibDataService.UpdateNodebInfoAndPublish(nodebInfo)

	if err != nil {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.handleUpdateAndPublishNodebInfo - RAN name: %s - Failed at UpdateNodebInfoAndPublish. error: %s", nodebInfo.RanName, err)
		return err
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.handleUpdateAndPublishNodebInfo - RAN name: %s - Successfully executed UpdateNodebInfoAndPublish", nodebInfo.RanName)
	return nil

}

func (h *E2SetupRequestNotificationHandler) handleNewRan(ranName string, e2tIpAddress string, setupRequest *models.E2SetupRequestMessage) (*entities.NodebInfo, error) {

	nodebInfo, err := h.buildNodebInfo(ranName, e2tIpAddress, setupRequest)
	if err != nil {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.handleNewRan - RAN name: %s - failed building nodebInfo. Error: %s", ranName, err)
		return nil, err
	}

	err = h.rNibDataService.SaveNodeb(nodebInfo)
	if err != nil {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.handleNewRan - RAN name: %s - failed saving nodebInfo. Error: %s", ranName, err)
		return nil, err
	}

	nbIdentity := h.buildNbIdentity(ranName, setupRequest)

	err = h.ranListManager.AddNbIdentity(nodebInfo.GetNodeType(), nbIdentity)

	if err != nil {
		return nil, err
	}

	return nodebInfo, nil
}

func (h *E2SetupRequestNotificationHandler) handleExistingRan(ranName string, nodebInfo *entities.NodebInfo, setupRequest *models.E2SetupRequestMessage) (bool, error) {
	if nodebInfo.GetConnectionStatus() == entities.ConnectionStatus_DISCONNECTED {
		delta_in_nano := uint64(time.Now().UnixNano()) - nodebInfo.StatusUpdateTimeStamp
		//The duration from last Disconnection for which a new request is to be rejected (currently 10 sec)
		if delta_in_nano < cleanUpDurationNanoSec {
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s, connection status: %s - nodeB entity disconnection in progress", ranName, nodebInfo.ConnectionStatus)
			return false, errors.New("nodeB entity disconnection in progress")
		}
		h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in disconnected state", ranName, nodebInfo.ConnectionStatus)
	} else if nodebInfo.GetConnectionStatus() == entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in incorrect state", ranName, nodebInfo.ConnectionStatus)
		return false, errors.New("nodeB entity in incorrect state")
	}

	nodebInfo.SetupFromNetwork = true

	e2NodeConfig := setupRequest.ExtractE2NodeConfigList()
	if e2NodeConfig == nil {
		return false, errors.New("Empty E2nodeComponentConfigAddition-List")
	}

	if nodebInfo.NodeType == entities.Node_ENB {
		if len(e2NodeConfig) == 0 && len(nodebInfo.GetEnb().GetNodeConfigs()) == 0 {
			return false, errors.New("Empty E2nodeComponentConfigAddition-List")
		}
		nodebInfo.GetEnb().NodeConfigs = e2NodeConfig

		return false, nil
	}

	if len(e2NodeConfig) == 0 && len(nodebInfo.GetGnb().GetNodeConfigs()) == 0 {
		return false, errors.New("Empty E2nodeComponentConfigAddition-List")
	}
	nodebInfo.GetGnb().NodeConfigs = e2NodeConfig

	setupMessageRanFuncs := setupRequest.ExtractRanFunctionsList()

	if setupMessageRanFuncs == nil || (len(setupMessageRanFuncs) == 0 && len(nodebInfo.GetGnb().RanFunctions) == 0) {
		return false, nil
	}

	nodebInfo.GetGnb().RanFunctions = setupMessageRanFuncs
	return true, nil
}

func (h *E2SetupRequestNotificationHandler) handleUnsuccessfulResponse(ranName string, req *models.NotificationRequest, cause models.Cause, setupRequest *models.E2SetupRequestMessage) {
	failureResponse := models.NewE2SetupFailureResponseMessage(models.TimeToWaitEnum.V60s, cause, setupRequest)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - E2_SETUP_RESPONSE has been built successfully %+v", failureResponse)

	responsePayload, err := xml.Marshal(&failureResponse.E2APPDU)
	if err != nil {
		h.logger.Warnf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_SETUP_RESP. Payload: %s", ranName, responsePayload)
	}

	responsePayload = utils.ReplaceEmptyTagsWithSelfClosing(responsePayload, emptyTagsToReplaceToSelfClosingTags)

	h.logger.Infof("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - payload: %s", responsePayload)
	msg := models.NewRmrMessage(rmrCgo.RIC_E2_SETUP_FAILURE, ranName, responsePayload, req.TransactionId, req.GetMsgSrc())
	h.logger.Infof("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - RAN name: %s - RIC_E2_SETUP_RESP message has been built successfully. Message: %x", ranName, msg)
	_ = h.rmrSender.WhSend(msg)

}

func (h *E2SetupRequestNotificationHandler) handleSuccessfulResponse(ranName string, req *models.NotificationRequest, setupRequest *models.E2SetupRequestMessage) {

	plmnId := buildPlmnId(h.config.GlobalRicId.Mcc, h.config.GlobalRicId.Mnc)

	ricNearRtId, err := convertTo20BitString(h.config.GlobalRicId.RicId)
	if err != nil {
		return
	}
	successResponse := models.NewE2SetupSuccessResponseMessage(plmnId, ricNearRtId, setupRequest)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - E2_SETUP_RESPONSE has been built successfully %+v", successResponse)

	responsePayload, err := xml.Marshal(&successResponse.E2APPDU)
	if err != nil {
		h.logger.Warnf("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_SETUP_RESP. Payload: %s", ranName, responsePayload)
	}

	responsePayload = utils.ReplaceEmptyTagsWithSelfClosing(responsePayload, emptyTagsToReplaceToSelfClosingTags)

	h.logger.Infof("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - payload: %s", responsePayload)

	msg := models.NewRmrMessage(rmrCgo.RIC_E2_SETUP_RESP, ranName, responsePayload, req.TransactionId, req.GetMsgSrc())
	h.logger.Infof("#E2SetupRequestNotificationHandler.handleSuccessfulResponse - RAN name: %s - RIC_E2_SETUP_RESP message has been built successfully. Message: %x", ranName, msg)
	_ = h.rmrSender.Send(msg)
}

func buildPlmnId(mmc string, mnc string) string {
	var b strings.Builder

	b.WriteByte(mmc[1])
	b.WriteByte(mmc[0])
	if len(mnc) == 2 {
		b.WriteString("F")
	} else {
		b.WriteByte(mnc[2])
	}
	b.WriteByte(mmc[2])
	b.WriteByte(mnc[1])
	b.WriteByte(mnc[0])

	return b.String()
}

func convertTo20BitString(ricNearRtId string) (string, error) {
	r, err := strconv.ParseUint(ricNearRtId, 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%020b", r)[:20], nil
}

func (h *E2SetupRequestNotificationHandler) parseSetupRequest(payload []byte) (*models.E2SetupRequestMessage, string, error) {

	pipInd := bytes.IndexByte(payload, '|')
	if pipInd < 0 {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
	}

	e2tIpAddress := string(payload[:pipInd])
	if len(e2tIpAddress) == 0 {
		return nil, "", errors.New("#E2SetupRequestNotificationHandler.parseSetupRequest - Empty E2T Address received")
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.parseSetupRequest - payload: %s", payload[pipInd+1:])

	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: %x", payload))
	}

	return setupRequest, e2tIpAddress, nil
}

func (h *E2SetupRequestNotificationHandler) buildNodebInfo(ranName string, e2tAddress string, request *models.E2SetupRequestMessage) (*entities.NodebInfo, error) {
	nodebInfo := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tAddress,
		RanName:                      ranName,
		GlobalNbId:                   h.buildGlobalNbId(request),
		SetupFromNetwork:             true,
	}
	err := h.setNodeTypeAndConfiguration(nodebInfo)
	if err != nil {
		return nil, err
	}

	e2NodeConfig := request.ExtractE2NodeConfigList()
	if e2NodeConfig == nil {
		return nil, errors.New("Empty E2nodeComponentConfigAddition-List")
	}

	if nodebInfo.NodeType == entities.Node_ENB {
		if len(e2NodeConfig) == 0 && len(nodebInfo.GetEnb().GetNodeConfigs()) == 0 {
			return nil, errors.New("Empty E2nodeComponentConfigAddition-List")
		}
		nodebInfo.GetEnb().NodeConfigs = e2NodeConfig

		return nodebInfo, nil
	}

	if len(e2NodeConfig) == 0 && len(nodebInfo.GetGnb().GetNodeConfigs()) == 0 {
		return nil, errors.New("Empty E2nodeComponentConfigAddition-List")
	}
	nodebInfo.GetGnb().NodeConfigs = e2NodeConfig

	if nodebInfo.NodeType == entities.Node_GNB {
		h.logger.Debugf("#E2SetupRequestNotificationHandler buildNodebInfo - entities.Node_GNB %d", entities.Node_GNB)

		gnbNodetype := h.setGnbNodeType(request)
		h.logger.Debugf("#E2SetupRequestNotificationHandler buildNodebInfo -gnbNodetype %s", gnbNodetype)
		nodebInfo.GnbNodeType = gnbNodetype
		nodebInfo.CuUpId = request.GetCuupId()
		nodebInfo.DuId = request.GetDuId()
		h.logger.Debugf("#E2SetupRequestNotificationHandler buildNodebInfo -cuupid%s", request.GetCuupId())
		h.logger.Debugf("#E2SetupRequestNotificationHandler buildNodebInfo -duid %s", request.GetDuId())
	}

	ranFuncs := request.ExtractRanFunctionsList()

	if ranFuncs != nil {
		nodebInfo.GetGnb().RanFunctions = ranFuncs
	}

	return nodebInfo, nil
}


func (h *E2SetupRequestNotificationHandler) setGnbNodeType(setupRequest *models.E2SetupRequestMessage) string {
	gnbNodetype := "gNB"
	if setupRequest.GetCuupId() != "" && setupRequest.GetDuId() != "" {
		gnbNodetype = "gNB"
	} else if setupRequest.GetCuupId() != "" {
		gnbNodetype = "gNB_CU_UP"
	} else if setupRequest.GetDuId() != "" {
		gnbNodetype = "gNB_DU"
	}
	return gnbNodetype
}

func (h *E2SetupRequestNotificationHandler) setNodeTypeAndConfiguration(nodebInfo *entities.NodebInfo) error {
	for k, v := range gnbTypesMap {
		if strings.HasPrefix(nodebInfo.RanName, k) {
			nodebInfo.NodeType = entities.Node_GNB
			nodebInfo.Configuration = &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{GnbType: v}}
			return nil
		}
	}
	for k, v := range enbTypesMap {
		if strings.HasPrefix(nodebInfo.RanName, k) {
			nodebInfo.NodeType = entities.Node_ENB
			nodebInfo.Configuration = &entities.NodebInfo_Enb{Enb: &entities.Enb{EnbType: v}}
			return nil
		}
	}

	return e2managererrors.NewUnknownSetupRequestRanNameError(nodebInfo.RanName)
}

func (h *E2SetupRequestNotificationHandler) buildGlobalNbId(setupRequest *models.E2SetupRequestMessage) *entities.GlobalNbId {
	return &entities.GlobalNbId{
		PlmnId: setupRequest.GetPlmnId(),
		NbId:   setupRequest.GetNbId(),
	}
}

func (h *E2SetupRequestNotificationHandler) buildNbIdentity(ranName string, setupRequest *models.E2SetupRequestMessage) *entities.NbIdentity {
	return &entities.NbIdentity{
		InventoryName: ranName,
		GlobalNbId:    h.buildGlobalNbId(setupRequest),
	}
}

func (h *E2SetupRequestNotificationHandler) fillCauseAndSendUnsuccessfulResponse(nodebInfo *entities.NodebInfo, request *models.NotificationRequest, setupRequest *models.E2SetupRequestMessage) {
	if nodebInfo.GetConnectionStatus() == entities.ConnectionStatus_DISCONNECTED {
		cause := models.Cause{Misc: &models.CauseMisc{ControlProcessingOverload: &struct{}{}}}
		h.handleUnsuccessfulResponse(nodebInfo.RanName, request, cause, setupRequest)
	}
}
