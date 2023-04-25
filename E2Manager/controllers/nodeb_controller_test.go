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
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unsafe"
)

const (
	RanName                      = "test"
	AssociatedE2TInstanceAddress = "10.0.2.15:38000"
	CorruptedJson                = "{\"errorCode\":401,\"errorMessage\":\"corrupted json\"}"
	ValidationFailureJson        = "{\"errorCode\":402,\"errorMessage\":\"Validation error\"}"
	ResourceNotFoundJson         = "{\"errorCode\":404,\"errorMessage\":\"Resource not found\"}"
	NodebExistsJson              = "{\"errorCode\":406,\"errorMessage\":\"Nodeb already exists\"}"
	RnibErrorJson                = "{\"errorCode\":500,\"errorMessage\":\"RNIB error\"}"
	InternalErrorJson            = "{\"errorCode\":501,\"errorMessage\":\"Internal Server Error. Please try again later\"}"
	AddEnbUrl                    = "/nodeb/enb"
)

var (
	ServedNrCellInformationRequiredFields = []string{"cellId", "choiceNrMode", "nrMode", "servedPlmns"}
	NrNeighbourInformationRequiredFields  = []string{"nrCgi", "choiceNrMode", "nrMode"}
	AddEnbRequestRequiredFields           = []string{"ranName", "enb", "globalNbId"}
	UpdateEnbRequestRequiredFields        = []string{"enb"}
	GlobalIdRequiredFields                = []string{"plmnId", "nbId"}
	EnbRequiredFields                     = []string{"enbType", "servedCells"}
	ServedCellRequiredFields              = []string{"broadcastPlmns", "cellId", "choiceEutraMode", "eutraMode", "tac"}
)

type controllerGetNodebTestContext struct {
	ranName              string
	nodebInfo            *entities.NodebInfo
	rnibError            error
	expectedStatusCode   int
	expectedJsonResponse string
}

type controllerGetNodebIdListTestContext struct {
	nodebIdList          []*entities.NbIdentity
	rnibError            error
	expectedStatusCode   int
	expectedJsonResponse string
}

type getNodebInfoResult struct {
	nodebInfo *entities.NodebInfo
	rnibError error
}

type updateGnbCellsParams struct {
	err error
}

type updateEnbCellsParams struct {
	err error
}

type addEnbParams struct {
	err error
}

type addNbIdentityParams struct {
	err error
}

type removeServedCellsParams struct {
	servedCellInfo []*entities.ServedCellInfo
	err            error
}

type removeServedNrCellsParams struct {
	servedNrCells []*entities.ServedNRCell
	err           error
}

type controllerUpdateEnbTestContext struct {
	getNodebInfoResult      *getNodebInfoResult
	removeServedCellsParams *removeServedCellsParams
	updateEnbCellsParams    *updateEnbCellsParams
	requestBody             map[string]interface{}
	expectedStatusCode      int
	expectedJsonResponse    string
}

type controllerUpdateGnbTestContext struct {
	getNodebInfoResult        *getNodebInfoResult
	removeServedNrCellsParams *removeServedNrCellsParams
	updateGnbCellsParams      *updateGnbCellsParams
	requestBody               map[string]interface{}
	expectedStatusCode        int
	expectedJsonResponse      string
}

type controllerAddEnbTestContext struct {
	getNodebInfoResult   *getNodebInfoResult
	addEnbParams         *addEnbParams
	addNbIdentityParams  *addNbIdentityParams
	requestBody          map[string]interface{}
	expectedStatusCode   int
	expectedJsonResponse string
}

type controllerDeleteEnbTestContext struct {
	getNodebInfoResult   *getNodebInfoResult
	expectedStatusCode   int
	expectedJsonResponse string
}

func generateServedNrCells(cellIds ...string) []*entities.ServedNRCell {

	servedNrCells := []*entities.ServedNRCell{}

	for _, v := range cellIds {
		servedNrCells = append(servedNrCells, &entities.ServedNRCell{ServedNrCellInformation: &entities.ServedNRCellInformation{
			CellId: v,
			ChoiceNrMode: &entities.ServedNRCellInformation_ChoiceNRMode{
				Fdd: &entities.ServedNRCellInformation_ChoiceNRMode_FddInfo{},
			},
			NrMode:      entities.Nr_FDD,
			NrPci:       5,
			ServedPlmns: []string{"whatever"},
		}})
	}

	return servedNrCells
}

func generateServedCells(cellIds ...string) []*entities.ServedCellInfo {

	var servedCells []*entities.ServedCellInfo

	for i, v := range cellIds {
		servedCells = append(servedCells, &entities.ServedCellInfo{
			CellId: v,
			ChoiceEutraMode: &entities.ChoiceEUTRAMode{
				Fdd: &entities.FddInfo{},
			},
			Pci:            uint32(i + 1),
			BroadcastPlmns: []string{"whatever"},
		})
	}

	return servedCells
}

func buildNrNeighbourInformation(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"nrCgi": "whatever",
		"choiceNrMode": map[string]interface{}{
			"tdd": map[string]interface{}{},
		},
		"nrMode": 1,
		"nrPci":  1,
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func buildServedNrCellInformation(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"cellId": "whatever",
		"choiceNrMode": map[string]interface{}{
			"fdd": map[string]interface{}{},
		},
		"nrMode": 1,
		"nrPci":  1,
		"servedPlmns": []interface{}{
			"whatever",
		},
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func buildServedCell(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"cellId": "whatever",
		"choiceEutraMode": map[string]interface{}{
			"fdd": map[string]interface{}{},
		},
		"eutraMode": 1,
		"pci":       1,
		"tac":       "whatever3",
		"broadcastPlmns": []interface{}{
			"whatever",
		},
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func getUpdateEnbRequest(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"enb": buildEnb(propToOmit),
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func getAddEnbRequest(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"ranName":    RanName,
		"globalNbId": buildGlobalNbId(""),
		"enb":        buildEnb(""),
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func buildEnb(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"enbType": 1,
		"servedCells": []interface{}{
			buildServedCell(""),
		}}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func buildGlobalNbId(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"plmnId": "whatever",
		"nbId":   "whatever2",
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func setupControllerTest(t *testing.T) (*NodebController, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, managers.RanListManager) {
	log := initLog(t)
	config := configuration.ParseConfiguration()

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	rmrSender := getRmrSender(rmrMessengerMock, log)
	e2tInstancesManager := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	ranListManager := managers.NewRanListManager(log, rnibDataService)
	ranAlarmService := &mocks.RanAlarmServiceMock{}
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(log, rnibDataService, ranListManager, ranAlarmService)
	nodebValidator := managers.NewNodebValidator()
	updateEnbManager := managers.NewUpdateEnbManager(log, rnibDataService, nodebValidator)
	updateGnbManager := managers.NewUpdateGnbManager(log, rnibDataService, nodebValidator)
	handlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(log, rmrSender, config, rnibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager, nodebValidator, updateEnbManager, updateGnbManager, ranListManager)
	controller := NewNodebController(log, handlerProvider)
	return controller, readerMock, writerMock, rmrMessengerMock, e2tInstancesManager, ranListManager
}

func setupDeleteEnbControllerTest(t *testing.T, preAddNbIdentity bool) (*NodebController, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *entities.NbIdentity) {
	log := initLog(t)
	config := configuration.ParseConfiguration()

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	rmrSender := getRmrSender(rmrMessengerMock, log)
	e2tInstancesManager := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	ranListManager := managers.NewRanListManager(log, rnibDataService)
	var nbIdentity *entities.NbIdentity
	if preAddNbIdentity {
		nbIdentity = &entities.NbIdentity{InventoryName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}
		writerMock.On("AddNbIdentity", entities.Node_ENB, nbIdentity).Return(nil)
		ranListManager.AddNbIdentity(entities.Node_ENB, nbIdentity)
	}
	ranAlarmService := &mocks.RanAlarmServiceMock{}
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(log, rnibDataService, ranListManager, ranAlarmService)
	nodebValidator := managers.NewNodebValidator()
	updateEnbManager := managers.NewUpdateEnbManager(log, rnibDataService, nodebValidator)
	updateGnbManager := managers.NewUpdateGnbManager(log, rnibDataService, nodebValidator)

	handlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(log, rmrSender, config, rnibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager, nodebValidator, updateEnbManager, updateGnbManager, ranListManager)
	controller := NewNodebController(log, handlerProvider)
	return controller, readerMock, writerMock, nbIdentity
}

func TestShutdownHandlerRnibError(t *testing.T) {
	controller, _, _, _, e2tInstancesManagerMock, _ := setupControllerTest(t)
	e2tInstancesManagerMock.On("GetE2TAddresses").Return([]string{}, e2managererrors.NewRnibDbError())

	writer := httptest.NewRecorder()

	controller.Shutdown(writer, tests.GetHttpRequest())

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, e2managererrors.NewRnibDbError().Code)
}

func TestSetGeneralConfigurationHandlerRnibError(t *testing.T) {
	controller, readerMock, _, _, _, _ := setupControllerTest(t)

	configuration := &entities.GeneralConfiguration{}
	readerMock.On("GetGeneralConfiguration").Return(configuration, e2managererrors.NewRnibDbError())

	writer := httptest.NewRecorder()

	httpRequest, _ := http.NewRequest("PUT", "https://localhost:3800/v1/nodeb/parameters", strings.NewReader("{\"enableRic\":false}"))

	controller.SetGeneralConfiguration(writer, httpRequest)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, e2managererrors.NewRnibDbError().Code, errorResponse.Code)
}

func TestSetGeneralConfigurationInvalidJson(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()

	httpRequest, _ := http.NewRequest("PUT", "https://localhost:3800/v1/nodeb/parameters", strings.NewReader("{\"enableRic\":false, \"someValue\":false}"))

	controller.SetGeneralConfiguration(writer, httpRequest)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	assert.Equal(t, e2managererrors.NewInvalidJsonError().Code, errorResponse.Code)
}

func controllerGetNodebTestExecuter(t *testing.T, context *controllerGetNodebTestContext) {
	controller, readerMock, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	readerMock.On("GetNodeb", context.ranName).Return(context.nodebInfo, context.rnibError)
	req, _ := http.NewRequest(http.MethodGet, "/nodeb", nil)
	req = mux.SetURLVars(req, map[string]string{"ranName": context.ranName})
	controller.GetNodeb(writer, req)
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
}

func controllerGetNodebIdListTestExecuter(t *testing.T, context *controllerGetNodebIdListTestContext) {
	controller, readerMock, _, _, _, ranListManager := setupControllerTest(t)
	writer := httptest.NewRecorder()
	readerMock.On("GetListNodebIds").Return(context.nodebIdList, context.rnibError)

	err := ranListManager.InitNbIdentityMap()
	if err != nil {
		t.Errorf("Error cannot init identity")
	}

	req, _ := http.NewRequest(http.MethodGet, "/nodeb/states", nil)
	controller.GetNodebIdList(writer, req)
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Contains(t, context.expectedJsonResponse, string(bodyBytes))
}

func activateControllerUpdateEnbMocks(context *controllerUpdateEnbTestContext, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock, updateEnbRequest *models.UpdateEnbRequest) {
	if context.getNodebInfoResult != nil {
		readerMock.On("GetNodeb", RanName).Return(context.getNodebInfoResult.nodebInfo, context.getNodebInfoResult.rnibError)
	}

	if context.removeServedCellsParams != nil {
		writerMock.On("RemoveServedCells", RanName, context.removeServedCellsParams.servedCellInfo).Return(context.removeServedCellsParams.err)
	}

	if context.updateEnbCellsParams != nil {
		updatedNodebInfo := *context.getNodebInfoResult.nodebInfo

		if context.getNodebInfoResult.nodebInfo.SetupFromNetwork {
			updateEnbRequest.Enb.EnbType = context.getNodebInfoResult.nodebInfo.GetEnb().EnbType
		}

		updatedNodebInfo.Configuration = &entities.NodebInfo_Enb{Enb: updateEnbRequest.Enb}

		writerMock.On("UpdateEnb", &updatedNodebInfo, updateEnbRequest.Enb.ServedCells).Return(context.updateEnbCellsParams.err)
	}
}

func activateControllerUpdateGnbMocks(context *controllerUpdateGnbTestContext, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	if context.getNodebInfoResult != nil {
		readerMock.On("GetNodeb", RanName).Return(context.getNodebInfoResult.nodebInfo, context.getNodebInfoResult.rnibError)
	}

	if context.removeServedNrCellsParams != nil {
		writerMock.On("RemoveServedNrCells", RanName, context.removeServedNrCellsParams.servedNrCells).Return(context.removeServedNrCellsParams.err)
	}

	if context.updateGnbCellsParams != nil {
		updatedNodebInfo := *context.getNodebInfoResult.nodebInfo
		gnb := entities.Gnb{}
		_ = jsonpb.Unmarshal(getJsonRequestAsBuffer(context.requestBody), &gnb)
		updatedGnb := *updatedNodebInfo.GetGnb()
		updatedGnb.ServedNrCells = gnb.ServedNrCells
		writerMock.On("UpdateGnbCells", &updatedNodebInfo, gnb.ServedNrCells).Return(context.updateGnbCellsParams.err)
	}
}

func assertControllerUpdateGnb(t *testing.T, context *controllerUpdateGnbTestContext, writer *httptest.ResponseRecorder, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func assertControllerUpdateEnb(t *testing.T, context *controllerUpdateEnbTestContext, writer *httptest.ResponseRecorder, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func assertControllerAddEnb(t *testing.T, context *controllerAddEnbTestContext, writer *httptest.ResponseRecorder, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func assertControllerDeleteEnb(t *testing.T, context *controllerDeleteEnbTestContext, writer *httptest.ResponseRecorder, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func buildUpdateEnbRequest(context *controllerUpdateEnbTestContext) *http.Request {
	updateEnbUrl := fmt.Sprintf("/nodeb/enb/%s", RanName)
	requestBody := getJsonRequestAsBuffer(context.requestBody)
	req, _ := http.NewRequest(http.MethodPut, updateEnbUrl, requestBody)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"ranName": RanName})
	return req
}

func buildUpdateGnbRequest(context *controllerUpdateGnbTestContext) *http.Request {
	updateGnbUrl := fmt.Sprintf("/nodeb/gnb/%s", RanName)
	requestBody := getJsonRequestAsBuffer(context.requestBody)
	req, _ := http.NewRequest(http.MethodPut, updateGnbUrl, requestBody)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"ranName": RanName})
	return req
}

func buildAddEnbRequest(context *controllerAddEnbTestContext) *http.Request {
	requestBody := getJsonRequestAsBuffer(context.requestBody)
	req, _ := http.NewRequest(http.MethodPost, AddEnbUrl, requestBody)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func controllerUpdateEnbTestExecuter(t *testing.T, context *controllerUpdateEnbTestContext) {
	controller, readerMock, writerMock, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()

	r := buildUpdateEnbRequest(context)
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, LimitRequest))

	updateEnbRequest := models.UpdateEnbRequest{}
	_ = json.Unmarshal(body, &updateEnbRequest)

	activateControllerUpdateEnbMocks(context, readerMock, writerMock, &updateEnbRequest)
	r = buildUpdateEnbRequest(context)
	defer r.Body.Close()

	controller.UpdateEnb(writer, r)

	assertControllerUpdateEnb(t, context, writer, readerMock, writerMock)
}

func controllerUpdateGnbTestExecuter(t *testing.T, context *controllerUpdateGnbTestContext) {
	controller, readerMock, writerMock, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()

	activateControllerUpdateGnbMocks(context, readerMock, writerMock)
	req := buildUpdateGnbRequest(context)
	controller.UpdateGnb(writer, req)
	assertControllerUpdateGnb(t, context, writer, readerMock, writerMock)
}

func activateControllerAddEnbMocks(context *controllerAddEnbTestContext, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock, addEnbRequest *models.AddEnbRequest) {
	if context.getNodebInfoResult != nil {
		readerMock.On("GetNodeb", RanName).Return(context.getNodebInfoResult.nodebInfo, context.getNodebInfoResult.rnibError)
	}

	if context.addEnbParams != nil {
		nodebInfo := entities.NodebInfo{
			RanName:          addEnbRequest.RanName,
			Ip:               addEnbRequest.Ip,
			Port:             addEnbRequest.Port,
			NodeType:         entities.Node_ENB,
			GlobalNbId:       addEnbRequest.GlobalNbId,
			Configuration:    &entities.NodebInfo_Enb{Enb: addEnbRequest.Enb},
			ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,
		}

		writerMock.On("AddEnb", &nodebInfo).Return(context.addEnbParams.err)
	}

	if context.addNbIdentityParams != nil {
		nbIdentity := entities.NbIdentity{InventoryName: addEnbRequest.RanName, GlobalNbId: addEnbRequest.GlobalNbId, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
		writerMock.On("AddNbIdentity", entities.Node_ENB, &nbIdentity).Return(context.addNbIdentityParams.err)
	}
}

func controllerAddEnbTestExecuter(t *testing.T, context *controllerAddEnbTestContext) {
	controller, readerMock, writerMock, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	r := buildAddEnbRequest(context)
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, LimitRequest))

	addEnbRequest := models.AddEnbRequest{}

	_ = json.Unmarshal(body, &addEnbRequest)
	activateControllerAddEnbMocks(context, readerMock, writerMock, &addEnbRequest)
	r = buildAddEnbRequest(context)
	defer r.Body.Close()
	controller.AddEnb(writer, r)
	assertControllerAddEnb(t, context, writer, readerMock, writerMock)
}

func controllerDeleteEnbTestExecuter(t *testing.T, context *controllerDeleteEnbTestContext, preAddNbIdentity bool) {
	controller, readerMock, writerMock, nbIdentity := setupDeleteEnbControllerTest(t, preAddNbIdentity)
	readerMock.On("GetNodeb", RanName).Return(context.getNodebInfoResult.nodebInfo, context.getNodebInfoResult.rnibError)
	if context.getNodebInfoResult.rnibError == nil && context.getNodebInfoResult.nodebInfo.GetNodeType() == entities.Node_ENB &&
		!context.getNodebInfoResult.nodebInfo.SetupFromNetwork {
		writerMock.On("RemoveEnb", context.getNodebInfoResult.nodebInfo).Return(nil)
		if preAddNbIdentity {
			writerMock.On("RemoveNbIdentity", entities.Node_ENB, nbIdentity).Return(nil)
		}
	}
	writer := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodDelete, AddEnbUrl+"/"+RanName, nil)
	r.Header.Set("Content-Type", "application/json")
	r = mux.SetURLVars(r, map[string]string{"ranName": RanName})
	controller.DeleteEnb(writer, r)
	assertControllerDeleteEnb(t, context, writer, readerMock, writerMock)
}

/*
UpdateGnb UTs
*/

// BEGIN - UpdateGnb Validation Failure UTs

func TestControllerUpdateGnbEmptyServedNrCells(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingServedNrCellInformation(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": nil,
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingServedNrCellRequiredProp(t *testing.T) {

	for _, v := range ServedNrCellInformationRequiredFields {
		context := controllerUpdateGnbTestContext{
			getNodebInfoResult: nil,
			requestBody: map[string]interface{}{
				"servedNrCells": []interface{}{
					map[string]interface{}{
						"servedNrCellInformation": buildServedNrCellInformation(v),
					},
				},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerUpdateGnbTestExecuter(t, &context)
	}
}

func TestControllerUpdateGnbMissingServedNrCellFddOrTdd(t *testing.T) {

	servedNrCellInformation := buildServedNrCellInformation("")
	servedNrCellInformation["choiceNrMode"] = map[string]interface{}{}

	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": servedNrCellInformation,
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingNeighbourInfoFddOrTdd(t *testing.T) {

	nrNeighbourInfo := buildNrNeighbourInformation("")
	nrNeighbourInfo["choiceNrMode"] = map[string]interface{}{}

	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						nrNeighbourInfo,
					},
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingNrNeighbourInformationRequiredProp(t *testing.T) {

	for _, v := range NrNeighbourInformationRequiredFields {
		context := controllerUpdateGnbTestContext{
			getNodebInfoResult: nil,
			requestBody: map[string]interface{}{
				"servedNrCells": []interface{}{
					map[string]interface{}{
						"servedNrCellInformation": buildServedNrCellInformation(""),
						"nrNeighbourInfos": []interface{}{
							buildNrNeighbourInformation(v),
						},
					},
				},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerUpdateGnbTestExecuter(t, &context)
	}
}

// END - UpdateGnb Validation Failure UTs

func TestControllerUpdateGnbValidServedNrCellInformationGetNodebNotFound(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
				},
			},
		},
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbValidServedNrCellInformationGetNodebInternalError(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessInvalidGnbConfiguration(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessRemoveServedNrCellsFailure(t *testing.T) {
	oldServedNrCells := generateServedNrCells("whatever1", "whatever2")
	context := controllerUpdateGnbTestContext{
		removeServedNrCellsParams: &removeServedNrCellsParams{
			err:           common.NewInternalError(errors.New("#writer.UpdateGnbCells - Internal Error")),
			servedNrCells: oldServedNrCells,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
				NodeType:                     entities.Node_GNB,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessUpdateGnbCellsFailure(t *testing.T) {
	oldServedNrCells := generateServedNrCells("whatever1", "whatever2")
	context := controllerUpdateGnbTestContext{
		removeServedNrCellsParams: &removeServedNrCellsParams{
			err:           nil,
			servedNrCells: oldServedNrCells,
		},
		updateGnbCellsParams: &updateGnbCellsParams{
			err: common.NewInternalError(errors.New("#writer.UpdateGnbCells - Internal Error")),
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
				NodeType:                     entities.Node_GNB,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbExistingEmptyCellsSuccess(t *testing.T) {
	oldServedNrCells := []*entities.ServedNRCell{}

	context := controllerUpdateGnbTestContext{
		updateGnbCellsParams: &updateGnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
				NodeType:                     entities.Node_GNB,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"nodeType\":\"GNB\",\"gnb\":{\"servedNrCells\":[{\"servedNrCellInformation\":{\"nrPci\":1,\"cellId\":\"whatever\",\"servedPlmns\":[\"whatever\"],\"nrMode\":\"FDD\",\"choiceNrMode\":{\"fdd\":{}}},\"nrNeighbourInfos\":[{\"nrPci\":1,\"nrCgi\":\"whatever\",\"nrMode\":\"FDD\",\"choiceNrMode\":{\"tdd\":{}}}]}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\"}",
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbSuccess(t *testing.T) {
	oldServedNrCells := generateServedNrCells("whatever1", "whatever2")

	context := controllerUpdateGnbTestContext{
		removeServedNrCellsParams: &removeServedNrCellsParams{
			err:           nil,
			servedNrCells: oldServedNrCells,
		},
		updateGnbCellsParams: &updateGnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
				NodeType:                     entities.Node_GNB,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"nodeType\":\"GNB\",\"gnb\":{\"servedNrCells\":[{\"servedNrCellInformation\":{\"nrPci\":1,\"cellId\":\"whatever\",\"servedPlmns\":[\"whatever\"],\"nrMode\":\"FDD\",\"choiceNrMode\":{\"fdd\":{}}},\"nrNeighbourInfos\":[{\"nrPci\":1,\"nrCgi\":\"whatever\",\"nrMode\":\"FDD\",\"choiceNrMode\":{\"tdd\":{}}}]}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\"}",
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

/*
UpdateEnb UTs
*/

func TestControllerUpdateEnbInvalidRequest(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	invalidJson := strings.NewReader("{enb:\"whatever\"")

	updateEnbUrl := fmt.Sprintf("/nodeb/enb/%s", RanName)
	req, _ := http.NewRequest(http.MethodPut, updateEnbUrl, invalidJson)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"ranName": RanName})

	controller.UpdateEnb(writer, req)

	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, CorruptedJson, string(bodyBytes))
}

func TestControllerUpdateEnbEmptyEnbType(t *testing.T) {
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult:   nil,
		requestBody:          getUpdateEnbRequest(EnbRequiredFields[0]),
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbEmptyServedCells(t *testing.T) {
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult:   nil,
		requestBody:          getUpdateEnbRequest(EnbRequiredFields[1]),
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbMissingEnb(t *testing.T) {
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult:   nil,
		requestBody:          getUpdateEnbRequest(UpdateEnbRequestRequiredFields[0]),
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbMissingRequiredServedCellProps(t *testing.T) {

	r := getUpdateEnbRequest("")

	for _, v := range ServedCellRequiredFields {
		enb := r["enb"]

		enbMap, _ := enb.(map[string]interface{})

		enbMap["servedCells"] = []interface{}{
			buildServedCell(v),
		}

		context := controllerUpdateEnbTestContext{
			requestBody:          r,
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerUpdateEnbTestExecuter(t, &context)
	}
}

func TestControllerUpdateEnbValidServedCellsGetNodebNotFound(t *testing.T) {
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbValidServedCellsGetNodebInternalError(t *testing.T) {
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbGetNodebSuccessGnbTypeFailure(t *testing.T) {
	oldServedCells := generateServedCells("whatever1", "whatever2")
	context := controllerUpdateEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_GNB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells}},
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbGetNodebSuccessRemoveServedCellsFailure(t *testing.T) {
	oldServedCells := generateServedCells("whatever1", "whatever2")
	context := controllerUpdateEnbTestContext{
		removeServedCellsParams: &removeServedCellsParams{
			err:            common.NewInternalError(errors.New("#writer.RemoveServedCells - Internal Error")),
			servedCellInfo: oldServedCells,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells}},
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbGetNodebSuccessUpdateEnbFailure(t *testing.T) {
	oldServedCells := generateServedCells("whatever1", "whatever2")
	context := controllerUpdateEnbTestContext{
		removeServedCellsParams: &removeServedCellsParams{
			err:            nil,
			servedCellInfo: oldServedCells,
		},
		updateEnbCellsParams: &updateEnbCellsParams{
			err: common.NewInternalError(errors.New("#writer.UpdateEnb - Internal Error")),
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells, EnbType: entities.EnbType_MACRO_ENB}},
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbExistingEmptyCellsSuccess(t *testing.T) {
	oldServedCells := []*entities.ServedCellInfo{}
	context := controllerUpdateEnbTestContext{
		updateEnbCellsParams: &updateEnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells, EnbType: entities.EnbType_MACRO_ENB}},
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"nodeType\":\"ENB\",\"enb\":{\"enbType\":\"MACRO_ENB\",\"servedCells\":[{\"pci\":1,\"cellId\":\"whatever\",\"tac\":\"whatever3\",\"broadcastPlmns\":[\"whatever\"],\"choiceEutraMode\":{\"fdd\":{}},\"eutraMode\":\"FDD\"}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\"}",
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbNgEnbFailure(t *testing.T) {

	requestBody := map[string]interface{}{
		"enb": map[string]interface{}{
			"enbType": 3,
			"servedCells": []interface{}{
				buildServedCell(""),
			}},
	}

	oldServedCells := generateServedCells("whatever1", "whatever2")

	context := controllerUpdateEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells, EnbType: entities.EnbType_MACRO_NG_ENB}},
			},
			rnibError: nil,
		},
		requestBody:          requestBody,
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbSuccessSetupFromNwFalse(t *testing.T) {
	oldServedCells := generateServedCells("whatever1", "whatever2")
	context := controllerUpdateEnbTestContext{
		removeServedCellsParams: &removeServedCellsParams{
			err:            nil,
			servedCellInfo: oldServedCells,
		},
		updateEnbCellsParams: &updateEnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells, EnbType: entities.EnbType_LONG_MACRO_ENB}},
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"nodeType\":\"ENB\",\"enb\":{\"enbType\":\"MACRO_ENB\",\"servedCells\":[{\"pci\":1,\"cellId\":\"whatever\",\"tac\":\"whatever3\",\"broadcastPlmns\":[\"whatever\"],\"choiceEutraMode\":{\"fdd\":{}},\"eutraMode\":\"FDD\"}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\"}",
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

func TestControllerUpdateEnbSuccessSetupFromNwTrue(t *testing.T) {
	oldServedCells := generateServedCells("whatever1", "whatever2")
	context := controllerUpdateEnbTestContext{
		removeServedCellsParams: &removeServedCellsParams{
			err:            nil,
			servedCellInfo: oldServedCells,
		},
		updateEnbCellsParams: &updateEnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				NodeType:                     entities.Node_ENB,
				Configuration:                &entities.NodebInfo_Enb{Enb: &entities.Enb{ServedCells: oldServedCells, EnbType: entities.EnbType_LONG_MACRO_ENB}},
				SetupFromNetwork:             true,
			},
			rnibError: nil,
		},
		requestBody:          getUpdateEnbRequest(""),
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"nodeType\":\"ENB\",\"enb\":{\"enbType\":\"LONG_MACRO_ENB\",\"servedCells\":[{\"pci\":1,\"cellId\":\"whatever\",\"tac\":\"whatever3\",\"broadcastPlmns\":[\"whatever\"],\"choiceEutraMode\":{\"fdd\":{}},\"eutraMode\":\"FDD\"}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\",\"setupFromNetwork\":true}",
	}

	controllerUpdateEnbTestExecuter(t, &context)
}

/*
AddEnb UTs
*/

func TestControllerAddEnbGetNodebInternalError(t *testing.T) {
	context := controllerAddEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		},
		requestBody:          getAddEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerAddEnbNodebExistsFailure(t *testing.T) {
	context := controllerAddEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{},
			rnibError: nil,
		},
		requestBody:          getAddEnbRequest(""),
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: NodebExistsJson,
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerAddEnbSaveNodebFailure(t *testing.T) {
	context := controllerAddEnbTestContext{
		addEnbParams: &addEnbParams{
			err: common.NewInternalError(errors.New("#reader.AddEnb - Internal Error")),
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody:          getAddEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerAddEnbAddNbIdentityFailure(t *testing.T) {
	context := controllerAddEnbTestContext{
		addEnbParams: &addEnbParams{
			err: nil,
		},
		addNbIdentityParams: &addNbIdentityParams{
			err: common.NewInternalError(errors.New("#writer.addNbIdentity - Internal Error")),
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody:          getAddEnbRequest(""),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerAddEnbMissingRequiredRequestProps(t *testing.T) {

	for _, v := range AddEnbRequestRequiredFields {
		context := controllerAddEnbTestContext{
			requestBody:          getAddEnbRequest(v),
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerAddEnbTestExecuter(t, &context)
	}
}

func TestControllerAddEnbInvalidRequest(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()

	// Invalid json: attribute name without quotes (should be "cause":).
	invalidJson := strings.NewReader("{ranName:\"whatever\"")
	req, _ := http.NewRequest(http.MethodPost, AddEnbUrl, invalidJson)

	controller.AddEnb(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, CorruptedJson, string(bodyBytes))

}

func TestControllerAddEnbMissingRequiredGlobalNbIdProps(t *testing.T) {

	r := getAddEnbRequest("")

	for _, v := range GlobalIdRequiredFields {
		r["globalNbId"] = buildGlobalNbId(v)

		context := controllerAddEnbTestContext{
			requestBody:          r,
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerAddEnbTestExecuter(t, &context)
	}
}

func TestControllerAddEnbMissingRequiredEnbProps(t *testing.T) {

	r := getAddEnbRequest("")

	for _, v := range EnbRequiredFields {
		r["enb"] = buildEnb(v)

		context := controllerAddEnbTestContext{
			requestBody:          r,
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerAddEnbTestExecuter(t, &context)
	}
}

func TestControllerAddEnbMissingRequiredServedCellProps(t *testing.T) {

	r := getAddEnbRequest("")

	for _, v := range ServedCellRequiredFields {
		enb := r["enb"]

		enbMap, _ := enb.(map[string]interface{})

		enbMap["servedCells"] = []interface{}{
			buildServedCell(v),
		}

		context := controllerAddEnbTestContext{
			requestBody:          r,
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerAddEnbTestExecuter(t, &context)
	}
}

func TestControllerAddEnbNgEnbFailure(t *testing.T) {

	requestBody := map[string]interface{}{
		"ranName":    RanName,
		"globalNbId": buildGlobalNbId(""),
		"enb": map[string]interface{}{
			"enbType": 5,
			"servedCells": []interface{}{
				buildServedCell(""),
			},
		},
	}

	context := controllerAddEnbTestContext{
		requestBody:          requestBody,
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerAddEnbSuccess(t *testing.T) {
	context := controllerAddEnbTestContext{
		addEnbParams: &addEnbParams{
			err: nil,
		},
		addNbIdentityParams: &addNbIdentityParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody:          getAddEnbRequest(""),
		expectedStatusCode:   http.StatusCreated,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"DISCONNECTED\",\"globalNbId\":{\"plmnId\":\"whatever\",\"nbId\":\"whatever2\"},\"nodeType\":\"ENB\",\"enb\":{\"enbType\":\"MACRO_ENB\",\"servedCells\":[{\"pci\":1,\"cellId\":\"whatever\",\"tac\":\"whatever3\",\"broadcastPlmns\":[\"whatever\"],\"choiceEutraMode\":{\"fdd\":{}},\"eutraMode\":\"FDD\"}]}}",
	}

	controllerAddEnbTestExecuter(t, &context)
}

func TestControllerDeleteEnbGetNodebInternalError(t *testing.T) {
	context := controllerDeleteEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerDeleteEnbTestExecuter(t, &context, false)
}

/*
DeleteEnb UTs
*/

func TestControllerDeleteEnbNodebNotExistsFailure(t *testing.T) {
	context := controllerDeleteEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found"),
		},
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerDeleteEnbTestExecuter(t, &context, false)
}

func TestControllerDeleteEnbNodebNotEnb(t *testing.T) {
	context := controllerDeleteEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{RanName: "ran1", NodeType: entities.Node_GNB, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED},
			rnibError: nil,
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerDeleteEnbTestExecuter(t, &context, false)
}

func TestControllerDeleteEnbSetupFromNetworkTrueFailure(t *testing.T) {
	context := controllerDeleteEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{RanName: RanName, NodeType: entities.Node_ENB, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, SetupFromNetwork: true},
			rnibError: nil,
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}
	controllerDeleteEnbTestExecuter(t, &context, true)
}

func TestControllerDeleteEnbSuccess(t *testing.T) {
	context := controllerDeleteEnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{RanName: RanName, NodeType: entities.Node_ENB, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED},
			rnibError: nil,
		},
		expectedStatusCode:   http.StatusNoContent,
		expectedJsonResponse: "",
	}
	controllerDeleteEnbTestExecuter(t, &context, true)
}

func getJsonRequestAsBuffer(requestJson map[string]interface{}) *bytes.Buffer {
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(requestJson)
	return b
}

/*
GetNodeb UTs
*/

func TestControllerGetNodebSuccess(t *testing.T) {
	ranName := "test"
	var rnibError error
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            &entities.NodebInfo{RanName: ranName, Ip: "10.0.2.15", Port: 1234},
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: fmt.Sprintf("{\"ranName\":\"%s\",\"ip\":\"10.0.2.15\",\"port\":1234}", ranName),
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebNotFound(t *testing.T) {

	ranName := "test"
	var nodebInfo *entities.NodebInfo
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            nodebInfo,
		rnibError:            common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebInternal(t *testing.T) {
	ranName := "test"
	var nodebInfo *entities.NodebInfo
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            nodebInfo,
		rnibError:            common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebIdListSuccess(t *testing.T) {
	var rnibError error
	nodebIdList := []*entities.NbIdentity{
		{InventoryName: "test1", GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}},
		{InventoryName: "test2", GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}},
	}

	context := controllerGetNodebIdListTestContext{
		nodebIdList:          nodebIdList,
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "[{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}},{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}}][{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}},{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}}]",
	}

	controllerGetNodebIdListTestExecuter(t, &context)
}

func TestControllerGetNodebIdListEmptySuccess(t *testing.T) {
	var rnibError error
	var nodebIdList []*entities.NbIdentity

	context := controllerGetNodebIdListTestContext{
		nodebIdList:          nodebIdList,
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "[]",
	}

	controllerGetNodebIdListTestExecuter(t, &context)
}

func TestHeaderValidationFailed(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()

	header := &http.Header{}

	controller.handleRequest(writer, header, httpmsghandlerprovider.ShutdownRequest, nil, true, 0)

	var errorResponse = parseJsonRequest(t, writer.Body)
	err := e2managererrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestShutdownStatusNoContent(t *testing.T) {
	controller, readerMock, _, _, e2tInstancesManagerMock, _ := setupControllerTest(t)
	e2tInstancesManagerMock.On("GetE2TAddresses").Return([]string{}, nil)
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, nil)

	writer := httptest.NewRecorder()
	controller.Shutdown(writer, tests.GetHttpRequest())

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestHandleInternalError(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := e2managererrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleCommandAlreadyInProgressError(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	err := e2managererrors.NewCommandAlreadyInProgressError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleRoutingManagerError(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	err := e2managererrors.NewRoutingManagerError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusServiceUnavailable, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleE2TInstanceAbsenceError(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := e2managererrors.NewE2TInstanceAbsenceError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusServiceUnavailable, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestValidateHeaders(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	result := controller.validateRequestHeader(&header)

	assert.Nil(t, result)
}

func parseJsonRequest(t *testing.T, r io.Reader) models.ErrorResponse {

	var errorResponse models.ErrorResponse
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error cannot deserialize json request")
	}
	_ = json.Unmarshal(body, &errorResponse)

	return errorResponse
}

func initLog(t *testing.T) *logger.Logger {
        InfoLevel := int8(3)
	log, err := logger.InitLogger(InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}

func TestX2ResetHandleSuccessfulRequestedCause(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock, _, _ := setupControllerTest(t)

	ranName := "test1"
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	data4Req := map[string]interface{}{"cause": "protocol:transfer-syntax-error"}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data4Req)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)

}

func TestX2ResetHandleSuccessfulRequestedDefault(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock, _, _ := setupControllerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, true).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	// no body
	b := new(bytes.Buffer)
	data4Req := map[string]interface{}{}
	_ = json.NewEncoder(b).Encode(data4Req)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestX2ResetHandleFailureBodyReadError(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	ranName := "test1"
	writer := httptest.NewRecorder()

	// Fake reader to return reading error.
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", errReader(0))
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

}

func TestX2ResetHandleFailureInvalidBody(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	ranName := "test1"

	writer := httptest.NewRecorder()

	// Invalid json: attribute name without quotes (should be "cause":).
	b := strings.NewReader("{cause:\"protocol:transfer-syntax-error\"")
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

}

/*
func TestControllerHealthCheckRequestSuccess(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock, _, _ := setupControllerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, true).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	data4Req := map[string]interface{}{"ranList": []string{"abcd"}}
    b := new(bytes.Buffer)
    _ = json.NewEncoder(b).Encode(data4Req)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/v1/nodeb/health", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.HealthCheckRequest(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}
*/

func TestHandleErrorResponse(t *testing.T) {
	controller, _, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRnibDbError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewCommandAlreadyInProgressError(), writer)
	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewHeaderValidationError(), writer)
	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewWrongStateError("", ""), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRequestValidationError(), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRmrError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewNoConnectedRanError(), writer)
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewResourceNotFoundError(), writer)
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(fmt.Errorf("ErrorError"), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}

func getRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
