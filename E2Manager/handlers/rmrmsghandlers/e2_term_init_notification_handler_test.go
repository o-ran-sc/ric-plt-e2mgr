package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/mock"
	"testing"
)

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, E2TermInitNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock) {

	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rnibDataService := services.NewRnibDataService(logger, config, rnibReaderProvider, rnibWriterProvider)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	ranReconnectionManager := managers.NewRanReconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, ranSetupManager)
	handler := NewE2TermInitNotificationHandler(logger, ranReconnectionManager, rnibDataService)
	return logger, handler, readerMock, writerMock, rmrMessengerMock
}

func TestE2TerminInitHandlerSuccessOneRan(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}}
	readerMock.On("GetListNodebIds").Return(ids, rnibErr)

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ids[0].InventoryName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg, nil)

	handler.Handle(nil)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TerminInitHandlerSuccessTwoRans(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}}
	readerMock.On("GetListNodebIds").Return(ids, rnibErr)

	var initialNodeb0 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, rnibErr)

	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(ids[0].InventoryName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg, nil)

	handler.Handle(nil)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestE2TerminInitHandlerSuccessThreeRansFirstRmrFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}
	readerMock.On("GetListNodebIds").Return(ids, rnibErr)

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, rnibErr)
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Fail = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 0}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Fail).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg0, fmt.Errorf("RMR Error"))

	handler.Handle(nil)

	//test1 (before send +1, after failure +1), test2 (0) test3 (0)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	//test1 failure (+1), test2  (0). test3 (0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TerminInitHandlerSuccessThreeRansSecondNotFoundFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}
	readerMock.On("GetListNodebIds").Return(ids, rnibErr)

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, common.NewResourceNotFoundError("not found"))
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Success = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Success).Return(rnibErr)

	var argNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb2Success = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb2).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb2Success).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg0, nil)

	handler.Handle(nil)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	//test1 (+1), test2 failure (0) test3 (+1)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	//test1 success (+1), test2  (0). test3 (+1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestE2TerminInitHandlerSuccessThreeRansSecondRnibInternalErrorFailure(t *testing.T) {
	log, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	ids := []*entities.NbIdentity{{InventoryName: "test1"}, {InventoryName: "test2"}, {InventoryName: "test3"}}
	readerMock.On("GetListNodebIds").Return(ids, rnibErr)

	var initialNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb1 = &entities.NodebInfo{RanName: ids[1].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var initialNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", ids[0].InventoryName).Return(initialNodeb0, rnibErr)
	readerMock.On("GetNodeb", ids[1].InventoryName).Return(initialNodeb1, common.NewInternalError(fmt.Errorf("internal error")))
	readerMock.On("GetNodeb", ids[2].InventoryName).Return(initialNodeb2, rnibErr)

	var argNodeb0 = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb0Success = &entities.NodebInfo{RanName: ids[0].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb0).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb0Success).Return(rnibErr)

	var argNodeb2 = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	var argNodeb2Success = &entities.NodebInfo{RanName: ids[2].InventoryName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, ConnectionAttempts: 1}
	writerMock.On("UpdateNodebInfo", argNodeb2).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodeb2Success).Return(rnibErr)

	payload := models.NewE2RequestMessage(ids[0].InventoryName /*tid*/, "", 0, ids[0].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	xaction := []byte(ids[0].InventoryName)
	msg0 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[0].InventoryName, &payload, &xaction)

	// Cannot use Mock because request MBuf contains pointers
	//payload =models.NewE2RequestMessage(ids[1].InventoryName /*tid*/, "", 0,ids[1].InventoryName, e2pdus.PackedX2setupRequest).GetMessageAsBytes(log)
	//xaction = []byte(ids[1].InventoryName)
	//msg1 := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ids[1].InventoryName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", mock.Anything).Return(msg0, nil)

	handler.Handle(nil)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 2)
	//test1 (+1), test2 failure (0) test3 (0)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	//test1 success (+1), test2  (0). test3 (+1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestE2TerminInitHandlerSuccessZeroRans(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)
	var rnibErr error

	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, rnibErr)

	handler.Handle(nil)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestE2TerminInitHandlerFailureGetListNodebIds(t *testing.T) {
	_, handler, readerMock, writerMock, rmrMessengerMock := initRanLostConnectionTest(t)

	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, common.NewInternalError(fmt.Errorf("internal error")))

	handler.Handle(nil)

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 0)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

// TODO: extract to test_utils
func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, &rmrMessenger)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}
