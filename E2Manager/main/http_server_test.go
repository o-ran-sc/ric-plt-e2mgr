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
	"e2mgr/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupRouterAndMocks() (*mux.Router, *mocks.ControllerMock, *mocks.NodebControllerMock) {
	controllerMock := &mocks.ControllerMock{}
	controllerMock.On("ShutdownHandler").Return(nil)
	controllerMock.On("X2ResetHandler").Return(nil)
	controllerMock.On("X2SetupHandler").Return(nil)
	controllerMock.On("EndcSetupHandler").Return(nil)

	nodebControllerMock := &mocks.NodebControllerMock{}
	nodebControllerMock.On("GetNodebIdList").Return(nil)
	nodebControllerMock.On("GetNodeb").Return(nil)
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

	controllerMock.AssertNumberOfCalls(t,"EndcSetupHandler", 1)
}

func TestRoutePostX2Setup(t *testing.T) {
	router, controllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("POST", "/v1/nodeb/x2-setup", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t,"X2SetupHandler", 1)
}

func TestRouteGetNodebIds(t *testing.T) {
	router, _, nodebControllerMock := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ids", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "GetNodebIdList", 1)
}

func TestRouteGetNodebRanName(t *testing.T) {
	router, _, nodebControllerMock := setupRouterAndMocks()

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

	controllerMock.AssertNumberOfCalls(t, "ShutdownHandler", 1)
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
	controllerMock.AssertNumberOfCalls(t, "X2ResetHandler", 1)
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

func TestParseConfigurationSuccess(t *testing.T) {
	config := configuration.ParseConfiguration()
	assert.Equal(t, 3800, config.Http.Port)
	assert.Equal(t, 3801, config.Rmr.Port)
	assert.Equal(t, 4096, config.Rmr.MaxMsgSize)
	assert.Equal(t, "info", config.Logging.LogLevel)
	assert.Equal(t, 100, config.NotificationResponseBuffer)
	assert.Equal(t, 5, config.BigRedButtonTimeoutSec)
}

func TestParseConfigurationFileNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	assert.Panics(t, func() { configuration.ParseConfiguration() })
}

func TestRmrConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 3800},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#http_server.fillRmrConfig - failed to fill RMR configuration: The entry 'rmr' not found\n", func() { configuration.ParseConfiguration() })
}

func TestLoggingConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestLoggingConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestLoggingConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":  map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"http": map[string]interface{}{"port": 3800},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#http_server_test.TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#http_server.fillLoggingConfig - failed to fill logging configuration: The entry 'logging' not found\n",
		func() { configuration.ParseConfiguration() })
}

func TestHttpConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestHttpConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestHttpConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#http_server_test.TestHttpConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#http_server_test.TestHttpConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#http_server.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n",
		func() { configuration.ParseConfiguration() })
}
