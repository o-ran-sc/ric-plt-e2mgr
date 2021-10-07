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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"os"
	"strconv"
)

const GeneralKeyDefaultValue = "{\"enableRic\":true}"

func initKeys(logger *logger.Logger, sdl *sdlgo.SyncStorage) error {
	ok, err := sdl.SetIfNotExists(common.GetRNibNamespace(), common.BuildGeneralConfigurationKey(), GeneralKeyDefaultValue)

	if err != nil {
		logger.Errorf("#app.main - Failed setting GENERAL key")
		return err
	}

	if ok {
		logger.Infof("#app.main - Successfully set GENERAL key")
	} else {
		logger.Infof("#app.main - GENERAL key exists, no need to set")
	}

	return nil

}

func main() {
	config := configuration.ParseConfiguration()
	logLevel, _ := logger.LogLevelTokenToLevel(config.Logging.LogLevel)
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}
	logger.Infof("#app.main - Configuration %s", config)
	sdl := sdlgo.NewSyncStorage()
	err = initKeys(logger, sdl)

	if err != nil {
		os.Exit(1)
	}

	defer sdl.Close()
	rnibDataService := services.NewRnibDataService(logger, config, reader.GetNewRNibReader(sdl), rNibWriter.GetRNibWriter(sdl, config.RnibWriter))

	ranListManager := managers.NewRanListManager(logger, rnibDataService)

	err = ranListManager.InitNbIdentityMap()

	if err != nil {
		logger.Errorf("#app.main - quit")
		os.Exit(1)
	}

	var msgImpl *rmrCgo.Context
	rmrMessenger := msgImpl.Init("tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, clients.NewHttpClient())
	ranAlarmService := services.NewRanAlarmService(logger, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(logger, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager)
	e2tShutdownManager := managers.NewE2TShutdownManager(logger, config, rnibDataService, e2tInstancesManager, e2tAssociationManager, ranConnectStatusChangeManager)
	e2tKeepAliveWorker := managers.NewE2TKeepAliveWorker(logger, rmrSender, e2tInstancesManager, e2tShutdownManager, config)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider()
	rmrNotificationHandlerProvider.Init(logger, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager)

	notificationManager := notificationmanager.NewNotificationManager(logger, rmrNotificationHandlerProvider)
	rmrReceiver := rmrreceiver.NewRmrReceiver(logger, rmrMessenger, notificationManager)
	nodebValidator := managers.NewNodebValidator()
	updateEnbManager := managers.NewUpdateEnbManager(logger, rnibDataService, nodebValidator)
	updateGnbManager := managers.NewUpdateGnbManager(logger, rnibDataService, nodebValidator)

	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()

	defer rmrMessenger.Close()

	go rmrReceiver.ListenAndHandle()
	go e2tKeepAliveWorker.Execute()

	httpMsgHandlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(logger, rmrSender, config, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager, nodebValidator, updateEnbManager, updateGnbManager, ranListManager)
	rootController := controllers.NewRootController(rnibDataService)
	nodebController := controllers.NewNodebController(logger, httpMsgHandlerProvider)
	e2tController := controllers.NewE2TController(logger, httpMsgHandlerProvider)
	symptomController := controllers.NewSymptomdataController(logger, httpMsgHandlerProvider, rnibDataService, ranListManager)
	_ = httpserver.Run(logger, config.Http.Port, rootController, nodebController, e2tController, symptomController)
}
