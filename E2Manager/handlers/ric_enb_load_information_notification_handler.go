package handlers

import "C"
import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rNibWriter"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"time"
)

type RicEnbLoadInformationNotificationHandler struct{}

func elapsed(startTime time.Time) float64 {
	return float64(time.Since(startTime)) / float64(time.Millisecond)
}

func (src RicEnbLoadInformationNotificationHandler) Handle(logger *logger.Logger, e2Sessions sessions.E2Sessions, request *models.NotificationRequest, messageChannel chan<- *models.NotificationResponse) {

	pdu, err := unpackX2apPdu(logger, MaxAsn1CodecAllocationBufferSize, request.Len, request.Payload, MaxAsn1CodecMessageBufferSize)

	//packedExampleString := "004c07080004001980da0100075bde017c148003d5a8205000017c180003d5a875555403331420000012883a0003547400cd20002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be0c4001ead4016e007ab50100002f8320067ab5005b8c1ead5070190c00001d637805f220000f56a081400005f020000f56a1d555400ccc508002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be044001ead4016e007ab50100002f8120067ab5005b8c1ead5070190c00000"
	//
	//var packedExampleByteSlice []byte
	//
	//_, err := fmt.Sscanf(packedExampleString, "%x", &packedExampleByteSlice)
	//
	//pdu, err := unpackX2apPduUPer(logger, MaxAsn1CodecAllocationBufferSize, len(packedExampleByteSlice), packedExampleByteSlice, MaxAsn1CodecMessageBufferSize)


	if err != nil {
		logger.Errorf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpack failed. Error: %v", request.RanName, err)
		return
	}

	logger.Debugf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpacked message successfully", request.RanName)

	ranLoadInformation := &entities.RanLoadInformation{LoadTimestamp: uint64(request.StartTime.UnixNano())}

	err = extractAndBuildRanLoadInformation(pdu, ranLoadInformation)

	if (err != nil) {
		logger.Errorf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed at extractAndBuildRanLoadInformation. Error: %v", request.RanName, err)
		return
	}

	logger.Debugf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully done with extracting and building RAN load information. elapsed: %f ms", request.RanName, elapsed(request.StartTime))

	rnibErr := rNibWriter.GetRNibWriter().SaveRanLoadInformation(request.RanName, ranLoadInformation) // TODO: Should inject RnibWriter

	if rnibErr != nil {
		logger.Errorf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed saving RAN load information. Error: %v", request.RanName, rnibErr)
		return
	}

	logger.Debugf("#RicEnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully saved RAN load information to RNIB. elapsed: %f ms", request.RanName, elapsed(request.StartTime))
}
