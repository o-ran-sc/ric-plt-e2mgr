package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"fmt"
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

func TestGetE2TInstanceSuccess(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address := "10.10.2.15:9800"
	e2tInstance := entities.NewE2TInstance(address)
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, nil)
	res, err := e2tInstancesManager.GetE2TInstance(address)
	assert.Nil(t, err)
	assert.Equal(t, e2tInstance, res)
}

func TestGetE2TInstanceFailure(t *testing.T) {
	rnibReaderMock, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address := "10.10.2.15:9800"
	var e2tInstance *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", address).Return(e2tInstance, common.NewInternalError(fmt.Errorf("for test")))
	res, err := e2tInstancesManager.GetE2TInstance(address)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAssociateRanSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	e2tInfo3 := entities.NewE2TInstanceInfo(address1)
	e2tInfo3.AssociatedRanCount = 1;
	e2tInfoList2 := []*entities.E2TInstanceInfo{e2tInfo3, e2tInfo2}
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList2).Return(nil)

	e2tInstance1  := entities.NewE2TInstance(address1)
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	err := e2tInstancesManager.AssociateRan("test1", address1)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestAssociateRanGetListFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"

	var e2tInfoList []*entities.E2TInstanceInfo
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInfoList")
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
	rnibReaderMock.AssertNotCalled(t, "GetE2TInstance")
}

func TestAssociateRanSaveListFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
	rnibReaderMock.AssertNotCalled(t, "GetE2TInstance")
}

func TestAssociateRanGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(nil)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestAssociateRanSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	e2tInfo3 := entities.NewE2TInstanceInfo(address1)
	e2tInfo3.AssociatedRanCount = 1;
	e2tInfoList2 := []*entities.E2TInstanceInfo{e2tInfo3, e2tInfo2}
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList2).Return(nil)

	e2tInstance1  := entities.NewE2TInstance(address1)
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.AssociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestDeassociateRanSuccess(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 1;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 0;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	e2tInfo3 := entities.NewE2TInstanceInfo(address1)
	e2tInfo3.AssociatedRanCount = 0;
	e2tInfoList2 := []*entities.E2TInstanceInfo{e2tInfo3, e2tInfo2}
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList2).Return(nil)

	e2tInstance1  := entities.NewE2TInstance(address1)
	e2tInstance1.AssociatedRanList = append(e2tInstance1.AssociatedRanList, "test0", "test1")
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(nil)

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.Nil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}

func TestDeassociateRanNoInstanceFound(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfoList := []*entities.E2TInstanceInfo{}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.Nil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInfoList")
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
	rnibReaderMock.AssertNotCalled(t, "GetE2TInstance")
}

func TestDeassociateRanGetListFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"

	var e2tInfoList []*entities.E2TInstanceInfo
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInfoList")
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
	rnibReaderMock.AssertNotCalled(t, "GetE2TInstance")
}

func TestDeassociateRanSaveListFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
	rnibReaderMock.AssertNotCalled(t, "GetE2TInstance")
}

func TestDeassociateRanGetInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 0;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 1;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList).Return(nil)

	var e2tInstance1 *entities.E2TInstance
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibWriterMock.AssertNotCalled(t, "SaveE2TInstance")
}

func TestDeassociateRanSaveInstanceFailure(t *testing.T) {
	rnibReaderMock, rnibWriterMock, e2tInstancesManager := initE2TInstancesManagerTest(t)

	address1 := "10.10.2.15:9800"
	e2tInfo1 := entities.NewE2TInstanceInfo(address1)
	e2tInfo1.AssociatedRanCount = 1;
	address2 := "10.10.2.15:9801"
	e2tInfo2 := entities.NewE2TInstanceInfo(address2)
	e2tInfo2.AssociatedRanCount = 0;
	e2tInfoList := []*entities.E2TInstanceInfo{e2tInfo1, e2tInfo2}
	rnibReaderMock.On("GetE2TInfoList").Return(e2tInfoList, nil)

	e2tInfo3 := entities.NewE2TInstanceInfo(address1)
	e2tInfo3.AssociatedRanCount = 0;
	e2tInfoList2 := []*entities.E2TInstanceInfo{e2tInfo3, e2tInfo2}
	rnibWriterMock.On("SaveE2TInfoList", e2tInfoList2).Return(nil)

	e2tInstance1  := entities.NewE2TInstance(address1)
	rnibReaderMock.On("GetE2TInstance", address1).Return(e2tInstance1, nil)
	rnibWriterMock.On("SaveE2TInstance", mock.Anything).Return(common.NewInternalError(fmt.Errorf("for test")))

	err := e2tInstancesManager.DeassociateRan("test1", address1)
	assert.NotNil(t, err)
	rnibReaderMock.AssertExpectations(t)
	rnibWriterMock.AssertExpectations(t)
}


func TestRemoveE2TInstance(t *testing.T) {
	_, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address1 := "10.10.2.15:9800"
	e2tInstance1  := entities.NewE2TInstance(address1)
	err := e2tInstancesManager.RemoveE2TInstance(e2tInstance1)
	assert.Nil(t, err)
}

func TestSelectE2TInstance(t *testing.T) {
	_, _, e2tInstancesManager := initE2TInstancesManagerTest(t)
	address1 := "10.10.2.15:9800"
	e2tInstance1  := entities.NewE2TInstance(address1)
	addr, err := e2tInstancesManager.SelectE2TInstance(e2tInstance1)
	assert.Nil(t, err)
	assert.Equal(t, "", addr)
}


