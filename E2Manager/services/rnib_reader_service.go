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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type RnibReaderService struct {
	rnibReaderProvider func() reader.RNibReader
}

func NewRnibReaderService(rnibReaderProvider func() reader.RNibReader) *RnibReaderService{
	return &RnibReaderService{rnibReaderProvider}
}

func (s RnibReaderService) GetNodeb(ranName string) (*entities.NodebInfo, common.IRNibError) {
	return s.rnibReaderProvider().GetNodeb(ranName)
}

func (s  RnibReaderService) GetNodebIdList()([]*entities.NbIdentity, common.IRNibError) {
	return s.rnibReaderProvider().GetListNodebIds()
}


