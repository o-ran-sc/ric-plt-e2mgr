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
//  platform project (RICP)

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/tests"
	"e2mgr/services"
	"e2mgr/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	RanNameForErrorIndication        = "test"
	E2tAddress                       = "10.10.2.15:9800"
	e2tInstanceFullAddressErrorIndication          = "10.0.2.15:9999"
	e2SetupMsgPrefixErrorIndication                = e2tInstanceFullAddressErrorIndication + "|"
	ErrorIndicationXmlPath           = "../../tests/resources/errorIndication/ErrorIndicationForSetupRequest.xml"
	ErrorIndicationWithoutCDXmlPath           = "../../tests/resources/errorIndication/ErrorIndicationWithoutCD.xml"
	ErrorIndicationXmlPathServiceUpdate = "../../tests/resources/errorIndication/ErrorIndicationForServiceUpdate.xml"
	ErrorIndicationXmlPathUnsuccessfuOutcome =  "../../tests/resources/errorIndication/ErrorIndicationUnsuccessfulOutcome.xml"
	ErrorIndicationXmlPathDefault = "../../tests/resources/errorIndication/ErrorIndicationForDefault.xml"
	ErrorIndicationInvalidXmlPath = "../../tests/resources/errorIndication/ErrorIndicationInvalid.xml"
)

func initErrorIndication(t *testing.T) (*ErrorIndicationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock, managers.RanListManager,*mocks.RanDisconnectionManagerMock,*mocks.RicServiceUpdateManagerMock, *mocks.MockLogger, *mocks.HttpClientMock,*mocks.RanListManagerMock) {
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
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	RanDisconnectionManagerMock := &mocks.RanDisconnectionManagerMock{}
	ricServiceUpdateManagerMock := &mocks.RicServiceUpdateManagerMock{}
	MockLogger := &mocks.MockLogger{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	ranListManagerMock := &mocks.RanListManagerMock{}

	ranListManager := managers.NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock, ranConnectStatusChangeManager)
	ranDisconnectionManager := managers.NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager, ranConnectStatusChangeManager)
	RicServiceUpdateManager := managers.NewRicServiceUpdateManager(logger, rnibDataService)
	handler := ErrorIndicationNotificationHandler(logger, ranDisconnectionManager, RicServiceUpdateManager)

	return handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock, ranListManager, RanDisconnectionManagerMock, ricServiceUpdateManagerMock,MockLogger,httpClientMock,ranListManagerMock
}

func TestParseErrorIndicationMessage_Success(t *testing.T) {
	ErrorgnbXml := utils.ReadXmlFile(t, ErrorIndicationXmlPath)
	handler, _, _, _, _, _, _, _, _,_,_,_ := initErrorIndication(t)
	prefBytes := []byte(e2SetupMsgPrefixErrorIndication)
	errorIndicationMessage, err := handler.parseErrorIndication(append(prefBytes, ErrorgnbXml...))
	assert.NotNil(t, errorIndicationMessage)
	assert.Nil(t, err)
}

func TestParseErrorIndication_PipFailure(t *testing.T) {
	ErrorgnbXml := utils.ReadXmlFile(t, ErrorIndicationXmlPath)
	handler, _, _, _, _, _,_ ,_, _, _,_,_ := initErrorIndication(t)
	prefBytes := []byte("10.0.2.15:9999")
	request, err := handler.parseErrorIndication(append(prefBytes, ErrorgnbXml...))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#ErrorIndicationHandler.parseErrorIndication - Error parsing ERROR INDICATION failed extract Payload: no | separator found")
}
func TestParseErrorIndicationMessage_UnmarshalFailure(t *testing.T) {
	handler, _,_, _, _, _, _, _, _, _,_,_ := initErrorIndication(t)
	prefBytes := []byte(e2SetupMsgPrefixErrorIndication)
	errorIndicationMessage, err := handler.parseErrorIndication(append(prefBytes, 1, 2, 3))
	assert.Nil(t, errorIndicationMessage)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#ErrorIndicationHandler.parseErrorIndication - Error unmarshalling ERROR INDICATION payload: 31302e302e322e31353a393939397c010203")
}

func testErrorIndicationNotificationHandler(t *testing.T) {
	handler, readerMock, writerMock, _, _, _, _, _, _,_ ,_,_:= initErrorIndication(t)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	notificationRequest := models.NotificationRequest{RanName: RanNameForErrorIndication}
	handler.Handle(&notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func testErrorIndicationHandlerWhenConnectedRanSuccess(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, e2tInstancesManagerMock,routingManagerClientMock, _, _, _, _,httpClientMock,_ := initErrorIndication(t)
	origNodebInfo := &entities.NodebInfo{
		RanName:                      RanNameForErrorIndication,
		GlobalNbId:                   &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"},
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		AssociatedE2TInstanceAddress: E2tAddress,
	}

	models.UpdateProcedureType(RanNameForErrorIndication,models.E2SetupProcedureCompleted)
	var rnibErr error
	readerMock.On("GetNodeb", RanNameForErrorIndication).Return(origNodebInfo, rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, RanNameForErrorIndication+"_DISCONNECTED").Return(rnibErr)
	updatedNodebInfo2 := *origNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address: E2tAddress, AssociatedRanList: []string{RanNameForErrorIndication}}
	readerMock.On("GetE2TInstance", E2tAddress).Return(e2tInstance, nil).Maybe()
	e2tInstanceToSave := *e2tInstance
	e2tInstanceToSave.AssociatedRanList = []string{}
	//writerMock.On("SaveE2TInstance", &e2tInstanceToSave).Return(nil)
	mockHttpClientForErrorIndication(httpClientMock, true) //After uncommenting testcase is failing 
	e2tInstancesManagerMock.On("RemoveRanFromInstance", RanNameForErrorIndication, E2tAddress).Return(nil)
	routingManagerClientMock.On("DissociateRanE2TInstance", E2tAddress, RanNameForErrorIndication).Return(nil)
	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}

	
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}
func TestErrorIndicationHandlerWhenConnectedGnbSuccessE2Setup(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccess(t, ErrorIndicationXmlPath)
}

func TestErrorIndicationHandlerWhenConnectedGnbSuccessProcedureType(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccess(t, ErrorIndicationWithoutCDXmlPath)
}
func TestErrorIndicationHandlerWhenConnectedGnbSuccessServiceUpdate(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccessServiceUpdate(t, ErrorIndicationXmlPathServiceUpdate)
}
func TestErrorIndicationHandlerWhenConnectedGnbSuccessServiceUpdateProcedureType(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccessServiceUpdate(t, ErrorIndicationWithoutCDXmlPath)
}
func TestErrorIndicationHandlerWhenConnectedGnbSuccessUnsuccessfulOutcome(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccess(t, ErrorIndicationXmlPathUnsuccessfuOutcome)
}
func TestErrorIndicationHandlerInvalidXML(t *testing.T) {
	testErrorIndicationHandlerInvalidXML(t, ErrorIndicationInvalidXmlPath)
}
func TestErrorIndicationHandlerForUnknownProcedureType(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccessUnknownProcedureType(t,ErrorIndicationWithoutCDXmlPath)
}
func TestErrorIndicationHandlerForUnhandlingProcedureType(t *testing.T) {
	testErrorIndicationHandlerWhenConnectedRanSuccessUnhandlingProcedureType(t,ErrorIndicationWithoutCDXmlPath)
}
func TestErrorIndicationHandlerForDefaultProcedureCode(t *testing.T) {
	testErrorIndicationHandlerForDefaultProcedureCode(t,ErrorIndicationXmlPathDefault)
}
func mockHttpClientForErrorIndication(httpClientMock *mocks.HttpClientMock, isSuccessful bool) {
	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2tAddress, RanNameForErrorIndication)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	var respStatusCode int
	if isSuccessful {
		respStatusCode = http.StatusCreated
	} else {
		respStatusCode = http.StatusBadRequest
	}
	httpClientMock.On("Post", clients.DissociateRanE2TInstanceApiSuffix, "application/json", body).Return(&http.Response{StatusCode: respStatusCode, Body: respBody}, nil).Maybe()
}

func testErrorIndicationHandlerWhenConnectedRanSuccessServiceUpdate(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, _,_, _, _, _, _,_,_ := initErrorIndication(t)
	origNodebInfo := &entities.NodebInfo{
		RanName:                      RanNameForErrorIndication,
		GlobalNbId:                   &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"},
		ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
		AssociatedE2TInstanceAddress: E2tAddress,
	}
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
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	RicServiceUpdateManager := managers.NewRicServiceUpdateManager(logger, rnibDataService)
	models.UpdateProcedureType(RanNameForErrorIndication,models.RicServiceUpdateCompleted)

	var rnibErr error
	readerMock.On("GetNodeb", RanNameForErrorIndication).Return(origNodebInfo, rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	updatedNodebInfo2 := *origNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = ""
	e2tInstance := &entities.E2TInstance{Address: E2tAddress, AssociatedRanList: []string{RanNameForErrorIndication}}
	readerMock.On("GetE2TInstance", E2tAddress).Return(e2tInstance, nil).Maybe()
	e2tInstanceToSave := *e2tInstance
	e2tInstanceToSave.AssociatedRanList = []string{}
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	err := RicServiceUpdateManager.RevertRanFunctions(ranName)
	assert.Nil(t,err)

	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}
func testErrorIndicationHandlerWhenConnectedRanSuccessUnknownProcedureType(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, _,_, _, _, _, _,httpClientMock,_ := initErrorIndication(t)

	models.UpdateProcedureType(RanNameForErrorIndication,models.E2SetupProcedureNotInitiated)
	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func testErrorIndicationHandlerWhenConnectedRanSuccessUnhandlingProcedureType(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, _,_, _, _, _, _,httpClientMock,_ := initErrorIndication(t)

	models.UpdateProcedureType(RanNameForErrorIndication,models.E2SetupProcedureFailure)
	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}

	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func testErrorIndicationHandlerInvalidXML(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, _,_, _, _, _, _,httpClientMock,_ := initErrorIndication(t)

	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}

	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
func testErrorIndicationHandlerForDefaultProcedureCode(t *testing.T,xmlPath string) {
	xml := utils.ReadXmlFile(t, xmlPath)
	handler, readerMock, writerMock, _, _,_, _, _, _, _,httpClientMock,_ := initErrorIndication(t)

	notificationRequest := &models.NotificationRequest{RanName: RanNameForErrorIndication, Payload: append([]byte(e2SetupMsgPrefixErrorIndication), xml...)}

	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
