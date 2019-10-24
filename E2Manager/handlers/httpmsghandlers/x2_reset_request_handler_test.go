package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupX2ResetRequestHandlerTest(t *testing.T) (*X2ResetRequestHandler, *mocks.RmrMessengerMock, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	rnibDataService := services.NewRnibDataService(log, config, readerProvider, writerProvider)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := getRmrSender(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(log, rmrSender, rnibDataService)

	return handler, rmrMessengerMock, readerMock
}
func TestHandleSuccessfulDefaultCause(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)

	rmrMessengerMock.On("SendMsg", msg).Return(msg, nil)

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.Nil(t, actual)
}

func TestHandleSuccessfulRequestedCause(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", msg).Return(msg, nil)

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName, Cause: "protocol:transfer-syntax-error"})

	assert.Nil(t, actual)
}

func TestHandleFailureUnknownCause(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName, Cause: "XXX"})

	assert.IsType(t, e2managererrors.NewRequestValidationError(), actual)

}

func TestHandleFailureWrongState(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewWrongStateError(X2_RESET_ACTIVITY_NAME, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)]), actual)
}

func TestHandleFailureRanNotFound(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"

	readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError("nodeb not found"))

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewResourceNotFoundError(), actual)
}

func TestHandleFailureRnibError(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"

	readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{}, common.NewInternalError(fmt.Errorf("internal error")))

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewRnibDbError(), actual)
}

func TestHandleFailureRmrError(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", msg).Return(&rmrCgo.MBuf{}, fmt.Errorf("rmr error"))

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewRmrError(), actual)
}
