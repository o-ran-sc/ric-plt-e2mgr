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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupSetGeneralConfigurationHandlerTest(t *testing.T) (*SetGeneralConfigurationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	handler := NewSetGeneralConfigurationHandler(log, rnibDataService)
	return handler, readerMock, writerMock
}

func TestSetGeneralConfigurationFalse_Success(t *testing.T) {
	handler, readerMock, writerMock := setupSetGeneralConfigurationHandlerTest(t)

	configuration := &entities.GeneralConfiguration{EnableRic: true}
	readerMock.On("GetGeneralConfiguration").Return(configuration, nil)

	updated := &entities.GeneralConfiguration{EnableRic: false}
	writerMock.On("SaveGeneralConfiguration", updated).Return(nil)

	response, err := handler.Handle(models.GeneralConfigurationRequest{EnableRic: false})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.GeneralConfigurationResponse{}, response)
}

func TestSetGeneralConfigurationTrue_Success(t *testing.T) {
	handler, readerMock, writerMock := setupSetGeneralConfigurationHandlerTest(t)

	configuration := &entities.GeneralConfiguration{EnableRic: false}
	readerMock.On("GetGeneralConfiguration").Return(configuration, nil)

	updated := &entities.GeneralConfiguration{EnableRic: true}
	writerMock.On("SaveGeneralConfiguration", updated).Return(nil)

	response, err := handler.Handle(models.GeneralConfigurationRequest{EnableRic: true})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.GeneralConfigurationResponse{}, response)
}

func TestSetGeneralConfigurationIgnore_Success(t *testing.T) {
	handler, readerMock, _ := setupSetGeneralConfigurationHandlerTest(t)

	configuration := &entities.GeneralConfiguration{EnableRic: false}
	readerMock.On("GetGeneralConfiguration").Return(configuration, nil)

	response, err := handler.Handle(models.GeneralConfigurationRequest{EnableRic: false})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.GeneralConfigurationResponse{}, response)
}

func TestSetGeneralConfigurationHandlerRnibError(t *testing.T) {
	handler, readerMock, writerMock := setupSetGeneralConfigurationHandlerTest(t)

	configuration := &entities.GeneralConfiguration{EnableRic: false}
	readerMock.On("GetGeneralConfiguration").Return(configuration, nil)

	updated := &entities.GeneralConfiguration{EnableRic: true}
	writerMock.On("SaveGeneralConfiguration", updated).Return(common.NewInternalError(errors.New("error")))

	response, err := handler.Handle(models.GeneralConfigurationRequest{EnableRic: true})

	assert.NotNil(t, err)
	assert.Nil(t, response)
}
