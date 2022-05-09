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
	"e2mgr/configuration"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"e2mgr/utils"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	E2nodeConfigUpdateXmlPath             = "../../tests/resources/configurationUpdate/e2NodeConfigurationUpdate.xml"
	E2nodeConfigUpdateOnlyAdditionXmlPath = "../../tests/resources/configurationUpdate/e2NodeConfigurationUpdateOnlyAddition.xml"
)

func initE2nodeConfigMocks(t *testing.T) (*E2nodeConfigUpdateNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock) {
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
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	handler := NewE2nodeConfigUpdateNotificationHandler(logger, rnibDataService, rmrSender)
	return handler, readerMock, writerMock, rmrMessengerMock
}

func TestE2nodeConfigUpdatetNotificationHandler(t *testing.T) {
	e2NodeConfigUpdateXml := utils.ReadXmlFile(t, E2nodeConfigUpdateOnlyAdditionXmlPath)
	handler, readerMock, writerMock, _ := initE2nodeConfigMocks(t)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      gnbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_DISCONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{
			Gnb: &entities.Gnb{},
		},
	}
	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(""), e2NodeConfigUpdateXml...)}
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestE2nodeConfigUpdatetParseFail(t *testing.T) {
	handler, _, _, _ := initE2nodeConfigMocks(t)
	badxml := []byte("abc")
	e2nodeConfig, err := handler.parseE2NodeConfigurationUpdate(badxml)

	var expected *models.E2nodeConfigurationUpdateMessage
	assert.Equal(t, expected, e2nodeConfig)
	assert.NotNil(t, err)
}

func TestHandleAddConfig(t *testing.T) {
	e2NodeConfigUpdateXml := utils.ReadXmlFile(t, E2nodeConfigUpdateXmlPath)

	handler, readerMock, writerMock, rmrMessengerMock := initE2nodeConfigMocks(t)
	var nodebInfo = &entities.NodebInfo{
		RanName:                      gnbNodebRanName,
		AssociatedE2TInstanceAddress: e2tInstanceFullAddress,
		ConnectionStatus:             entities.ConnectionStatus_DISCONNECTED,
		NodeType:                     entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{
			Gnb: &entities.Gnb{},
		},
	}
	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	var errEmpty error
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(&rmrCgo.MBuf{}, errEmpty)

	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(""), e2NodeConfigUpdateXml...)}

	handler.Handle(notificationRequest)

	t.Logf("len of addtionList : %d", len(nodebInfo.GetGnb().NodeConfigs))

	assert.Equal(t, 5, len(nodebInfo.GetGnb().NodeConfigs))
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)
}
