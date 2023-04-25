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


package logger

import (
	mdclog "gerrit.o-ran-sc.org/r/com/golog"
	"time"
)

type Logger struct {
	logger *mdclog.MdcLogger
}

func InitLogger(loglevel int8 ) (*Logger, error) {
       name := "e2mgr"
       log ,err:= NewLogger(name)
       return log,err
}

func NewLogger(name string) (*Logger, error) {
	l,err:= mdclog.InitLogger(name)
	return &Logger{
		logger: l,
	},err
}

func (l *Logger) SetFormat(logMonitor int) {
    l.logger.Mdclog_format_initialize(logMonitor)
}

func (l *Logger) SetLevel(level int) {
	l.logger.LevelSet(mdclog.Level(level))
}

func (l *Logger) SetMdc(key string, value string) {
	l.logger.MdcAdd(key, value)
}

func (l *Logger) Errorf(pattern string, args ...interface{}) {
	l.SetMdc("time", time.Now().Format(time.RFC3339))
	l.logger.Error(pattern, args...)
}

func (l *Logger) Warnf(pattern string, args ...interface{}) {
	l.SetMdc("time", time.Now().Format(time.RFC3339))
	l.logger.Warning(pattern, args...)
}

func (l *Logger) Infof(pattern string, args ...interface{}) {
	l.SetMdc("time", time.Now().Format(time.RFC3339))
	l.logger.Info(pattern, args...)
}

func (l *Logger) Debugf(pattern string, args ...interface{}) {
	l.SetMdc("time", time.Now().Format(time.RFC3339))
	l.logger.Debug(pattern, args...)
}
