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

package handlers

import "C"
import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/stateMachine"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"time"
)

type DeleteAllRequestHandler struct {
	readerProvider func() reader.RNibReader
	writerProvider func() rNibWriter.RNibWriter
	rmrService *services.RmrService
	config *configuration.Configuration
}

func NewDeleteAllRequestHandler(rmrService *services.RmrService, config *configuration.Configuration, writerProvider func() rNibWriter.RNibWriter,
	readerProvider func() reader.RNibReader) *DeleteAllRequestHandler {
	return &DeleteAllRequestHandler {
		readerProvider: readerProvider,
		writerProvider: writerProvider,
		rmrService: rmrService,
		config: config,
	}
}

func (handler *DeleteAllRequestHandler) Handle(logger *logger.Logger, request models.Request) error {

	err, continueFlow := handler.updateNodebStates(logger, false)
	if err != nil {
		return err
	}

	if continueFlow == false{
		return nil
	}

	//TODO change to rmr_request
	response := models.NotificationResponse{MgsType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	if err:= handler.rmrService.SendRmrMessage(&response); err != nil {
		logger.Errorf("#DeleteAllRequestHandler.Handle - failed to send sctp clear all message to RMR: %s", err)
		return  e2managererrors.NewRmrError()
	}

	time.Sleep(time.Duration(handler.config.BigRedButtonTimeoutSec) * time.Second)
	logger.Infof("#DeleteAllRequestHandler.Handle - timer expired")

	err, _ = handler.updateNodebStates(logger, true)
	if err != nil {
		return err
	}

	return nil
}

func (handler *DeleteAllRequestHandler) updateNodebStates(logger *logger.Logger, timeoutExpired bool) (error, bool){
	nbIdentityList, err := handler.readerProvider().GetListNodebIds()

	if err != nil {
		logger.Errorf("#DeleteAllRequestHandler.updateNodebStates - failed to get nodes list from RNIB. Error: %s", err.Error())
		return e2managererrors.NewRnibDbError(), false
	}

	if len(nbIdentityList) == 0 {
		return nil, false
	}

	numOfRanToShutDown := 0
	for _,nbIdentity := range nbIdentityList{

		node, err := handler.readerProvider().GetNodeb((*nbIdentity).GetInventoryName())

		if err != nil {
			logger.Errorf("#DeleteAllRequestHandler.updateNodebStates - failed to get nodeB entity for ran name: %v from RNIB. Error: %s",
				(*nbIdentity).GetInventoryName(), err.Error())
			continue
		}

		if timeoutExpired{

			if handler.saveNodebShutDownState(logger, nbIdentity, node){
				numOfRanToShutDown++
			}
			continue
		}
		if handler.saveNodebNextState(logger, nbIdentity, node){
			numOfRanToShutDown++
		}
	}

	if numOfRanToShutDown > 0{
		logger.Infof("#DeleteAllRequestHandler.updateNodebStates - update nodebs states in RNIB completed")
	}else {
		logger.Infof("#DeleteAllRequestHandler.updateNodebStates - nodebs states are not updated ")
		return nil, false
	}

	return nil, true
}

func (handler *DeleteAllRequestHandler) saveNodebNextState(logger *logger.Logger, nbIdentity *entities.NbIdentity, node *entities.NodebInfo) bool{

	if node.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN{
		return true
	}

	nextStatus, res := stateMachine.NodeNextStateDeleteAll(node.ConnectionStatus)
	if res == false {
		return false
	}

	node.ConnectionStatus = nextStatus

	err := handler.writerProvider().SaveNodeb(nbIdentity, node)

	if err != nil {
		logger.Errorf("#DeleteAllRequestHandler.saveNodebNextState - failed to save nodeB entity for inventory name: %v to RNIB. Error: %s",
			(*nbIdentity).GetInventoryName(), err.Error())
		return false
	}

	if logger.DebugEnabled() {
		logger.Debugf("#DeleteAllRequestHandler.saveNodebNextState - connection status of inventory name: %v changed to %v",
			(*nbIdentity).GetInventoryName(), nextStatus.String())
	}
	return true
}

func (handler *DeleteAllRequestHandler) saveNodebShutDownState(logger *logger.Logger, nbIdentity *entities.NbIdentity, node *entities.NodebInfo) bool{

	if node.ConnectionStatus == entities.ConnectionStatus_SHUT_DOWN{
		return false
	}

	if node.ConnectionStatus != entities.ConnectionStatus_SHUTTING_DOWN {
		logger.Errorf("#DeleteAllRequestHandler.saveNodebShutDownState - ignore, status is not Shutting Down, inventory name: %v ", (*nbIdentity).GetInventoryName())
		return false
	}

	node.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN

	err := handler.writerProvider().SaveNodeb(nbIdentity, node)

	if err != nil {
		logger.Errorf("#DeleteAllRequestHandler.saveNodebShutDownState - failed to save nodeB entity for inventory name: %v to RNIB. Error: %s",
			(*nbIdentity).GetInventoryName(), err.Error())
		return false
	}

	logger.Errorf("#DeleteAllRequestHandler.saveNodebShutDownState - Shut Down , inventory name: %v ", (*nbIdentity).GetInventoryName())
	return true
}
