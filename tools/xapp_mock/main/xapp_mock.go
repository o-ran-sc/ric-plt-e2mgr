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
package main

import (
	"../frontend"
	"../rmr"
	"../sender"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
)

const (
	ENV_RMR_PORT = "RMR_PORT"
	RMR_PORT_DEFAULT = 5001
)

var rmrService *rmr.Service

func main() {
	var rmrContext *rmr.Context

	var rmrConfig rmr.Config = rmr.Config{Port: RMR_PORT_DEFAULT, MaxMsgSize: rmr.RMR_MAX_MSG_SIZE, MaxRetries: 3, Flags: 0}
	if port, err := strconv.ParseUint(os.Getenv(ENV_RMR_PORT), 10, 16); err == nil {
		rmrConfig.Port = int(port)
	} else {
		log.Printf("%s: %s, using default (%d).", ENV_RMR_PORT, err,RMR_PORT_DEFAULT)
	}
	rmrService = rmr.NewService(rmrConfig, rmrContext)

	/* Load configuration file*/
	err := frontend.ProcessConfigurationFile("resources","conf",  ".json",
		func(data []byte) error {
			return  frontend.JsonCommandsDecoder(data,jsonCommandsDecoderCB)
		})
	if err != nil {
		log.Fatalf("processing Error: %s", err)
	}

	log.Print("xapp_mock is up and running.")

	cmd:= flag.Arg(0) /*first remaining argument after flags have been processed*/
	if err :=  frontend.JsonCommandDecoder([]byte(cmd),jsonCommandDecoderCB); err != nil {
		log.Printf("command processing Error: %s", err)
	}

	rmrService.CloseContext()

	log.Print("xapp_mock is down.")
}


// TODO: move callbacks to Dispatcher.
func jsonCommandsDecoderCB(command *frontend.JsonCommand) error {
	if len(command.Id) == 0{
		return errors.New(fmt.Sprintf("invalid command, no Id"))
	}
	frontend.Configuration[command.Id] = command
	if rmrMsgId, err := rmr.MessageIdToUint(command.WaitForRmrMessageType); err != nil {
		return errors.New(fmt.Sprintf("invalid rmr message id: %s",command.WaitForRmrMessageType))
	} else {
		frontend.WaitedForRmrMessageType[int(rmrMsgId)] = command
	}
	return nil
}

// TODO: merge command with configuration
func jsonCommandDecoderCB(command *frontend.JsonCommand) error {
	if len(command.Id) == 0{
		return errors.New(fmt.Sprintf("invalid command, no Id"))
	}
	switch command.Action {
	case frontend.SendRmrMessage:
		if  err := sender.SendJsonRmrMessage(*command, nil, rmrService); err != nil {
			return err
		}
		if len(command.WaitForRmrMessageType) > 0 {
			rmrService.ListenAndHandle() //TODO: handle error
		}
	case frontend.ReceiveRmrMessage:
		if rmrMsgId, err := rmr.MessageIdToUint(command.RmrMessageType); err != nil {
			return errors.New(fmt.Sprintf("invalid rmr message id: %s",command.WaitForRmrMessageType))
		} else {
			frontend.WaitedForRmrMessageType[int(rmrMsgId)] = command
		}
		rmrService.ListenAndHandle() //TODO: handle error
	default:
		return errors.New(fmt.Sprintf("invalid command action %s", command.Action))
	}
	return nil
}