//
// Copyright 2022 Samsung Electronics Co.
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
	"e2mgr/utils"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	e2NodeConfigurationUpdateOnlyAdditionXmlPath = "../tests/resources/configurationUpdate/e2NodeConfigurationUpdateOnlyAddition.xml"
	e2NodeConfigurationUpdateXmlPath             = "../tests/resources/configurationUpdate/e2NodeConfigurationUpdate.xml"
)

func getTestE2NodeConfigurationUpdateMessage(t *testing.T, reqXmlPath string) *models.E2nodeConfigurationUpdateMessage {
	xmlConfUpdate := utils.ReadXmlFile(t, reqXmlPath)
	confUpdateMsg := &models.E2nodeConfigurationUpdateMessage{}
	err := xml.Unmarshal(utils.NormalizeXml(xmlConfUpdate), &confUpdateMsg.E2APPDU)
	assert.Nil(t, err)
	return confUpdateMsg
}

func TestParseE2NodeConfigurationUpdateSuccessAdditionOnly(t *testing.T) {
	configurationUpdate := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateOnlyAdditionXmlPath)
	assert.NotEqual(t, nil, configurationUpdate, "xml is not parsed correctly")
	assert.Equal(t, "6", configurationUpdate.E2APPDU.InitiatingMessage.ProcedureCode)
	assert.Equal(t, 1, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs))
	assert.Equal(t, 1, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer))
	assert.Equal(t, false, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)
}

func TestParseE2NodeConfigurationUpdateSuccess(t *testing.T) {
	configurationUpdate := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	assert.NotEqual(t, nil, configurationUpdate, "xml is not parsed correctly")
	assert.Equal(t, "6", configurationUpdate.E2APPDU.InitiatingMessage.ProcedureCode)
	assert.Equal(t, 3, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs))

	assert.Equal(t, 7, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer))
	assert.Equal(t, false, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[0].Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)

	updateIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[1]
	assert.Equal(t, 7, len(updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer))
	assert.Equal(t, false, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.E1 == nil)

	removalIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[2]
	assert.Equal(t, 7, len(removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer))
	assert.Equal(t, false, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.E1 == nil)
}
