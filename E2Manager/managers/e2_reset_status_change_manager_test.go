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

package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initE2ResetStatusChangeTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *RanResetManager) {
	DebugLevel := int8(4)
	logger, err := logger.InitLogger(DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranListManager := NewRanListManager(logger, rnibDataService)
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2ResetStatusChangeManager := NewRanResetManager(logger, rnibDataService, ranConnectStatusChangeManager)
	return logger, rmrMessengerMock, readerMock, writerMock, e2ResetStatusChangeManager
}

func TestE2ResetStatusChangeSucceeds(t *testing.T) {
	logger, _, readerMock, writerMock, e2ResetStatusChangeManager := initE2ResetStatusChangeTest(t)
	logger.Infof("#TestRanResetManager.ResetRan - RAN name")
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, mock.Anything).Return(nil)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_UNDER_RESET
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	_, err := e2ResetStatusChangeManager.ResetRan(ranName)
	assert.Nil(t, err)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}
