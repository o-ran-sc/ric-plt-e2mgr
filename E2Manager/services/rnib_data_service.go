//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	PingRnib() bool
}

type rNibDataService struct {
	logger        *logger.Logger
	rnibReader    reader.RNibReader
	rnibWriter    rNibWriter.RNibWriter
	maxAttempts   int
	retryInterval time.Duration
}

func NewRnibDataService(logger *logger.Logger, config *configuration.Configuration, rnibReader reader.RNibReader, rnibWriter rNibWriter.RNibWriter) *rNibDataService {
	return &rNibDataService{
		logger:        logger,
		rnibReader:    rnibReader,
		rnibWriter:    rnibWriter,
		maxAttempts:   config.MaxRnibConnectionAttempts,
		retryInterval: time.Duration(config.RnibRetryIntervalMs) * time.Millisecond,
	}
}

func (w *rNibDataService) UpdateNodebInfo(nodebInfo *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.UpdateNodebInfo - nodebInfo: %s", nodebInfo)

	err := w.retry("UpdateNodebInfo", func() (err error) {
		err = w.rnibWriter.UpdateNodebInfo(nodebInfo)
		return
	})

	return err
}

func (w *rNibDataService) SaveNodeb(nbIdentity *entities.NbIdentity, nb *entities.NodebInfo) error {
	w.logger.Infof("#RnibDataService.SaveNodeb - nbIdentity: %s, nodebInfo: %s", nbIdentity, nb)

	err := w.retry("SaveNodeb", func() (err error) {
		err = w.rnibWriter.SaveNodeb(nbIdentity, nb)
		return
	})

	return err
}

func (w *rNibDataService) SaveRanLoadInformation(inventoryName string, ranLoadInformation *entities.RanLoadInformation) error {
	w.logger.Infof("#RnibDataService.SaveRanLoadInformation - inventoryName: %s, ranLoadInformation: %s", inventoryName, ranLoadInformation)

	err := w.retry("SaveRanLoadInformation", func() (err error) {
		err = w.rnibWriter.SaveRanLoadInformation(inventoryName, ranLoadInformation)
		return
	})

	return err
}

func (w *rNibDataService) GetNodeb(ranName string) (*entities.NodebInfo, error) {
	w.logger.Infof("#RnibDataService.GetNodeb - RAN name: %s", ranName)

	var nodeb *entities.NodebInfo = nil

	err := w.retry("GetNodeb", func() (err error) {
		nodeb, err = w.rnibReader.GetNodeb(ranName)
		return
	})

	return nodeb, err
}

func (w *rNibDataService) GetListNodebIds() ([]*entities.NbIdentity, error) {
	w.logger.Infof("#RnibDataService.GetListNodebIds")

	var nodeIds []*entities.NbIdentity = nil

	err := w.retry("GetListNodebIds", func() (err error) {
		nodeIds, err = w.rnibReader.GetListNodebIds()
		return
	})

	return nodeIds, err
}

func (w *rNibDataService) PingRnib() bool {
	err := w.retry("GetListNodebIds", func() (err error) {
		_, err = w.rnibReader.GetListNodebIds()
		return
	})

	return !isRnibConnectionError(err)
}

func (w *rNibDataService) retry(rnibFunc string, f func() error) (err error) {
	attempts := w.maxAttempts

	for i := 1; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if !isRnibConnectionError(err) {
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

func isRnibConnectionError(err error) bool {
	internalErr, ok := err.(*common.InternalError)
	if !ok {
		return false
	}
	_, ok = internalErr.Err.(*net.OpError)

	return ok
}
