package managers

import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2SetupFailureResponseManager struct{}

func NewX2SetupFailureResponseManager() *X2SetupFailureResponseManager {
	return &X2SetupFailureResponseManager{}
}

func (m *X2SetupFailureResponseManager) SetNodeb(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	failureResponse, err := converters.UnpackX2SetupFailureResponseAndExtract(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		logger.Errorf("#X2SetupFailureResponseManager.SetNodeb - RAN name: %s - Unpack & extract failed. Error: %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#X2SetupFailureResponseManager.SetNodeb - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED_SETUP_FAILED
	nodebInfo.SetupFailure = failureResponse
	nodebInfo.FailureType = entities.Failure_X2_SETUP_FAILURE
	return nil
}
