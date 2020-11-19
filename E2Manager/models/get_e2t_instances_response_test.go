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

func TestGetE2tInstancesResponseMarshalSuccess(t *testing.T) {
	E2TAddress := "10.0.2.15:38000"
	E2TAddress2 := "10.0.2.16:38001"
	ranNames1 := []string{"test1", "test2", "test3"}
	expectedE2TInstancesResponse := models.E2TInstancesResponse{
		&models.E2TInstanceResponseModel{
			E2TAddress: E2TAddress,
			RanNames:   ranNames1,
		},
		&models.E2TInstanceResponseModel{
			E2TAddress: E2TAddress2,
			RanNames:   []string{},
		},
	}
	e2tInstanceResponseModel1 := models.NewE2TInstanceResponseModel(E2TAddress, ranNames1)
	e2tInstanceResponseModel2 := models.NewE2TInstanceResponseModel(E2TAddress2, []string{})
	e2tInstancesResponse := models.E2TInstancesResponse{e2tInstanceResponseModel1, e2tInstanceResponseModel2}
	expectedData, _ := json.Marshal(expectedE2TInstancesResponse)

	resp, err := e2tInstancesResponse.Marshal()
	assert.Nil(t, err)
	assert.Equal(t, expectedData, resp)
}
