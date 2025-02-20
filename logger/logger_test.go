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
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogExample(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	_ = buf
	logger := InitLogger(LoggerConfig{
		Output: os.Stdout,
		Level:  "debug",
	})

	// use logger
	logger.Debugf("debug-log-%d", 1)
	logger.WithField("field-key-1", "field-val-1").Infof("info-log-%d", 1)

	// with context, set traceId
	ctx := context.WithValue(context.Background(), TraceIdKey{}, "TRACE-ID")
	ctxlog := logger.WithContext(ctx)
	ctxlog.Debugf("debug-log-%d", 2)
	fieldlog := ctxlog.WithFields(map[string]interface{}{
		"field-key-2": "field-val-2",
		"field-key-3": "field-val-3",
	})
	// use the same field logger to avoid allocte new Entry
	fieldlog.Infof("info-log-%d", 2)
	fieldlog.Infof("info-log-%s", "2.1")

	// use logrus
	logrus.Debugf("debug-log-%d", 3)
	logrus.WithField("field-key-3", "field-val-3").Infof("info-log-%d", 3)
	fmt.Printf("%s", buf.Bytes())
}

func TestLogFile(t *testing.T) {
	InitLogger(LoggerConfig{
		Output:     nil,
		Level:      "debug",
		Filename:   "../tests/test.log",
		MaxSize:    10, // 10M
		MaxAge:     3,  // 3days
		MaxBackups: 3,
		LocalTime:  false,
		Compress:   false,
	})

	// use logrus
	logrus.Debugf("debug-log-%d", 1)
	logrus.WithField("field-key-1", "field-val-1").Infof("info-log-%d", 1)
}
