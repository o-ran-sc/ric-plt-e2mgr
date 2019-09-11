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
//

package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"testing"
	"time"
)

func initX2ResetResponseHandlerTest(t *testing.T) (*logger.Logger, X2ResetResponseHandler, *mocks.RnibReaderMock) {
	log, err := logger.InitLogger(logger.DebugLevel)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	if err!=nil{
		t.Errorf("#sctp_errors_notification_handler_test.TestHandleInSession - failed to initialize logger, error: %s", err)
	}
	readerMock :=&mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	rnibDataService := services.NewRnibDataService(log, config, rnibReaderProvider, nil)

	h := NewX2ResetResponseHandler(rnibDataService)
	return log, h, readerMock
}

func TestX2ResetResponseSuccess(t *testing.T) {
	log, h, readerMock := initX2ResetResponseHandlerTest(t)

	payload, err := tests.BuildPackedX2ResetResponse()
	if err != nil {
		t.Errorf("#x2_reset_response_handler_test.TestX2resetResponse - failed to build and pack X2ResetResponse. Error %x", err)
	}

	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}
	var messageChannel chan<- *models.NotificationResponse

	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED_SETUP_FAILED,}
	var rnibErr error
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log, &notificationRequest, messageChannel)

	//TODO:Nothing to verify
}

func TestX2ResetResponseReaderFailure(t *testing.T) {
	log, h, readerMock := initX2ResetResponseHandlerTest(t)

	var payload []byte
	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}
	var messageChannel chan<- *models.NotificationResponse

	var nb *entities.NodebInfo
	rnibErr  := common.NewResourceNotFoundError("nodeb not found")
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log, &notificationRequest, messageChannel)

	//TODO:Nothing to verify
}

func TestX2ResetResponseUnpackFailure(t *testing.T) {
	log, h, readerMock := initX2ResetResponseHandlerTest(t)

	payload := []byte("not valid payload")
	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}
	var messageChannel chan<- *models.NotificationResponse

	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED_SETUP_FAILED,}
	var rnibErr error
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log, &notificationRequest, messageChannel)

	//TODO:Nothing to verify
}
