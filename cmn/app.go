//
// Copyright 2019 Abdulkadir DILSIZ - TransferChain
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

package cmn

import (
	"github.com/akdilsiz/release-agent/model"
	"os"
)

// App structure
type App struct {
	Database	*Database
	Channel		chan os.Signal
	Config		*model.Config
	Logger 		*Logger
	RabbitMq	model.RabbitMq
	Redis		model.Redis
	Job			*Job
	Mode		model.MODE
}

// NewApp application with config structure and logger package
func NewApp(config *model.Config, logger *Logger) *App {
	app := &App{
		Config:		config,
		RabbitMq:	model.RabbitMq{},
		Redis:		model.Redis{},
		Logger:		logger,
	}

	if app.Config.RabbitMqHost != "" {
		app.Config.RabbitMq = true
	} else if app.Config.RedisHost != "" {
		app.Config.Redis = true
	}

	app.Job = NewJob(app)
	app.Job.Start()

	app.Logger.LogInfo("Started application")

	return app
}
