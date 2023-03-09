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

package configuration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParseConfigurationSuccess(t *testing.T) {
	config := ParseConfiguration()
	assert.Equal(t, 3800, config.Http.Port)
	assert.Equal(t, 3801, config.Rmr.Port)
	assert.Equal(t, 65536, config.Rmr.MaxMsgSize)
	assert.Equal(t, "info", config.Logging.LogLevel)
	assert.Equal(t, 100, config.NotificationResponseBuffer)
	assert.Equal(t, 5, config.BigRedButtonTimeoutSec)
	assert.Equal(t, 4500, config.KeepAliveResponseTimeoutMs)
	assert.Equal(t, 1500, config.KeepAliveDelayMs)
	assert.Equal(t, 15000, config.E2TInstanceDeletionTimeoutMs)
	assert.Equal(t, 10, config.E2ResetTimeOutSec)
	assert.NotNil(t, config.GlobalRicId)
	assert.Equal(t, "AACCE", config.GlobalRicId.RicId)
	assert.Equal(t, "310", config.GlobalRicId.Mcc)
	assert.Equal(t, "411", config.GlobalRicId.Mnc)
	assert.Equal(t, "RAN_CONNECTION_STATUS_CHANGE", config.RnibWriter.StateChangeMessageChannel)
	assert.Equal(t, "RAN_MANIPULATION", config.RnibWriter.RanManipulationMessageChannel)
}

func TestStringer(t *testing.T) {
	config := ParseConfiguration().String()
	assert.NotEmpty(t, config)
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
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
		"globalRicId":    map[string]interface{}{"plmnId": "131014", "ricNearRtId": "556670"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.populateRmrConfig - failed to populate RMR configuration: The entry 'rmr' not found\n", func() { ParseConfiguration() })
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
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"http":           map[string]interface{}{"port": 3800},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
		"globalRicId":    map[string]interface{}{"plmnId": "131014", "ricNearRtId": "556670"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestRmrConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.populateLoggingConfig - failed to populate logging configuration: The entry 'logging' not found\n",
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
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
		"globalRicId":    map[string]interface{}{"plmnId": "131014", "ricNearRtId": "556670"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestHttpConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestHttpConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.populateHttpConfig - failed to populate HTTP configuration: The entry 'http' not found\n",
		func() { ParseConfiguration() })
}

func TestRoutingManagerConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestRoutingManagerConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestRoutingManagerConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":         map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":     map[string]interface{}{"logLevel": "info"},
		"http":        map[string]interface{}{"port": 3800},
		"globalRicId": map[string]interface{}{"mcc": 327, "mnc": 94, "ricId": "AACCE"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestRoutingManagerConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestRoutingManagerConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.populateRoutingManagerConfig - failed to populate Routing Manager configuration: The entry 'routingManager' not found\n",
		func() { ParseConfiguration() })
}

func TestGlobalRicIdConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateGlobalRicIdConfig - failed to populate Global RicId configuration: The entry 'globalRicId' not found\n",
		func() { ParseConfiguration() })
}

func TestRnibWriterConfigNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
		"globalRicId":    map[string]interface{}{"mcc": 327, "mnc": 94, "ricId": "AACCE"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestGlobalRicIdConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.populateRnibWriterConfig - failed to populate Rnib Writer configuration: The entry 'rnibWriter' not found\n",
		func() { ParseConfiguration() })
}

func TestEmptyRicIdFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestEmptyRicIdFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "327", "mnc": "94", "ricId": ""},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateRicId - ricId is missing or empty\n",
		func() { ParseConfiguration() })
}

func TestMissingRicIdFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestEmptyRicIdFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "327", "mnc": "94"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestEmptyRicIdFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateRicId - ricId is missing or empty\n",
		func() { ParseConfiguration() })
}

func TestNonHexRicIdFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestNonHexRicIdFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestNonHexRicIdFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "327", "mnc": "94", "ricId": "TEST1"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestNonHexRicIdFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestNonHexRicIdFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateRicId - ricId is not hex number\n",
		func() { ParseConfiguration() })
}

func TestWrongRicIdLengthFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestWrongRicIdLengthFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestWrongRicIdLengthFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "327", "mnc": "94", "ricId": "AA43"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestWrongRicIdLengthFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestWrongRicIdLengthFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateRicId - ricId length should be 5 hex characters\n",
		func() { ParseConfiguration() })
}

func TestMccNotThreeDigitsFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestMccNotThreeDigitsFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestMccNotThreeDigitsFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "31", "mnc": "94", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestMccNotThreeDigitsFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestMccNotThreeDigitsFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMcc - mcc is not 3 digits\n",
		func() { ParseConfiguration() })
}

func TestMncLengthIsGreaterThanThreeDigitsFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestMncLengthIsGreaterThanThreeDigitsFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestMncLengthIsGreaterThanThreeDigitsFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "310", "mnc": "6794", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestMncLengthIsGreaterThanThreeDigitsFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestMncLengthIsGreaterThanThreeDigitsFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is not 2 or 3 digits\n",
		func() { ParseConfiguration() })
}

func TestMncLengthIsLessThanTwoDigitsFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestMncLengthIsLessThanTwoDigitsFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestMncLengthIsLessThanTwoDigitsFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "310", "mnc": "4", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestMncLengthIsLessThanTwoDigitsFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestMncLengthIsLessThanTwoDigitsFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is not 2 or 3 digits\n",
		func() { ParseConfiguration() })
}

func TestNegativeMncFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestNegativeMncFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "310", "mnc": "-2", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is negative\n",
		func() { ParseConfiguration() })
}

func TestNegativeMccFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestNegativeMncFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "-31", "mnc": "222", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestNegativeMncFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMcc - mcc is negative\n",
		func() { ParseConfiguration() })
}

func TestAlphaNumericMccFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestAlphaNumericMccFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestAlphaNumericMccFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "1W2", "mnc": "222", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestAlphaNumericMccFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestAlphaNumericMccFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMcc - mcc is not a number\n",
		func() { ParseConfiguration() })
}

func TestAlphaNumericMncFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestAlphaNumericMncFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestAlphaNumericMncFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "111", "mnc": "2A8", "ricId": "AA443"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestAlphaNumericMncFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestAlphaNumericMncFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is not a number\n",
		func() { ParseConfiguration() })
}

func TestMissingMmcFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestMissingMmcFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestMissingMmcFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mnc": "94", "ricId": "AABB3"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestMissingMmcFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestMissingMmcFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMcc - mcc is missing or empty\n",
		func() { ParseConfiguration() })
}

func TestEmptyMmcFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestEmptyMmcFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestEmptyMmcFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "", "mnc": "94", "ricId": "AABB3"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestEmptyMmcFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestEmptyMmcFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMcc - mcc is missing or empty\n",
		func() { ParseConfiguration() })
}

func TestEmptyMncFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestEmptyMncFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestEmptyMncFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "111", "mnc": "", "ricId": "AABB3"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestEmptyMncFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestEmptyMncFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is missing or empty\n",
		func() { ParseConfiguration() })
}

func TestMissingMncFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#TestMissingMncFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#TestMissingMncFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()
	yamlMap := map[string]interface{}{
		"rmr":            map[string]interface{}{"port": 3801, "maxMsgSize": 4096},
		"logging":        map[string]interface{}{"logLevel": "info"},
		"http":           map[string]interface{}{"port": 3800},
		"globalRicId":    map[string]interface{}{"mcc": "111", "ricId": "AABB3"},
		"routingManager": map[string]interface{}{"baseUrl": "http://localhost:8080/ric/v1/handles/"},
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#TestMissingMncFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#TestMissingMncFailure - failed to write configuration file: %s\n", configPath)
	}
	assert.PanicsWithValue(t, "#configuration.validateMnc - mnc is missing or empty\n",
		func() { ParseConfiguration() })
}
