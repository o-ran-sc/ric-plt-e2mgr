//
// Copyright (c) 2022 Samsung Electronics Co., Ltd. All Rights Reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewE2nodeConfigurationUpdateSuccessResponseMessage(t *testing.T) {
	configurationUpdate := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	ack := models.NewE2nodeConfigurationUpdateSuccessResponseMessage(configurationUpdate)
	successOutcome := ack.Outcome.(models.E2nodeConfigurationUpdateAcknowledgeSuccessfulOutcome)

	assert.Equal(t, models.ProcedureCode_id_E2nodeConfigurationUpdate, successOutcome.ProcedureCode)
	assert.Equal(t, 4, len(successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs))

	txIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[0]
	assert.Equal(t, models.ProtocolIE_ID_id_TransactionID, txIE.ID)
	assert.Equal(t, "1234", txIE.Value.(models.E2nodeConfigurationUpdateAcknowledgeTransID).TransactionID)

	additionIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[1]
	assert.Equal(t, models.ProtocolIE_ID_id_E2nodeComponentConfigAdditionAck, additionIE.ID)
	assert.Equal(t, 5, len(additionIE.Value.(models.E2nodeComponentConfigAdditionAckList).E2nodeComponentConfigAdditionAckList.ProtocolIESingleContainer))

	updateIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[2]
	assert.Equal(t, models.ProtocolIE_ID_id_E2nodeComponentConfigUpdateAck, updateIE.ID)
	assert.Equal(t, 5, len(updateIE.Value.(models.E2nodeComponentConfigUpdateAckList).E2nodeComponentConfigUpdateAckList.ProtocolIESingleContainer))

	removalIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[3]
	assert.Equal(t, models.ProtocolIE_ID_id_E2nodeComponentConfigRemovalAck, removalIE.ID)
	assert.Equal(t, 5, len(removalIE.Value.(models.E2nodeComponentConfigRemovalAckList).E2nodeComponentConfigRemovalAckList.ProtocolIESingleContainer))
}

func TestNewE2nodeConfigurationUpdateSuccessResponseMessageAdditionOnly(t *testing.T) {
	configurationUpdate := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateOnlyAdditionXmlPath)
	ack := models.NewE2nodeConfigurationUpdateSuccessResponseMessage(configurationUpdate)
	successOutcome := ack.Outcome.(models.E2nodeConfigurationUpdateAcknowledgeSuccessfulOutcome)

	assert.Equal(t, models.ProcedureCode_id_E2nodeConfigurationUpdate, successOutcome.ProcedureCode)
	assert.Equal(t, 2, len(successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs))

	txIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[0]
	assert.Equal(t, models.ProtocolIE_ID_id_TransactionID, txIE.ID)
	assert.Equal(t, "1234", txIE.Value.(models.E2nodeConfigurationUpdateAcknowledgeTransID).TransactionID)

	additionIE := successOutcome.Value.E2nodeConfigurationUpdateAcknowledge.ProtocolIEs.E2nodeConfigurationUpdateAcknowledgeIEs[1]
	assert.Equal(t, models.ProtocolIE_ID_id_E2nodeComponentConfigAdditionAck, additionIE.ID)
	assert.Equal(t, 1, len(additionIE.Value.(models.E2nodeComponentConfigAdditionAckList).E2nodeComponentConfigAdditionAckList.ProtocolIESingleContainer))
}
