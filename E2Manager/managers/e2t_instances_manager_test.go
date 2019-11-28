package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const E2TAddress = "10.0.2.15"

func initE2TInstancesManagerTest(t *testing.T) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *E2TInstancesManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)
	return readerMock, writerMock, e2tInstancesManager
}

func TestAddNewE2TInstanceEmptyAddress(t *testing.T) {
	_, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	err := e2tInstancesManager.AddE2TInstance("")
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestAddNewE2TInstanceSaveE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "GetE2TInfoList")
}

func TestAddNewE2TInstanceGetE2TInfoListInternalFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tInfoList := []*entities.E2TInstanceInfo{}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, common.NewInternalError(errors.New("Error")))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
	rnibReaderMock.AssertNotCalled(t, "SaveE2TInfoList")
}

func TestAddNewE2TInstanceNoE2TInfoList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tInfoList := []*entities.E2TInstanceInfo{}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, common.NewResourceNotFoundError(""))
	e2tInfoList = append(e2tInfoList, &entities.E2TInstanceInfo{Address: E2TAddress})
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TInfoList", e2tInfoList)
}

func TestAddNewE2TInstanceEmptyE2TInfoList(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tInfoList := []*entities.E2TInstanceInfo{}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)
	e2tInfoList = append(e2tInfoList, &entities.E2TInstanceInfo{Address: E2TAddress})
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(nil)
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.Nil(t, err)
	rnibWriterMock.AssertCalled(t, "SaveE2TInfoList", e2tInfoList)
}

func TestAddNewE2TInstanceSaveE2TInfoListFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)
	e2tInfoList := []*entities.E2TInstanceInfo{}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)
	e2tInfoList = append(e2tInfoList, &entities.E2TInstanceInfo{Address: E2TAddress})
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(common.NewResourceNotFoundError(""))
	err := e2tInstancesManager.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
}
