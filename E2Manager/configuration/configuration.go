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
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type RnibWriterConfig struct {
	StateChangeMessageChannel     string
	RanManipulationMessageChannel string
}

type Configuration struct {
	Logging struct {
		LogLevel string
	}
	Http struct {
		Port int
	}
	Rmr struct {
		Port       int
		MaxMsgSize int
	}
	RoutingManager struct {
		BaseUrl string
	}

	NotificationResponseBuffer   int
	BigRedButtonTimeoutSec       int
	MaxRnibConnectionAttempts    int
	RnibRetryIntervalMs          int
	KeepAliveResponseTimeoutMs   int
	KeepAliveDelayMs             int
	E2TInstanceDeletionTimeoutMs int
	E2ResetTimeOutSec            int
	GlobalRicId                  struct {
		RicId string
		Mcc   string
		Mnc   string
	}
	RnibWriter RnibWriterConfig
}

func ParseConfiguration() *Configuration {
	viper.SetConfigType("yaml")
	viper.SetConfigName("configuration")
	viper.AddConfigPath("E2Manager/resources/")
	viper.AddConfigPath("./resources/")     //For production
	viper.AddConfigPath("../resources/")    //For test under Docker
	viper.AddConfigPath("../../resources/") //For test under Docker
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("#configuration.ParseConfiguration - failed to read configuration file: %s\n", err))
	}

	config := Configuration{}
	config.populateRmrConfig(viper.Sub("rmr"))
	config.populateHttpConfig(viper.Sub("http"))
	config.populateLoggingConfig(viper.Sub("logging"))
	config.populateRoutingManagerConfig(viper.Sub("routingManager"))
	config.NotificationResponseBuffer = viper.GetInt("notificationResponseBuffer")
	config.BigRedButtonTimeoutSec = viper.GetInt("bigRedButtonTimeoutSec")
	config.MaxRnibConnectionAttempts = viper.GetInt("maxRnibConnectionAttempts")
	config.RnibRetryIntervalMs = viper.GetInt("rnibRetryIntervalMs")
	config.KeepAliveResponseTimeoutMs = viper.GetInt("keepAliveResponseTimeoutMs")
	config.KeepAliveDelayMs = viper.GetInt("KeepAliveDelayMs")
	config.E2TInstanceDeletionTimeoutMs = viper.GetInt("e2tInstanceDeletionTimeoutMs")
	//E2ResetTimeOutSec : timeout expiry threshold required for handling reset and thus the time for which the nodeb is under reset connection state.
	config.E2ResetTimeOutSec = viper.GetInt("e2ResetTimeOutSec")
	config.populateGlobalRicIdConfig(viper.Sub("globalRicId"))
	config.populateRnibWriterConfig(viper.Sub("rnibWriter"))
	return &config
}

func (c *Configuration) populateLoggingConfig(logConfig *viper.Viper) {
	if logConfig == nil {
		panic(fmt.Sprintf("#configuration.populateLoggingConfig - failed to populate logging configuration: The entry 'logging' not found\n"))
	}
	c.Logging.LogLevel = logConfig.GetString("logLevel")
}

func (c *Configuration) populateHttpConfig(httpConfig *viper.Viper) {
	if httpConfig == nil {
		panic(fmt.Sprintf("#configuration.populateHttpConfig - failed to populate HTTP configuration: The entry 'http' not found\n"))
	}
	c.Http.Port = httpConfig.GetInt("port")
}

func (c *Configuration) populateRmrConfig(rmrConfig *viper.Viper) {
	if rmrConfig == nil {
		panic(fmt.Sprintf("#configuration.populateRmrConfig - failed to populate RMR configuration: The entry 'rmr' not found\n"))
	}
	c.Rmr.Port = rmrConfig.GetInt("port")
	c.Rmr.MaxMsgSize = rmrConfig.GetInt("maxMsgSize")
}

func (c *Configuration) populateRoutingManagerConfig(rmConfig *viper.Viper) {
	if rmConfig == nil {
		panic(fmt.Sprintf("#configuration.populateRoutingManagerConfig - failed to populate Routing Manager configuration: The entry 'routingManager' not found\n"))
	}
	c.RoutingManager.BaseUrl = rmConfig.GetString("baseUrl")
}

func (c *Configuration) populateRnibWriterConfig(rnibWriterConfig *viper.Viper) {
	if rnibWriterConfig == nil {
		panic(fmt.Sprintf("#configuration.populateRnibWriterConfig - failed to populate Rnib Writer configuration: The entry 'rnibWriter' not found\n"))
	}
	c.RnibWriter.StateChangeMessageChannel = rnibWriterConfig.GetString("stateChangeMessageChannel")
	c.RnibWriter.RanManipulationMessageChannel = rnibWriterConfig.GetString("ranManipulationMessageChannel")
}

func (c *Configuration) populateGlobalRicIdConfig(globalRicIdConfig *viper.Viper) {
	err := validateGlobalRicIdConfig(globalRicIdConfig)
	if err != nil {
		panic(err.Error())
	}
	c.GlobalRicId.RicId = globalRicIdConfig.GetString("ricId")
	c.GlobalRicId.Mcc = globalRicIdConfig.GetString("mcc")
	c.GlobalRicId.Mnc = globalRicIdConfig.GetString("mnc")
}

func validateGlobalRicIdConfig(globalRicIdConfig *viper.Viper) error {
	if globalRicIdConfig == nil {
		return errors.New("#configuration.validateGlobalRicIdConfig - failed to populate Global RicId configuration: The entry 'globalRicId' not found\n")
	}

	err := validateRicId(globalRicIdConfig.GetString("ricId"))

	if err != nil {
		return err
	}

	err = validateMcc(globalRicIdConfig.GetString("mcc"))

	if err != nil {
		return err
	}

	err = validateMnc(globalRicIdConfig.GetString("mnc"))

	if err != nil {
		return err
	}

	return nil
}

func validateMcc(mcc string) error {

	if len(mcc) == 0 {
		return errors.New("#configuration.validateMcc - mcc is missing or empty\n")
	}

	if len(mcc) != 3 {
		return errors.New("#configuration.validateMcc - mcc is not 3 digits\n")
	}

	mccInt, err := strconv.Atoi(mcc)

	if err != nil {
		return errors.New("#configuration.validateMcc - mcc is not a number\n")
	}

	if mccInt < 0 {
		return errors.New("#configuration.validateMcc - mcc is negative\n")
	}
	return nil
}

func validateMnc(mnc string) error {

	if len(mnc) == 0 {
		return errors.New("#configuration.validateMnc - mnc is missing or empty\n")
	}

	if len(mnc) < 2 || len(mnc) > 3 {
		return errors.New("#configuration.validateMnc - mnc is not 2 or 3 digits\n")
	}

	mncAsInt, err := strconv.Atoi(mnc)

	if err != nil {
		return errors.New("#configuration.validateMnc - mnc is not a number\n")
	}

	if mncAsInt < 0 {
		return errors.New("#configuration.validateMnc - mnc is negative\n")
	}

	return nil
}

func validateRicId(ricId string) error {

	if len(ricId) == 0 {
		return errors.New("#configuration.validateRicId - ricId is missing or empty\n")
	}

	if len(ricId) != 5 {
		return errors.New("#configuration.validateRicId - ricId length should be 5 hex characters\n")
	}

	_, err := strconv.ParseUint(ricId, 16, 64)
	if err != nil {
		return errors.New("#configuration.validateRicId - ricId is not hex number\n")
	}

	return nil
}

func (c *Configuration) String() string {
	return fmt.Sprintf("{logging.logLevel: %s, http.port: %d, rmr: { port: %d, maxMsgSize: %d}, routingManager.baseUrl: %s, "+
		"notificationResponseBuffer: %d, bigRedButtonTimeoutSec: %d, maxRnibConnectionAttempts: %d, "+
		"rnibRetryIntervalMs: %d, keepAliveResponseTimeoutMs: %d, keepAliveDelayMs: %d, e2tInstanceDeletionTimeoutMs: %d,e2ResetTimeOutSec: %d,"+
		"globalRicId: { ricId: %s, mcc: %s, mnc: %s}, rnibWriter: { stateChangeMessageChannel: %s, ranManipulationChannel: %s}",
		c.Logging.LogLevel,
		c.Http.Port,
		c.Rmr.Port,
		c.Rmr.MaxMsgSize,
		c.RoutingManager.BaseUrl,
		c.NotificationResponseBuffer,
		c.BigRedButtonTimeoutSec,
		c.MaxRnibConnectionAttempts,
		c.RnibRetryIntervalMs,
		c.KeepAliveResponseTimeoutMs,
		c.KeepAliveDelayMs,
		c.E2TInstanceDeletionTimeoutMs,
		c.E2ResetTimeOutSec,
		c.GlobalRicId.RicId,
		c.GlobalRicId.Mcc,
		c.GlobalRicId.Mnc,
		c.RnibWriter.StateChangeMessageChannel,
		c.RnibWriter.RanManipulationMessageChannel,
	)
}
