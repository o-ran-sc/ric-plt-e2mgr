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


package httpserver

import (
	"e2mgr/controllers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run(port int, controller controllers.IRootController, newController controllers.INodebController) {

	router := mux.NewRouter();
	initializeRoutes(router, controller, newController)

	addr := fmt.Sprintf(":%d", port)

	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatalf("#http_server.Run - Fail initiating HTTP server. Error: %v", err)
	}
}

func initializeRoutes(router *mux.Router, rootController controllers.IRootController, nodebController controllers.INodebController) {
	r := router.PathPrefix("/v1").Subrouter()
	r.HandleFunc("/health", rootController.HandleHealthCheckRequest).Methods("GET")

	rr := r.PathPrefix("/nodeb").Subrouter()
	rr.HandleFunc("/ids", nodebController.GetNodebIdList).Methods("GET")
	rr.HandleFunc("/{ranName}", nodebController.GetNodeb).Methods("GET")
	rr.HandleFunc("/shutdown", nodebController.Shutdown).Methods("PUT")
	rr.HandleFunc("/{ranName}/reset", nodebController.X2Reset).Methods("PUT")
	rr.HandleFunc("/x2-setup", nodebController.X2Setup).Methods("POST")
	rr.HandleFunc("/endc-setup", nodebController.EndcSetup).Methods("POST")
}
