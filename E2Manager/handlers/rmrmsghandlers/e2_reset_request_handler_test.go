//
// Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.
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

	// "e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"e2mgr/utils"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	E2ResetXmlPath = "../../tests/resources/reset/reset-request.xml"
)

func initE2ResetMocks(t *testing.T) (*E2ResetRequestNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{
		RnibRetryIntervalMs:       10,
		MaxRnibConnectionAttempts: 3,
		E2ResetTimeOutSec:         10,
		RnibWriter: configuration.RnibWriterConfig{
			StateChangeMessageChannel: StateChangeMessageChannel,
		}}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	handler := NewE2ResetRequestNotificationHandler(logger, rnibDataService, config, rmrSender)
	return handler, readerMock, writerMock, rmrMessengerMock
}

func TestE2ResettNotificationHandler(t *testing.T) {
	e2ResetXml := utils.ReadXmlFile(t, E2ResetXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock := initE2ResetMocks(t)
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
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(""), e2ResetXml...)}
	handler.Handle(notificationRequest)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestE2ResettNotificationHandler_UpdateStatus_Connected(t *testing.T) {
	e2ResetXml := utils.ReadXmlFile(t, E2ResetXmlPath)
	handler, readerMock, writerMock, rmrMessengerMock := initE2ResetMocks(t)
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
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	readerMock.On("GetNodeb", gnbNodebRanName).Return(nodebInfo, nil)

	var errEmpty error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, errEmpty)
	notificationRequest := &models.NotificationRequest{RanName: gnbNodebRanName, Payload: append([]byte(""), e2ResetXml...)}
	handler.Handle(notificationRequest)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	writerMock.AssertCalled(t, "UpdateNodebInfoAndPublish", mock.Anything)
	readerMock.AssertCalled(t, "GetNodeb", mock.Anything)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}
