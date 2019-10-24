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
	"e2mgr/configuration"
	"e2mgr/controllers"
	"e2mgr/converters"
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
	"os"
	"strconv"
)

const MAX_RNIB_POOL_INSTANCES = 4

func main() {
	config := configuration.ParseConfiguration()
	logLevel, _ := logger.LogLevelTokenToLevel(config.Logging.LogLevel)
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}
	rNibWriter.Init("e2Manager", MAX_RNIB_POOL_INSTANCES)
	defer rNibWriter.Close()
	reader.Init("e2Manager", MAX_RNIB_POOL_INSTANCES)
	defer reader.Close()
	rnibDataService := services.NewRnibDataService(logger, config, reader.GetRNibReader, rNibWriter.GetRNibWriter)
	var msgImpl *rmrCgo.Context
	rmrMessenger := msgImpl.Init("tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	ranReconnectionManager := managers.NewRanReconnectionManager(logger, config, rnibDataService, ranSetupManager)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(logger, rmrSender)
	x2SetupResponseConverter := converters.NewX2SetupResponseConverter(logger)
	x2SetupResponseManager := managers.NewX2SetupResponseManager(x2SetupResponseConverter)
	x2SetupFailureResponseConverter := converters.NewX2SetupFailureResponseConverter(logger)
	x2SetupFailureResponseManager := managers.NewX2SetupFailureResponseManager(x2SetupFailureResponseConverter)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider(logger, rnibDataService, ranReconnectionManager, ranStatusChangeManager, rmrSender, x2SetupResponseManager, x2SetupFailureResponseManager)

	notificationManager := notificationmanager.NewNotificationManager(logger, rmrNotificationHandlerProvider)
	rmrReceiver := rmrreceiver.NewRmrReceiver(logger, rmrMessenger, notificationManager)

	defer (*rmrMessenger).Close()

	go rmrReceiver.ListenAndHandle()

	httpMsgHandlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(logger, rmrSender, config, rnibDataService, ranSetupManager)
	rootController := controllers.NewRootController(rnibDataService)
	nodebController := controllers.NewNodebController(logger, httpMsgHandlerProvider)
	httpserver.Run(config.Http.Port, rootController, nodebController)
}
