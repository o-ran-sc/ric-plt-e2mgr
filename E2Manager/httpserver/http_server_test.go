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

package httpserver

import (
	"e2mgr/logger"
	"e2mgr/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouterAndMocks() (*mux.Router, *mocks.ControllerMock, *mocks.NodebControllerMock) {
	controllerMock := &mocks.ControllerMock{}
	controllerMock.On("Shutdown").Return(nil)
	controllerMock.On("X2Reset").Return(nil)
	controllerMock.On("X2Setup").Return(nil)
	controllerMock.On("EndcSetup").Return(nil)
	controllerMock.On("GetNodeb").Return(nil)
	controllerMock.On("GetNodebIdList").Return(nil)



	nodebControllerMock := &mocks.NodebControllerMock{}
	nodebControllerMock.On("GetNodebIdList").Return(nil)
	nodebControllerMock.On("GetNodeb").Return(nil) // TODO: remove
	nodebControllerMock.On("HandleHealthCheckRequest").Return(nil)

	router := mux.NewRouter()
	initializeRoutes(router, nodebControllerMock, controllerMock)
	return router, controllerMock, nodebControllerMock
}

func TestRoutePostEndcSetup(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("POST", "/v1/nodeb/endc-setup", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t,"EndcSetup", 1)
}

func TestRoutePostX2Setup(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("POST", "/v1/nodeb/x2-setup", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t,"X2Setup", 1)
}

func TestRouteGetNodebIds(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ids", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t, "GetNodebIdList", 1)
}

func TestRouteGetNodebRanName(t *testing.T) {
	router, controllerMock,_ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "ran1", rr.Body.String(), "handler returned wrong body")
	controllerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
}

func TestRouteGetHealth(t *testing.T) {
	router, _, nodebControllerMock := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "HandleHealthCheckRequest", 1)
}

func TestRoutePutNodebShutdown(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t, "Shutdown", 1)
}

func TestRoutePutNodebResetRanName(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/ran1/reset", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "ran1", rr.Body.String(), "handler returned wrong body")
	controllerMock.AssertNumberOfCalls(t, "X2Reset", 1)
}

func TestRouteNotFound(t *testing.T) {
	router, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/no/such/route", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestRunError(t *testing.T) {
	log := initLog(t)
	err := Run(log, 1234567, &mocks.NodebControllerMock{}, &mocks.ControllerMock{})
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	log := initLog(t)
	_, controllerMock, nodebControllerMock := setupRouterAndMocks()
	go Run(log,11223, nodebControllerMock, controllerMock)

	time.Sleep(time.Millisecond * 100)
	resp, err := http.Get("http://localhost:11223/v1/health")
	if err != nil {
		t.Fatalf("failed to perform GET to http://localhost:11223/v1/health")
	}
	assert.Equal(t, 200, resp.StatusCode)
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}