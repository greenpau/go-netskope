// Copyright 2020 Paul Greenberg greenpau@outlook.com
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

package netskope

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(opts map[string]interface{}) (*zap.Logger, error) {
	logLevel := "info"
	logEncoder := "console"
	if v, exists := opts["log_level"]; exists {
		logLevel = v.(string)
	}
	if v, exists := opts["log_encoder"]; exists {
		logEncoder = v.(string)
	}

	logAtom := zap.NewAtomicLevel()
	switch logLevel {
	case "info", "INFO":
		logAtom.SetLevel(zapcore.InfoLevel)
	case "warn", "WARN", "warning", "WARNING":
		logAtom.SetLevel(zapcore.WarnLevel)
	case "debug", "DEBUG", "dbg":
		logAtom.SetLevel(zapcore.DebugLevel)
	case "error", "ERROR":
		logAtom.SetLevel(zapcore.ErrorLevel)
	case "fatal", "FATAL":
		logAtom.SetLevel(zapcore.FatalLevel)
	default:
		return nil, fmt.Errorf("unsupported log level %s", logLevel)
	}

	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logEncoderConfig.TimeKey = "time"

	var core zapcore.Core

	if logEncoder == "console" {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(logEncoderConfig),
			zapcore.Lock(os.Stdout),
			logAtom,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(logEncoderConfig),
			zapcore.Lock(os.Stdout),
			logAtom,
		)
	}

	logger := zap.New(core)

	return logger, nil
}
