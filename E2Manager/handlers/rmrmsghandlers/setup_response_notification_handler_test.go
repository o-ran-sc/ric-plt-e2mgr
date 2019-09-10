package rmrmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const (
	RanName                           = "test"
	X2SetupResponsePackedPdu          = "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	EndcSetupResponsePackedPdu        = "202400808e00000100f600808640000200fc00090002f829504a952a0a00fd007200010c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000211e148033e4e5e4c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a00021a0044033e4e5e000000002c001e3f271f2e3d4ff0031e3f274400010000150400000a00020000"
	X2SetupFailureResponsePackedPdu   = "4006001a0000030005400200000016400100001140087821a00000008040"
	EndcSetupFailureResponsePackedPdu = "4024001a0000030005400200000016400100001140087821a00000008040"
)

type setupResponseTestContext struct {
	logger               *logger.Logger
	readerMock           *mocks.RnibReaderMock
	writerMock           *mocks.RnibWriterMock
	rnibReaderProvider   func() reader.RNibReader
	rnibWriterProvider   func() rNibWriter.RNibWriter
	setupResponseManager managers.ISetupResponseManager
}

func NewSetupResponseTestContext(manager managers.ISetupResponseManager) *setupResponseTestContext {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	return &setupResponseTestContext{
		logger:     logger,
		readerMock: readerMock,
		writerMock: writerMock,
		rnibReaderProvider: func() reader.RNibReader {
			return readerMock
		},
		rnibWriterProvider: func() rNibWriter.RNibWriter {
			return writerMock
		},
		setupResponseManager: manager,
	}
}

func TestSetupResponseGetNodebFailure(t *testing.T) {
	notificationRequest := models.NotificationRequest{RanName: RanName}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.rnibReaderProvider, testContext.rnibWriterProvider, &managers.X2SetupResponseManager{}, "X2 Setup Response")
	testContext.readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewInternalError(errors.New("Error")))
	handler.Handle(testContext.logger, &notificationRequest, nil)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
}

func TestSetupResponseInvalidConnectionStatus(t *testing.T) {
	ranName := "test"
	notificationRequest := models.NotificationRequest{RanName: ranName}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.rnibReaderProvider, testContext.rnibWriterProvider, &managers.X2SetupResponseManager{}, "X2 Setup Response")
	var rnibErr error
	testContext.readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}, rnibErr)
	handler.Handle(testContext.logger, &notificationRequest, nil)
	testContext.readerMock.AssertCalled(t, "GetNodeb", ranName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
}

func executeHandleSuccessSetupResponse(t *testing.T, packedPdu string, setupResponseManager managers.ISetupResponseManager, notificationType string, saveNodebMockReturnValue error) (*setupResponseTestContext, *entities.NodebInfo) {
	var payload []byte
	_, err := fmt.Sscanf(packedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	notificationRequest := models.NotificationRequest{RanName: RanName, Payload: payload}
	testContext := NewSetupResponseTestContext(setupResponseManager)

	handler := NewSetupResponseNotificationHandler(testContext.rnibReaderProvider, testContext.rnibWriterProvider, testContext.setupResponseManager, notificationType)

	var rnibErr error

	nodebInfo := &entities.NodebInfo{
		ConnectionStatus:   entities.ConnectionStatus_CONNECTING,
		ConnectionAttempts: 1,
		RanName:            RanName,
		Ip:                 "10.0.2.2",
		Port:               1231,
	}

	testContext.readerMock.On("GetNodeb", RanName).Return(nodebInfo, rnibErr)
	testContext.writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(saveNodebMockReturnValue)
	handler.Handle(testContext.logger, &notificationRequest, nil)

	return testContext, nodebInfo
}

func TestX2SetupResponse(t *testing.T) {
	var rnibErr error
	testContext, nodebInfo := executeHandleSuccessSetupResponse(t, X2SetupResponsePackedPdu, &managers.X2SetupResponseManager{}, "X2 Setup Response", rnibErr)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, 0, nodebInfo.ConnectionAttempts)
	assert.EqualValues(t, entities.Node_ENB, nodebInfo.NodeType)

	assert.IsType(t, &entities.NodebInfo_Enb{}, nodebInfo.Configuration)
	i, _ := nodebInfo.Configuration.(*entities.NodebInfo_Enb)
	assert.NotNil(t, i.Enb)
}

func TestX2SetupFailureResponse(t *testing.T) {
	var rnibErr error
	testContext, nodebInfo := executeHandleSuccessSetupResponse(t, X2SetupFailureResponsePackedPdu, &managers.X2SetupFailureResponseManager{}, "X2 Setup Failure Response", rnibErr)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, 0, nodebInfo.ConnectionAttempts)
	assert.EqualValues(t, entities.Failure_X2_SETUP_FAILURE, nodebInfo.FailureType)
	assert.NotNil(t, nodebInfo.SetupFailure)
}

func TestEndcSetupResponse(t *testing.T) {
	var rnibErr error
	testContext, nodebInfo := executeHandleSuccessSetupResponse(t, EndcSetupResponsePackedPdu, &managers.EndcSetupResponseManager{}, "ENDC Setup Response", rnibErr)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, 0, nodebInfo.ConnectionAttempts)
	assert.EqualValues(t, entities.Node_GNB, nodebInfo.NodeType)
	assert.IsType(t, &entities.NodebInfo_Gnb{}, nodebInfo.Configuration)

	i, _ := nodebInfo.Configuration.(*entities.NodebInfo_Gnb)
	assert.NotNil(t, i.Gnb)
}
func TestEndcSetupFailureResponse(t *testing.T) {
	var rnibErr error
	testContext, nodebInfo := executeHandleSuccessSetupResponse(t, EndcSetupFailureResponsePackedPdu, &managers.EndcSetupFailureResponseManager{}, "ENDC Setup Failure Response", rnibErr)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
	assert.EqualValues(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.EqualValues(t, 0, nodebInfo.ConnectionAttempts)
	assert.EqualValues(t, entities.Failure_ENDC_X2_SETUP_FAILURE, nodebInfo.FailureType)
	assert.NotNil(t, nodebInfo.SetupFailure)
}

func TestSetupResponseInvalidPayload(t *testing.T) {
	ranName := "test"
	notificationRequest := models.NotificationRequest{RanName: ranName, Payload: []byte("123")}
	testContext := NewSetupResponseTestContext(nil)
	handler := NewSetupResponseNotificationHandler(testContext.rnibReaderProvider, testContext.rnibWriterProvider, &managers.X2SetupResponseManager{}, "X2 Setup Response")
	var rnibErr error
	testContext.readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, ConnectionAttempts: 1}, rnibErr)
	handler.Handle(testContext.logger, &notificationRequest, nil)
	testContext.readerMock.AssertCalled(t, "GetNodeb", ranName)
	testContext.writerMock.AssertNotCalled(t, "SaveNodeb")
}

func TestSetupResponseSaveNodebFailure(t *testing.T) {
	rnibErr := common.NewInternalError(errors.New("Error"))
	testContext, nodebInfo := executeHandleSuccessSetupResponse(t, X2SetupResponsePackedPdu, &managers.X2SetupResponseManager{}, "X2 Setup Response", rnibErr)
	testContext.readerMock.AssertCalled(t, "GetNodeb", RanName)
	testContext.writerMock.AssertCalled(t, "SaveNodeb", mock.Anything, nodebInfo)
}
