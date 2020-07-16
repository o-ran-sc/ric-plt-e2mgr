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
	"encoding/xml"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"strconv"
	"strings"
)

var (
	emptyTagsToReplaceToSelfClosingTags = []string{"reject", "ignore", "transport-resource-unavailable", "om-intervention",
		"v60s", "v20s", "v10s", "v5s", "v2s", "v1s"}
)

type E2SetupRequestNotificationHandler struct {
	logger                        *logger.Logger
	config                        *configuration.Configuration
	e2tInstancesManager           managers.IE2TInstancesManager
	rmrSender                     *rmrsender.RmrSender
	rNibDataService               services.RNibDataService
	e2tAssociationManager         *managers.E2TAssociationManager
	ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager
	ranListManager managers.RanListManager
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
		ranListManager: ranListManager,
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

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - got general configuration from rnib - enableRic: %t", generalConfiguration.EnableRic)

	if !generalConfiguration.EnableRic {
		cause := models.Cause{Misc: &models.CauseMisc{OmIntervention: &struct{}{}}}
		h.handleUnsuccessfulResponse(ranName, request, cause)
		return
	}

	setupRequest, e2tIpAddress, err := h.parseSetupRequest(request.Payload)
	if err != nil {
		h.logger.Errorf(err.Error())
		return
	}

	h.logger.Infof("#E2SetupRequestNotificationHandler.Handle - E2T Address: %s - handling E2_SETUP_REQUEST", e2tIpAddress)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.Handle - E2_SETUP_REQUEST has been parsed successfully %+v", setupRequest)

	_, err = h.e2tInstancesManager.GetE2TInstance(e2tIpAddress)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
		return
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(ranName)

	if err != nil {

		if _, ok := err.(*common.ResourceNotFoundError); !ok {
			h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to retrieve nodebInfo entity. Error: %s", ranName, err)
			return

		}

		if nodebInfo, err = h.handleNewRan(ranName, e2tIpAddress, setupRequest); err != nil {
			return
		}

	} else {
		if err = h.handleExistingRan(ranName, nodebInfo, setupRequest); err != nil {
			return
		}
	}

	err = h.e2tAssociationManager.AssociateRan(e2tIpAddress, nodebInfo)

	if err != nil {

		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s - failed to associate E2T to nodeB entity. Error: %s", ranName, err)
		if _, ok := err.(*e2managererrors.RoutingManagerError); ok {
			cause := models.Cause{Transport: &models.CauseTransport{TransportResourceUnavailable: &struct{}{}}}
			h.handleUnsuccessfulResponse(nodebInfo.RanName, request, cause)
		}
		return
	}

	h.handleSuccessfulResponse(ranName, request, setupRequest)
}

func (h *E2SetupRequestNotificationHandler) handleNewRan(ranName string, e2tIpAddress string, setupRequest *models.E2SetupRequestMessage) (*entities.NodebInfo, error) {

	nodebInfo := h.buildNodebInfo(ranName, e2tIpAddress, setupRequest)
	err := h.rNibDataService.SaveNodeb(nodebInfo)

	if err != nil {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.handleNewRan - RAN name: %s - failed saving nodebInfo. Error: %s", ranName, err)
		return nil, err
	}

	nbIdentity := h.buildNbIdentity(ranName, setupRequest)

	err = h.ranListManager.AddNbIdentity(entities.Node_GNB, nbIdentity)

	if err != nil {
		return nil, err
	}

	return nodebInfo, nil
}

func (h *E2SetupRequestNotificationHandler) setGnbFunctions(nodebInfo *entities.NodebInfo, setupRequest *models.E2SetupRequestMessage) {
	ranFunctions := setupRequest.ExtractRanFunctionsList()

	if ranFunctions != nil {
		nodebInfo.GetGnb().RanFunctions = ranFunctions
	}
}

func (h *E2SetupRequestNotificationHandler) handleExistingRan(ranName string, nodebInfo *entities.NodebInfo, setupRequest *models.E2SetupRequestMessage) error {
	if nodebInfo.GetConnectionStatus() == entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Errorf("#E2SetupRequestNotificationHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in incorrect state", ranName, nodebInfo.ConnectionStatus)
		return errors.New("nodeB entity in incorrect state")
	}

	h.setGnbFunctions(nodebInfo, setupRequest)

	return h.rNibDataService.UpdateNodebInfo(nodebInfo)
}

func (h *E2SetupRequestNotificationHandler) handleUnsuccessfulResponse(ranName string, req *models.NotificationRequest, cause models.Cause) {
	failureResponse := models.NewE2SetupFailureResponseMessage(models.TimeToWaitEnum.V60s, cause)
	h.logger.Debugf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - E2_SETUP_RESPONSE has been built successfully %+v", failureResponse)

	responsePayload, err := xml.Marshal(&failureResponse.E2APPDU)
	if err != nil {
		h.logger.Warnf("#E2SetupRequestNotificationHandler.handleUnsuccessfulResponse - RAN name: %s - Error marshalling RIC_E2_SETUP_RESP. Payload: %s", ranName, responsePayload)
	}

	responsePayload = replaceEmptyTagsWithSelfClosing(responsePayload)

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

	responsePayload = replaceEmptyTagsWithSelfClosing(responsePayload)

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

func replaceEmptyTagsWithSelfClosing(responsePayload []byte) []byte {

	emptyTagVsSelfClosingTagPairs := make([]string, len(emptyTagsToReplaceToSelfClosingTags)*2)

	j := 0

	for i := 0; i < len(emptyTagsToReplaceToSelfClosingTags); i++ {
		emptyTagVsSelfClosingTagPairs[j] = fmt.Sprintf("<%[1]s></%[1]s>", emptyTagsToReplaceToSelfClosingTags[i])
		emptyTagVsSelfClosingTagPairs[j+1] = fmt.Sprintf("<%s/>", emptyTagsToReplaceToSelfClosingTags[i])
		j += 2
	}
	responseString := strings.NewReplacer(emptyTagVsSelfClosingTagPairs...).Replace(string(responsePayload))
	return []byte(responseString)
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
	err := xml.Unmarshal(normalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)
	if err != nil {
		return nil, "", errors.New(fmt.Sprintf("#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: %x", payload))
	}

	return setupRequest, e2tIpAddress, nil
}

func normalizeXml(payload []byte) []byte {
	xmlStr := string(payload)
	normalized := strings.NewReplacer("&lt;", "<", "&gt;", ">").Replace(xmlStr)
	return []byte(normalized)
}

func (h *E2SetupRequestNotificationHandler) buildNodebInfo(ranName string, e2tAddress string, request *models.E2SetupRequestMessage) *entities.NodebInfo {
	nodebInfo := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tAddress,
		RanName:                      ranName,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
		GlobalNbId:                   h.buildGlobalNbId(request),
	}

	h.setGnbFunctions(nodebInfo, request)
	return nodebInfo
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