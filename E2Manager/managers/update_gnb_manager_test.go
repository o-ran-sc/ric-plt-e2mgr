package managers

import (
        "e2mgr/configuration"
        "e2mgr/logger"
        "e2mgr/services"
        "e2mgr/mocks"
        //"e2mgr/models"
        "testing"
        "github.com/stretchr/testify/assert"
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)
func initUpdateGnbManagerTest(t *testing.T) (*UpdateGnbManager, *logger.Logger, services.RNibDataService, *NodebValidator) {
        logger, err := logger.InitLogger(logger.DebugLevel)
        if err != nil {
                t.Errorf("#... - failed to initialize logger, error: %s", err)
        }
        config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
        //rmrMessengerMock := &mocks.RmrMessengerMock{}
        readerMock := &mocks.RnibReaderMock{}
        writerMock := &mocks.RnibWriterMock{}
        nodebValidator := NewNodebValidator()
        rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
        UpdateGnbManager := NewUpdateGnbManager(logger, rnibDataService, nodebValidator)
        return UpdateGnbManager, logger, rnibDataService, nodebValidator
}

func TestRemoveNodebCellsGnb(t *testing.T) {
      UpdateGnbManager,_,_, _ := initUpdateGnbManagerTest(t)
      nodebInfo := &entities.NodebInfo{}
      //updateEnbRequest := &models.UpdateEnbRequest{}
      res :=UpdateGnbManager.RemoveNodebCells(nodebInfo)
      assert.NotNil(t, res)
}

func TestValidateNodebGnb(t *testing.T) {
      UpdateGnbManager,_,_, _ := initUpdateGnbManagerTest(t)
      nodebInfo := &entities.NodebInfo{}
      res := UpdateGnbManager.ValidateNodeb(nodebInfo)
      assert.Nil(t, res)
}



