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

package managers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const RanName = "test"
const BaseRMUrl = "http://10.10.2.15:12020/routingmanager"

func initE2TAssociationManagerTest(t *testing.T) (*E2TAssociationManager, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	config.RoutingManager.BaseUrl = BaseRMUrl

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	manager := NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient)

	return manager, readerMock, writerMock, httpClientMock
}

func mockHttpClientAssociateRan(httpClientMock *mocks.HttpClientMock, isSuccessful bool) {
	data := models.NewRoutingManagerE2TData(E2TAddress, RanName)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := BaseRMUrl + clients.AssociateRanToE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	var respStatusCode int
	if isSuccessful {
		respStatusCode = http.StatusCreated
	} else {
		respStatusCode = http.StatusBadRequest
	}
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: respStatusCode, Body: respBody}, nil)
}

func TestAssociateRanSuccess(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionAttempts: 1}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = append(updatedE2tInstance.AssociatedRanList, RanName)
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanRoutingManagerError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, false)

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanGetNodebError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, true)
	var nb *entities.NodebInfo
	readerMock.On("GetNodeb", RanName).Return(nb, e2managererrors.NewRnibDbError())

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanUpdateNodebError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionAttempts: 1}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(e2managererrors.NewRnibDbError())

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanGetE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionAttempts: 1}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	var e2tInstance *entities.E2TInstance
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, errors.New("test"))

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanSaveE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientAssociateRan(httpClientMock, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionAttempts: 1}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.ConnectionAttempts = 0
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = append(updatedE2tInstance.AssociatedRanList, RanName)
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(errors.New("test"))

	err := manager.AssociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
