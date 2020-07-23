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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package e2managererrors

import "fmt"

type UnknownSetupRequestRanNameError struct {
	*BaseError
	ranName string
}

func NewUnknownSetupRequestRanNameError(ranName string) *UnknownSetupRequestRanNameError {
	return &UnknownSetupRequestRanNameError{
		&BaseError{
			Code:    599,
			Message: fmt.Sprintf("Can't extract node-type out of RAN name %s", ranName),
		},
		ranName,
	}
}

func (e *UnknownSetupRequestRanNameError) Error() string {
	return e.Message
}
