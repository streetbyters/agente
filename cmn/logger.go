//
// Copyright 2019 Abdulkadir DILSIZ
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cmn provides in-app modules
package cmn

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"strings"
	"time"
)

// Logger Custom Logger Structure
type Logger struct {
	zerolog.Logger
}

// NewLogger We start custom logger according to our working environment.
func NewLogger(mode string) *Logger {
	var logger zerolog.Logger

	if mode == "test" || mode == "dev" {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("Msg: %v |", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}

		logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	return &Logger{Logger: logger}
}

// LogError Error log
func (l *Logger) LogError(err error, msg string) {
	_, fn, line, _ := runtime.Caller(1)

	l.Error().Str("fn", fn).Int("line", line).Err(err).Msg(msg)
}

// LogFatal Fatal log
func (l *Logger) LogFatal(err error) {
	_, fn, line, _ := runtime.Caller(1)

	l.Fatal().Str("fn", fn).Int("line", line).Err(err).Msg("")
}

// LogInfo Info log
func (l *Logger) LogInfo(msg string) {
	l.Info().Msg(msg)
}

// LogDebug Debug log
func (l *Logger) LogDebug(msg string) {
	l.Debug().Msg(msg)
}

