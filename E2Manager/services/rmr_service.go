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
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/sessions"
	"strconv"
	"sync"
)


type RmrConfig struct {
	Port       int
	MaxMsgSize int
	Flags      int
	Logger     *logger.Logger
}

func NewRmrConfig(port int, maxMsgSize int, flags int, logger *logger.Logger) *RmrConfig {
	return &RmrConfig{port, maxMsgSize, flags, logger}
}

// RmrService holds an instance of RMR messenger as well as its configuration
type RmrService struct {
	config     *RmrConfig
	messenger  *rmrCgo.RmrMessenger
	e2sessions sessions.E2Sessions
	nManager   *managers.NotificationManager
	rmrResponse chan *models.NotificationResponse
}

// NewRmrService instantiates a new Rmr service instance
func NewRmrService(rmrConfig *RmrConfig, msrImpl rmrCgo.RmrMessenger, e2sessions sessions.E2Sessions, nManager *managers.NotificationManager,
	rmrResponse chan *models.NotificationResponse) *RmrService {

	return &RmrService{
		config:     rmrConfig,
		messenger:  msrImpl.Init("tcp:"+strconv.Itoa(rmrConfig.Port), rmrConfig.MaxMsgSize, rmrConfig.Flags, rmrConfig.Logger),
		e2sessions: e2sessions,
		nManager:   nManager,
		rmrResponse: rmrResponse,
	}
}

func (r *RmrService) SendMessage(messageType int, messageChannel chan *models.E2RequestMessage, errorChannel chan error,
	wg sync.WaitGroup) {

	wg.Add(1)
	setupRequestMessage := <-messageChannel
	e2Message := setupRequestMessage.GetMessageAsBytes(r.config.Logger)

	transactionId := []byte(setupRequestMessage.TransactionId())

	msg := rmrCgo.NewMBuf(messageType, len(e2Message)/*r.config.MaxMsgSize*/, setupRequestMessage.RanName(), &e2Message, &transactionId)

	r.config.Logger.Debugf("#rmr_service.SendMessage - Going to send the message: %#v\n", msg)
	_, err := (*r.messenger).SendMsg(msg, r.config.MaxMsgSize)

	errorChannel <- err
	wg.Done()
}

func (r *RmrService) SendRmrMessage(response *models.NotificationResponse) {

	msgAsBytes := response.GetMessageAsBytes(r.config.Logger)
	transactionIdByteArr := []byte(response.RanName)

	msg := rmrCgo.NewMBuf(response.MgsType, len(msgAsBytes), response.RanName, &msgAsBytes, &transactionIdByteArr)

	r.config.Logger.Debugf("#rmr_service.SendRmrMessage - Going to send the message: %#v\n", msg)

	_, err := (*r.messenger).SendMsg(msg, r.config.MaxMsgSize)

	if err != nil {
		r.config.Logger.Errorf("#rmr_service.SendRmrMessage - error: %#v\n", err)
		return
	}
}

// ListenAndHandle waits for messages coming from rmr_rcv_msg and sends it to a designated message handler
func (r *RmrService) ListenAndHandle() {

	for {
		mbuf, err := (*r.messenger).RecvMsg()
		r.config.Logger.Debugf("#rmr_service.ListenAndHandle - Going to handle received message: %#v\n", mbuf)

		// TODO: one mbuf received immediately execute goroutine
		if err != nil {
			continue	//TODO log error
		}

		r.nManager.HandleMessage(r.config.Logger, r.e2sessions, mbuf, r.rmrResponse)
	}
}

func (r *RmrService) SendResponse(){
	for{

		response, ok := <-r.rmrResponse
		if !ok {

			r.config.Logger.Errorf("#rmr_service.SendResponse - channel closed")
			break
		}

		r.config.Logger.Debugf("#rmr_service.SendResponse - Going to send message: %#v\n", response)
		r.SendRmrMessage(response)
	}
}

func (r *RmrService) CloseContext() {
	if r.config.Logger.DebugEnabled(){
		r.config.Logger.Debugf("#rmr_service.CloseContext - RMR is ready: %v", (*r.messenger).IsReady())
		(*r.messenger).Close()
		r.config.Logger.Debugf("#rmr_service.CloseContext - RMR is ready: %v", (*r.messenger).IsReady())
	}
}

