////
//// Copyright 2019 AT&T Intellectual Property
//// Copyright 2019 Nokia
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////      http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
////
//
package handlers
//
//import (
//	"e2mgr/logger"
//	"e2mgr/mocks"
//	"e2mgr/models"
//	"e2mgr/rNibWriter"
//	"e2mgr/rmrCgo"
//	"e2mgr/sessions"
//	"e2mgr/tests"
//	"fmt"
//	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
//	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
//	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
//	"github.com/stretchr/testify/mock"
//	"testing"
//	"time"
//)
//
///*
// * Test an error response while in an x2 setup request session
// */
//func TestHandleInSession(t *testing.T){
//	log, err := logger.InitLogger(logger.InfoLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleInSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	xaction := []byte(fmt.Sprintf("%32s", "1234"))
//	e2Sessions[string(xaction)] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
//		StartTime: time.Now(), TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//
//	if _, ok := e2Sessions[string(xaction)]; ok {
//		t.Errorf("want: no session entry, got: session entry for: %s", string(xaction) )
//	}
//}
//
///*
// * Test an error response triggered by the E2 Term
// */
//
//func TestHandleNoSession(t *testing.T){
//	log, err := logger.InitLogger(logger.InfoLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//			TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//
//	if _, ok := e2Sessions[transactionId]; !ok {
//		t.Errorf("want: session entry for %s, got: no session entry", transactionId )
//	}
//}
///*
// * Test an error response triggered by the E2 Term
// */
//func TestHandleUnsolicitedDisconnectionConnectedSuccess(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
//
//func TestHandleUnsolicitedDisconnectionNotConnectedSuccess(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
//
//func TestHandleUnsolicitedDisconnectionShuttingDownSuccess(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_SHUT_DOWN,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
//
//func TestHandleUnsolicitedDisconnectionShutDownSuccess(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_SHUT_DOWN,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_SHUT_DOWN,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(rnibErr)
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
//
//func TestHandleUnsolicitedDisconnectionReaderFailure(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	var nb *entities.NodebInfo
//	rnibErr := common.RNibError{}
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
//
//func TestHandleUnsolicitedDisconnectionWriterFailure(t *testing.T){
//	log, err := logger.InitLogger(logger.DebugLevel)
//	if err!=nil{
//		t.Errorf("#sctp_errors_notification_handler_test.TestHandleNoSession - failed to initialize logger, error: %s", err)
//	}
//
//	readerMock :=&mocks.RnibReaderMock{}
//	rnibReaderProvider := func() reader.RNibReader {
//		return readerMock
//	}
//	writerMock := &mocks.RnibWriterMock{}
//	rnibWriterProvider := func() rNibWriter.RNibWriter {
//		return writerMock
//	}
//	h := NewRanLostConnectionHandler(rnibReaderProvider,rnibWriterProvider)
//
//	e2Sessions := make(sessions.E2Sessions)
//	transactionId := "1234"
//	xaction := []byte(fmt.Sprintf("%32s", transactionId+"6"))
//	e2Sessions[transactionId] = sessions.E2SessionDetails{SessionStart: time.Now()}
//	payload := []byte("Error")
//	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
//	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(),
//		TransactionId: string(xaction)}
//	var messageChannel chan<- *models.NotificationResponse
//
//	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED,}
//	var rnibErr common.IRNibError
//	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)
//	updatedNb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
//	writerMock.On("SaveNodeb", mock.Anything, updatedNb).Return(common.RNibError{})
//
//	h.Handle(log,e2Sessions, &notificationRequest, messageChannel)
//}
