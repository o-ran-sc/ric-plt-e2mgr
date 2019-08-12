package handlers

import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/sessions"
)

type RicEnbLoadInformationNotificationHandler struct{}


func (src RicEnbLoadInformationNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions,
	request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	notification, err := unpackX2apPduAndRefine(logger, MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize /*message buffer*/)

	if err != nil {
		logger.Errorf("#ric_enb_load_information_notification_handler.Handle - unpack failed. Error: %v", err)
	}

	logger.Infof("#ric_enb_load_information_notification_handler.handle - Enb load information notification message received")
	logger.Debugf("#ric_enb_load_information_notification_handler.handle - Enb load information notification message payload: %s", notification.pduPrint)
}
