package services

import (
        "e2mgr/configuration"
        "e2mgr/logger"
        "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
        "testing"
        "github.com/stretchr/testify/assert"
)


func RanAlarmServiceTest(t *testing.T) (RanAlarmService, *logger.Logger, *configuration.Configuration) {
    DebugLevel := int8(4)
    logger, err := logger.InitLogger(DebugLevel)
        if err != nil {
                t.Errorf("#... - failed to initialize logger, error: %s", err)
        }
    config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
    ranAlarmServiceInstance := NewRanAlarmService(logger , config)
    return ranAlarmServiceInstance,logger, config
}


func TestSetConnectivityChangeAlarmTest(t *testing.T){
     ranAlarmServiceInstance,_,_ := RanAlarmServiceTest(t)
     nodebInfo := &entities.NodebInfo{}
     err := ranAlarmServiceInstance.SetConnectivityChangeAlarm(nodebInfo)
     assert.Nil(t,err)
}

