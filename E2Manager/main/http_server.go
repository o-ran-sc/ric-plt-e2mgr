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
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/managers/notificationmanager"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/receivers"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const MAX_RNIB_POOL_INSTANCES = 4

func main() {
	config := configuration.ParseConfiguration()
	logLevel, _ := logger.LogLevelTokenToLevel(config.Logging.LogLevel)
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#http_server.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}
	rmrConfig := services.NewRmrConfig(config.Rmr.Port, config.Rmr.MaxMsgSize, 0, logger)
	var msgImpl *rmrCgo.Context
	rNibWriter.Init("e2Manager", MAX_RNIB_POOL_INSTANCES)
	defer rNibWriter.Close()
	reader.Init("e2Manager", MAX_RNIB_POOL_INSTANCES)
	defer reader.Close()

	rmrResponseChannel := make(chan *models.NotificationResponse, config.NotificationResponseBuffer)
	rmrService := services.NewRmrService(rmrConfig, msgImpl, controllers.E2Sessions, rmrResponseChannel)
	var ranSetupManager = managers.NewRanSetupManager(logger, rmrService, rNibWriter.GetRNibWriter)
	var ranReconnectionManager = managers.NewRanReconnectionManager(logger, config, reader.GetRNibReader, rNibWriter.GetRNibWriter, ranSetupManager)
	var nManager = notificationmanager.NewNotificationManager(reader.GetRNibReader, rNibWriter.GetRNibWriter, ranReconnectionManager)

	rmrServiceReceiver := receivers.NewRmrServiceReceiver(*rmrService, nManager)
	defer rmrService.CloseContext()
	go rmrServiceReceiver.ListenAndHandle()
	go rmrService.SendResponse()

	controller := controllers.NewNodebController(logger, rmrService, reader.GetRNibReader, rNibWriter.GetRNibWriter)
	newController := controllers.NewController(logger, rmrService, reader.GetRNibReader, rNibWriter.GetRNibWriter, config, ranSetupManager)
	runServer(config.Http.Port, controller, newController)
}

func runServer(port int, controller controllers.INodebController, newController controllers.IController) {

	router := mux.NewRouter();
	initializeRoutes(router, controller, newController)

	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("#runNodebServer - fail to start http server. Error: %v", err)
	}
}

func initializeRoutes(router *mux.Router, controller controllers.INodebController, newController controllers.IController) {
	r := router.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/health", controller.HandleHealthCheckRequest).Methods("GET")

	rr := r.PathPrefix("/nodeb").Subrouter()
	rr.HandleFunc("/ids", controller.GetNodebIdList).Methods("GET") // nodeb/ids
	rr.HandleFunc("/{ranName}", controller.GetNodeb).Methods("GET")
	rr.HandleFunc("/shutdown", newController.ShutdownHandler).Methods("PUT")
	rr.HandleFunc("/{ranName}/reset", newController.X2ResetHandler).Methods("PUT") // nodeb/{ranName}/reset
	rr.HandleFunc("/x2-setup", newController.X2SetupHandler).Methods("POST")
	rr.HandleFunc("/endc-setup", newController.EndcSetupHandler).Methods("POST")
}
