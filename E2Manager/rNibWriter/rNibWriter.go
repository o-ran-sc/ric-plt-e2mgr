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

package rNibWriter

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"github.com/golang/protobuf/proto"
)



var writerPool *common.Pool

type rNibWriterInstance struct {
	sdl       *common.ISdlInstance
	namespace string
}
/*
RNibWriter interface allows saving data to the redis DB
 */
type RNibWriter interface {
	SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) common.IRNibError
	SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) common.IRNibError
}
/*
Init initializes the infrastructure required for the RNibWriter instance
 */
func Init(namespace string, poolSize int) {
	initPool(poolSize,
		func() interface{} {
			var sdlI common.ISdlInstance = sdlgo.NewSdlInstance(namespace, sdlgo.NewDatabase())
			return &rNibWriterInstance{sdl: &sdlI, namespace: namespace}
		},
		func(obj interface{}) {
			(*obj.(*rNibWriterInstance).sdl).Close()
		})
}
/*
InitPool initializes the writer's instances pool
 */
func initPool(poolSize int, newObj func() interface{}, destroyObj func(interface{})) {
	writerPool = common.NewPool(poolSize, newObj, destroyObj)
}
/*
GetRNibWriter returns RNibWriter instance from the pool
 */
func GetRNibWriter() RNibWriter {
	return writerPool.Get().(RNibWriter)
}
/*
SaveNodeb saves nodeB entity data in the redis DB according to the specified data model
 */
func (w *rNibWriterInstance) SaveNodeb(nbIdentity *entities.NbIdentity, entity *entities.NodebInfo) common.IRNibError {

	isNotEmptyIdentity := isNotEmpty(nbIdentity)

	if isNotEmptyIdentity && entity.GetNodeType() == entities.Node_UNKNOWN{
		return common.NewValidationError(errors.New( fmt.Sprintf("#rNibWriter.saveNodeB - Unknown responding node type, entity: %v", entity)))
	}
	defer writerPool.Put(w)
	data, err := proto.Marshal(entity)
	if err != nil {
		return common.NewInternalError(err)
	}
	var pairs []interface{}
	key, rNibErr := common.ValidateAndBuildNodeBNameKey(nbIdentity.InventoryName)
	if rNibErr != nil{
		return rNibErr
	}
	pairs = append(pairs, key, data)

	if isNotEmptyIdentity {
		key, rNibErr = common.ValidateAndBuildNodeBIdKey(entity.GetNodeType().String(), nbIdentity.GlobalNbId.GetPlmnId(), nbIdentity.GlobalNbId.GetNbId())
		if rNibErr != nil{
			return rNibErr
		}
		pairs = append(pairs, key, data)
	}

	if entity.GetEnb() != nil {
		pairs, rNibErr = appendEnbCells(nbIdentity, entity.GetEnb().GetServedCells(), pairs)
		if rNibErr != nil{
			return rNibErr
		}
	}
	if entity.GetGnb() != nil {
		pairs, rNibErr = appendGnbCells(nbIdentity, entity.GetGnb().GetServedNrCells(), pairs)
		if rNibErr != nil{
			return rNibErr
		}
	}
	err = (*w.sdl).Set(pairs)
	if err != nil {
		return common.NewInternalError(err)
	}

	ranNameIdentity := &entities.NbIdentity{InventoryName: nbIdentity.InventoryName}

	if isNotEmptyIdentity{
		nbIdData, err := proto.Marshal(ranNameIdentity)
		if err != nil {
			return common.NewInternalError(err)
		}
		err = (*w.sdl).RemoveMember(entities.Node_UNKNOWN.String(), nbIdData)
		if err != nil {
			return common.NewInternalError(err)
		}
	} else {
		nbIdentity = ranNameIdentity
	}

	nbIdData, err := proto.Marshal(nbIdentity)
	if err != nil {
		return common.NewInternalError(err)
	}
	err = (*w.sdl).AddMember(entity.GetNodeType().String(), nbIdData)
	if err != nil {
		return common.NewInternalError(err)
	}
	return nil
}

/*
SaveRanLoadInformation stores ran load information for the provided ran
*/
func (w *rNibWriterInstance) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) common.IRNibError {

	defer writerPool.Put(w)

	key, rnibErr:= common.ValidateAndBuildRanLoadInformationKey(inventoryName)

	if rnibErr != nil {
		return rnibErr
	}

	data, err := proto.Marshal(ranLoadInformation)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = (*w.sdl).Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}

/*
Close closes writer's pool
 */
func Close(){
	writerPool.Close()
}

func appendEnbCells(nbIdentity *entities.NbIdentity, cells []*entities.ServedCellInfo, pairs []interface{}) ([]interface{}, common.IRNibError) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type:entities.Cell_LTE_CELL, Cell:&entities.Cell_ServedCellInfo{ServedCellInfo:cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildCellIdKey(cell.GetCellId())
		if rNibErr != nil{
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(nbIdentity.InventoryName, cell.GetPci())
		if rNibErr != nil{
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func appendGnbCells(nbIdentity *entities.NbIdentity, cells  []*entities.ServedNRCell, pairs []interface{}) ([]interface{}, common.IRNibError) {
	for _, cell := range cells {
		cellEntity := entities.Cell{Type:entities.Cell_NR_CELL, Cell:&entities.Cell_ServedNrCell{ServedNrCell:cell}}
		cellData, err := proto.Marshal(&cellEntity)
		if err != nil {
			return pairs, common.NewInternalError(err)
		}
		key, rNibErr := common.ValidateAndBuildNrCellIdKey(cell.GetServedNrCellInformation().GetCellId())
		if rNibErr != nil{
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
		key, rNibErr = common.ValidateAndBuildCellNamePciKey(nbIdentity.InventoryName, cell.GetServedNrCellInformation().GetNrPci())
		if rNibErr != nil{
			return pairs, rNibErr
		}
		pairs = append(pairs, key, cellData)
	}
	return pairs, nil
}

func isNotEmpty(nbIdentity *entities.NbIdentity) bool {
	return nbIdentity.GlobalNbId != nil && nbIdentity.GlobalNbId.PlmnId != "" && nbIdentity.GlobalNbId.NbId != ""
}
