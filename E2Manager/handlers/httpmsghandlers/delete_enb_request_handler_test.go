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


package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupDeleteEnbRequestHandlerTest(t *testing.T) (*DeleteEnbRequestHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	handler := NewDeleteEnbRequestHandler(log, rnibDataService)
	return handler, readerMock, writerMock
}

func TestHandleDeleteEnbSuccess(t *testing.T) {
	handler, readerMock, writerMock := setupDeleteEnbRequestHandlerTest(t)

	ranName := "test1"
	var rnibError error
	nodebInfo := &entities.NodebInfo{RanName: ranName, NodeType: entities.Node_ENB}
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, rnibError)
	writerMock.On("RemoveEnb", nodebInfo).Return(nil)
	result, err := handler.Handle(&models.DeleteEnbRequest{RanName: ranName})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, &models.DeleteEnbResponse{}, result)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestHandleDeleteEnbInternalGetNodebError(t *testing.T) {
	handler, readerMock, writerMock := setupDeleteEnbRequestHandlerTest(t)

	ranName := "test1"
	rnibError := errors.New("for test")
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, rnibError)
	result, err := handler.Handle(&models.DeleteEnbRequest{RanName: ranName})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestHandleDeleteEnbInternalRemoveEnbError(t *testing.T) {
	handler, readerMock, writerMock := setupDeleteEnbRequestHandlerTest(t)

	ranName := "test1"
	rnibError := errors.New("for test")
	nodebInfo  := &entities.NodebInfo{RanName: ranName, NodeType: entities.Node_ENB}
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, nil)
	writerMock.On("RemoveEnb", nodebInfo).Return(rnibError)
	result, err := handler.Handle(&models.DeleteEnbRequest{RanName: ranName})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestHandleDeleteEnbResourceNotFoundError(t *testing.T) {
	handler, readerMock, writerMock := setupDeleteEnbRequestHandlerTest(t)

	ranName := "test1"
	rnibError := common.NewResourceNotFoundError("for test")
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, rnibError)
	result, err := handler.Handle(&models.DeleteEnbRequest{RanName: ranName})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}

func TestHandleDeleteEnbNodeTypeNotEnbError(t *testing.T) {
	handler, readerMock, writerMock := setupDeleteEnbRequestHandlerTest(t)

	ranName := "test1"
	nodebInfo  := &entities.NodebInfo{RanName: ranName, NodeType: entities.Node_GNB}
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, nil)
	result, err := handler.Handle(&models.DeleteEnbRequest{RanName: ranName})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
}