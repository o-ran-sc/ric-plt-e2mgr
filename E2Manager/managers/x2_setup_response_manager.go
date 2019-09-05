package managers

import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2SetupResponseManager struct {}

func NewX2SetupResponseManager() *X2SetupResponseManager {
	return &X2SetupResponseManager{}
}

func (m *X2SetupResponseManager) SetNodeb(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {
		enbId, enb, err := converters.UnpackX2SetupResponseAndExtract(logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize)

		if err != nil || enbId == nil || enb == nil {
			logger.Errorf("#X2SetupResponseNotificationHandler.SetNodeb - Unpack failed. Error: %v", err)
			return err
		}

		nbIdentity.InventoryName = nodebInfo.RanName
		nbIdentity.GlobalNbId = enbId
		nodebInfo.GlobalNbId = nbIdentity.GlobalNbId
		nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
		nodebInfo.E2ApplicationProtocol = entities.E2ApplicationProtocol_X2_SETUP_REQUEST
		nodebInfo.NodeType = entities.Node_ENB
		nodebInfo.Configuration = &entities.NodebInfo_Enb{Enb: enb}

		return nil
}

