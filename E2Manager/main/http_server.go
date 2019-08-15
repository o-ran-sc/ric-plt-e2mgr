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
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/receivers"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

const MAX_RNIB_PULL_INSTANCES = 4

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
	rNibWriter.Init("e2Manager", MAX_RNIB_PULL_INSTANCES)
	defer rNibWriter.Close()
	reader.Init("e2Manager", MAX_RNIB_PULL_INSTANCES)
	defer reader.Close()
	var nManager = managers.NewNotificationManager(reader.GetRNibReader, rNibWriter.GetRNibWriter)

	rmrResponseChannel := make(chan *models.NotificationResponse, config.NotificationResponseBuffer)
	rmrService := services.NewRmrService(rmrConfig, msgImpl, controllers.E2Sessions, rmrResponseChannel)
	rmrServiceReceiver := receivers.NewRmrServiceReceiver(*rmrService, nManager)
	defer rmrService.CloseContext()
	go rmrServiceReceiver.ListenAndHandle()
	go rmrService.SendResponse()
	runServer(rmrService, logger, config, rmrResponseChannel)
}

func runServer(rmrService *services.RmrService, logger *logger.Logger, config *configuration.Configuration, rmrResponseChannel chan *models.NotificationResponse) {

	router := httprouter.New()
	controller := controllers.NewNodebController(logger, rmrService, reader.GetRNibReader, rNibWriter.GetRNibWriter)
	newController := controllers.NewController(logger, rmrService, reader.GetRNibReader, rNibWriter.GetRNibWriter, config, rmrResponseChannel)

	router.POST("/v1/nodeb/:messageType", controller.HandleRequest)
	router.GET("/v1/nodeb-ids", controller.GetNodebIdList)
	router.GET("/v1/nodeb/:ranName", controller.GetNodeb)
	router.GET("/v1/health", controller.HandleHealthCheckRequest)
	router.PUT("/v1/nodeb/shutdown", newController.ShutdownHandler)
	router.PUT("/v1/nodeb-reset/:ranName", newController.X2ResetHandler)

	port := fmt.Sprintf(":%d", config.Http.Port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("#runNodebServer - fail to start http server. Error: %v", err)
	}
}