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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package e2managererrors

import "fmt"

type WrongStateError struct {
	*BaseError
}

func NewWrongStateError(activityName string, state string) *WrongStateError {
	return &WrongStateError{
		&BaseError{
			Code:    403,
			Message: fmt.Sprintf("Activity %s rejected. RAN current state %s does not allow its execution ", activityName, state) ,
		},
	}
}

func (e *WrongStateError) Error() string {
	return e.Message
}
