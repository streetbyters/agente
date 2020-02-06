// Copyright 2019 StreetByters Community
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
	"github.com/streetbyters/agente/database"
	model2 "github.com/streetbyters/agente/database/model"
	"github.com/streetbyters/agente/model"
	"github.com/streetbyters/agente/utils"
	"os"
)

// App structure
type App struct {
	Database        *database.Database
	Channel         chan os.Signal
	Config          *model.Config
	Logger          *utils.Logger
	ChannelRabbitMq *ChannelRabbitMq
	ChannelRedis    *ChannelRedis
	Scheduler       *Scheduler
	Mode            model.MODE
	Job             *Job
	Node            *model2.Node
}

// NewApp building new app
func NewApp(config *model.Config, logger *utils.Logger) *App {
	app := &App{
		Config: config,
		Logger: logger,
	}

	if app.Config.RabbitMqHost != "" {
		app.Config.RabbitMq = true
		app.ChannelRabbitMq = NewRabbitMq(app)
		app.ChannelRabbitMq.Start()
		app.ChannelRabbitMq.Subscribe()
		go app.ChannelRabbitMq.Receive()
	}

	if app.Config.RedisHost != "" {
		app.Config.Redis = true
		app.ChannelRedis = NewRedis(app)
		app.ChannelRedis.Start()
		app.ChannelRedis.Subscribe()
		go app.ChannelRedis.Receive()
	}

	app.Job = NewJob(app)

	app.Logger.LogInfo("Started application")

	return app
}

// FailOnError panic error with logger
func FailOnError(logger *utils.Logger, err error) {
	if err != nil {
		logger.Panic().Err(err)
		panic(err)
	}
}
