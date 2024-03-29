//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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
	"e2mgr/logger"
	"e2mgr/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouterAndMocks() (*mux.Router, *mocks.RootControllerMock, *mocks.NodebControllerMock, *mocks.E2TControllerMock, *mocks.SymptomdataControllerMock) {
	rootControllerMock := &mocks.RootControllerMock{}
	rootControllerMock.On("HandleHealthCheckRequest").Return(nil)

	nodebControllerMock := &mocks.NodebControllerMock{}
	nodebControllerMock.On("Shutdown").Return(nil)
	nodebControllerMock.On("GetNodeb").Return(nil)
	nodebControllerMock.On("GetNodebIdList").Return(nil)
	nodebControllerMock.On("GetNodebId").Return(nil)
	nodebControllerMock.On("SetGeneralConfiguration").Return(nil)
	nodebControllerMock.On("DeleteEnb").Return(nil)
	nodebControllerMock.On("AddEnb").Return(nil)
	nodebControllerMock.On("UpdateEnb").Return(nil)
	nodebControllerMock.On("HealthCheckRequest").Return(nil)

	e2tControllerMock := &mocks.E2TControllerMock{}
	e2tControllerMock.On("GetE2TInstances").Return(nil)

	symptomdataControllerMock := &mocks.SymptomdataControllerMock{}
	symptomdataControllerMock.On("GetSymptomData").Return(nil)

	router := mux.NewRouter()
	initializeRoutes(router, rootControllerMock, nodebControllerMock, e2tControllerMock, symptomdataControllerMock)
	return router, rootControllerMock, nodebControllerMock, e2tControllerMock, symptomdataControllerMock
}

func TestRouteGetNodebIdList(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/states", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "GetNodebIdList", 1)
}

func TestRouteGetNodebId(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/states/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	nodebControllerMock.AssertNumberOfCalls(t, "GetNodebId", 1)
}

func TestRouteGetNodebRanName(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "ran1", rr.Body.String(), "handler returned wrong body")
	nodebControllerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
}

func TestRouteGetHealth(t *testing.T) {
	router, rootControllerMock, _, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	rootControllerMock.AssertNumberOfCalls(t, "HandleHealthCheckRequest", 1)
}

func TestRoutePutNodebShutdown(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "Shutdown", 1)
}

func TestHealthCheckRequest(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code, "handler returned wrong status code")
	nodebControllerMock.AssertNumberOfCalls(t, "HealthCheckRequest", 1)
}

func TestRoutePutNodebSetGeneralConfiguration(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/parameters", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "SetGeneralConfiguration", 1)
}

func TestRoutePutUpdateEnb(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/enb/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "UpdateEnb", 1)
}

func TestRouteNotFound(t *testing.T) {
	router, _, _, _, _ := setupRouterAndMocks()

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
	err := Run(log, 1234567, &mocks.RootControllerMock{}, &mocks.NodebControllerMock{}, &mocks.E2TControllerMock{}, &mocks.SymptomdataControllerMock{})
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	log := initLog(t)
	_, rootControllerMock, nodebControllerMock, e2tControllerMock, symptomdataControllerMock := setupRouterAndMocks()
	go Run(log, 11223, rootControllerMock, nodebControllerMock, e2tControllerMock, symptomdataControllerMock)

	time.Sleep(time.Millisecond * 100)
	resp, err := http.Get("http://localhost:11223/v1/health")
	if err != nil {
		t.Fatalf("failed to perform GET to http://localhost:11223/v1/health")
	}
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRouteAddEnb(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("POST", "/v1/nodeb/enb", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	nodebControllerMock.AssertNumberOfCalls(t, "AddEnb", 1)
}

func TestRouteDeleteEnb(t *testing.T) {
	router, _, nodebControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("DELETE", "/v1/nodeb/enb/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code, "handler returned wrong status code")
	nodebControllerMock.AssertNumberOfCalls(t, "DeleteEnb", 1)
}

func initLog(t *testing.T) *logger.Logger {
	InfoLevel := int8(3)
	log, err := logger.InitLogger(InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}
