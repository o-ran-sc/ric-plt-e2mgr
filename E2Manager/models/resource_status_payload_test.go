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
	"e2mgr/enums"
	"e2mgr/models"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceStatusPayloadSuccess(t *testing.T) {
	expectedResourceStatusPayload := &models.ResourceStatusPayload{
		NodeType:         entities.Node_ENB,
		MessageDirection: enums.RIC_TO_RAN,
	}

	payload := models.NewResourceStatusPayload(entities.Node_ENB, enums.RIC_TO_RAN)
	assert.Equal(t, expectedResourceStatusPayload, payload)
}
