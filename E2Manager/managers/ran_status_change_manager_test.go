package managers

import (
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func initRanStatusChangeManagerTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *rmrsender.RmrSender) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Fatalf("#initStatusChangeManagerTest - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	return logger, rmrMessengerMock, rmrSender
}

func TestMarshalFailure(t *testing.T) {
	logger, _, rmrSender := initRanStatusChangeManagerTest(t)
	m := NewRanStatusChangeManager(logger, rmrSender)

	nodebInfo := entities.NodebInfo{}
	err := m.Execute(123, 4, &nodebInfo)

	assert.NotNil(t, err)
}

func TestMarshalSuccess(t *testing.T) {
	logger, rmrMessengerMock, rmrSender := initRanStatusChangeManagerTest(t)
	m := NewRanStatusChangeManager(logger, rmrSender)

	nodebInfo := entities.NodebInfo{NodeType: entities.Node_ENB}
	var err error
	rmrMessengerMock.On("SendMsg", mock.Anything).Return(&rmrCgo.MBuf{}, err)
	err  = m.Execute(rmrCgo.RAN_CONNECTED, enums.RIC_TO_RAN, &nodebInfo)

	assert.Nil(t, err)
}
