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

package rNibWriter

import (
	"e2mgr/configuration"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/proto"
)

const (
	E2TAddressesKey = "E2TAddresses"
	RanAddedEvent   = "ADDED"
	RanUpdatedEvent = "UPDATED"
	RanDeletedEvent = "DELETED"
)

type rNibWriterInstance struct {
	sdl              common.ISdlSyncStorage
	rnibWriterConfig configuration.RnibWriterConfig
	ns               string
}

/*
RNibWriter interface allows saving data to the redis DB
*/
type RNibWriter interface {
	SaveNodeb(nodebInfo *entities.NodebInfo) error
	UpdateNodebInfo(nodebInfo *entities.NodebInfo) error
	UpdateNodebInfoAndPublish(nodebInfo *entities.NodebInfo) error
	SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error
	SaveE2TInstance(e2tInstance *entities.E2TInstance) error
	SaveE2TAddresses(addresses []string) error
	RemoveE2TInstance(e2tAddress string) error
	UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error
	RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error
	UpdateNodebInfoOnConnectionStatusInversion(nodebInfo *entities.NodebInfo, ent string) error
	SaveGeneralConfiguration(config *entities.GeneralConfiguration) error
	RemoveEnb(nodebInfo *entities.NodebInfo) error
	RemoveServedCells(inventoryName string, servedCells []*entities.ServedCellInfo) error
	UpdateEnb(nodebInfo *entities.NodebInfo, servedCells []*entities.ServedCellInfo) error
	AddNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error
	RemoveNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error
	AddEnb(nodebInfo *entities.NodebInfo) error
	UpdateNbIdentities(nodeType entities.Node_Type, oldNbIdentities []*entities.NbIdentity, newNbIdentities []*entities.NbIdentity) error
}

/*
GetRNibWriter returns reference to RNibWriter
*/

func GetRNibWriter(sdl common.ISdlSyncStorage, rnibWriterConfig configuration.RnibWriterConfig) RNibWriter {
	return &rNibWriterInstance{
		sdl:              sdl,
		rnibWriterConfig: rnibWriterConfig,
		ns:               common.GetRNibNamespace(),
	}
}

func getChannelsAndEventsPair(channel string, ranName string, event string) []string {
	return []string{channel, fmt.Sprintf("%s_%s", ranName, event)}
}

func (w *rNibWriterInstance) AddNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error {
	nbIdData, err := proto.Marshal(nbIdentity)

	if err != nil {
		return common.NewInternalError(err)
	}

	err = w.sdl.AddMember(w.ns, nodeType.String(), nbIdData)

	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

func (w *rNibWriterInstance) RemoveServedNrCells(inventoryName string, servedNrCells []*entities.ServedNRCell) error {
	cellKeysToRemove := buildServedNRCellKeysToRemove(inventoryName, servedNrCells)

	err := w.sdl.Remove(w.ns, cellKeysToRemove)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) RemoveServedCells(inventoryName string, servedCells []*entities.ServedCellInfo) error {
	cellKeysToRemove := buildServedCellInfoKeysToRemove(inventoryName, servedCells)

	err := w.sdl.Remove(w.ns, cellKeysToRemove)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) SaveGeneralConfiguration(config *entities.GeneralConfiguration) error {

	err := w.SaveWithKeyAndMarshal(common.BuildGeneralConfigurationKey(), config)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
SaveNodeb saves nodeB entity data in the redis DB according to the specified data model
*/
func (w *rNibWriterInstance) SaveNodeb(nodebInfo *entities.NodebInfo) error {

	data, err := proto.Marshal(nodebInfo)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	key, rNibErr := common.ValidateAndBuildNodeBNameKey(nodebInfo.RanName)

	if rNibErr != nil {
		return rNibErr
	}

	pairs = append(pairs, key, data)

	if nodebInfo.GlobalNbId != nil {

		key, rNibErr = common.ValidateAndBuildNodeBIdKey(nodebInfo.GetNodeType().String(), nodebInfo.GlobalNbId.GetPlmnId(), nodebInfo.GlobalNbId.GetNbId())
		if rNibErr != nil {
			return rNibErr
		}
		pairs = append(pairs, key, data)
	}

	if nodebInfo.GetEnb() != nil {
		pairs, rNibErr = appendEnbCells(nodebInfo.RanName, nodebInfo.GetEnb().GetServedCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}

	if nodebInfo.GetGnb() != nil {
		pairs, rNibErr = appendGnbCells(nodebInfo.RanName, nodebInfo.GetGnb().GetServedNrCells(), pairs)
		if rNibErr != nil {
			return rNibErr
		}
	}

	err = w.sdl.Set(w.ns, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) AddEnb(nodebInfo *entities.NodebInfo) error {

	data, err := proto.Marshal(nodebInfo)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	key, rNibErr := common.ValidateAndBuildNodeBNameKey(nodebInfo.RanName)

	if rNibErr != nil {
		return rNibErr
	}

	pairs = append(pairs, key, data)

	if nodebInfo.GlobalNbId != nil {

		key, rNibErr = common.ValidateAndBuildNodeBIdKey(nodebInfo.GetNodeType().String(), nodebInfo.GlobalNbId.GetPlmnId(), nodebInfo.GlobalNbId.GetNbId())
		if rNibErr != nil {
			return rNibErr
		}
		pairs = append(pairs, key, data)
	}

	pairs, rNibErr = appendEnbCells(nodebInfo.RanName, nodebInfo.GetEnb().GetServedCells(), pairs)
	if rNibErr != nil {
		return rNibErr
	}

	channelsAndEvents := getChannelsAndEventsPair(w.rnibWriterConfig.RanManipulationMessageChannel, nodebInfo.RanName, RanAddedEvent)
	err = w.sdl.SetAndPublish(w.ns, channelsAndEvents, pairs)
	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) UpdateNbIdentities(nodeType entities.Node_Type, oldNbIdentities []*entities.NbIdentity, newNbIdentities []*entities.NbIdentity) error {

	nbIdIdentitiesToRemove, err := w.buildNbIdentitiesMembers(oldNbIdentities)
	if err != nil {
		return err
	}

	err = w.sdl.RemoveMember(w.ns, nodeType.String(), nbIdIdentitiesToRemove[:]...)
	if err != nil {
		return err
	}

	nbIdIdentitiesToAdd, err := w.buildNbIdentitiesMembers(newNbIdentities)
	if err != nil {
		return err
	}

	err = w.sdl.AddMember(w.ns, nodeType.String(), nbIdIdentitiesToAdd[:]...)
	if err != nil {
		return err
	}

	return nil
}

func (w *rNibWriterInstance) UpdateGnbCells(nodebInfo *entities.NodebInfo, servedNrCells []*entities.ServedNRCell) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	pairs, err = appendGnbCells(nodebInfo.RanName, servedNrCells, pairs)

	if err != nil {
		return err
	}

	channelsAndEvents := getChannelsAndEventsPair(w.rnibWriterConfig.RanManipulationMessageChannel, nodebInfo.RanName, RanUpdatedEvent)
	err = w.sdl.SetAndPublish(w.ns, channelsAndEvents, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func buildServedNRCellKeysToRemove(inventoryName string, servedNrCellsToRemove []*entities.ServedNRCell) []string {

	var cellKeysToRemove []string

	for _, cell := range servedNrCellsToRemove {

		key, _ := common.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}

		key, _ = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetServedNrCellInformation().GetNrPci())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}
	}

	return cellKeysToRemove
}

func buildServedCellInfoKeysToRemove(inventoryName string, servedCellsToRemove []*entities.ServedCellInfo) []string {

	var cellKeysToRemove []string

	for _, cell := range servedCellsToRemove {

		key, _ := common.ValidateAndBuildCellIdKey(cell.GetCellId())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}

		key, _ = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetPci())

		if len(key) != 0 {
			cellKeysToRemove = append(cellKeysToRemove, key)
		}
	}

	return cellKeysToRemove
}

func buildUpdateNodebInfoPairs(nodebInfo *entities.NodebInfo) ([]interface{}, error) {
	nodebNameKey, rNibErr := common.ValidateAndBuildNodeBNameKey(nodebInfo.GetRanName())

	if rNibErr != nil {
		return []interface{}{}, rNibErr
	}

	nodebIdKey, buildNodebIdKeyError := common.ValidateAndBuildNodeBIdKey(nodebInfo.GetNodeType().String(), nodebInfo.GlobalNbId.GetPlmnId(), nodebInfo.GlobalNbId.GetNbId())

	data, err := proto.Marshal(nodebInfo)

	if err != nil {
		return []interface{}{}, common.NewInternalError(err)
	}

	pairs := []interface{}{nodebNameKey, data}

	if buildNodebIdKeyError == nil {
		pairs = append(pairs, nodebIdKey, data)
	}

	return pairs, nil
}

func (w *rNibWriterInstance) buildRemoveEnbKeys(nodebInfo *entities.NodebInfo) ([]string, error) {
	keys := buildServedCellInfoKeysToRemove(nodebInfo.GetRanName(), nodebInfo.GetEnb().GetServedCells())

	nodebNameKey, rNibErr := common.ValidateAndBuildNodeBNameKey(nodebInfo.GetRanName())

	if rNibErr != nil {
		return []string{}, rNibErr
	}

	keys = append(keys, nodebNameKey)

	nodebIdKey, buildNodebIdKeyError := common.ValidateAndBuildNodeBIdKey(nodebInfo.GetNodeType().String(), nodebInfo.GlobalNbId.GetPlmnId(), nodebInfo.GlobalNbId.GetNbId())

	if buildNodebIdKeyError == nil {
		keys = append(keys, nodebIdKey)
	}

	return keys, nil
}

func (w *rNibWriterInstance) RemoveNbIdentity(nodeType entities.Node_Type, nbIdentity *entities.NbIdentity) error {
	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		return common.NewInternalError(err)
	}
	err = w.sdl.RemoveMember(w.ns, nodeType.String(), nbIdData)
	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

func (w *rNibWriterInstance) updateNodebInfo(nodebInfo *entities.NodebInfo, publish bool) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	if publish {
		channelsAndEvents := getChannelsAndEventsPair(w.rnibWriterConfig.RanManipulationMessageChannel, nodebInfo.RanName, RanUpdatedEvent)
		err = w.sdl.SetAndPublish(w.ns, channelsAndEvents, pairs)
	} else {
		err = w.sdl.Set(w.ns, pairs)
	}

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
UpdateNodebInfo...
*/
func (w *rNibWriterInstance) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {
	return w.updateNodebInfo(nodebInfo, false)
}

/*
UpdateNodebInfoAndPublish...
*/
func (w *rNibWriterInstance) UpdateNodebInfoAndPublish(nodebInfo *entities.NodebInfo) error {
	return w.updateNodebInfo(nodebInfo, true)
}

/*
SaveRanLoadInformation stores ran load information for the provided ran
*/
func (w *rNibWriterInstance) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {

	key, rnibErr := common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if rnibErr != nil {
		return rnibErr
	}

	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = w.sdl.Set(w.ns, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) SaveE2TInstance(e2tInstance *entities.E2TInstance) error {

	key, rnibErr := common.ValidateAndBuildE2TInstanceKey(e2tInstance.Address)

	if rnibErr != nil {
		return rnibErr
	}

	data, err := json.Marshal(e2tInstance)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = w.sdl.Set(w.ns, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) SaveE2TAddresses(addresses []string) error {

	data, err := json.Marshal(addresses)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, E2TAddressesKey, data)

	err = w.sdl.Set(w.ns, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) RemoveE2TInstance(address string) error {
	key, rNibErr := common.ValidateAndBuildE2TInstanceKey(address)
	if rNibErr != nil {
		return rNibErr
	}
	err := w.sdl.Remove(w.ns, []string{key})

	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

func (w *rNibWriterInstance) SaveWithKeyAndMarshal(key string, entity interface{}) error {

	data, err := json.Marshal(entity)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = w.sdl.Set(w.ns, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
UpdateNodebInfoOnConnectionStatusInversion...
*/
func (w *rNibWriterInstance) UpdateNodebInfoOnConnectionStatusInversion(nodebInfo *entities.NodebInfo, event string) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	err = w.sdl.SetAndPublish(w.ns, []string{w.rnibWriterConfig.StateChangeMessageChannel, event}, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) RemoveEnb(nodebInfo *entities.NodebInfo) error {
	keysToRemove, err := w.buildRemoveEnbKeys(nodebInfo)
	if err != nil {
		return err
	}

	channelsAndEvents := getChannelsAndEventsPair(w.rnibWriterConfig.RanManipulationMessageChannel, nodebInfo.RanName, RanDeletedEvent)
	err = w.sdl.RemoveAndPublish(w.ns, channelsAndEvents, keysToRemove)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) UpdateEnb(nodebInfo *entities.NodebInfo, servedCells []*entities.ServedCellInfo) error {

	pairs, err := buildUpdateNodebInfoPairs(nodebInfo)

	if err != nil {
		return err
	}

	pairs, err = appendEnbCells(nodebInfo.RanName, servedCells, pairs)

	if err != nil {
		return err
	}

	channelsAndEvents := getChannelsAndEventsPair(w.rnibWriterConfig.RanManipulationMessageChannel, nodebInfo.RanName, RanUpdatedEvent)
	err = w.sdl.SetAndPublish(w.ns, channelsAndEvents, pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

func (w *rNibWriterInstance) buildNbIdentitiesMembers(nbIdentities []*entities.NbIdentity) ([]interface{}, error) {

	var nbIdIdentitiesMembers []interface{}
	for _, nbIdentity := range nbIdentities {

		nbIdData, err := proto.Marshal(nbIdentity)
		if err != nil {
			return nil, common.NewInternalError(err)
		}
		nbIdIdentitiesMembers = append(nbIdIdentitiesMembers, nbIdData)
	}

	return nbIdIdentitiesMembers, nil
}

/*
Close the writer
*/
func Close() {
	//Nothing to do
}

func appendEnbCells(inventoryName string, cells []*entities.ServedCellInfo, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type: entities.Cell_LTE_CELL, Cell: &entities.Cell_ServedCellInfo{ServedCellInfo: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildCellIdKey(cell.GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func appendGnbCells(inventoryName string, cells []*entities.ServedNRCell, pairs []interface{}) ([]interface{}, error) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type: entities.Cell_NR_CELL, Cell: &entities.Cell_ServedNrCell{ServedNrCell: cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(inventoryName, cell.GetServedNrCellInformation().GetNrPci())
		if rNibErr != nil {
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}
