//
// Copyright 2023 Nokia
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
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"testing"
	"github.com/stretchr/testify/assert"
	"e2mgr/tests"
	"github.com/stretchr/testify/mock"
)

const (
	serviceUpdateRANName1 = "gnb:TestRan1"
    E2tAddress = "10.10.2.15:9800"
)


func initRicServiceUpdateManagerTest(t *testing.T) (*logger.Logger,*mocks.RnibReaderMock, *mocks.RnibWriterMock,services.RNibDataService, *configuration.Configuration, *RicServiceUpdateManager) {
	logger := tests.InitLog(t)

	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	RicServiceUpdateManager := NewRicServiceUpdateManager(logger, rnibDataService)
	return logger, readerMock, writerMock, rnibDataService, config, RicServiceUpdateManager
}
func TestUpdateRevertRanFunctions(t *testing.T) {

	_,readerMock, writerMock, _, _, RicServiceUpdateManager := initRicServiceUpdateManagerTest(t)
	InvName := "test"
	nodebInfo := &entities.NodebInfo{
		RanName: InvName,
		NodeType:                     entities.Node_GNB,
		Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}
	gnb := nodebInfo.GetGnb()
	gnb.RanFunctions = []*entities.RanFunction{{RanFunctionId: 2, RanFunctionRevision: 2}}
	readerMock.On("GetNodeb", InvName).Return(nodebInfo, nil)
	writerMock.On("UpdateNodebInfoAndPublish", mock.Anything).Return(nil)
	err := RicServiceUpdateManager.StoreExistingRanFunctions(ranName)
	assert.Nil(t, err)
	err = RicServiceUpdateManager.RevertRanFunctions(ranName)
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)
	readerMock.AssertCalled(t, "GetNodeb", InvName)
}
