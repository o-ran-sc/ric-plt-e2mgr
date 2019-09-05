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
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"e2mgr/tests"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandleSuccessEnbConfigUpdate(t *testing.T){
		log, err := logger.InitLogger(logger.InfoLevel)
		if err!=nil{
			t.Errorf("#endc_configuration_update_handler_test.TestHandleSuccessEndcConfigUpdate - failed to initialize logger, error: %s", err)
		}
		h := X2EnbConfigurationUpdateHandler{}
		E2Sessions := make(sessions.E2Sessions)

		payload := tests.GetPackedPayload(t)
		mBuf := rmrCgo.NewMBuf(10370, len(payload),"RanName", &payload, &tests.DummyXAction)
		notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now()}
		messageChannel := make(chan *models.NotificationResponse)

		go h.Handle(log, E2Sessions, &notificationRequest, messageChannel)

		response := <-messageChannel

		assert.NotEmpty(t, response)
		assert.EqualValues(t, 10081, response.MgsType)
		assert.True(t, len(response.Payload) > 0)
}

func TestHandleFailureEnbConfigUpdate(t *testing.T){
	log, err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("#endc_configuration_update_handler_test.TestHandleFailureEndcConfigUpdate - failed to initialize logger, error: %s", err)
	}
	h := X2EnbConfigurationUpdateHandler{}
	E2Sessions := make(sessions.E2Sessions)

	mBuf := rmrCgo.NewMBuf(tests.MessageType, 4,"RanName", &tests.DummyPayload, &tests.DummyXAction)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now()}
	messageChannel := make(chan *models.NotificationResponse)

	go h.Handle(log, E2Sessions, &notificationRequest, messageChannel)

	response := <-messageChannel

	assert.NotEmpty(t, response)
	assert.EqualValues(t, 10082, response.MgsType)
	assert.True(t, len(response.Payload) > 0)
}
