//
// Copyright 2020 AT&T Intellectual Property
// Copyright 2020 Nokia
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

package models_test

import (
	"e2mgr/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheckResponseMarshalSuccess(t *testing.T) {
	healthMsg := "OK"
	expectedResponse := models.HealthCheckSuccessResponse{
		Message: healthMsg,
	}
	expectedData, _ := json.Marshal(expectedResponse)

	healthCheckSuccessResponse := models.NewHealthCheckSuccessResponse(healthMsg)
	resp, err := healthCheckSuccessResponse.Marshal()
	assert.Nil(t, err)
	assert.Equal(t, expectedData, resp)
}
