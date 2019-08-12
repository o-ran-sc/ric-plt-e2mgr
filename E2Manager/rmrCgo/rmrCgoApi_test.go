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

package rmrCgo

import (
	"e2mgr/logger"
	"e2mgr/tests"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

var (
	log     *logger.Logger
	msgr *RmrMessenger
)

func TestLogger(t *testing.T){
	var err error
	log, err = logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#rmrCgoApi_test.TestLogger - failed to initialize logger, error: %s", err)
	}
	data :=  map[string]interface{}{"messageType": 1001, "ranIp":"10.0.0.3", "ranPort": 879, "ranName":"test1"}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data)
	req := tests.GetHttpRequest()
	boo, _ := ioutil.ReadAll(req.Body)
	log.Debugf("#rmrCgoApi_test.TestLogger - request header: %v\n; request body: %s\n", req.Header, string(boo))
}


func TestNewMBufSuccess(t *testing.T) {
	var err error
	log, err = logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#rmrCgoApi_test.TestNewMBufSuccess - failed to initialize logger, error: %s", err)
	}
	msg := NewMBuf(tests.MessageType, len(tests.DummyPayload),"RanName", &tests.DummyPayload, &tests.DummyXAction)
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg.Payload)
	assert.NotEmpty(t, msg.XAction)
	assert.Equal(t, msg.MType, tests.MessageType)
	assert.Equal(t, msg.Meid, "RanName")
	assert.Equal(t, msg.Len, len(tests.DummyPayload))
}

func TestInitFailure(t *testing.T) {
	var err error
	log, err = logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#rmrCgoApi_test.TestInitFailure - failed to initialize logger, error: %s", err)
	}
	go initRmr(tests.GetPort(), tests.MaxMsgSize, tests.Flags, log)
	time.Sleep(time.Second)
	if msgr != nil {
		t.Errorf("The rmr router is ready, should be not ready")
	}
}

//func TestInitSuccess(t *testing.T) {
//	var err error
//	log, err = logger.InitLogger(true)
//	if err != nil {
//		t.Errorf("#rmrCgoApi_test.TestInitSuccess - failed to initialize logger, error: %s", err)
//	}
//	go initRmr(tests.GetPort(), tests.MaxMsgSize, tests.Flags, log)
//	time.Sleep(time.Second)
//	if msgr == nil {
//		t.Errorf("The rmr router is not ready, should be ready")
//	}
//}

func TestIsReadyFailure(t *testing.T) {
	var err error
	log, err = logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#rmrCgoApi_test.TestIsReadyFailure - failed to initialize logger, error: %s", err)
	}

	go initRmr(tests.GetPort(), tests.MaxMsgSize, tests.Flags, log)
	time.Sleep(time.Second)
	assert.True(t, msgr == nil || !(*msgr).IsReady())
}

//func TestSendRecvMsgSuccess(t *testing.T) {
//	var err error
//	log, err = logger.InitLogger(true)
//	if err != nil {
//		t.Errorf("#rmrCgoApi_test.TestSendRecvMsgSuccess - failed to initialize logger, error: %s", err)
//	}
//	go initRmr(tests.GetPort(), tests.MaxMsgSize, tests.Flags, log)
//	time.Sleep(time.Second)
//	if msgr == nil || !(*msgr).IsReady()  {
//		t.Errorf("#rmrCgoApi_test.TestSendRecvMsgSuccess - The rmr router is not ready")
//	}
//	msg := NewMBuf(1, tests.MaxMsgSize, &tests.DummyPayload, &tests.DummyXAction)
//	log.Debugf("#rmrCgoApi_test.TestSendRecvMsgSuccess - Going to send the message: %#v\n", msg)
//	msgR, _ := (*msgr).SendMsg(msg, tests.MaxMsgSize)
//	log.Debugf("#rmrCgoApi_test.TestSendRecvMsgSuccess - The message has been sent %#v\n", msgR)
//	log.Debugf("#rmrCgoApi_test.TestSendRecvMsgSuccess - The payload: %#v\n", msgR.Payload)
//	msgR = (*msgr).RecvMsg()
//	log.Debugf("#rmrCgoApi_test.TestSendRecvMsgSuccess - The message has been received: %#v\n", msgR)
//	log.Debugf("#rmrCgoApi_test.TestSendRecvMsgSuccess - The payload: %#v\n", msgR.Payload)
//	(*msgr).Close()
//}

//func TestIsReadySuccess(t *testing.T) {
//	var err error
//	log, err = logger.InitLogger(true)
//	if err != nil {
//		t.Errorf("#rmrCgoApi_test.TestIsReadySuccess - The rmr router is not ready")
//	}
//	go initRmr(tests.GetPort(), tests.MaxMsgSize, tests.Flags, log)
//	time.Sleep(time.Second)
//	if msgr == nil || !(*msgr).IsReady()  {
//		t.Errorf("#rmrCgoApi_test.TestIsReadySuccess - The rmr router is not ready")
//	}
//}

func initRmr(port string, maxMsgSize int, flags int, log *logger.Logger){
	var ctx *Context
	msgr = ctx.Init(port, maxMsgSize, flags, log)
}
