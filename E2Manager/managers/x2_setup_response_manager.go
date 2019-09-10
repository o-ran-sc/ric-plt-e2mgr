package managers

import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2SetupResponseManager struct{}

func NewX2SetupResponseManager() *X2SetupResponseManager {
	return &X2SetupResponseManager{}
}

func (m *X2SetupResponseManager) PopulateNodebByPdu(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	enbId, enb, err := converters.UnpackX2SetupResponseAndExtract(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		logger.Errorf("#X2SetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpack and extract failed. %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#X2SetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nbIdentity.GlobalNbId = enbId
	nodebInfo.GlobalNbId = enbId
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	nodebInfo.NodeType = entities.Node_ENB
	nodebInfo.Configuration = &entities.NodebInfo_Enb{Enb: enb}

	return nil
}
