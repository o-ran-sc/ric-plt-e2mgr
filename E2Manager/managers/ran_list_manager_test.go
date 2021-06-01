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
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initRanListManagerTest(t *testing.T) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, RanListManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Fatalf("#... - failed to initialize logger, error: %s", err)
	}

	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3,
		RnibWriter: configuration.RnibWriterConfig{StateChangeMessageChannel: "RAN_CONNECTION_STATUS_CHANGE"}}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranListManager := NewRanListManager(logger, rnibDataService)
	return readerMock, writerMock, ranListManager
}

func TestRanListManagerInstance_InitNbIdentityMapSuccess(t *testing.T) {
	readerMock, _, ranListManager := initRanListManagerTest(t)
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{{InventoryName: RanName, GlobalNbId: &entities.GlobalNbId{NbId: "asd", PlmnId: "efg"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED}}, nil)
	err := ranListManager.InitNbIdentityMap()
	assert.Nil(t, err)
}

func TestRanListManagerInstance_InitNbIdentityMapFailure(t *testing.T) {
	readerMock, _, ranListManager := initRanListManagerTest(t)
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, common.NewInternalError(errors.New("#reader.GetListNodebIds - Internal Error")))
	err := ranListManager.InitNbIdentityMap()
	assert.NotNil(t, err)
}

func TestRanListManagerInstance_AddNbIdentitySuccess(t *testing.T) {
        _,writerMock, ranListManager := initRanListManagerTest(t)
        //readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{{InventoryName: RanName, GlobalNbId: &entities.GlobalNbId{NbId: "asd", PlmnId: "efg"}, ConnectionStatus: entities.ConnectionStatus_CONNECTED}}, nil)
        nbIdentity :=  &entities.NbIdentity{}
        writerMock.On("AddNbIdentity", entities.Node_ENB, nbIdentity).Return(nil)
        nodetype := entities.Node_ENB
        err := ranListManager.AddNbIdentity(nodetype,nbIdentity)
        assert.Nil(t, err)
}


func TestRanListManagerInstance_RemoveNbIdentitySuccess(t *testing.T) {
        _,writerMock, ranListManager := initRanListManagerTest(t)
        ranName := "ran1"
        writerMock.On("RemoveNbIdentity", entities.Node_ENB,"ran1" ).Return(nil)
        nodetype := entities.Node_ENB
        err := ranListManager.RemoveNbIdentity(nodetype,ranName)
        assert.Nil(t, err)
}

func TestRanListManagerInstance_GetNbIdentity(t *testing.T) {
       _,writerMock,ranListManager := initRanListManagerTest(t)
       ranName := "ran1"
       nbIdentity := &entities.NbIdentity{}
       writerMock.On("GetNbIdentity", entities.Node_ENB, nbIdentity).Return(nil)
       err,nb := ranListManager.GetNbIdentity(ranName)
       assert.NotNil(t, nb)
       assert.Nil(t,err)
}

func TestRanListManagerInstance_GetNbIdentityList(t *testing.T) {
       _,writerMock,ranListManager := initRanListManagerTest(t)
       writerMock.On("GetNbIdentityList").Return(nil)
       Ids := ranListManager.GetNbIdentityList()
       assert.NotNil(t, Ids)
}


func TestRanListManagerInstance_UpdateNbIdentities(t *testing.T) {
        _,writerMock,ranListManager := initRanListManagerTest(t)
        nodeType := entities.Node_ENB
        oldNbIdentities := []*entities.NbIdentity{}
        newNbIdentities := []*entities.NbIdentity{}
        writerMock.On("UpdateNbIdentities",entities.Node_ENB,oldNbIdentities,newNbIdentities).Return(nil)
        res := ranListManager.UpdateNbIdentities(nodeType, oldNbIdentities, newNbIdentities)
        assert.Nil(t, res)
}
