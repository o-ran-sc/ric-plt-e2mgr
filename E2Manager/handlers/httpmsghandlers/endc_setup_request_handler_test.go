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

package httpmsghandlers

import (
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewEndcSetupRequestHandler(t *testing.T) {

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	h := NewEndcSetupRequestHandler(rnibWriterProvider)
	assert.NotNil(t, h)
}

func TestCreateEndcX2SetupMessageSuccess(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#setup_request_handler_test.TestCreateMessageSuccess - failed to initialize logger, error: %s", err)
	}
	messageChannel := make(chan *models.E2RequestMessage)
	assert.NotPanics(t, func() { createEndcX2SetupMsg(log, messageChannel) })
	assert.NotEmpty(t, <-messageChannel)
}

func createEndcX2SetupMsg(log *logger.Logger, messageChannel chan *models.E2RequestMessage) {
	h := EndcSetupRequestHandler{}
	E2Sessions := make(sessions.E2Sessions)
	var wg sync.WaitGroup
	var rd models.RequestDetails
	go h.CreateMessage(log, &rd, messageChannel, E2Sessions, time.Now(), wg)
	wg.Wait()
}
