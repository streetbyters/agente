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

package model

// Config Application config structure
type Config struct {
	Path         string `json:"path"`
	Mode         MODE   `json:"mode"`
	Port         int    `json:"port"`
	DB           DB     `json:"db"`
	DBPat        string `json:"db_path"`
	DBName       string `json:"db_name"`
	DBHost       string `json:"db_host"`
	DBPort       int    `json:"db_port"`
	DBUser       string `json:"db_user"`
	DBPass       string `json:"db_pass"`
	DBSsl        string `json:"db_ssl"`
	RabbitMq     bool   `json:"-"`
	RabbitMqHost string `json:"rabbitmq_host"`
	RabbitMqPort int    `json:"rabbitmq_port"`
	RabbitMqUser string `json:"rabbitmq_user"`
	RabbitMqPass string `json:"rabbitmq_pass"`
	Redis        bool   `json:"-"`
	RedisHost    string `json:"redis_host"`
	RedisPort    int    `json:"redis_port"`
	RedisPass    string `json:"redis_pass"`
	RedisDB      int    `json:"redis_db"`
	ChannelName  string `json:"channel_name"`
	Versioning   bool   `json:"versioning"`
	Scheduler    string `json:"scheduler"`
}
