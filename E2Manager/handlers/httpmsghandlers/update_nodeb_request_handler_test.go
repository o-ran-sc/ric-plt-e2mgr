package httpmsghandlers

import (
    "e2mgr/configuration"
    "e2mgr/managers"
    "e2mgr/mocks"
    "e2mgr/logger"
    "e2mgr/models"
    "e2mgr/services"
 //   "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
    "github.com/stretchr/testify/assert"
    "testing"
)

func setupUpdateNodebRequestHandlerTest(t *testing.T) ( *UpdateNodebRequestHandler,  *mocks.RnibReaderMock, *mocks.RnibWriterMock){
        logger, err := logger.InitLogger(logger.DebugLevel)
        if err != nil {
                t.Errorf("#... - failed to initialize logger, error: %s", err)
        }
        config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
        readerMock := &mocks.RnibReaderMock{}
        writerMock := &mocks.RnibWriterMock{}
        rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
        updateNodebManager := managers.IUpdateNodebManager(nil)
        handler := NewUpdateNodebRequestHandler(logger,rnibDataService,updateNodebManager)
        return handler,readerMock, writerMock

}

func TestGetRanName(t *testing.T) {
        handler,_,_ := setupUpdateNodebRequestHandlerTest(t)
        updateEnbRequest := models.UpdateEnbRequest{}
        ret := handler.getRanName(updateEnbRequest)
        assert.Equal(t, "", ret )
}


