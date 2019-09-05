package rmrmsghandlers

import "C"
import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"time"
)

type EnbLoadInformationNotificationHandler struct {
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewEnbLoadInformationNotificationHandler(rnibWriterProvider func() rNibWriter.RNibWriter) EnbLoadInformationNotificationHandler {
	return EnbLoadInformationNotificationHandler{
		rnibWriterProvider: rnibWriterProvider,
	}
}

func elapsed(startTime time.Time) float64 {
	return float64(time.Since(startTime)) / float64(time.Millisecond)
}

func (src EnbLoadInformationNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions, request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	pdu, err := converters.UnpackX2apPdu(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpack failed. Error: %v", request.RanName, err)
		return
	}

	logger.Debugf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpacked message successfully", request.RanName)

	ranLoadInformation := &entities.RanLoadInformation{LoadTimestamp: uint64(request.StartTime.UnixNano())}

	err = converters.ExtractAndBuildRanLoadInformation(pdu, ranLoadInformation)

	if (err != nil) {
		logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed at ExtractAndBuildRanLoadInformation. Error: %v", request.RanName, err)
		return
	}

	logger.Debugf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully done with extracting and building RAN load information. elapsed: %f ms", request.RanName, elapsed(request.StartTime))

	rnibErr := src.rnibWriterProvider().SaveRanLoadInformation(request.RanName, ranLoadInformation)

	if rnibErr != nil {
		logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed saving RAN load information. Error: %v", request.RanName, rnibErr)
		return
	}

	logger.Infof("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully saved RAN load information to RNIB. elapsed: %f ms", request.RanName, elapsed(request.StartTime))
}
