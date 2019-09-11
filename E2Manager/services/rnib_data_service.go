package services

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"net"
	"time"
)

type RNibDataService interface {
	SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error
	UpdateNodebInfo(nodebInfo *entities.NodebInfo) error
	SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error
	GetNodeb(ranName string) (*entities.NodebInfo, error)
	GetListNodebIds() ([]*entities.NbIdentity, error)
}

type rNibDataService struct {
	logger             *logger.Logger
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
	maxAttempts        int
	retryInterval      time.Duration
}

func NewRnibDataService(logger *logger.Logger, config *configuration.Configuration, rnibReaderProvider func() reader.RNibReader, rnibWriterProvider func() rNibWriter.RNibWriter) *rNibDataService {
	return &rNibDataService{
		logger:             logger,
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
		maxAttempts:        config.MaxRnibConnectionAttempts,
		retryInterval:      time.Duration(config.RnibRetryIntervalMs) * time.Millisecond,
	}
}

func (w *rNibDataService) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.UpdateNodebInfo - nodebInfo: %s", nodebInfo)

	err := w.retry("UpdateNodebInfo", func() (err error) {
		err = w.rnibWriterProvider().UpdateNodebInfo(nodebInfo)
		return
	})

	return err
}

func (w *rNibDataService) SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.SaveNodeb - nbIdentity: %s, nodebInfo: %s", nbIdentity, nb)

	err := w.retry("SaveNodeb", func() (err error) {
		err = w.rnibWriterProvider().SaveNodeb(nbIdentity, nb)
		return
	})

	return err
}

func (w *rNibDataService) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {
	w.logger.Infof("#RnibDataService.SaveRanLoadInformation - inventoryName: %s, ranLoadInformation: %s", inventoryName, ranLoadInformation)

	err := w.retry("SaveRanLoadInformation", func() (err error) {
		err = w.rnibWriterProvider().SaveRanLoadInformation(inventoryName, ranLoadInformation)
		return
	})

	return err
}

func (w *rNibDataService) GetNodeb(ranName string) (*entities.NodebInfo, error) {
	w.logger.Infof("#RnibDataService.GetNodeb - ranName: %s", ranName)

	var nodeb *entities.NodebInfo = nil

	err := w.retry("GetNodeb", func() (err error) {
		nodeb, err = w.rnibReaderProvider().GetNodeb(ranName)
		return
	})

	return nodeb, err
}

func (w *rNibDataService) GetListNodebIds() ([]*entities.NbIdentity, error) {
	w.logger.Infof("#RnibDataService.GetListNodebIds")

	var nodeIds []*entities.NbIdentity = nil

	err := w.retry("GetListNodebIds", func() (err error) {
		nodeIds, err = w.rnibReaderProvider().GetListNodebIds()
		return
	})

	return nodeIds, err
}

func (w *rNibDataService) retry(rnibFunc string, f func() error) (err error) {
	attempts := w.maxAttempts

	for i := 1; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if !w.isConnError(err) {
			return err
		}
		if i >= attempts {
			w.logger.Errorf("#RnibDataService.retry - after %d attempts of %s, last error: %s", attempts, rnibFunc, err)
			return err
		}
		time.Sleep(w.retryInterval)

		w.logger.Infof("#RnibDataService.retry - retrying %d %s after error: %s", i, rnibFunc, err)
	}
}

func (w *rNibDataService) isConnError(err error) bool {
	internalErr, ok := err.(common.InternalError)
	if !ok {
		return false
	}
	_, ok = internalErr.Err.(*net.OpError)

	return ok
}
