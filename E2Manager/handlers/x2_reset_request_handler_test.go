package handlers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)


func TestHandleSuccessfulDefaultCause(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	// o&m intervention
	payload:= []byte {0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg:= rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg",msg,mock.Anything).Return(msg,nil)

	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	var nodeb = &entities.NodebInfo{ConnectionStatus:  entities.ConnectionStatus_CONNECTED }
	readerMock.On("GetNodeb",ranName).Return(nodeb, nil)

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName})

	assert.Nil(t, actual)

}

func TestHandleSuccessfulRequestedCause(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	payload:= []byte {0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	xaction := []byte(ranName)
	msg:= rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg",msg,mock.Anything).Return(msg,nil)

	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	var nodeb = &entities.NodebInfo{ConnectionStatus:  entities.ConnectionStatus_CONNECTED }
	readerMock.On("GetNodeb",ranName).Return(nodeb, nil)

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName , Cause:"protocol:transfer-syntax-error"})

	assert.Nil(t, actual)
}

func TestHandleFailureUnknownCause(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}


	rmrMessengerMock := &mocks.RmrMessengerMock{}


	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	var nodeb = &entities.NodebInfo{ConnectionStatus:  entities.ConnectionStatus_CONNECTED }
	readerMock.On("GetNodeb",ranName).Return(nodeb, nil)

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName , Cause:"XXX"})

	assert.IsType(t, e2managererrors.NewRequestValidationError(), actual)

}

func TestHandleFailureWrongState(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}


	rmrMessengerMock := &mocks.RmrMessengerMock{}


	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	var nodeb = &entities.NodebInfo{ConnectionStatus:  entities.ConnectionStatus_DISCONNECTED }
	readerMock.On("GetNodeb",ranName).Return(nodeb, nil)

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName })

	assert.IsType(t, e2managererrors.NewWrongStateError(X2_RESET_ACTIVITY_NAME, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)]), actual)

}



func TestHandleFailureRanNotFound(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}


	rmrMessengerMock := &mocks.RmrMessengerMock{}


	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	readerMock.On("GetNodeb",ranName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(fmt.Errorf("nodeb not found")))

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName })

	assert.IsType(t, e2managererrors.NewResourceNotFoundError(), actual)

}


func TestHandleFailureRnibError(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}


	rmrMessengerMock := &mocks.RmrMessengerMock{}


	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	readerMock.On("GetNodeb",ranName).Return(&entities.NodebInfo{}, common.NewInternalError(fmt.Errorf("internal error")))

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName })

	assert.IsType(t, e2managererrors.NewRnibDbError(), actual)

}


func TestHandleFailureRmrError(t *testing.T){
	log := initLog(t)

	ranName := "test1"

	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	writerProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	// o&m intervention
	payload:= []byte {0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	xaction := []byte(ranName)
	msg:= rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("SendMsg",msg,mock.Anything).Return(&rmrCgo.MBuf{},fmt.Errorf("rmr error"))

	config := configuration.ParseConfiguration()
	rmrService:=getRmrService(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(rmrService, config, writerProvider, readerProvider)

	var nodeb = &entities.NodebInfo{ConnectionStatus:  entities.ConnectionStatus_CONNECTED }
	readerMock.On("GetNodeb",ranName).Return(nodeb, nil)

	actual := handler.Handle(log, models.ResetRequest{RanName: ranName })

	assert.IsType(t, e2managererrors.NewRmrError(), actual)

}
