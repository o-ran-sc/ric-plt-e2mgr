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

func initChangeStatusToConnectedRanTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *ChangeStatusToConnectedRanManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
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
	changeStatusToConnectedRanManager := NewChangeStatusToConnectedRanManager(logger, rnibDataService, ranConnectStatusChangeManager)
	return logger, rmrMessengerMock, readerMock, writerMock, changeStatusToConnectedRanManager
}

func TestChangeStatusToConnectedRanSucceeds(t *testing.T) {
	logger, _, readerMock, writerMock, changeStatusToConnectedRanManager := initChangeStatusToConnectedRanTest(t)
	logger.Infof("#TestChangeStatusToConnectedRanManager.ConnectedRan - RAN name")
	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	writerMock.On("UpdateNodebInfoOnConnectionStatusInversion", mock.Anything, "test_CONNECTED").Return(rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(rnibErr)
	_, err := changeStatusToConnectedRanManager.ChangeStatusToConnectedRan(ranName)
	assert.Nil(t, err)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}
