package managers

import (
       /* "e2mgr/configuration"
        "e2mgr/logger"
        "e2mgr/mocks"
        "e2mgr/services"
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
       */
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
        //"github.com/pkg/errors"
        "github.com/stretchr/testify/assert"
        "testing"
)

func initNodebValidatorTest(t *testing.T)(*NodebValidator){
      nodebValidator := NewNodebValidator()
      return nodebValidator
}

func TestIsGnbValid(t *testing.T){
     nodebValidator := initNodebValidatorTest(t)
     gnb := entities.Gnb{}
     res := nodebValidator.IsGnbValid(&gnb)
     assert.NotNil(t,res)
}


func TestIsEnbValid(t *testing.T){
     nodebValidator := initNodebValidatorTest(t)
     enb := entities.Enb{}
     res := nodebValidator.IsEnbValid(&enb)
     assert.NotNil(t,res)
}

func TestIsServedNrCellInformationValid(t *testing.T){
        servedNrCellInformation :=  entities.ServedNRCellInformation{}
        err := isServedNrCellInformationValid(&servedNrCellInformation)
        assert.NotNil(t,err)
}

func TestIsServedNrCellInfoTddValid(t *testing.T){
        tdd := entities.ServedNRCellInformation_ChoiceNRMode_TddInfo{}
        err := isServedNrCellInfoTddValid(&tdd)
        assert.Nil(t,err)
}

func TestIsServedNrCellInfoFddValid(t *testing.T){
        fdd := entities.ServedNRCellInformation_ChoiceNRMode_TddInfo{}
        err := isServedNrCellInfoTddValid(&fdd)
        assert.Nil(t,err)
}

func TestIsNrNeighbourInfoTddValid(t *testing.T){
        tdd := entities.NrNeighbourInformation_ChoiceNRMode_TddInfo{}
        err := isNrNeighbourInfoTddValid(&tdd)
        assert.Nil(t,err)
}

func TestIsNrNeighbourInfoFddValid(t *testing.T){
        fdd := entities.NrNeighbourInformation_ChoiceNRMode_FddInfo{}
        err := isNrNeighbourInfoFddValid(&fdd)
        assert.Nil(t,err)
}

func TestIsTddInfoValid(t *testing.T){
        tdd := entities.TddInfo{}
        err := isTddInfoValid(&tdd)
        assert.Nil(t,err)
}

func TestIsFddInfoValid(t *testing.T){
        fdd := entities.FddInfo{}
        err := isFddInfoValid(&fdd)
        assert.Nil(t,err)
}

func TestIsServedNrCellInfoChoiceNrModeValid(t *testing.T){
        choiceNrMode := entities.ServedNRCellInformation_ChoiceNRMode{}
        res := isServedNrCellInfoChoiceNrModeValid(&choiceNrMode)
        assert.NotNil(t,res)
}

func  TestIsNrNeighbourInformationValid(t *testing.T){
        nrNeighbourInformation := entities.NrNeighbourInformation{}
        res := isNrNeighbourInformationValid(&nrNeighbourInformation)
        assert.NotNil(t,res)
}

func TestIsNrNeighbourInfoChoiceNrModeValid(t *testing.T){
        choiceNrMode := entities.NrNeighbourInformation_ChoiceNRMode{}
        res := isNrNeighbourInfoChoiceNrModeValid(&choiceNrMode)
        assert.NotNil(t,res)
}

