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

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"e2mgr/utils"
	"encoding/xml"
	"fmt"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	serviceUpdateE2tInstanceAddress = "10.0.0.27:9999"
	serviceUpdateE2SetupMsgPrefix   = serviceUpdateE2tInstanceAddress + "|"
	serviceUpdateRANName            = "gnb:TestRan"
	RanManipulationMessageChannel   = "RAN_MANIPULATION"
	RICServiceUpdate_E2SetupReqPath = "../../tests/resources/serviceUpdate/RicServiceUpdate_SetupRequest.xml"
	RicServiceUpdateModifiedPath    = "../../tests/resources/serviceUpdate/RicServiceUpdate_ModifiedFunction.xml"
	RicServiceUpdateDeletePath      = "../../tests/resources/serviceUpdate/RicServiceUpdate_DeleteFunction.xml"
	RicServiceUpdateAddedPath       = "../../tests/resources/serviceUpdate/RicServiceUpdate_AddedFunction.xml"
	RicServiceUpdateEmptyPath       = "../../tests/resources/serviceUpdate/RicServiceUpdate_Empty.xml"
	RicServiceUpdateAckModifiedPath = "../../tests/resources/serviceUpdateAck/RicServiceUpdateAck_ModifiedFunction.xml"
	RicServiceUpdateAckAddedPath    = "../../tests/resources/serviceUpdateAck/RicServiceUpdateAck_AddedFunction.xml"
	RicServiceUpdateAckDeletePath   = "../../tests/resources/serviceUpdateAck/RicServiceUpdateAck_DeleteFunction.xml"
	RicServiceUpdateAckEmptyPath    = "../../tests/resources/serviceUpdateAck/RicServiceUpdateAck_Empty.xml"
)

func initRicServiceUpdateHandler(t *testing.T) (*RicServiceUpdateHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.RanListManagerMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{
		RnibRetryIntervalMs:       10,
		MaxRnibConnectionAttempts: 3,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel:     StateChangeMessageChannel,
			RanManipulationMessageChannel: RanManipulationMessageChannel,
		},
		GlobalRicId: struct {
			RicId string
			Mcc   string
			Mnc   string
		}{
			Mcc:   "337",
			Mnc:   "94",
			RicId: "AACCE",
		}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranListManagerMock := &mocks.RanListManagerMock{}
	handler := NewRicServiceUpdateHandler(logger, rmrSender, rnibDataService, ranListManagerMock)
	return handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock
}

func TestRICServiceUpdateModifiedFuncSuccess(t *testing.T) {
	testServiceUpdateSuccess(t, RicServiceUpdateModifiedPath, RicServiceUpdateAckModifiedPath)
}

func TestRICServiceUpdateAddedFuncSuccess(t *testing.T) {

	testServiceUpdateSuccess(t, RicServiceUpdateAddedPath, RicServiceUpdateAckAddedPath)
}

func TestRICServiceUpdateDeleteFuncSuccess(t *testing.T) {
	testServiceUpdateSuccess(t, RicServiceUpdateDeletePath, RicServiceUpdateAckDeletePath)
}

func TestRICServiceUpdateRnibFailure(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, RicServiceUpdateDeletePath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	readerMock.On("GetNodeb", serviceUpdateRANName).Return(&entities.NodebInfo{}, common.NewInternalError(fmt.Errorf("internal error")))
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName, Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func TestRICServiceUpdateRnibNotFound(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, RicServiceUpdateModifiedPath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	readerMock.On("GetNodeb", serviceUpdateRANName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError("nodeb not found"))
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName, Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func TestRICServiceUpdateNodeBInfoFailure(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, RicServiceUpdateDeletePath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	nb1 := createNbInfo(t, serviceUpdateRANName, entities.ConnectionStatus_CONNECTED)
	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName, Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(common.NewInternalError(fmt.Errorf("internal error")))

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func TestSendRICServiceUpdateAckFailure(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, RicServiceUpdateModifiedPath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	nb1 := createNbInfo(t, serviceUpdateRANName, entities.ConnectionStatus_CONNECTED)
	oldnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	newnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName, Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}
	ricServiceAckMsg := createRicServiceQueryAckRMRMbuf(t, RicServiceUpdateAckModifiedPath, notificationRequest)
	ranListManagerMock.On("UpdateHealthcheckTimeStampReceived", nb1.RanName).Return(oldnbIdentity, newnbIdentity)
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	rmrMessengerMock.On("SendMsg", ricServiceAckMsg, true).Return(&rmrCgo.MBuf{}, fmt.Errorf("rmr send failure"))
	ranListManagerMock.On("UpdateNbIdentities", nb1.NodeType, []*entities.NbIdentity{oldnbIdentity}, []*entities.NbIdentity{newnbIdentity}).Return(nil)

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func TestRICServiceUpdateUpdateNbIdentitiesFailure(t *testing.T) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, RicServiceUpdateDeletePath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	nb1 := createNbInfo(t, serviceUpdateRANName, entities.ConnectionStatus_CONNECTED)
	oldnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	newnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName, Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}
	ranListManagerMock.On("UpdateHealthcheckTimeStampReceived", nb1.RanName).Return(oldnbIdentity, newnbIdentity)
	ranListManagerMock.On("UpdateNbIdentities", nb1.NodeType, []*entities.NbIdentity{oldnbIdentity}, []*entities.NbIdentity{newnbIdentity}).Return(common.NewInternalError(fmt.Errorf("internal error")))
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func TestRICServiceUpdateParseRequest_PipFailure(t *testing.T) {
	xmlGnb := utils.ReadXmlFile(t, RICServiceUpdate_E2SetupReqPath)
	handler, _, _, _, _ := initRicServiceUpdateHandler(t)
	prefBytes := []byte(serviceUpdateE2tInstanceAddress)
	ricServiceUpdate, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Nil(t, ricServiceUpdate)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#RicServiceUpdateHandler.parseSetupRequest - Error parsing RIC SERVICE UPDATE failed extract Payload: no | separator found")
}

func TestRICServiceUppdateParseRequest_UnmarshalFailure(t *testing.T) {
	handler, _, _, _, _ := initRicServiceUpdateHandler(t)
	prefBytes := []byte(serviceUpdateE2SetupMsgPrefix)
	ricServiceUpdate, err := handler.parseSetupRequest(append(prefBytes, 1, 2, 3))
	assert.Nil(t, ricServiceUpdate)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#RicServiceUpdateHandler.parseSetupRequest - Error unmarshalling RIC SERVICE UPDATE payload: 31302e302e302e32373a393939397c010203")
}

func testServiceUpdateSuccess(t *testing.T, servicepdatePath string, serviceUpdateAckPath string) {
	handler, readerMock, writerMock, rmrMessengerMock, ranListManagerMock := initRicServiceUpdateHandler(t)
	xmlserviceUpdate := utils.ReadXmlFile(t, servicepdatePath)
	xmlserviceUpdate = utils.CleanXML(xmlserviceUpdate)
	nb1 := createNbInfo(t, serviceUpdateRANName, entities.ConnectionStatus_CONNECTED)
	oldnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	newnbIdentity := &entities.NbIdentity{InventoryName: nb1.RanName, ConnectionStatus: nb1.ConnectionStatus}
	readerMock.On("GetNodeb", nb1.RanName).Return(nb1, nil)
	notificationRequest := &models.NotificationRequest{RanName: serviceUpdateRANName,
		Payload: append([]byte(serviceUpdateE2SetupMsgPrefix), xmlserviceUpdate...)}
	ricServiceAckMsg := createRicServiceQueryAckRMRMbuf(t, serviceUpdateAckPath, notificationRequest)
	ranListManagerMock.On("UpdateHealthcheckTimeStampReceived", nb1.RanName).Return(oldnbIdentity, newnbIdentity)
	ranListManagerMock.On("UpdateNbIdentities", nb1.NodeType, []*entities.NbIdentity{oldnbIdentity},
		[]*entities.NbIdentity{newnbIdentity}).Return(nil)
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	rmrMessengerMock.On("SendMsg", ricServiceAckMsg, true).Return(&rmrCgo.MBuf{}, nil)

	handler.Handle(notificationRequest)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
	readerMock.AssertExpectations(t)
	ranListManagerMock.AssertExpectations(t)
}

func createRicServiceQueryAckRMRMbuf(t *testing.T, xmlFile string, req *models.NotificationRequest) *rmrCgo.MBuf {
	ricServiceQueryAckXml := utils.ReadXmlFile(t, xmlFile)
	ricServiceQueryAckXml = utils.CleanXML(ricServiceQueryAckXml)
	payLoad := utils.NormalizeXml(ricServiceQueryAckXml)

	xAction := req.TransactionId
	msgsrc := req.GetMsgSrc()

	rmrMessage := models.NewRmrMessage(rmrCgo.RIC_SERVICE_UPDATE_ACK, serviceUpdateRANName, payLoad, xAction, msgsrc)
	return rmrCgo.NewMBuf(rmrMessage.MsgType, len(rmrMessage.Payload), rmrMessage.RanName, &rmrMessage.Payload, &rmrMessage.XAction, rmrMessage.GetMsgSrc())
}

func createNbInfo(t *testing.T, RanName string, connectionStatus entities.ConnectionStatus) *entities.NodebInfo {
	xmlgnb := utils.ReadXmlFile(t, RICServiceUpdate_E2SetupReqPath)
	xmlgnb = utils.CleanXML(xmlgnb)
	payload := append([]byte(serviceUpdateE2SetupMsgPrefix), xmlgnb...)
	pipInd := bytes.IndexByte(payload, '|')
	setupRequest := &models.E2SetupRequestMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(payload[pipInd+1:]), &setupRequest.E2APPDU)
	if err != nil {
		t.Fatal(err)
	}

	nodeb := &entities.NodebInfo{
		AssociatedE2TInstanceAddress: serviceUpdateE2tInstanceAddress,
		RanName:                      RanName,
		SetupFromNetwork:             true,
		NodeType:                     entities.Node_GNB,
		ConnectionStatus:             connectionStatus,
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
