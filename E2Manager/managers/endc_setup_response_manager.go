package managers

import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type EndcSetupResponseManager struct{}

func NewEndcSetupResponseManager() *EndcSetupResponseManager {
	return &EndcSetupResponseManager{}
}

func (m *EndcSetupResponseManager) SetNodeb(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	gnbId, gnb, err := converters.UnpackEndcX2SetupResponseAndExtract(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		logger.Errorf("#EndcSetupResponseManager.SetNodeb - RAN name: %s - Unpack & extract failed. Error: %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#EndcSetupResponseManager.SetNodeb - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nbIdentity.GlobalNbId = gnbId
	nodebInfo.GlobalNbId = gnbId
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	nodebInfo.NodeType = entities.Node_GNB
	nodebInfo.Configuration = &entities.NodebInfo_Gnb{Gnb: gnb}

	return nil
}
