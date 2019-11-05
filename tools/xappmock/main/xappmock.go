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
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
	"xappmock/dispatcher"
	"xappmock/frontend"
	"xappmock/rmr"
)

const (
	ENV_RMR_PORT     = "RMR_PORT"
	RMR_PORT_DEFAULT = 5001
)

var rmrService *rmr.Service

func main() {
	var rmrContext *rmr.Context
	var rmrConfig = rmr.Config{Port: RMR_PORT_DEFAULT, MaxMsgSize: rmr.RMR_MAX_MSG_SIZE, MaxRetries: 10, Flags: 0}

	if port, err := strconv.ParseUint(os.Getenv(ENV_RMR_PORT), 10, 16); err == nil {
		rmrConfig.Port = int(port)
	} else {
		log.Printf("#main - %s: %s, using default (%d).", ENV_RMR_PORT, err, RMR_PORT_DEFAULT)
	}

	rmrService = rmr.NewService(rmrConfig, rmrContext)
	dispatcherDesc := dispatcher.New(rmrService)

	/* Load configuration file*/
	err := frontend.ProcessConfigurationFile("resources", "conf", ".json",
		func(data []byte) error {
			return frontend.JsonCommandsDecoder(data, dispatcherDesc.JsonCommandsDecoderCB)
		})

	if err != nil {
		log.Fatalf("#main - processing error: %s", err)
	}

	log.Print("#main - xApp Mock is up and running...")

	flag.Parse()
	cmd := flag.Arg(0) /*first remaining argument after flags have been processed*/

	command, err := frontend.DecodeJsonCommand([]byte(cmd))

	if err != nil {
		log.Printf("#main - command decoding error: %s", err)
		rmrService.CloseContext()
		log.Print("#main - xApp Mock is down")
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
		rmrService.CloseContext()
	}()

	processStartTime := time.Now()
	dispatcherDesc.ProcessJsonCommand(ctx, command)
	pr := dispatcherDesc.GetProcessResult()

	if pr.Err != nil {
		log.Printf("#main - command processing Error: %s", err)
	}

	processElapsedTimeInMs := float64(time.Since(processStartTime)) / float64(time.Millisecond)

	log.Printf("#main - processing (sending/receiving) messages took %.2f ms", processElapsedTimeInMs)
	log.Printf("#main - process result: %s", pr)

	rmrService.CloseContext() // TODO: called twice
	log.Print("#main - xApp Mock is down")
}
