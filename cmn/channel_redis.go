// Copyright 2019 Abdulkadir Dilsiz
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
	"github.com/akdilsiz/agente/model"
	"github.com/go-redis/redis/v7"
	"strconv"
	"strings"
)

// ChannelRedis queuing structure
type ChannelRedis struct {
	ChannelInterface
	App    *App
	Client *redis.Client
	PubSub *redis.PubSub
}

// NewRedis building redis queuing
func NewRedis(app *App) *ChannelRedis {
	return &ChannelRedis{App: app}
}

// Start Redis Conn
func (r *ChannelRedis) Start() {
	client := redis.NewClient(&redis.Options{
		Network:   "tcp",
		Addr:      strings.Join([]string{r.App.Config.RedisHost, strconv.Itoa(r.App.Config.RedisPort)}, ":"),
		Dialer:    nil,
		OnConnect: nil,
		Password:  r.App.Config.RedisPass,
		DB:        r.App.Config.RedisDB,
		PoolSize:  5,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	r.Client = client
	r.App.Logger.LogInfo("Start Redis connection")
}

// Subscribe redis channel
func (r *ChannelRedis) Subscribe() {
	r.PubSub = r.Client.Subscribe(r.App.Config.ChannelName)
}

// Receive redis channel
func (r *ChannelRedis) Receive() {
	_, err := r.PubSub.Receive()
	if err != nil {
		r.App.Logger.LogError(err, "redis error channel receive")
	}

	ch := r.PubSub.Channel()

	for received := range ch {
		r.App.Logger.LogInfo("Receive redis message: " + received.Payload)
		r.App.Job.Run(model.NewReceivedMessage(received.Payload))
	}

	defer r.Client.Close()
}
