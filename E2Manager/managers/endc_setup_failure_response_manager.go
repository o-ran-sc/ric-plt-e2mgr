package managers

import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type EndcSetupFailureResponseManager struct{}

func NewEndcSetupFailureResponseManager() *EndcSetupFailureResponseManager {
	return &EndcSetupFailureResponseManager{}
}

func (m *EndcSetupFailureResponseManager) PopulateNodebByPdu(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	failureResponse, err := converters.UnpackEndcX2SetupFailureResponseAndExtract(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		logger.Errorf("#EndcSetupFailureResponseManager.PopulateNodebByPdu - RAN name: %s - Unpack and extract failed. Error: %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#EndcSetupFailureResponseManager.PopulateNodebByPdu - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED_SETUP_FAILED
	nodebInfo.SetupFailure = failureResponse
	nodebInfo.FailureType = entities.Failure_ENDC_X2_SETUP_FAILURE
	return nil
}
