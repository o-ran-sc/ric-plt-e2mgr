package managers

import (
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type RanSetupManager struct {
	logger             *logger.Logger
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
	rmrService         *services.RmrService
}

func NewRanSetupManager(logger *logger.Logger, rmrService *services.RmrService, rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) *RanSetupManager {
	return &RanSetupManager{
		logger:             logger,
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
	}
}

func (m *RanSetupManager) ExecuteSetup(nodebInfo *entities.NodebInfo) error {
	return nil
}
