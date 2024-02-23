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
	//"fmt"
    "flag"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
    "gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
    "github.com/spf13/viper"
    "github.com/fsnotify/fsnotify"
	"os"
	"strconv"
)

const GeneralKeyDefaultValue = "{\"enableRic\":true}"
const DEFAULT_CONFIG_FILE = "../resources/configuration.yaml"
const DEFAULT_PORT = "8080"
var Log *logger.Logger

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
/**Dynamic log-level changes **/

func loadConfig() {
        viper.SetConfigFile(parseCmd())
        if err := viper.ReadInConfig(); err != nil {
                Log.Errorf("Error reading config file, %s", err)
        }
        Log.Infof("Using config file: %s\n", viper.ConfigFileUsed())
        // Watch for config file changes and re-read data ...
        watch()
}
func parseCmd() string {
        var fileName *string
        fileName = flag.String("f", DEFAULT_CONFIG_FILE, "Specify the configuration file.")
	flag.String("port", DEFAULT_PORT, "Specify the port file.")
        flag.Parse()

        return *fileName
}

func watch() {
        viper.WatchConfig()
        viper.OnConfigChange(func(e fsnotify.Event) {
                Log.Infof("config file changed  %s", e.Name)
        	setLoglevel()
        })
}

/*********MDC LOG CHNAGES ********/
func setLoglevel() {
    var loglevel int
    if err := viper.UnmarshalKey("loglevel", &loglevel); err != nil {
        Log.Errorf("Unmarshalling failed while reading %d", loglevel)
    }

    switch loglevel {
    case 1:
        Log.Infof("LOGLEVEL is set to ERROR\n")
    case 2:
        Log.Infof("LOGLEVEL is set to WARNING\n")
    case 3:
        Log.Infof("LOGLEVEL is set to INFO\n")
    case 4:
        Log.Infof("LOGLEVEL is set to DEBUG\n")
    }
    Log.SetLevel(loglevel)
}


func main() {
	config := configuration.ParseConfiguration()
        level := int8(4)
	Log, _ = logger.InitLogger(level)
	Log.SetFormat(0)
        Log.SetMdc("e2mgr", "0.2.2")
	/*if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}*/
	Log.Infof("#app.main - Configuration %s", config)
	loadConfig()
	
	setLoglevel()
	sdl := sdlgo.NewSyncStorage()
	err := initKeys(Log, sdl)

	if err != nil {
		os.Exit(1)
	}

	defer sdl.Close()
	rnibDataService := services.NewRnibDataService(Log, config, reader.GetNewRNibReader(sdl), rNibWriter.GetRNibWriter(sdl, config.RnibWriter))

	ranListManager := managers.NewRanListManager(Log, rnibDataService)
	RicServiceUpdateManager := managers.NewRicServiceUpdateManager(Log, rnibDataService)

	err = ranListManager.InitNbIdentityMap()

	if err != nil {
		Log.Errorf("#app.main - quit")
		os.Exit(1)
	}

	var msgImpl *rmrCgo.Context
	rmrMessenger := msgImpl.Init("tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, Log)
	rmrSender := rmrsender.NewRmrSender(Log, rmrMessenger)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, Log)
	routingManagerClient := clients.NewRoutingManagerClient(Log, config, clients.NewHttpClient())
	ranAlarmService := services.NewRanAlarmService(Log, config)
	ranConnectStatusChangeManager := managers.NewRanConnectStatusChangeManager(Log, rnibDataService, ranListManager, ranAlarmService)
	e2tAssociationManager := managers.NewE2TAssociationManager(Log, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager)
	e2tShutdownManager := managers.NewE2TShutdownManager(Log, config, rnibDataService, e2tInstancesManager, e2tAssociationManager, ranConnectStatusChangeManager)
	e2tKeepAliveWorker := managers.NewE2TKeepAliveWorker(Log, rmrSender, e2tInstancesManager, e2tShutdownManager, config)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider()
	rmrNotificationHandlerProvider.Init(Log, config, rnibDataService, rmrSender, e2tInstancesManager, routingManagerClient, e2tAssociationManager, ranConnectStatusChangeManager, ranListManager, RicServiceUpdateManager)

	notificationManager := notificationmanager.NewNotificationManager(Log, rmrNotificationHandlerProvider)
	rmrReceiver := rmrreceiver.NewRmrReceiver(Log, rmrMessenger, notificationManager)
	nodebValidator := managers.NewNodebValidator()
	updateEnbManager := managers.NewUpdateEnbManager(Log, rnibDataService, nodebValidator)
	updateGnbManager := managers.NewUpdateGnbManager(Log, rnibDataService, nodebValidator)

	e2tInstancesManager.ResetKeepAliveTimestampsForAllE2TInstances()

	defer rmrMessenger.Close()

	go rmrReceiver.ListenAndHandle()
	go e2tKeepAliveWorker.Execute()

	httpMsgHandlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(Log, rmrSender, config, rnibDataService, e2tInstancesManager, routingManagerClient, ranConnectStatusChangeManager, nodebValidator, updateEnbManager, updateGnbManager, ranListManager)
	rootController := controllers.NewRootController(rnibDataService)
	nodebController := controllers.NewNodebController(Log, httpMsgHandlerProvider)
	e2tController := controllers.NewE2TController(Log, httpMsgHandlerProvider)
	symptomController := controllers.NewSymptomdataController(Log, httpMsgHandlerProvider, rnibDataService, ranListManager)
        //fmt.Println("loadconfig called at last")
        //loadConfig()
	_ = httpserver.Run(Log, config.Http.Port, rootController, nodebController, e2tController, symptomController)
	//fmt.Println("loadconfig called at last")
	//loadConfig()
}
