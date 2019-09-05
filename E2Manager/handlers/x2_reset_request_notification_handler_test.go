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

package handlers

import (
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/tests"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestX2ResetRequestNotifSuccess(t *testing.T) {
	log := initLog(t)
	payload := []byte("payload")
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}

	h := NewX2ResetRequestNotificationHandler(readerProvider)

	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}

	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_CONNECTED,}
	var rnibErr error
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	messageChannel := make(chan *models.NotificationResponse)

	go h.Handle(log,nil, &notificationRequest, messageChannel)

	result := <-messageChannel
	assert.Equal(t, result.RanName, mBuf.Meid)
	assert.Equal(t, result.MgsType, rmrCgo.RIC_X2_RESET_RESP)
}

func TestHandleX2ResetRequestNotifShuttingDownStatus(t *testing.T) {
	log := initLog(t)
	var payload []byte
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}

	h := NewX2ResetRequestNotificationHandler(readerProvider)

	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}

	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN,}
	var rnibErr error

	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log,nil, &notificationRequest, nil)
}

func TestHandleX2ResetRequestNotifDisconnectStatus(t *testing.T) {
	log := initLog(t)
	var payload []byte
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}

	h := NewX2ResetRequestNotificationHandler(readerProvider)

	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}

	nb := &entities.NodebInfo{RanName:mBuf.Meid, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED,}
	var rnibErr error

	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log,nil, &notificationRequest, nil)
}

func TestHandleX2ResetRequestNotifGetNodebFailed(t *testing.T){

	log := initLog(t)
	var payload []byte
	readerMock := &mocks.RnibReaderMock{}
	readerProvider := func() reader.RNibReader {
		return readerMock
	}

	h := NewX2ResetRequestNotificationHandler(readerProvider)
	xaction := []byte("RanName")
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload),"RanName", &payload, &xaction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: string(xaction)}

	var nb *entities.NodebInfo
	rnibErr  := &common.ResourceNotFoundError{}
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(log,nil, &notificationRequest, nil)
}
