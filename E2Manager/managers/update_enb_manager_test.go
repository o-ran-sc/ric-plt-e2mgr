package managers

import (
        "e2mgr/configuration"
        "e2mgr/logger"
        "e2mgr/services"
        "e2mgr/mocks"
        "e2mgr/models"
        "testing"
        "github.com/stretchr/testify/assert"
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)
func initUpdateEnbManagerTest(t *testing.T) (*UpdateEnbManager, *logger.Logger, services.RNibDataService, *NodebValidator) {
 	DebugLevel := int8(4)      
	 logger, err := logger.InitLogger(DebugLevel)
        if err != nil {
                t.Errorf("#... - failed to initialize logger, error: %s", err)
        }
        config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
        //rmrMessengerMock := &mocks.RmrMessengerMock{}
        readerMock := &mocks.RnibReaderMock{}
        writerMock := &mocks.RnibWriterMock{}
        nodebValidator := NewNodebValidator()
        rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
        UpdateEnbManager := NewUpdateEnbManager(logger, rnibDataService, nodebValidator)
        return UpdateEnbManager, logger, rnibDataService, nodebValidator
}

func TestSuccessfulSetNodeb(t *testing.T) {
        UpdateEnbManager,_,_, _ := initUpdateEnbManagerTest(t)
        nodebInfo := &entities.NodebInfo{}
        writerMock := &mocks.RnibWriterMock{}
        writerMock.On("SetNodeb", nodebInfo).Return(nil)
        updateEnbRequest := &models.UpdateEnbRequest{}
        UpdateEnbManager.SetNodeb(nodebInfo,updateEnbRequest)
        //writerMock.AssertNumberOfCalls(t, "SaveNodeb", 1)
}

func TestValidateRequestBody(t *testing.T) {
        UpdateEnbManager,_,_, _ := initUpdateEnbManagerTest(t)
        //nodebInfo := &entities.NodebInfo{}
        //writerMock := &mocks.RnibWriterMock{}
        //writerMock.On("UpdateNodeb", nodebInfo).Return(nil)
        updateEnbRequest := &models.UpdateEnbRequest{}
        res := UpdateEnbManager.validateRequestBody(updateEnbRequest)
        assert.NotNil(t, res)
}


func TestValidateNodeb(t *testing.T) {
      UpdateEnbManager,_,_, _ := initUpdateEnbManagerTest(t)
      nodebInfo := &entities.NodebInfo{}
      res := UpdateEnbManager.ValidateNodeb(nodebInfo)
      assert.Nil(t, res) 
}

func TestValidate(t *testing.T) {
      UpdateEnbManager,_,_, _ := initUpdateEnbManagerTest(t)
      //nodebInfo := &entities.NodebInfo{}
      updateEnbRequest := &models.UpdateEnbRequest{}
      res :=UpdateEnbManager.Validate(updateEnbRequest)
      assert.NotNil(t, res)
}

func TestRemoveNodebCells(t *testing.T) {
      UpdateEnbManager,_,_, _ := initUpdateEnbManagerTest(t)
      nodebInfo := &entities.NodebInfo{}
      //updateEnbRequest := &models.UpdateEnbRequest{}
      res :=UpdateEnbManager.RemoveNodebCells(nodebInfo)
      assert.NotNil(t, res)
     
}

