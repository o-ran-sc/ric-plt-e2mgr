//
// Copyright 2021 AT&T Intellectual Property
// Copyright 2021 Nokia
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

package mocks

import "github.com/stretchr/testify/mock"

type MockSdlSyncStorage struct {
	mock.Mock
}

func (m *MockSdlSyncStorage) SubscribeChannel(ns string, cb func(string, ...string), channels ...string) error {
	a := m.Called(ns, cb, channels)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) UnsubscribeChannel(ns string, channels ...string) error {
	a := m.Called(ns, channels)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) SetAndPublish(ns string, channelsAndEvents []string, pairs ...interface{}) error {
	a := m.Called(ns, channelsAndEvents, pairs)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) SetIfAndPublish(ns string, channelsAndEvents []string, key string, oldData, newData interface{}) (bool, error) {
	a := m.Called(ns, channelsAndEvents, key, oldData, newData)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlSyncStorage) SetIfNotExistsAndPublish(ns string, channelsAndEvents []string, key string, data interface{}) (bool, error) {
	a := m.Called(ns, channelsAndEvents, key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlSyncStorage) RemoveAndPublish(ns string, channelsAndEvents []string, keys []string) error {
	a := m.Called(ns, channelsAndEvents, keys)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) RemoveIfAndPublish(ns string, channelsAndEvents []string, key string, data interface{}) (bool, error) {
	a := m.Called(ns, channelsAndEvents, key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlSyncStorage) RemoveAllAndPublish(ns string, channelsAndEvents []string) error {
	a := m.Called(ns, channelsAndEvents)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) Set(ns string, pairs ...interface{}) error {
	a := m.Called(ns, pairs)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) Get(ns string, keys []string) (map[string]interface{}, error) {
	a := m.Called(ns, keys)
	return a.Get(0).(map[string]interface{}), a.Error(1)
}

func (m *MockSdlSyncStorage) GetAll(ns string) ([]string, error) {
	a := m.Called(ns)
	return a.Get(0).([]string), a.Error(1)
}

func (m *MockSdlSyncStorage) ListKeys(ns string, pattern string) ([]string, error) {
	a := m.Called(ns, pattern)
	return a.Get(0).([]string), a.Error(1)
}

func (m *MockSdlSyncStorage) Close() error {
	a := m.Called()
	return a.Error(0)
}

func (m *MockSdlSyncStorage) Remove(ns string, keys []string) error {
	a := m.Called(ns, keys)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) RemoveAll(ns string) error {
	a := m.Called(ns)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) SetIf(ns string, key string, oldData, newData interface{}) (bool, error) {
	a := m.Called(ns, key, oldData, newData)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlSyncStorage) SetIfNotExists(ns string, key string, data interface{}) (bool, error) {
	a := m.Called(ns, key, data)
	return a.Bool(0), a.Error(1)
}
func (m *MockSdlSyncStorage) RemoveIf(ns string, key string, data interface{}) (bool, error) {
	a := m.Called(ns, key, data)
	return a.Bool(0), a.Error(1)
}

func (m *MockSdlSyncStorage) AddMember(ns string, group string, member ...interface{}) error {
	a := m.Called(ns, group, member)
	return a.Error(0)
}

func (m *MockSdlSyncStorage) RemoveMember(ns string, group string, member ...interface{}) error {
	a := m.Called(ns, group, member)
	return a.Error(0)
}
func (m *MockSdlSyncStorage) RemoveGroup(ns string, group string) error {
	a := m.Called(ns, group)
	return a.Error(0)
}
func (m *MockSdlSyncStorage) GetMembers(ns string, group string) ([]string, error) {
	a := m.Called(ns, group)
	return a.Get(0).([]string), a.Error(1)
}
func (m *MockSdlSyncStorage) IsMember(ns string, group string, member interface{}) (bool, error) {
	a := m.Called(ns, group, member)
	return a.Bool(0), a.Error(1)
}
func (m *MockSdlSyncStorage) GroupSize(ns string, group string) (int64, error) {
	a := m.Called(ns, group)
	return int64(a.Int(0)), a.Error(1)
}
