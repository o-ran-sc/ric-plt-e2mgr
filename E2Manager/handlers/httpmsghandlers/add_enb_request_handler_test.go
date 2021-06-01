package httpmsghandlers

import (
    "e2mgr/configuration"
    "e2mgr/managers"
    "e2mgr/mocks"
    "e2mgr/models"
    "e2mgr/services"
    "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
    "github.com/stretchr/testify/assert"
    "testing"
)

func setupAddEnbRequestHandlerTest(t *testing.T) (*AddEnbRequestHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock) {
    log := initLog(t)
    config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
    readerMock := &mocks.RnibReaderMock{}
    writerMock := &mocks.RnibWriterMock{}
    rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
    ranListManager := managers.NewRanListManager(log, rnibDataService)
    nodebValidator := managers.NewNodebValidator()
    handler := NewAddEnbRequestHandler(log,rnibDataService,nodebValidator, ranListManager)
    return handler, readerMock, writerMock
}

func CreateNodebInfoTest(t *testing.T,RanName string,  connectionStatus entities.ConnectionStatus) *entities.NodebInfo{

        setupRequest := &models.E2SetupRequestMessage{}
        nodebInfo := &entities.NodebInfo{
                RanName:                      RanName,
                SetupFromNetwork:             true,
                NodeType:                     entities.Node_GNB,
                ConnectionStatus:                         connectionStatus,
                Configuration: &entities.NodebInfo_Gnb{
                        Gnb: &entities.Gnb{
                                GnbType:      entities.GnbType_GNB,
                                RanFunctions: setupRequest.ExtractRanFunctionsList(),
                        },
                },
                GlobalNbId: &entities.GlobalNbId{
                        PlmnId: setupRequest.GetPlmnId(),
                        NbId:   setupRequest.GetNbId(),
                },
        }
        return nodebInfo
}

func TestGlobalNbIdValid(t *testing.T){
        globalNbId := &entities.GlobalNbId{}
        res :=  isGlobalNbIdValid(globalNbId)
        assert.NotNil(t,res)

}

func TestHandleAddEnbSuccess(t *testing.T) {
        handler, readerMock, writerMock := setupAddEnbRequestHandlerTest(t)
        ranName := "ran1"
        var rnibError error
        nodebInfo := &entities.NodebInfo{RanName: ranName, NodeType: entities.Node_ENB}
        readerMock.On("GetNodeb", ranName).Return(nodebInfo, rnibError)
        //writerMock.On("AddNbIdentity", entities.Node_ENB, nbIdentity).Return(nil)
        writerMock.On("AddEnb", nodebInfo).Return(nil)
        writerMock.On("AddNbIdentity", entities.Node_ENB, &entities.NbIdentity{InventoryName: "ran1", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}}).Return(nil)
        addEnbRequest := &models.AddEnbRequest{RanName: ranName}
        result, err := handler.Handle(addEnbRequest)
        assert.NotNil(t, err)
        assert.Nil(t, result)
}

func TestValidateRequestBody(t *testing.T){
     handler, _,_  := setupAddEnbRequestHandlerTest(t)
     
     ranName := "ran1"
     addEnbRequest := &models.AddEnbRequest{RanName: ranName}
     err := handler.validateRequestBody(addEnbRequest)
     assert.NotNil(t,err)
     }



