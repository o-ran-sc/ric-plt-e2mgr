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

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/stateMachine"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"time"

	"e2mgr/models"
)

type DeleteAllRequestHandler struct {
	readerProvider func() reader.RNibReader
	writerProvider func() rNibWriter.RNibWriter
	config *configuration.Configuration
}

func NewDeleteAllRequestHandler(config *configuration.Configuration, writerProvider func() rNibWriter.RNibWriter,
	readerProvider func() reader.RNibReader) *DeleteAllRequestHandler {
	return &DeleteAllRequestHandler {
		readerProvider: readerProvider,
		writerProvider: writerProvider,
		config: config,
	}
}

func (handler *DeleteAllRequestHandler) Handle(logger *logger.Logger, request models.Request, rmrResponseChannel chan<- *models.NotificationResponse) error {

	err, continueFlow := handler.updateNodebStates(logger, false)
	if err != nil {
		return err
	}

	if continueFlow == false{
		return nil
	}

	//TODO change to rmr_request
	response := models.NotificationResponse{MgsType: rmrCgo.RIC_SCTP_CLEAR_ALL}
	rmrResponseChannel <- &response

	time.Sleep(time.Duration(handler.config.BigRedButtonTimeoutSec) * time.Second)
	logger.Infof("#delete_all_request_handler.Handle - timer expired")

	err, _ = handler.updateNodebStates(logger, true)
	if err != nil {
		return err
	}

	return nil
}

func (handler *DeleteAllRequestHandler) updateNodebStates(logger *logger.Logger, timeoutExpired bool) (error, bool){
	nbIdentityList, err := handler.readerProvider().GetListNodebIds()

	if err != nil {
		logger.Errorf("#delete_all_request_handler.updateNodebStates - failed to get nodes list from RNIB. Error: %s", err.Error())
		return e2managererrors.NewRnibDbError(), false
	}

	if len(nbIdentityList) == 0 {
		return nil, false
	}

	for _,nbIdentity := range nbIdentityList{

		node, err := handler.readerProvider().GetNodeb((*nbIdentity).GetInventoryName())

		if err != nil {
			logger.Errorf("#delete_all_request_handler.updateNodebStates - failed to get nodeB entity for ran name: %v from RNIB. Error: %s",
				(*nbIdentity).GetInventoryName(), err.Error())
			continue
		}

		if timeoutExpired{

			handler.saveNodebShutDownState(logger, nbIdentity, node)
			continue
		}
		handler.saveNodebNextState(logger, nbIdentity, node)
	}

	logger.Infof("#delete_all_request_handler.updateNodebStates - update nodeb states in RNIB completed")
	return nil, true
}

func (handler *DeleteAllRequestHandler) saveNodebNextState(logger *logger.Logger, nbIdentity *entities.NbIdentity, node *entities.NodebInfo) {

	nextStatus, res := stateMachine.NodeNextStateDeleteAll(node.ConnectionStatus)
	if res == false {
		return
	}

	node.ConnectionStatus = nextStatus

	err := handler.writerProvider().SaveNodeb(nbIdentity, node)

	if err != nil {
		logger.Errorf("#delete_all_request_handler.saveNodebNextState - failed to save nodeB entity for inventory name: %v to RNIB. Error: %s",
			(*nbIdentity).GetInventoryName(), err.Error())
		return
	}

	if logger.DebugEnabled() {
		logger.Debugf("#delete_all_request_handler.saveNodebNextState - connection status of inventory name: %v changed to %v",
			(*nbIdentity).GetInventoryName(), nextStatus.String())
	}
}

func (handler *DeleteAllRequestHandler) saveNodebShutDownState(logger *logger.Logger, nbIdentity *entities.NbIdentity, node *entities.NodebInfo) {

	if node.ConnectionStatus == entities.ConnectionStatus_SHUT_DOWN{
		return
	}

	if node.ConnectionStatus != entities.ConnectionStatus_SHUTTING_DOWN {
		logger.Errorf("#delete_all_request_handler.saveNodebShutDownState - ignore, status is not Shutting Down, inventory name: %v ", (*nbIdentity).GetInventoryName())
		return
	}

	node.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN

	err := handler.writerProvider().SaveNodeb(nbIdentity, node)

	if err != nil {
		logger.Errorf("#delete_all_request_handler.saveNodebShutDownState - failed to save nodeB entity for inventory name: %v to RNIB. Error: %s",
			(*nbIdentity).GetInventoryName(), err.Error())
		return
	}

	logger.Errorf("#delete_all_request_handler.saveNodebShutDownState - Shut Down , inventory name: %v ", (*nbIdentity).GetInventoryName())
}
