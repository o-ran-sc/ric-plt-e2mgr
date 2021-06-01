package managers

import (
        "e2mgr/converters"
        "e2mgr/tests"
        "fmt"
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
        "github.com/stretchr/testify/assert"
        "testing"
)

func TestPopulateX2NodebPduSuccess(t *testing.T) {
        logger := tests.InitLog(t)
        nodebInfo := &entities.NodebInfo{}
        nodebIdentity := &entities.NbIdentity{}
        handler := NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger))
        err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createRandomPayload())
        assert.NotNil(t, err)
}

func TestPopulateX2NodebPduFailure(t *testing.T) {
        logger := tests.InitLog(t)
        nodebInfo := &entities.NodebInfo{}
        nodebIdentity := &entities.NbIdentity{}
        handler := NewX2SetupResponseManager(converters.NewX2SetupResponseConverter(logger))
        err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createX2SetupResponsePayload(t))
        assert.Nil(t, err)
}
func createX2SetupResponsePayload(t *testing.T) []byte {
        packedPdu := "4006001a0000030005400200000016400100001140087821a00000008040"
        var payload []byte
        _, err := fmt.Sscanf(packedPdu, "%x", &payload)
        if err != nil {
                t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
        }
        return payload
}


