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
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"encoding/xml"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	e2tInstanceFullAddress                   = "10.0.2.15:9999"
	e2SetupMsgPrefix                         = e2tInstanceFullAddress + "|"
	gnbNodebRanName                          = "gnb:310-410-b5c67788"
	enbNodebRanName                          = "enB_macro:P310-410-b5c67788"
	GnbSetupRequestXmlPath                   = "../../tests/resources/setupRequest_gnb.xml"
	EnGnbSetupRequestXmlPath                 = "../../tests/resources/setupRequest_en-gNB.xml"
	NgEnbSetupRequestXmlPath                 = "../../tests/resources/setupRequest_ng-eNB.xml"
	EnbSetupRequestXmlPath                   = "../../tests/resources/setupRequest_enb.xml"
	GnbWithoutFunctionsSetupRequestXmlPath   = "../../tests/resources/setupRequest_gnb_without_functions.xml"
	E2SetupFailureResponseWithMiscCause      = "<E2AP-PDU><unsuccessfulOutcome><procedureCode>1</procedureCode><criticality><reject/></criticality><value><E2setupFailure><protocolIEs><E2setupFailureIEs><id>1</id><criticality><ignore/></criticality><value><Cause><misc><om-intervention/></misc></Cause></value></E2setupFailureIEs><E2setupFailureIEs><id>31</id><criticality><ignore/></criticality><value><TimeToWait><v60s/></TimeToWait></value></E2setupFailureIEs></protocolIEs></E2setupFailure></value></unsuccessfulOutcome></E2AP-PDU>"
	E2SetupFailureResponseWithTransportCause = "<E2AP-PDU><unsuccessfulOutcome><procedureCode>1</procedureCode><criticality><reject/></criticality><value><E2setupFailure><protocolIEs><E2setupFailureIEs><id>1</id><criticality><ignore/></criticality><value><Cause><transport><transport-resource-unavailable/></transport></Cause></value></E2setupFailureIEs><E2setupFailureIEs><id>31</id><criticality><ignore/></criticality><value><TimeToWait><v60s/></TimeToWait></value></E2setupFailureIEs></protocolIEs></E2setupFailure></value></unsuccessfulOutcome></E2AP-PDU>"
	E2SetupFailureResponseWithRicCause       = "<E2AP-PDU><unsuccessfulOutcome><procedureCode>1</procedureCode><criticality><reject/></criticality><value><E2setupFailure><protocolIEs><E2setupFailureIEs><id>1</id><criticality><ignore/></criticality><value><Cause><ricRequest><request-id-unknown/></ricRequest></Cause></value></E2setupFailureIEs><E2setupFailureIEs><id>31</id><criticality><ignore/></criticality><value><TimeToWait><v60s/></TimeToWait></value></E2setupFailureIEs></protocolIEs></E2setupFailure></value></unsuccessfulOutcome></E2AP-PDU>"
	StateChangeMessageChannel                = "RAN_CONNECTION_STATUS_CHANGE"
)

func initMocks(t *testing.T) (*E2SetupRequestNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{
		RnibRetryIntervalMs:       10,
		MaxRnibConnectionAttempts: 3,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel: StateChangeMessageChannel,
		},
		GlobalRicId: struct {
			RicId string
			Mcc   string
			Mnc   string
		}{Mcc: "327", Mnc: "94", RicId: "AACCE"}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock, ranConnectStatusChangeManager)
	handler := NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManagerMock, rmrSender, rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager)
	return handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock
}

func getMbuf(ranName string, msgType int, payloadStr string, request *models.NotificationRequest) *rmrCgo.MBuf {
	payload := []byte(payloadStr)
	mbuf := rmrCgo.NewMBuf(msgType, len(payload), ranName, &payload, &request.TransactionId, request.GetMsgSrc())
	return mbuf
}

func readXmlFile(t *testing.T, xmlPath string) []byte {
	path, err := filepath.Abs(xmlPath)
	if err != nil {
		t.Fatal(err)
	}
	xmlAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return xmlAsBytes
}

func TestParseGnbSetupRequest_Success(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(e2SetupMsgPrefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Equal(t, "02F829", request.GetPlmnId())
	assert.Equal(t, "001100000011000000110000", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseEnGnbSetupRequest_Success(t *testing.T) {
	enGnbXml := readXmlFile(t, EnGnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(e2SetupMsgPrefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, enGnbXml...))
	assert.Equal(t, "131014", request.GetPlmnId())
	assert.Equal(t, "11000101110001101100011111111000", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseNgEnbSetupRequest_Success(t *testing.T) {
	ngEnbXml := readXmlFile(t, NgEnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(e2SetupMsgPrefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, ngEnbXml...))
	assert.Equal(t, "131014", request.GetPlmnId())
	assert.Equal(t, "101010101010101010", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseEnbSetupRequest_Success(t *testing.T) {
	enbXml := readXmlFile(t, EnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(e2SetupMsgPrefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, enbXml...))
	assert.Equal(t, "6359AB", request.GetPlmnId())
	assert.Equal(t, "101010101010101010", request.GetNbId())
	assert.Nil(t, err)
}

func TestParseSetupRequest_PipFailure(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte("10.0.2.15:9999")
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
}

func TestParseSetupRequest_UnmarshalFailure(t *testing.T) {
	handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(e2SetupMsgPrefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, 1, 2, 3))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: 31302e302e322e31353a393939397c010203")
}

func TestE2SetupRequestNotificationHandler_HandleParseError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte("invalid_prefix"), xmlGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleUnmarshalError(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), "xmlGnb"...)}
	handler.Handle(notificationRequest)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_ConvertTo20BitStringError(t *testing.T) {
	xmlEnGnb := readXmlFile(t, EnGnbSetupRequestXmlPath)
	logger := tests.InitLog(t)
	config := &configuration.Configuration{
		RnibRetryIntervalMs:       10,
		MaxRnibConnectionAttempts: 3,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel: StateChangeMessageChannel,
		},
		GlobalRicId: struct {
			RicId string
			Mcc   string
			Mnc   string
		}{Mcc: "327", Mnc: "94", RicId: "10011001101010101011"}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)

	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock, ranConnectStatusChangeManager)
	handler := NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManagerMock, rmrSender, rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	writerMock.On("AddNbIdentity", entities.Node_GNB, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	prefBytes := []byte(e2SetupMsgPrefix)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append(prefBytes, xmlEnGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestExtractionOfNodeTypeFromRanName(t *testing.T) {
	handler, _, _, _, _, _ := initMocks(t)
	validRanNames := []string{"gnb_P310_410_b5c67788", "en_gnb_P310_410_b5c67788", "ng_enB_macro_P310_410_b5c67788", "ng_enB_shortmacro_P310_410_b5c67788", "ng_enB_longmacro_P310_410_b5c67788", "enB_macro_P310_410_b5c67788", "enB_home_P310_410_b5c67788", "enB_shortmacro_P310_410_b5c67788", "enB_longmacro_P310_410_b5c67788"}
	for _, v := range validRanNames {
		nodeb := &entities.NodebInfo{RanName: v}
		err := handler.setNodeTypeAndConfiguration(nodeb)
		assert.Nil(t, err)
	}
	inValidRanNames := []string{"P310_410_b5c67788", "blabla_P310_410_b5c67788", "ng_enB-macro_P310_410_b5c67788", "ng_enb_shortmacro_P310_410_b5c67788", "ng_enB-longmacro_P310_410_b5c67788", "enB_new_macro_P310_410_b5c67788"}
	for _, v := range inValidRanNames {
		nodeb := &entities.NodebInfo{RanName: v}
		err := handler.setNodeTypeAndConfiguration(nodeb)
		assert.NotNil(t, err, v)
	}
}

func TestE2SetupRequestNotificationHandler_GetGeneralConfigurationFailure(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{}, common.NewInternalError(errors.New("some error")))
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	handler.Handle(notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance")
	readerMock.AssertNotCalled(t, "GetNodeb")
	writerMock.AssertNotCalled(t, "SaveNodeb")
}

func TestE2SetupRequestNotificationHandler_EnableRicFalse(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: false}, nil)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	mbuf := getMbuf(gnbNodebRanName, rmrCgo.RIC_E2_SETUP_FAILURE, E2SetupFailureResponseWithMiscCause, notificationRequest)
	rmrMessengerMock.On("WhSendMsg", mbuf, true).Return(&rmrCgo.MBuf{}, nil)
	handler.Handle(notificationRequest)
	rmrMessengerMock.AssertCalled(t, "WhSendMsg", mbuf, true)
	e2tInstancesManagerMock.AssertNotCalled(t, "GetE2TInstance")
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance")
	readerMock.AssertNotCalled(t, "GetNodeb")
	writerMock.AssertNotCalled(t, "SaveNodeb")
}

func TestE2SetupRequestNotificationHandler_HandleGetE2TInstanceError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, common.NewResourceNotFoundError("Not found"))
	prefBytes := []byte(e2SetupMsgPrefix)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	readerMock.AssertNotCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleGetNodebError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, routingManagerClientMock, e2tInstancesManagerMock, rmrMessengerMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewInternalError(errors.New("some error")))
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "GetE2TInstance", e2tInstanceFullAddress)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertNotCalled(t, "SaveNodeb", mock.Anything, mock.Anything)
	routingManagerClientMock.AssertNotCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo", mock.Anything)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance", mock.Anything, mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

/* New Ran UTs - BEGIN */

func TestE2SetupRequestNotificationHandler_HandleNewRanSaveNodebFailure(t *testing.T) {
	xml := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, _, e2tInstancesManagerMock, _ := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", gnbNodebRanName).Return(gnb, common.NewResourceNotFoundError("Not found"))
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xml...)}
	nodebInfo := getExpectedGnbNodebForNewRan(notificationRequest.Payload)
	writerMock.On("SaveNodeb", nodebInfo).Return(common.NewInternalError(errors.New("error")))
	handler.Handle(notificationRequest)

	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewRan_invalidRanName(t *testing.T) {
	xml := readXmlFile(t, EnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, _ := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var enb *entities.NodebInfo
	invalidEnbRanName := "enB-macro:P310-410-b5c67788"
	readerMock.On("GetNodeb", invalidEnbRanName).Return(enb, common.NewResourceNotFoundError("Not found"))
	notificationRequest := &models.NotificationRequest{RanName: invalidEnbRanName, Payload: append([]byte(e2SetupMsgPrefix), xml...)}
	mbuf := getMbuf(invalidEnbRanName, rmrCgo.RIC_E2_SETUP_FAILURE, E2SetupFailureResponseWithRicCause, notificationRequest)
	rmrMessengerMock.On("WhSendMsg", mbuf, true).Return(&rmrCgo.MBuf{}, nil)
	handler.Handle(notificationRequest)

	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "WhSendMsg", mbuf, true)
	e2tInstancesManagerMock.AssertExpectations(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewRanAddNbIdentityFailure(t *testing.T) {
	xml := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, _, e2tInstancesManagerMock, _ := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", gnbNodebRanName).Return(gnb, common.NewResourceNotFoundError("Not found"))
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xml...)}
	nodebInfo := getExpectedGnbNodebForNewRan(notificationRequest.Payload)
	writerMock.On("SaveNodeb", nodebInfo).Return(nil)
	nbIdentity := &entities.NbIdentity{InventoryName: gnbNodebRanName, GlobalNbId: nodebInfo.GlobalNbId}
	writerMock.On("AddNbIdentity", entities.Node_GNB, nbIdentity).Return(common.NewInternalError(errors.New("error")))
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewRanRoutingManagerError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)

	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	readerMock.On("GetNodeb", gnbNodebRanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError("Not found"))
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	nodebInfo := getExpectedGnbNodebForNewRan(notificationRequest.Payload)
	writerMock.On("SaveNodeb", nodebInfo).Return(nil)
	nbIdentity := &entities.NbIdentity{InventoryName: gnbNodebRanName, GlobalNbId: nodebInfo.GlobalNbId}
	writerMock.On("AddNbIdentity", entities.Node_GNB, nbIdentity).Return(nil)
	updatedNodebInfo := *nodebInfo
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(errors.New("association error"))
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(nil)
	var errEmpty error
	mbuf := getMbuf(gnbNodebRanName, rmrCgo.RIC_E2_SETUP_FAILURE, E2SetupFailureResponseWithTransportCause, notificationRequest)
	rmrMessengerMock.On("WhSendMsg", mbuf, true).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	routingManagerClientMock.AssertCalled(t, "AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "WhSendMsg", mbuf, true)
}

func testE2SetupRequestNotificationHandler_HandleNewRanSuccess(t *testing.T, xmlPath string, nodeType entities.Node_Type) {
	xml := readXmlFile(t, xmlPath)
	var ranName string
	if nodeType == entities.Node_GNB {
		ranName = gnbNodebRanName
	} else {
		ranName = enbNodebRanName
	}
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(gnb, common.NewResourceNotFoundError("Not found"))
	notificationRequest := &models.NotificationRequest{RanName: ranName, Payload: append([]byte(e2SetupMsgPrefix), xml...)}
	var expectedNodebInfo *entities.NodebInfo
	if nodeType == entities.Node_GNB {
		expectedNodebInfo = getExpectedGnbNodebForNewRan(notificationRequest.Payload)
	} else {
		expectedNodebInfo = getExpectedEnbNodebForNewRan(notificationRequest.Payload)
	}
	writerMock.On("SaveNodeb", expectedNodebInfo).Return(nil)
	nbIdentity := &entities.NbIdentity{InventoryName: ranName, GlobalNbId: expectedNodebInfo.GlobalNbId}
	writerMock.On("AddNbIdentity", nodeType, nbIdentity).Return(nil)
	updatedNodebInfo := *expectedNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &updatedNodebInfo, ranName+"_CONNECTED").Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	updatedNodebInfo2 := *expectedNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = e2tInstanceFullAddress
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo2).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", e2tInstanceFullAddress, []string{ranName}).Return(nil)
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(&rmrCgo.MBuf{}, nil)

	handler.Handle(notificationRequest)

	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleNewRanSuccess(t, GnbSetupRequestXmlPath, entities.Node_GNB)
}

func TestE2SetupRequestNotificationHandler_HandleNewGnbWithoutFunctionsSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleNewRanSuccess(t, GnbWithoutFunctionsSetupRequestXmlPath, entities.Node_GNB)
}

func TestE2SetupRequestNotificationHandler_HandleNewEnGnbSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleNewRanSuccess(t, EnGnbSetupRequestXmlPath, entities.Node_GNB)
}

func TestE2SetupRequestNotificationHandler_HandleNewNgEnbSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleNewRanSuccess(t, NgEnbSetupRequestXmlPath, entities.Node_ENB)
}

func getExpectedGnbNodebForNewRan(payload []byte) *entities.NodebInfo {
	pipInd := bytes.IndexByte(payload, '|')
	setupRequest := &models.E2SetupRequestMessage{}
	_ = xml.Unmarshal(normalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)

	nodeb := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		RanName:                      gnbNodebRanName,
		SetupFromNetwork:             true,
		NodeType:                     entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{
			Gnb: &entities.Gnb{
				GnbType:      entities.GnbType_GNB,
				RanFunctions: setupRequest.ExtractRanFunctionsList(),
			},
		},
		GlobalNbId: &entities.GlobalNbId{
			PlmnId: setupRequest.GetPlmnId(),
			NbId:   setupRequest.GetNbId(),
		},
	}

	return nodeb
}

func getExpectedEnbNodebForNewRan(payload []byte) *entities.NodebInfo {
	pipInd := bytes.IndexByte(payload, '|')
	setupRequest := &models.E2SetupRequestMessage{}
	_ = xml.Unmarshal(normalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)

	nodeb := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		RanName:                      enbNodebRanName,
		SetupFromNetwork:             true,
		NodeType:                     entities.Node_ENB,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				EnbType: entities.EnbType_MACRO_ENB,
			},
		},
		GlobalNbId: &entities.GlobalNbId{
			PlmnId: setupRequest.GetPlmnId(),
			NbId:   setupRequest.GetNbId(),
		},
	}

	return nodeb
}

/* New Ran UTs - END */

/* Existing Ran UTs - BEGIN */

func getExpectedNodebForExistingRan(nodeb entities.NodebInfo, payload []byte) *entities.NodebInfo {
	pipInd := bytes.IndexByte(payload, '|')
	setupRequest := &models.E2SetupRequestMessage{}
	_ = xml.Unmarshal(normalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)

	if nodeb.NodeType == entities.Node_ENB {
		return &nodeb
	}

	if ranFuncs := setupRequest.ExtractRanFunctionsList(); ranFuncs != nil {
		nodeb.GetGnb().RanFunctions = ranFuncs
	}
	return &nodeb
}

func TestE2SetupRequestNotificationHandler_HandleExistingConnectedEnbSuccess(t *testing.T) {
	xmlEnb := readXmlFile(t, EnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      enbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		NodeType:                     entities.Node_ENB,
		Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{}},
	}

	readerMock.On("GetNodeb", enbNodebRanName).Return(nodebInfo, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)

	notificationRequest := &models.NotificationRequest{RanName: enbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlEnb...)}
	enbToUpdate := getExpectedNodebForExistingRan(*nodebInfo, notificationRequest.Payload)
	enbToUpdate.SetupFromNetwork = true
	writerMock.On("UpdateNodebInfo", enbToUpdate).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", e2tInstanceFullAddress, []string{enbNodebRanName}).Return(nil)
	var errEmpty error
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mock.Anything, true)
}

func TestE2SetupRequestNotificationHandler_HandleExistingDisconnectedEnbSuccess(t *testing.T) {
	xmlEnb := readXmlFile(t, EnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      enbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_DISCONNECTED,
		NodeType:                     entities.Node_ENB,
		Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{}},
	}

	readerMock.On("GetNodeb", enbNodebRanName).Return(nodebInfo, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)

	notificationRequest := &models.NotificationRequest{RanName: enbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlEnb...)}
	enbToUpdate := getExpectedNodebForExistingRan(*nodebInfo, notificationRequest.Payload)

	enbToUpdate2 := *enbToUpdate
	enbToUpdate2.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	enbToUpdate2.SetupFromNetwork = true
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &enbToUpdate2, enbNodebRanName+"_CONNECTED").Return(nil)
	enbToUpdate3 := *enbToUpdate
	enbToUpdate3.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	enbToUpdate3.SetupFromNetwork = true
	enbToUpdate3.AssociatedE2TInstanceAddress = e2tInstanceFullAddress
	writerMock.On("UpdateNodebInfo", &enbToUpdate3).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", e2tInstanceFullAddress, []string{enbNodebRanName}).Return(nil)
	var errEmpty error
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mock.Anything, true)
}

func testE2SetupRequestNotificationHandler_HandleExistingConnectedGnbSuccess(t *testing.T, withFunctions bool) {
	var xmlGnb []byte

	if withFunctions {
		xmlGnb = readXmlFile(t, GnbSetupRequestXmlPath)

	} else {
		xmlGnb = readXmlFile(t, GnbWithoutFunctionsSetupRequestXmlPath)
	}
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      gnbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}

	if withFunctions {
		gnb := nodebInfo.GetGnb()
		gnb.RanFunctions = []*entities.RanFunction{{RanFunctionId: 2, RanFunctionRevision: 2}}
	}

	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)

	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	gnbToUpdate := getExpectedNodebForExistingRan(*nodebInfo, notificationRequest.Payload)
	gnbToUpdate.SetupFromNetwork = true
	writerMock.On("UpdateNodebInfo", gnbToUpdate).Return(nil)
	if withFunctions {
		writerMock.On("UpdateNodebInfoAndPublish", gnbToUpdate).Return(nil)
	}
	e2tInstancesManagerMock.On("AddRansToInstance", e2tInstanceFullAddress, []string{gnbNodebRanName}).Return(nil)
	var errEmpty error
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mock.Anything, true)
}

func TestE2SetupRequestNotificationHandler_HandleExistingConnectedGnbRoutingManagerError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      gnbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}

	gnb := nodebInfo.GetGnb()
	gnb.RanFunctions = []*entities.RanFunction{{RanFunctionId: 2, RanFunctionRevision: 2}}

	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(errors.New("association error"))

	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	gnbToUpdate := getExpectedNodebForExistingRan(*nodebInfo, notificationRequest.Payload)

	gnbToUpdate2 := *gnbToUpdate
	gnbToUpdate2.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	gnbToUpdate2.SetupFromNetwork = true
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &gnbToUpdate2, gnbNodebRanName+"_DISCONNECTED").Return(nil)

	var errEmpty error
	rmrMessengerMock.On("WhSendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertCalled(t, "WhSendMsg", mock.Anything, true)
}

func TestE2SetupRequestNotificationHandler_HandleExistingGnbInvalidConnectionStatusError(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, _, e2tInstancesManagerMock, rmrMessengerMock := initMocks(t)
	var gnb = &entities.NodebInfo{RanName: gnbNodebRanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	readerMock.On("GetNodeb", gnbNodebRanName).Return(gnb, nil)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	prefBytes := []byte(e2SetupMsgPrefix)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
}

func TestE2SetupRequestNotificationHandler_HandleExistingConnectedGnbWithoutFunctionsSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleExistingConnectedGnbSuccess(t, false)
}

func TestE2SetupRequestNotificationHandler_HandleExistingConnectedGnbWithFunctionsSuccess(t *testing.T) {
	testE2SetupRequestNotificationHandler_HandleExistingConnectedGnbSuccess(t, true)
}

func TestE2SetupRequestNotificationHandler_HandleExistingDisconnectedGnbSuccess(t *testing.T) {
	xmlGnb := readXmlFile(t, GnbSetupRequestXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	readerMock.On("GetGeneralConfiguration").Return(&entities.GeneralConfiguration{EnableRic: true}, nil)
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(&entities.E2TInstance{}, nil)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      gnbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_DISCONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}

	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)

	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(e2SetupMsgPrefix), xmlGnb...)}
	gnbToUpdate := getExpectedNodebForExistingRan(*nodebInfo, notificationRequest.Payload)
	gnbToUpdate2 := *gnbToUpdate
	gnbToUpdate2.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	gnbToUpdate2.SetupFromNetwork = true
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", &gnbToUpdate2, gnbNodebRanName+"_CONNECTED").Return(nil)
	gnbToUpdate3 := *gnbToUpdate
	gnbToUpdate3.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	gnbToUpdate3.SetupFromNetwork = true
	gnbToUpdate3.AssociatedE2TInstanceAddress = e2tInstanceFullAddress
	writerMock.On("UpdateNodebInfo", &gnbToUpdate3).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", e2tInstanceFullAddress, []string{gnbNodebRanName}).Return(nil)
	var errEmpty error
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(&rmrCgo.MBuf{}, errEmpty)
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
}

/* Existing Ran UTs - END */
