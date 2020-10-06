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
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package httpmsghandlers

import (
	"e2mgr/mocks"
	"e2mgr/models"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"e2mgr/e2managererrors"
	"testing"
)

func setupGetNodebIdRequestHandlerTest(t *testing.T) (*GetNodebIdRequestHandler, *mocks.RanListManagerMock) {
	log := initLog(t)
	ranListManagerMock := &mocks.RanListManagerMock{}

	handler := NewGetNodebIdRequestHandler(log, ranListManagerMock)
	return handler, ranListManagerMock
}

func TestHandleGetNodebIdSuccess(t *testing.T) {
	handler, ranListManagerMock := setupGetNodebIdRequestHandlerTest(t)
	nbIdentity := &entities.NbIdentity{
		InventoryName:    "test",
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		HealthCheckTimestampSent: 12345678,
		HealthCheckTimestampReceived: 12346548,
	}
	ranListManagerMock.On("GetNbIdentity",nbIdentity.InventoryName).Return(nbIdentity, nil)
	response, err := handler.Handle(models.GetNodebIdRequest{RanName: nbIdentity.InventoryName})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.NodebIdResponse{}, response)
}

func TestHandleGetNodebIdNotFoundFailure(t *testing.T) {
	handler, ranListManagerMock := setupGetNodebIdRequestHandlerTest(t)
	nbIdentity := &entities.NbIdentity{
		InventoryName:    "test",
	}

	ranListManagerMock.On("GetNbIdentity",nbIdentity.InventoryName).Return(nbIdentity, e2managererrors.NewResourceNotFoundError())
	_, err := handler.Handle(models.GetNodebIdRequest{RanName: nbIdentity.InventoryName})

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.ResourceNotFoundError{}, err)
}