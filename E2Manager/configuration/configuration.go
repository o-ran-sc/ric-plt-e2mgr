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
	"fmt"
	"github.com/spf13/viper"
)

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
	NotificationResponseBuffer   int
	BigRedButtonTimeoutSec       int
	MaxConnectionAttempts        int
	MaxRnibConnectionAttempts    int
	RnibRetryIntervalMs          int
	KeepAliveResponseTimeoutMs 	 int
	KeepAliveDelayMs             int
	RoutingManagerBaseUrl		 string
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
	config.fillRmrConfig(viper.Sub("rmr"))
	config.fillHttpConfig(viper.Sub("http"))
	config.fillLoggingConfig(viper.Sub("logging"))

	config.NotificationResponseBuffer = viper.GetInt("notificationResponseBuffer")
	config.BigRedButtonTimeoutSec = viper.GetInt("bigRedButtonTimeoutSec")
	config.MaxConnectionAttempts = viper.GetInt("maxConnectionAttempts")
	config.MaxRnibConnectionAttempts = viper.GetInt("maxRnibConnectionAttempts")
	config.RnibRetryIntervalMs = viper.GetInt("rnibRetryIntervalMs")
	config.KeepAliveResponseTimeoutMs = viper.GetInt("keepAliveResponseTimeoutMs")
	config.KeepAliveDelayMs = viper.GetInt("KeepAliveDelayMs")
	config.RoutingManagerBaseUrl = viper.GetString("routingManagerBaseUrl")
	return &config
}

func (c *Configuration) fillLoggingConfig(logConfig *viper.Viper) {
	if logConfig == nil {
		panic(fmt.Sprintf("#configuration.fillLoggingConfig - failed to fill logging configuration: The entry 'logging' not found\n"))
	}
	c.Logging.LogLevel = logConfig.GetString("logLevel")
}

func (c *Configuration) fillHttpConfig(httpConfig *viper.Viper) {
	if httpConfig == nil {
		panic(fmt.Sprintf("#configuration.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n"))
	}
	c.Http.Port = httpConfig.GetInt("port")
}

func (c *Configuration) fillRmrConfig(rmrConfig *viper.Viper) {
	if rmrConfig == nil {
		panic(fmt.Sprintf("#configuration.fillRmrConfig - failed to fill RMR configuration: The entry 'rmr' not found\n"))
	}
	c.Rmr.Port = rmrConfig.GetInt("port")
	c.Rmr.MaxMsgSize = rmrConfig.GetInt("maxMsgSize")
}
