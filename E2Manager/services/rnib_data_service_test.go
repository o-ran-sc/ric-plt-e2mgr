package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rNibWriter"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

func setupTest(t *testing.T) (*rNibDataService, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}

	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}

	rnibDataService := NewRnibDataService(logger, config, rnibReaderProvider, rnibWriterProvider)
	assert.NotNil(t, rnibDataService)

	return rnibDataService, readerMock, writerMock
}

func TestSuccessfulSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	mockErr := common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 3)
}

func TestNonConnFailureSaveNodeb(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	nodebInfo := &entities.NodebInfo{}
	nbIdentity := &entities.NbIdentity{}
	mockErr := common.InternalError{Err: fmt.Errorf("non connection failure")}
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(mockErr)

	rnibDataService.SaveNodeb(nbIdentity, nodebInfo)
	writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestSuccessfulUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	nodebInfo := &entities.NodebInfo{}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(nil)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnFailureUpdateNodebInfo(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	nodebInfo := &entities.NodebInfo{}
	mockErr := common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("UpdateNodebInfo", nodebInfo).Return(mockErr)

	rnibDataService.UpdateNodebInfo(nodebInfo)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 3)
}

func TestSuccessfulSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(nil)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 1)
}

func TestConnFailureSaveRanLoadInformation(t *testing.T) {
	rnibDataService, _, writerMock := setupTest(t)

	var ranName string = "abcd"
	ranLoadInformation := &entities.RanLoadInformation{}
	mockErr := common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	writerMock.On("SaveRanLoadInformation", ranName, ranLoadInformation).Return(mockErr)

	rnibDataService.SaveRanLoadInformation(ranName, ranLoadInformation)
	writerMock.AssertNumberOfCalls(t, "SaveRanLoadInformation", 3)
}

func TestSuccessfulGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupTest(t)

	invName := "abcd"
	nodebInfo := &entities.NodebInfo{}
	readerMock.On("GetNodeb", invName).Return(nodebInfo, nil)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
	assert.Equal(t, nodebInfo, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodeb(t *testing.T) {
	rnibDataService, readerMock, _ := setupTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	mockErr := common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeb, res)
}

func TestSuccessfulGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupTest(t)

	nodeIds := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.Equal(t, nodeIds, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock, _ := setupTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error", ))
	assert.Equal(t, nodeIds, res)
}

//func TestConnFailureThenSuccessGetNodebIdList(t *testing.T) {
//	rnibDataService, readerMock, _ := setupTest(t)
//
//	var nilNodeIds []*entities.NbIdentity = nil
//	nodeIds := []*entities.NbIdentity{}
//	mockErr := common.InternalError{Err: &net.OpError{Err:fmt.Errorf("connection error")}}
//	//readerMock.On("GetListNodebIds").Return(nilNodeIds, mockErr)
//	//readerMock.On("GetListNodebIds").Return(nodeIds, nil)
//
//	res, err := rnibDataService.GetListNodebIds()
//	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 2)
//	assert.True(t, strings.Contains(err.Error(),"connection failure", ))
//	assert.Equal(t, nodeIds, res)
//}
