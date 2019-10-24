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
package rmr

import (
	"../frontend"
	"log"
	"strconv"
)
// RmrService holds an instance of RMR messenger as well as its configuration
type Service struct {
	messenger  *Messenger
}

// NewRmrService instantiates a new Rmr service instance
func NewService(rmrConfig Config, messenger Messenger) *Service {
	return &Service{
		messenger: messenger.Init("tcp:"+strconv.Itoa(rmrConfig.Port), rmrConfig.MaxMsgSize, rmrConfig.MaxRetries, rmrConfig.Flags),
	}
}

func (r *Service) SendMessage(messageType int, msg []byte, transactionId []byte) (*MBuf, error){
	log.Printf( "SendMessage (type: %d, tid: %s, msg: %v", messageType, transactionId, msg)
	mbuf := NewMBuf(messageType, len(msg), msg, transactionId)
	return (*r.messenger).SendMsg(mbuf)
}

// ListenAndHandle waits for messages coming from rmr_rcv_msg and sends it to a designated message handler
func (r *Service) ListenAndHandle() error {
	for {
		mbuf, err := (*r.messenger).RecvMsg()

		if err != nil {
			return err
		}

		if _, ok := frontend.WaitedForRmrMessageType[mbuf.MType]; ok {
			log.Printf( "ListenAndHandle Expected msg: %s", mbuf)
			break
		} else {
			log.Printf( "ListenAndHandle Unexpected msg: %s", mbuf)
		}
	}
	return nil
}


func (r *Service) CloseContext() {
	(*r.messenger).Close()

}



