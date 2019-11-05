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

package models

import "fmt"

type ProcessStats struct {
	SentCount               int
	SentErrorCount          int
	ReceivedExpectedCount   int
	ReceivedUnexpectedCount int
	ReceivedErrorCount      int
}

type ProcessResult struct {
	Stats ProcessStats
	Err   error
}

func (pr ProcessResult) String() string {
	return fmt.Sprintf("\nNumber of sent messages: %d\nNumber of send errors: %d\n" +
		"Number of expected received messages: %d\nNumber of unexpected received messages: %d\n" +
		"Number of receive errors: %d\n", pr.Stats.SentCount, pr.Stats.SentErrorCount, pr.Stats.ReceivedExpectedCount, pr.Stats.ReceivedUnexpectedCount, pr.Stats.ReceivedErrorCount)
}
