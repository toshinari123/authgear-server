// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"context"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	loggers map[string]*logrus.Logger
	lock    sync.Mutex
)

func init() {
	loggers = map[string]*logrus.Logger{}
	loggers[""] = logrus.StandardLogger()
}

func Logger(name string) *logrus.Logger {
	lock.Lock()
	defer lock.Unlock()

	logger, ok := loggers[name]
	if !ok {
		logger = logrus.New()

		if logger == nil {
			panic("logrus.New() returns nil")
		}

		loggers[name] = logger
	}

	return logger
}

func Loggers() map[string]*logrus.Logger {
	lock.Lock()
	defer lock.Unlock()

	ret := map[string]*logrus.Logger{}
	for loggerName, logger := range loggers {
		ret[loggerName] = logger
	}
	return ret
}

func SetFormatter(formatter logrus.Formatter) {
	lock.Lock()
	defer lock.Unlock()

	for _, logger := range loggers {
		logger.Formatter = formatter
	}
}

func SetLevel(level logrus.Level) {
	lock.Lock()
	defer lock.Unlock()

	for _, logger := range loggers {
		logger.Level = level
	}
}

func SetOutput(out io.Writer) {
	lock.Lock()
	defer lock.Unlock()

	for _, logger := range loggers {
		logger.Out = out
	}
}

func AddHook(hook logrus.Hook) {
	lock.Lock()
	defer lock.Unlock()

	for _, logger := range loggers {
		logger.Hooks.Add(hook)
	}
}

func LoggerEntry(name string) *logrus.Entry {
	return LoggerEntryWithTag(name, name)
}

func LoggerEntryWithTag(name string, tag string) *logrus.Entry {
	logger := Logger(name)
	fields := logrus.Fields{}
	if name != "" {
		fields["logger"] = name
	}
	if tag != "" {
		fields["tag"] = tag
	}
	return logger.WithFields(fields)
}

func CreateLogger(ctx context.Context, logger string) *logrus.Entry {
	var requestTag string
	fields := logrus.Fields{}
	if ctx != nil {
		if tag, ok := ctx.Value("RequestTag").(string); ok {
			requestTag = tag
		}

		if requestID, ok := ctx.Value("RequestID").(string); ok {
			fields["request_id"] = requestID
		}
	}
	return LoggerEntryWithTag(logger, requestTag).WithFields(fields)
}
