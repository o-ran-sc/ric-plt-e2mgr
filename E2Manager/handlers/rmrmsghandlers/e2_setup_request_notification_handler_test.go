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

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	prefix = "10.0.2.15:9999|"
	logFilePath = "./loggerTest.txt"
	e2tInstanceFullAddress = "10.0.2.15:9999"
	nodebRanName = "gnb:310-410-b5c67788"
)

func TestParseSetupRequest_Success(t *testing.T){
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Equal(t, request.GetPlmnId(), "131014")
	assert.Equal(t, request.GetNbId(), "10110101110001100111011110001000")
}

func TestParseSetupRequest_PipFailure(t *testing.T){
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte("10.0.2.15:9999")
	request, _, err := handler.parseSetupRequest(append(prefBytes, xmlGnb...))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error parsing E2 Setup Request failed extract Payload: no | separator found")
}

func TestParseSetupRequest_UnmarshalFailure(t *testing.T){
	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte(prefix)
	request, _, err := handler.parseSetupRequest(append(prefBytes, 1,2,3))
	assert.Nil(t, request)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "#E2SetupRequestNotificationHandler.parseSetupRequest - Error unmarshalling E2 Setup Request payload: 31302e302e322e31353a393939397c010203")
}

func TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertSuccessFlowNewNodebLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewEnGnbSuccess(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_en-gNB.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlEnGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewEnGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlEnGnb...)}
	handler.Handle(notificationRequest)
	assertSuccessFlowNewNodebLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleNewNgEnbSuccess(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_ng-eNB.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlEnGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewNgEnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	handler := stubMockSuccessFlowNewNodeb(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlEnGnb...)}
	handler.Handle(notificationRequest)
	assertSuccessFlowNewNodebLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleExistingGnbSuccess(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	handler := stubMockSuccessFlowExistingNodeb(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertSuccessFlowExistingNodebLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleParseError(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	_, handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte("invalid_prefix")
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertParseErrorFlowLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleUnmarshalError(t *testing.T) {
	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	_, handler, _, _, _, _, _ := initMocks(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, "xmlGnb"...)}
	handler.Handle(notificationRequest)
	assertUnmarshalErrorFlowLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleGetE2TInstanceError(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	_, handler, _, _, _, e2tInstancesManagerMock, _ := initMocks(t)
	var e2tInstance * entities.E2TInstance
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, common.NewResourceNotFoundError("Not found"))
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertGetE2TInstanceErrorLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleGetNodebError(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile
	_, handler, readerMock, _, _, e2tInstancesManagerMock, _ := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewInternalError(errors.New("Some error")))
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertGetNodebErrorLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleAssociationError(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	_, handler, readerMock, writerMock, _, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(errors.New("association error"))

	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertAssociationErrorLogs(t)
}

func TestE2SetupRequestNotificationHandler_HandleExistingGnbInvalidStatusError(t *testing.T) {
	path, err :=filepath.Abs("../../tests/resources/setupRequest_gnb.xml")
	if err != nil {
		t.Fatal(err)
	}
	xmlGnb, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil{
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleNewGnbSuccess - failed to create file, error: %s", err)
	}
	oldStdout := os.Stdout
	defer changeStdout(oldStdout)
	defer removeLogFile(t)
	os.Stdout = logFile

	handler := stubMockInvalidStatusFlowExistingNodeb(t)
	prefBytes := []byte(prefix)
	notificationRequest := &models.NotificationRequest{RanName: nodebRanName, Payload: append(prefBytes, xmlGnb...)}
	handler.Handle(notificationRequest)
	assertInvalidNodebStatusLogs(t)
}

func assertInvalidNodebStatusLogs(t *testing.T){
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertInvalidNodebStatusLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertInvalidNodebStatusLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "#RnibDataService.GetNodeb")
	assert.Contains(t, record, "connection status: SHUTTING_DOWN")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#E2SetupRequestNotificationHandler.Handle")
	assert.Contains(t, record, "connection status: SHUTTING_DOWN - nodeB entity in incorrect state")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#E2SetupRequestNotificationHandler.Handle")
	assert.Contains(t, record, "Summary: elapsed time for receiving and handling setup request message from E2 terminator")
}

func assertAssociationErrorLogs(t *testing.T){
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertNewNodebSavedLog(buf, t)
	assertAssociationErrorLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertAssociationErrorLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "#E2TAssociationManager.AssociateRan - Associating RAN")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#E2TAssociationManager.AssociateRan - RoutingManager failure: Failed to associate RAN")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#E2SetupRequestNotificationHandler.Handle - RAN name:")
	assert.Contains(t, record, "failed to associate E2T to nodeB entity")
}

func assertGetNodebErrorLogs(t *testing.T) {
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertGetNodebErrorLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertGetNodebErrorLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "failed to retrieve nodebInfo entity")
}

func assertGetE2TInstanceErrorLogs(t *testing.T) {
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertGetE2TInstanceErrorLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertGetE2TInstanceErrorLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "Failed retrieving E2TInstance")
}

func removeLogFile(t *testing.T) {
	err := os.Remove(logFilePath)
	if err != nil {
		t.Errorf("e2_setup_request_notification_handler_test.TestE2SetupRequestNotificationHandler_HandleGnbSuccess - failed to remove file, error: %s", err)
	}
}

func assertParseErrorFlowLogs(t *testing.T) {
	buf := getLogFileBuffer(t)
	assertReceivedAndFailedParseLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertUnmarshalErrorFlowLogs(t *testing.T) {
	buf := getLogFileBuffer(t)
	assertReceivedAndFailedUnmarshalLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertSuccessFlowNewNodebLogs(t *testing.T){
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertNewNodebSavedLog(buf, t)
	assertAssociatedLog(buf, t)
	assertRequestBuiltLog(buf, t)
	assertRequestSentLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertSuccessFlowExistingNodebLogs(t *testing.T){
	buf := getLogFileBuffer(t)
	assertReceivedAndParsedLog(buf, t)
	assertExistingNodebRetrievedLog(buf, t)
	assertAssociatedLog(buf, t)
	assertRequestBuiltLog(buf, t)
	assertRequestSentLog(buf, t)
	assertNoMoreRecordsLog(buf, t)
}

func assertReceivedAndParsedLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "received E2_SETUP_REQUEST")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "handling E2_SETUP_REQUEST")
}

func assertReceivedAndFailedParseLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "received E2_SETUP_REQUEST")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "Error parsing E2 Setup Request")
}

func assertReceivedAndFailedUnmarshalLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "received E2_SETUP_REQUEST")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "Error unmarshalling E2 Setup Request")
}

func assertNewNodebSavedLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "#RnibDataService.SaveNodeb - nbIdentity:")
}

func assertExistingNodebRetrievedLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "#RnibDataService.GetNodeb - RAN name:")
}

func assertAssociatedLog(buf *bytes.Buffer, t *testing.T){
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "#E2TAssociationManager.AssociateRan - Associating RAN")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#RnibDataService.UpdateNodebInfo")
	record, _ = buf.ReadString('\n')
	assert.Contains(t, record, "#E2TAssociationManager.AssociateRan - successfully associated RAN")
}

func assertRequestSentLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "uccessfully sent RMR message")
}
func assertRequestBuiltLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Contains(t, record, "RIC_E2_SETUP_RESP message has been built successfully")
}

func assertNoMoreRecordsLog(buf *bytes.Buffer, t *testing.T) {
	record, _ := buf.ReadString('\n')
	assert.Empty(t, record)
}

func stubMockSuccessFlowNewNodeb(t *testing.T) E2SetupRequestNotificationHandler{
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb *entities.NodebInfo
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, common.NewResourceNotFoundError("Not found"))
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var err error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, err)
	return handler
}

func stubMockSuccessFlowExistingNodeb(t *testing.T) E2SetupRequestNotificationHandler{
	_, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb = &entities.NodebInfo{
		RanName: nodebRanName,
		AssociatedE2TInstanceAddress: e2tAddress,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		NodeType: entities.Node_GNB,
		Configuration: &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
	}
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, nil)
	routingManagerClientMock.On("AssociateRanToE2TInstance", e2tInstanceFullAddress, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	var err error
	rmrMessage := &rmrCgo.MBuf{}
	rmrMessengerMock.On("SendMsg", mock.Anything, mock.Anything).Return(rmrMessage, err)
	return handler
}

func stubMockInvalidStatusFlowExistingNodeb(t *testing.T) E2SetupRequestNotificationHandler{
	_, handler, readerMock, _, _, e2tInstancesManagerMock, _ := initMocks(t)
	var e2tInstance = &entities.E2TInstance{}
	e2tInstancesManagerMock.On("GetE2TInstance", e2tInstanceFullAddress).Return(e2tInstance, nil)
	var gnb = &entities.NodebInfo{RanName: nodebRanName, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN}
	readerMock.On("GetNodeb", mock.Anything).Return(gnb, nil)
	return handler
}

func initMocks(t *testing.T) (*logger.Logger, E2SetupRequestNotificationHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock, *mocks.RoutingManagerClientMock) {
	logger := tests.InitLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := tests.InitRmrSender(rmrMessengerMock, logger)
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	routingManagerClientMock := &mocks.RoutingManagerClientMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManagerMock, routingManagerClientMock)
	handler := NewE2SetupRequestNotificationHandler(logger, config, e2tInstancesManagerMock, rmrSender, rnibDataService, e2tAssociationManager)
	return logger, handler, readerMock, writerMock, rmrMessengerMock, e2tInstancesManagerMock, routingManagerClientMock
}

func changeStdout(old *os.File) {
	os.Stdout = old
}

func getLogFileBuffer(t *testing.T) *bytes.Buffer {
	logFile, err := os.Open(logFilePath)
	if err != nil {
		t.Errorf("e2_setup_request_notification_handler_test.assertSuccessFlowNewNodebLogRecords - failed to open file, error: %s", err)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, logFile)
	if err != nil {
		t.Errorf("e2_setup_request_notification_handler_test.assertSuccessFlowNewNodebLogRecords - failed to copy bytes, error: %s", err)
	}
	return &buf
}

