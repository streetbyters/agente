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

// Package cmn provides in-app modules
package cmn

import (
	"github.com/akdilsiz/release-agent/model"
	"github.com/streadway/amqp"
	"net/url"
	"strconv"
)

// RabbitMq Conn Structure
type ChannelRabbitMq struct {
	ChannelInterface
	App 		*App
	Conn		*amqp.Connection
	Channel		*amqp.Channel
}

// NewRabbitMq
func NewRabbitMq(app *App) *ChannelRabbitMq {
	return &ChannelRabbitMq{App: app}
}

// Start RabbitMQ Conn
func (r *ChannelRabbitMq) Start() {
	uri := url.URL{
		Scheme:     "amqp",
		User:      	url.UserPassword(r.App.Config.RabbitMqUser, r.App.Config.RabbitMqPass),
		Host:       r.App.Config.RabbitMqHost + ":" + strconv.Itoa(r.App.Config.RabbitMqPort),
	}
	conn, err := amqp.Dial(uri.String())
	if err != nil {
		panic(err)
	}
	r.Conn = conn
	r.App.Logger.LogInfo("Start RabbitMQ Connection")
}

func (r *ChannelRabbitMq) Subscribe() {
	channel, err := r.Conn.Channel()
	if err != nil {
		panic(err)
	}

	r.Channel = channel

	err = r.Channel.ExchangeDeclare(
		r.App.Config.ChannelName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel exchange declare")
	}
}

func (r *ChannelRabbitMq) Receive() {
	q, err := r.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel queue declare")
	}

	err = r.Channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		r.App.Config.ChannelName, // exchange
		false,
		nil)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel queue bind")
	}

	received, err := r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for receive := range received {
			r.App.Logger.LogInfo("Receive rabbitmq message: " + string(receive.Body))
			r.App.Job.Run(model.NewReceivedMessage(string(receive.Body)))
		}
	}()

	<- r.App.Channel

	defer r.Channel.Close()
	defer r.Conn.Close()
}
