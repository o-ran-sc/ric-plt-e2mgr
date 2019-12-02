package notificationmanager

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initNotificationManagerTest(t *testing.T) (*logger.Logger, *mocks.RnibReaderMock, *NotificationManager) {
	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	httpClient := &mocks.HttpClientMock{}

	rmrSender := initRmrSender(&mocks.RmrMessengerMock{}, logger)
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClient)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider()
	rmrNotificationHandlerProvider.Init(logger, config, rnibDataService, rmrSender, ranSetupManager, e2tInstancesManager,routingManagerClient)
	notificationManager := NewNotificationManager(logger, rmrNotificationHandlerProvider )
	return logger, readerMock, notificationManager
}

func TestHandleMessageUnexistingMessageType(t *testing.T) {
	_, _, nm := initNotificationManagerTest(t)

	mbuf := &rmrCgo.MBuf{MType: 1234}

	err := nm.HandleMessage(mbuf)
	assert.NotNil(t, err)
}

func TestHandleMessageExistingMessageType(t *testing.T) {
	_, readerMock, nm := initNotificationManagerTest(t)
	payload := []byte("123")
	xaction := []byte("test")
	mbuf := &rmrCgo.MBuf{MType: rmrCgo.RIC_X2_SETUP_RESP, Meid: "test", Payload: &payload, XAction: &xaction}
	readerMock.On("GetNodeb", "test").Return(&entities.NodebInfo{}, fmt.Errorf("Some error"))
	err := nm.HandleMessage(mbuf)
	assert.Nil(t, err)
}

// TODO: extract to test_utils
func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
