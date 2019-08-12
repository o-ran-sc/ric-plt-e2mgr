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

package services

import (
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	"time"
)

func TestNewRmrConfig(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#rmr_service_test.TestNewRmrConfig - failed to initialize logger, error: %s", err)
	}
	assert.NotNil(t, NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log))
}

func TestSendMessage(t *testing.T){
	log, err := logger.InitLogger(logger.DebugLevel)
	if err!=nil{
		t.Errorf("#rmr_service_test.TestSendMessage - failed to initialize logger, error: %s", err)
	}
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	mbuf := rmrCgo.NewMBuf(tests.MessageType, tests.MaxMsgSize,"RanName" , &tests.DummyPayload, &tests.DummyXAction)
	rmrMessengerMock.On("SendMsg", mock.AnythingOfType(fmt.Sprintf("%T", mbuf)), tests.MaxMsgSize).Return(mbuf, nil)

	errorChannel := make(chan error)
	var wg sync.WaitGroup
	go getRmrService( rmrMessengerMock, log).SendMessage(tests.MessageType, make(chan *models.E2RequestMessage), errorChannel, wg)
	wg.Wait()

	assert.Empty (t, errorChannel)
}

func TestListenAndHandle(t *testing.T){
	log, err := logger.InitLogger(logger.DebugLevel)
	if err!=nil{
		t.Errorf("#rmr_service_test.TestListenAndHandle - failed to initialize logger, error: %s", err)
	}
	rmrMessengerMock := &mocks.RmrMessengerMock{}

	var buf *rmrCgo.MBuf
	e := fmt.Errorf("test error")
	rmrMessengerMock.On("RecvMsg").Return(buf, e)

	go  getRmrService(rmrMessengerMock,log).ListenAndHandle()

	time.Sleep(time.Microsecond*10)
}

func TestCloseContext(t *testing.T){
	log, err := logger.InitLogger(logger.DebugLevel)
	if err!=nil{
		t.Errorf("#rmr_service_test.TestCloseContext - failed to initialize logger, error: %s", err)
	}
	rmrMessengerMock := &mocks.RmrMessengerMock{}

	rmrMessengerMock.On("IsReady").Return(true)
	rmrMessengerMock.On("Close").Return()

	getRmrService(rmrMessengerMock, log).CloseContext()
	time.Sleep(time.Microsecond*10)
}

func getRmrService(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *RmrService {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	readerMock :=&mocks.RnibReaderMock{}
	rnibReaderProvider := func() reader.RNibReader {
		return readerMock
	}
	writerMock := &mocks.RnibWriterMock{}
	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return writerMock
	}
	nManager := managers.NewNotificationManager(rnibReaderProvider, rnibWriterProvider)
	messageChannel := make(chan *models.NotificationResponse)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return NewRmrService(NewRmrConfig(tests.Port, tests.MaxMsgSize, tests.Flags, log), rmrMessenger,  make(sessions.E2Sessions), nManager, messageChannel)
}