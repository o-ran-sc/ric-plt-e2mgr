package models

import (
	"e2mgr/enums"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type ResourceStatusPayload struct {
	NodeType         entities.Node_Type     `json:"nodeType"`
	MessageDirection enums.MessageDirection `json:"messageDirection"`
}

func NewResourceStatusPayload(nodeType entities.Node_Type, messageDirection enums.MessageDirection) *ResourceStatusPayload {
	return &ResourceStatusPayload{
		NodeType:         nodeType,
		MessageDirection: messageDirection,
	}
}
