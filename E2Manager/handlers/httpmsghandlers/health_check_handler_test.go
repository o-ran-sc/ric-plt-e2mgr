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
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupHealthCheckHandlerTest(t *testing.T) (*HealthCheckRequestHandler, services.RNibDataService, *mocks.RnibReaderMock, *mocks.RanListManagerMock) {
	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	ranListManagerMock := &mocks.RanListManagerMock{}

	rmrSender := getRmrSender(rmrMessengerMock, logger)
	handler := NewHealthCheckRequestHandler(logger, rnibDataService, ranListManagerMock, rmrSender)

	return handler, rnibDataService, readerMock, ranListManagerMock
}

func TestHealthCheckRequestHandlerArguementHasRanNameSuccess(t *testing.T) {
	handler, _, readerMock, _ := setupHealthCheckHandlerTest(t)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	ranNames := []string{"RanName_1"}

	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	_, err := handler.Handle(models.HealthCheckRequest{ranNames})

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
}

func TestHealthCheckRequestHandlerArguementHasNoRanNameSuccess(t *testing.T) {
	handler, _, readerMock, ranListManagerMock := setupHealthCheckHandlerTest(t)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED},
		{InventoryName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}}

	ranListManagerMock.On("GetNbIdentityList").Return(nbIdentityList)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)

	_, err := handler.Handle(models.HealthCheckRequest{[]string{}})

	assert.Nil(t, err)

}

func TestHealthCheckRequestHandlerArguementHasNoRanConnectedFailure(t *testing.T) {
	handler, _, readerMock, ranListManagerMock := setupHealthCheckHandlerTest(t)

	nbIdentityList := []*entities.NbIdentity{{InventoryName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED},
		{InventoryName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}}
	ranListManagerMock.On("GetNbIdentityList").Return(nbIdentityList)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)

	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)

	_, err := handler.Handle(models.HealthCheckRequest{[]string{}})

	assert.IsType(t, &e2managererrors.NoConnectedRanError{}, err)

}

func TestHealthCheckRequestHandlerArguementHasRanNameDBErrorFailure(t *testing.T) {
	handler, _, readerMock, _ := setupHealthCheckHandlerTest(t)

	ranNames := []string{"RanName_1"}
	readerMock.On("GetNodeb", "RanName_1").Return(&entities.NodebInfo{}, errors.New("error"))

	_, err := handler.Handle(models.HealthCheckRequest{ranNames})

	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
}
