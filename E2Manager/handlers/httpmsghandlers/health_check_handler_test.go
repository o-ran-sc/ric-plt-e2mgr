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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package httpmsghandlers

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/utils"
	"encoding/xml"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"unsafe"
)

const (
	e2tInstanceFullAddress                   = "10.0.2.15:9999"
	e2SetupMsgPrefix                         = e2tInstanceFullAddress + "|"
	GnbSetupRequestXmlPath                   = "../../tests/resources/setupRequest/setupRequest_gnb.xml"
)

func setupHealthCheckHandlerTest(t *testing.T) (*HealthCheckRequestHandler, services.RNibDataService, *mocks.RnibReaderMock, *mocks.RanListManagerMock, *mocks.RmrMessengerMock) {
	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	ranListManagerMock := &mocks.RanListManagerMock{}

	rmrSender := getRmrSender(rmrMessengerMock, logger)
	handler := NewHealthCheckRequestHandler(logger, rnibDataService, ranListManagerMock, rmrSender)

	return handler, rnibDataService, readerMock, ranListManagerMock, rmrMessengerMock
}

func TestHealthCheckRequestHandlerArguementHasRanNameSuccess(t *testing.T) {
	handler, _, readerMock, ranListManagerMock, rmrMessengerMock := setupHealthCheckHandlerTest(t)
	ranNames := []string{"RanName_1"}

	nb1:= createNbIdentity(t,"RanName_1", entities.ConnectionStatus_CONNECTED)
	oldnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	newnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}

	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)

	mbuf:= createRMRMbuf(t, nb1)
	rmrMessengerMock.On("SendMsg",mbuf,true).Return(mbuf,nil)
	ranListManagerMock.On("UpdateHealthcheckTimeStampSent",nb1.RanName).Return(oldnbIdentity, newnbIdentity)
	ranListManagerMock.On("UpdateNbIdentities",nb1.NodeType, []*entities.NbIdentity{oldnbIdentity}, []*entities.NbIdentity{newnbIdentity}).Return(nil)

	resp, err := handler.Handle(models.HealthCheckRequest{ranNames})

	assert.IsType(t, &models.HealthCheckSuccessResponse{}, resp)
	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
}

func TestHealthCheckRequestHandlerArguementHasNoRanNameSuccess(t *testing.T) {
	handler, _, readerMock, ranListManagerMock, rmrMessengerMock := setupHealthCheckHandlerTest(t)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED},
		{InventoryName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}}

	ranListManagerMock.On("GetNbIdentityList").Return(nbIdentityList)

	nb1:= createNbIdentity(t,"RanName_1", entities.ConnectionStatus_CONNECTED)
	oldnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	newnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}

	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)

	mbuf:= createRMRMbuf(t, nb1)
	rmrMessengerMock.On("SendMsg",mbuf,true).Return(mbuf,nil)
	ranListManagerMock.On("UpdateHealthcheckTimeStampSent",nb1.RanName).Return(oldnbIdentity, newnbIdentity)
	ranListManagerMock.On("UpdateNbIdentities",nb1.NodeType, []*entities.NbIdentity{oldnbIdentity}, []*entities.NbIdentity{newnbIdentity}).Return(nil)

	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)

	resp, err := handler.Handle(models.HealthCheckRequest{[]string{}})

	assert.Nil(t, err)
	assert.IsType(t, &models.HealthCheckSuccessResponse{}, resp)

}

func TestHealthCheckRequestHandlerArguementHasNoRanConnectedFailure(t *testing.T) {
	handler, _, readerMock, ranListManagerMock, rmrMessengerMock := setupHealthCheckHandlerTest(t)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED},
		{InventoryName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}}
	ranListManagerMock.On("GetNbIdentityList").Return(nbIdentityList)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)

	_, err := handler.Handle(models.HealthCheckRequest{[]string{}})

	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
	ranListManagerMock.AssertNotCalled(t,"UpdateHealthcheckTimeStampSent",mock.Anything)
	ranListManagerMock.AssertNotCalled(t,"UpdateNbIdentities",mock.Anything, mock.Anything, mock.Anything)
	assert.IsType(t, &e2managererrors.NoConnectedRanError{}, err)

}

func TestHealthCheckRequestHandlerArguementHasRanNameDBErrorFailure(t *testing.T) {
	handler, _, readerMock, ranListManagerMock, rmrMessengerMock := setupHealthCheckHandlerTest(t)

	ranNames := []string{"RanName_1"}
	readerMock.On("GetNodeb", "RanName_1").Return(&entities.NodebInfo{}, errors.New("error"))

	_, err := handler.Handle(models.HealthCheckRequest{ranNames})

	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
	ranListManagerMock.AssertNotCalled(t,"UpdateHealthcheckTimeStampSent",mock.Anything)
	ranListManagerMock.AssertNotCalled(t,"UpdateNbIdentities",mock.Anything, mock.Anything, mock.Anything)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
}

func createRMRMbuf(t *testing.T, nodebInfo *entities.NodebInfo) *rmrCgo.MBuf{
	serviceQuery := models.NewRicServiceQueryMessage(nodebInfo.GetGnb().RanFunctions)
	payLoad, err := xml.Marshal(&serviceQuery.E2APPDU)
	payLoad = utils.NormalizeXml(payLoad)
	tagsToReplace := []string{"reject","ignore","protocolIEs"}
	payLoad = utils.ReplaceEmptyTagsWithSelfClosing(payLoad, tagsToReplace)

	if err != nil {
		t.Fatal(err)
	}

	var xAction []byte
	var msgSrc unsafe.Pointer

	rmrMessage := models.NewRmrMessage(rmrCgo.RIC_SERVICE_QUERY, nodebInfo.RanName, payLoad, xAction, msgSrc)
	return rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
}

func createNbIdentity(t *testing.T, RanName string,  connectionStatus entities.ConnectionStatus) *entities.NodebInfo {
	xmlgnb := utils.ReadXmlFile(t, GnbSetupRequestXmlPath)
	payload := append([]byte(e2SetupMsgPrefix), xmlgnb...)
	pipInd := bytes.IndexByte(payload, '|')
	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)
	if err != nil {
		t.Fatal(err)
	}

	nodeb := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		RanName:                      RanName,
		SetupFromNetwork:             true,
		NodeType:                     entities.Node_GNB,
		ConnectionStatus: 			  connectionStatus,
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

func normalizeXml(payload []byte) []byte {
	xmlStr := string(payload)
	normalized := strings.NewReplacer("&lt;", "<", "&gt;", ">",
		"<reject></reject>","<reject/>","<ignore></ignore>","<ignore/>",
		"<protocolIEs></protocolIEs>","<protocolIEs/>").Replace(xmlStr)
	return []byte(normalized)
}
