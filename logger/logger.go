/**
 * Copyright 2025 OceanBase
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const defaultTimestampFormat = "2006-01-02T15:04:05.99999-07:00"
const INIT_TRACEID = "0000000000000000"

var textFormatter = &TextFormatter{
	TimestampFormat:        "2006-01-02T15:04:05.99999-07:00", // log timestamp format
	FullTimestamp:          true,
	DisableLevelTruncation: true,
	FieldMap: map[string]string{
		"WARNING": "WARN", // log level string, use WARN
	},
	// log caller, filename:line callFunction
	CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
		n := 0
		filename := frame.File
		// 获取包名
		for i := len(filename) - 1; i > 0; i-- {
			if filename[i] == '/' {
				n++
				if n >= 2 {
					filename = filename[i+1:]
					break
				}
			}
		}

		name := frame.Function
		idx := strings.LastIndex(name, ".")
		return name[idx+1:], fmt.Sprintf("%s:%d", filename, frame.Line)
	},
}

type LoggerConfig struct {
	Output     io.Writer
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"maxage"`
	MaxBackups int    `yaml:"maxbackups"`
	LocalTime  bool   `yaml:"localtime"`
	Compress   bool   `yaml:"compress"`
}

func InitLogger(config LoggerConfig) *logrus.Logger {
	logger := logrus.StandardLogger()
	// log output
	if config.Output == nil {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		})
	} else {
		logger.SetOutput(config.Output)
	}

	// log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		panic(fmt.Sprintf("parse log level: %+v", err))
	}
	logger.SetLevel(level)

	// log format
	logger.SetFormatter(textFormatter)
	logger.SetReportCaller(true)

	return logger
}
