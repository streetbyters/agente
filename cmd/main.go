//
// Copyright 2019 Abdulkadir DILSIZ <TransferChain>
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
	"flag"
	"github.com/akdilsiz/release-agent/api"
	"github.com/akdilsiz/release-agent/cmn"
	"github.com/akdilsiz/release-agent/model"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

var configFile string
var devMode bool

func main() {
	var mode model.MODE
	var dbPath string
	appPath, _ := os.Getwd()
	dirs := strings.SplitAfter(appPath, "release-agent")

	flag.BoolVar(&devMode, "dev", false, "Development Mode")
	flag.StringVar(&configFile, "config", "release-agent.env", "Config file")
	flag.Parse()

	if devMode {
		mode = model.Dev
		appPath = path.Join(dirs[0])
		dbPath = appPath
	} else {
		mode = model.Prod
		appPath = path.Join("etc", "tc-release-agent")
		dbPath = path.Join("var", "lib", "tc-release-agent")
	}

	logger := cmn.NewLogger(string(mode))

	viper.SetConfigName(configFile)
	viper.AddConfigPath(appPath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic().Err(err)
	}

	config := &model.Config{
		Path:			appPath,
		Port:			viper.GetInt("PORT"),
		DB:				model.DB(viper.GetString("DB")),
		DBPath:			dbPath,
		DBName:			viper.GetString("DB_NAME"),
		DBHost: 		viper.GetString("DB_HOST"),
		DBPort:			viper.GetInt("DB_PORT"),
		DBUser:			viper.GetString("DB_USER"),
		DBPass:			viper.GetString("DB_PASS"),
		DBSsl:			viper.GetString("DB_SSL"),
		RabbitMqHost:	viper.GetString("RABBITMQ_HOST"),
		RabbitMqPort:	viper.GetInt("RABBITMQ_PORT"),
		RabbitMqUser:	viper.GetString("RABBITMQ_USER"),
		RabbitMqPass:	viper.GetString("RABBITMQ_PASS"),
		RedisHost:		viper.GetString("REDIS_HOST"),
		RedisPort:		viper.GetInt("REDIS_PORT"),
		RedisPass:		viper.GetString("REDIS_PASS"),
		RedisDB:		viper.GetString("REDIS_DB"),
		Versioning:		viper.GetBool("VERSIONING"),
	}

	database, err := cmn.NewDB(config, logger)
	if err != nil {
		logger.Panic().Err(err)
	}
	//
	//ch := make(chan os.Signal)
	//signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	newApp := cmn.NewApp(config, logger)
	//newApp.Channel = ch
	newApp.Database = database

	newApi := api.NewApi(newApp)
	logger.LogFatal(newApi.Router.Server.ListenAndServe(newApi.Router.Addr))
}

