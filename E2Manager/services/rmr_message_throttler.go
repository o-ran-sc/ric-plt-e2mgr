//
// Copyright 2026 O-RAN SC
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

package services

import (
	"sync"
	"time"
)

type RmrMessageThrottler struct {
	messageCounters      map[string]*messageCounter
	mutex                sync.RWMutex
	maxMessagesPerSecond int
	windowSize           time.Duration
}

type messageCounter struct {
	count  int
	window time.Time
}

func NewRmrMessageThrottler(maxMessagesPerSecond int) *RmrMessageThrottler {
	return &RmrMessageThrottler{
		messageCounters:      make(map[string]*messageCounter),
		maxMessagesPerSecond: maxMessagesPerSecond,
		windowSize:           time.Second,
	}
}

func (t *RmrMessageThrottler) AllowMessage(ranName string) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()
	counter, exists := t.messageCounters[ranName]

	if !exists || now.Sub(counter.window) > t.windowSize {
		t.messageCounters[ranName] = &messageCounter{count: 1, window: now}
		return true
	}

	if counter.count >= t.maxMessagesPerSecond {
		return false
	}

	counter.count++
	return true
}
