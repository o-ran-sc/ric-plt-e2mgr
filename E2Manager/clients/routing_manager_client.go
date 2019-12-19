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

package clients

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AddE2TInstanceApiSuffix = "e2t"
)

type RoutingManagerClient struct {
	logger     *logger.Logger
	config     *configuration.Configuration
	httpClient HttpClient
}

type IRoutingManagerClient interface {
	AddE2TInstance(e2tAddress string) error
}

func NewRoutingManagerClient(logger *logger.Logger, config *configuration.Configuration, httpClient HttpClient) *RoutingManagerClient {
	return &RoutingManagerClient{
		logger:     logger,
		config:     config,
		httpClient: httpClient,
	}
}

func (c *RoutingManagerClient) AddE2TInstance(e2tAddress string) error {
	data := models.NewRoutingManagerE2TData(e2tAddress)

	marshaled, err := json.Marshal(data)

	if err != nil {
		return e2managererrors.NewRoutingManagerError(err)
	}

	body := bytes.NewBuffer(marshaled)
	c.logger.Infof("[E2M -> Routing Manager] #RoutingManagerClient.AddE2TInstance - request body: %+v", body)

	url := c.config.RoutingManager.BaseUrl + AddE2TInstanceApiSuffix
	resp, err := c.httpClient.Post(url, "application/json", body)

	if err != nil {
		c.logger.Errorf("#RoutingManagerClient.AddE2TInstance - failed sending request. error: %s", err)
		return e2managererrors.NewRoutingManagerError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		c.logger.Infof("[Routing Manager -> E2M] #RoutingManagerClient.AddE2TInstance - success. http status code: %d", resp.StatusCode)
		return nil
	}

	c.logger.Errorf("[Routing Manager -> E2M] #RoutingManagerClient.AddE2TInstance - failure. http status code: %d", resp.StatusCode)
	return e2managererrors.NewRoutingManagerError(fmt.Errorf("invalid data"))
}
