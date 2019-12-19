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
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/controllers"
	"e2mgr/httpserver"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/managers/notificationmanager"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrreceiver"
	"e2mgr/services/rmrsender"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"net/http"
	"os"
	"strconv"
)

func main() {
	config := configuration.ParseConfiguration()
	logLevel, _ := logger.LogLevelTokenToLevel(config.Logging.LogLevel)
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}
	db := sdlgo.NewDatabase()
	sdl := sdlgo.NewSdlInstance("e2Manager", db)
	defer sdl.Close()
	rnibDataService := services.NewRnibDataService(logger, config, reader.GetRNibReader(sdl), rNibWriter.GetRNibWriter( sdl))
	var msgImpl *rmrCgo.Context
	rmrMessenger := msgImpl.Init("tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	e2tShutdownManager := managers.NewE2TShutdownManager(logger, rnibDataService, e2tInstancesManager)
	e2tKeepAliveWorker := managers.NewE2TKeepAliveWorker(logger, rmrSender, e2tInstancesManager, e2tShutdownManager, config)
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, &http.Client{})
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider()
	rmrNotificationHandlerProvider.Init(logger, config, rnibDataService, rmrSender, ranSetupManager, e2tInstancesManager, routingManagerClient)

	notificationManager := notificationmanager.NewNotificationManager(logger, rmrNotificationHandlerProvider)
	rmrReceiver := rmrreceiver.NewRmrReceiver(logger, rmrMessenger, notificationManager)

	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()

	defer rmrMessenger.Close()

	go rmrReceiver.ListenAndHandle()
	go e2tKeepAliveWorker.Execute()

	httpMsgHandlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(logger, rmrSender, config, rnibDataService, ranSetupManager, e2tInstancesManager)
	rootController := controllers.NewRootController(rnibDataService)
	nodebController := controllers.NewNodebController(logger, httpMsgHandlerProvider)
	_ = httpserver.Run(logger, config.Http.Port, rootController, nodebController)
}