// Copyright 2019 Street Byters Community
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

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/forgolang/agente/api"
	"github.com/forgolang/agente/cmn"
	"github.com/forgolang/agente/database"
	model2 "github.com/forgolang/agente/database/model"
	"github.com/forgolang/agente/model"
	"github.com/forgolang/agente/utils"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

// configFile App config file name string
var configFile string

// devMode Development mode flag
var devMode string

// migrate
var migrate bool

// reset
var reset bool

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	var mode model.MODE
	var dbPath string
	var appPath string
	var libPath string

	flag.StringVar(&devMode, "mode", "dev", "Development Mode")
	flag.BoolVar(&migrate, "migrate", false, "Run migrations")
	flag.BoolVar(&reset, "reset", false, "Reset database")
	flag.StringVar(&configFile, "config", "", "Config file")
	flag.StringVar(&dbPath, "dbPath", "", "Database path")
	flag.StringVar(&appPath, "appPath", "", "Application path")
	flag.Parse()

	if appPath == "" {
		appPath, _ = os.Getwd()
	}
	dirs := strings.SplitAfter(appPath, "agente")

	if devMode == "dev" || devMode == "test" {
		mode = model.MODE(devMode)
		appPath = path.Join(dirs[0])
		dbPath = appPath
		libPath = path.Join(appPath, "files")
	} else {
		mode = model.Prod
		appPath = path.Join("etc", "agente")
		dbPath = path.Join("var", "lib", "agente", "db")
		libPath = path.Join("var", "lib", "agente", "files")
	}

	if configFile == "" {
		configFile = "agente." + string(mode) + ".env"
	}

	logger := utils.NewLogger(string(mode))

	viper.SetConfigName(configFile)
	viper.AddConfigPath(appPath)
	err := viper.ReadInConfig()
	cmn.FailOnError(logger, err)

	config := &model.Config{
		NodeType:     model.Node(viper.GetString("TYPE")),
		Path:         appPath,
		LibPath:      libPath,
		Port:         viper.GetInt("PORT"),
		SecretKey:    viper.GetString("SECRET_KEY"),
		DB:           model.DB(viper.GetString("DB")),
		DBPath:       dbPath,
		DBName:       viper.GetString("DB_NAME"),
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetInt("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPass:       viper.GetString("DB_PASS"),
		DBSsl:        viper.GetString("DB_SSL"),
		RabbitMqHost: viper.GetString("RABBITMQ_HOST"),
		RabbitMqPort: viper.GetInt("RABBITMQ_PORT"),
		RabbitMqUser: viper.GetString("RABBITMQ_USER"),
		RabbitMqPass: viper.GetString("RABBITMQ_PASS"),
		RedisHost:    viper.GetString("REDIS_HOST"),
		RedisPort:    viper.GetInt("REDIS_PORT"),
		RedisPass:    viper.GetString("REDIS_PASS"),
		RedisDB:      viper.GetInt("REDIS_DB"),
		Versioning:   viper.GetBool("VERSIONING"),
		ChannelName:  viper.GetString("CHANNEL_NAME"),
		Scheduler:    viper.GetString("SCHEDULER"),
	}

	if config.DB == "" {
		panic(errors.New("enter DB conf"))
	}
	if config.Port == 0 {
		panic(errors.New("enter PORT conf"))
	}

	if config.ChannelName == "" {
		panic(errors.New("enter CHANNEL_NAME conf"))
	}

	if config.RedisHost == "" && config.RabbitMqHost == "" {
		panic(errors.New("enter one of the redis or rabbitMQ configurations"))
	}

	if config.RedisHost != "" && config.RabbitMqHost != "" {
		panic(errors.New("you can only work on one queue(redis or rabbitMQ) system"))
	}

	db, err := database.NewDB(config)
	cmn.FailOnError(logger, err)
	db.Logger = logger

	newApp := cmn.NewApp(config, logger)
	newApp.Channel = ch
	newApp.Database = db
	newApp.Mode = mode

	if migrate {
		db.Reset = reset
		database.InstallDB(db)
		return
	}

	genNode(newApp)

	newApp.Scheduler = cmn.NewScheduler(newApp)

	newAPI := api.NewAPI(newApp)
	go func() {
		err := newAPI.Router.Server.ListenAndServe(newAPI.Router.Addr)
		cmn.FailOnError(logger, err)
	}()

	<-newApp.Channel
}

func genNode(app *cmn.App) {
	app.Logger.LogInfo("Generating node information")

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	hostname = fmt.Sprintf("node_%s@%s", string(app.Mode), hostname)
	node := model2.NewNode()
	res := app.Database.QueryRowWithModel(fmt.Sprintf(`SELECT * FROM %s `+
		`WHERE code = $1`+
		`ORDER BY id DESC LIMIT 1`,
		node.TableName()),
		node,
		hostname)

	app.Config.NodeName = hostname

	if res.Error != nil {
		node.Name = hostname
		node.Code = hostname
		err := app.Database.Insert(new(model2.Node), node, "id", "inserted_at")
		if err != nil {
			panic(errors.New("node information could not be created on the database, "+err.Error()))
		}
	}

	app.Node = node

	app.Logger.LogInfo("Node information was created")
}
