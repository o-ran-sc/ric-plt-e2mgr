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
	Config      *RmrConfig
	Messenger   *rmrCgo.RmrMessenger
	E2sessions  sessions.E2Sessions
	RmrResponse chan *models.NotificationResponse
}

// NewRmrService instantiates a new Rmr service instance
func NewRmrService(rmrConfig *RmrConfig, msrImpl rmrCgo.RmrMessenger, e2sessions sessions.E2Sessions, rmrResponse chan *models.NotificationResponse) *RmrService {

	return &RmrService{
		Config:      rmrConfig,
		Messenger:   msrImpl.Init("tcp:"+strconv.Itoa(rmrConfig.Port), rmrConfig.MaxMsgSize, rmrConfig.Flags, rmrConfig.Logger),
		E2sessions:  e2sessions,
		RmrResponse: rmrResponse,
	}
}

func (r *RmrService) SendMessage(messageType int, messageChannel chan *models.E2RequestMessage, errorChannel chan error,
	wg sync.WaitGroup) {

	wg.Add(1)
	setupRequestMessage := <-messageChannel
	e2Message := setupRequestMessage.GetMessageAsBytes(r.Config.Logger)

	transactionId := []byte(setupRequestMessage.TransactionId())

	msg := rmrCgo.NewMBuf(messageType, len(e2Message) /*r.config.MaxMsgSize*/, setupRequestMessage.RanName(), &e2Message, &transactionId)

	r.Config.Logger.Debugf("#rmr_service.SendMessage - Going to send the message: %#v\n", msg)
	_, err := (*r.Messenger).SendMsg(msg, r.Config.MaxMsgSize)

	errorChannel <- err
	wg.Done()
}

func (r *RmrService) SendRmrMessage(response *models.NotificationResponse) error {

	msgAsBytes := response.GetMessageAsBytes(r.Config.Logger)
	transactionIdByteArr := []byte(response.RanName)

	msg := rmrCgo.NewMBuf(response.MgsType, len(msgAsBytes), response.RanName, &msgAsBytes, &transactionIdByteArr)

	_, err := (*r.Messenger).SendMsg(msg, r.Config.MaxMsgSize)

	if err != nil {
		return err
	}
	return nil
}

func (r *RmrService) SendResponse() {
	for {

		response, ok := <-r.RmrResponse
		if !ok {

			r.Config.Logger.Errorf("#rmr_service.SendResponse - channel closed")
			break
		}

		r.Config.Logger.Debugf("#rmr_service.SendResponse - Going to send message: %#v\n", response)
		if err := r.SendRmrMessage(response); err != nil {
			r.Config.Logger.Errorf("#rmr_service.SendResponse - error: %#v\n", err)
		}
	}
}

func (r *RmrService) CloseContext() {
	if r.Config.Logger.DebugEnabled() {
		r.Config.Logger.Debugf("#rmr_service.CloseContext - RMR is ready: %v", (*r.Messenger).IsReady())
		(*r.Messenger).Close()
		r.Config.Logger.Debugf("#rmr_service.CloseContext - RMR is ready: %v", (*r.Messenger).IsReady())
	}
}
