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

package configuration

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigurationSuccess(t *testing.T) {
	config := ParseConfiguration()
	assert.Equal(t, 3800, config.Http.Port)
	assert.Equal(t, 3801, config.Rmr.Port)
	assert.Equal(t, 65536, config.Rmr.MaxMsgSize)
	assert.Equal(t, "info", config.Logging.LogLevel)
	assert.Equal(t, 100, config.NotificationResponseBuffer)
	assert.Equal(t, 5, config.BigRedButtonTimeoutSec)
	assert.Equal(t, 1500, config.KeepAliveResponseTimeoutMs)
	assert.Equal(t, 500, config.KeepAliveDelayMs)
}

func TestParseConfigurationFileNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	assert.Panics(t, func() { ParseConfiguration() })
}

func TestRmrConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestRmrConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 3800},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.fillRmrConfig - failed to fill RMR configuration: The entry 'rmr' not found\n", func() { ParseConfiguration() })
}

func TestLoggingConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestLoggingConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestLoggingConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":  map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"http": map[string]interface{}{"port": 3800},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.fillLoggingConfig - failed to fill logging configuration: The entry 'logging' not found\n",
		func() { ParseConfiguration() })
}

func TestHttpConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestHttpConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestHttpConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestHttpConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestHttpConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n",
		func() { ParseConfiguration() })
}
