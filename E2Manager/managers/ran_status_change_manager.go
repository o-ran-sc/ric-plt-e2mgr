package managers

import (
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services/rmrsender"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type RanStatusChangeManager struct {
	logger    *logger.Logger
	rmrSender *rmrsender.RmrSender
}

func NewRanStatusChangeManager(logger *logger.Logger, rmrSender *rmrsender.RmrSender) *RanStatusChangeManager {
	return &RanStatusChangeManager{
		logger:    logger,
		rmrSender: rmrSender,
	}
}

type IRanStatusChangeManager interface {
	Execute(msgType int, msgDirection enums.MessageDirection, nodebInfo *entities.NodebInfo) error
}

func (m *RanStatusChangeManager) Execute(msgType int, msgDirection enums.MessageDirection, nodebInfo *entities.NodebInfo) error {

	resourceStatusPayload := models.NewResourceStatusPayload(nodebInfo.NodeType, msgDirection)
	resourceStatusJson, err := json.Marshal(resourceStatusPayload)

	if err != nil {
		m.logger.Errorf("#RanStatusChangeManager.Execute - RAN name: %s - Error marshaling resource status payload: %v", nodebInfo.RanName, err)
		return err
	}

	rmrMessage := models.NewRmrMessage(msgType, nodebInfo.RanName, resourceStatusJson)
	return m.rmrSender.Send(rmrMessage)
}
