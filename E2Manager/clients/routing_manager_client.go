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
	"io/ioutil"
	"net/http"
)

const (
	AddE2TInstanceApiSuffix            = "e2t"
	AssociateRanToE2TInstanceApiSuffix = "associate-ran-to-e2t"
	DissociateRanE2TInstanceApiSuffix  = "dissociate-ran"
)

type RoutingManagerClient struct {
	logger     *logger.Logger
	config     *configuration.Configuration
	httpClient HttpClient
}

type IRoutingManagerClient interface {
	AddE2TInstance(e2tAddress string) error
	AssociateRanToE2TInstance(e2tAddress string, ranName string) error
	DissociateRanE2TInstance(e2tAddress string, ranName string) error
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
	url := c.config.RoutingManagerBaseUrl + AddE2TInstanceApiSuffix

	return c.PostMessage(data, url)
}

func (c *RoutingManagerClient) AssociateRanToE2TInstance(e2tAddress string, ranName string) error {

	data := models.NewRoutingManagerE2TData(e2tAddress, ranName)
	url := c.config.RoutingManagerBaseUrl + AssociateRanToE2TInstanceApiSuffix

	return c.PostMessage(data, url)
}

func (c *RoutingManagerClient) DissociateRanE2TInstance(e2tAddress string, ranName string) error {

	data := models.NewRoutingManagerE2TData(e2tAddress, ranName)
	url := c.config.RoutingManagerBaseUrl + DissociateRanE2TInstanceApiSuffix

	return c.PostMessage(data, url)
}

func (c *RoutingManagerClient) PostMessage(data *models.RoutingManagerE2TData, url string) error {
	marshaled, err := json.Marshal(data)

	if err != nil {
		return e2managererrors.NewRoutingManagerError(err)
	}

	body := bytes.NewBuffer(marshaled)
	c.logger.Infof("[E2M -> Routing Manager] #RoutingManagerClient.PostMessage - request body: %+v", body)

	resp, err := c.httpClient.Post(url, "application/json", body)

	if err != nil {
		return e2managererrors.NewRoutingManagerError(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return e2managererrors.NewRoutingManagerError(err)
	}

	if resp.StatusCode != http.StatusOK { // TODO: shall we check for != 201?
		c.logger.Errorf("[Routing Manager -> E2M] #RoutingManagerClient.PostMessage - failure. http status code: %d, response body: %s", resp.StatusCode, string(respBody))
		return e2managererrors.NewRoutingManagerError(fmt.Errorf("Invalid data")) // TODO: which error shall we return?
	}

	c.logger.Infof("[Routing Manager -> E2M] #RoutingManagerClient.PostMessage - success. http status code: %d, response body: %s", resp.StatusCode, string(respBody))
	return nil
}