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
	assert.Equal(t, 2, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs))

	additionIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[1]
	assert.Equal(t, 1, len(additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer))
	assert.Equal(t, false, additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)
}

func TestParseE2NodeConfigurationUpdateSuccess(t *testing.T) {
	configurationUpdate := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	assert.NotEqual(t, nil, configurationUpdate, "xml is not parsed correctly")
	assert.Equal(t, "6", configurationUpdate.E2APPDU.InitiatingMessage.ProcedureCode)
	assert.Equal(t, 4, len(configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs))

	additionIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[1]
	assert.Equal(t, 7, len(additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer))
	assert.Equal(t, false, additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, additionIE.Value.E2nodeComponentConfigAdditionList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigAdditionItem.E2nodeComponentInterfaceType.E1 == nil)

	updateIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[2]
	assert.Equal(t, 7, len(updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer))
	assert.Equal(t, false, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, updateIE.Value.E2nodeComponentConfigUpdateList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigUpdateItem.E2nodeComponentInterfaceType.E1 == nil)

	removalIE := configurationUpdate.E2APPDU.InitiatingMessage.Value.E2nodeConfigurationUpdate.ProtocolIEs.E2nodeConfigurationUpdateIEs[3]
	assert.Equal(t, 7, len(removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer))
	assert.Equal(t, false, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.Ng == nil)
	assert.Equal(t, true, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.E1 == nil)
	assert.Equal(t, true, removalIE.Value.E2nodeComponentConfigRemovalList.ProtocolIESingleContainer[0].Value.E2nodeComponentConfigRemovalItem.E2nodeComponentInterfaceType.E1 == nil)
}

func TestExtractAdditionConfigList(t *testing.T) {
	configurationUpdate1 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	additionList := configurationUpdate1.ExtractConfigAdditionList()

	assert.Equal(t, 5, len(additionList), "Addtion List is not matching")

	configurationUpdate2 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateOnlyAdditionXmlPath)
	additionList2 := configurationUpdate2.ExtractConfigAdditionList()

	assert.Equal(t, 1, len(additionList2), "Addtion List is not matching")
}

func TestExtractUpdateConfigList(t *testing.T) {
	configurationUpdate1 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	updateList1 := configurationUpdate1.ExtractConfigUpdateList()

	assert.Equal(t, 5, len(updateList1), "Update List is not matching")

	configurationUpdate2 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateOnlyAdditionXmlPath)
	updateList2 := configurationUpdate2.ExtractConfigUpdateList()

	assert.Equal(t, 0, len(updateList2), "Update List is not matching")
}

func TestExtractDeleteConfigList(t *testing.T) {
	configurationRemoval1 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateXmlPath)
	removalList1 := configurationRemoval1.ExtractConfigDeletionList()

	assert.Equal(t, 5, len(removalList1), "Removal List is not matching")

	configurationRemoval2 := getTestE2NodeConfigurationUpdateMessage(t, e2NodeConfigurationUpdateOnlyAdditionXmlPath)
	removalList2 := configurationRemoval2.ExtractConfigDeletionList()

	assert.Equal(t, 0, len(removalList2), "Removal List is not matching")
}
