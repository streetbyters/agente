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

package cmn

import (
	"github.com/forgolang/agente/model"
	"github.com/streadway/amqp"
	"net/url"
	"strconv"
)

// ChannelRabbitMq queuing structure
type ChannelRabbitMq struct {
	ChannelInterface
	App     *App
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMq building rabbitMQ queuing
func NewRabbitMq(app *App) *ChannelRabbitMq {
	return &ChannelRabbitMq{App: app}
}

// Start RabbitMQ Conn
func (r *ChannelRabbitMq) Start() {
	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(r.App.Config.RabbitMqUser, r.App.Config.RabbitMqPass),
		Host:   r.App.Config.RabbitMqHost + ":" + strconv.Itoa(r.App.Config.RabbitMqPort),
	}
	conn, err := amqp.Dial(uri.String())
	FailOnError(r.App.Logger, err)

	r.Conn = conn
	r.App.Logger.LogInfo("Start RabbitMQ Connection")
}

// Subscribe rabbitMQ channel
func (r *ChannelRabbitMq) Subscribe() {
	channel, err := r.Conn.Channel()
	FailOnError(r.App.Logger, err)

	r.Channel = channel

	err = r.Channel.ExchangeDeclare(
		r.App.Config.ChannelName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel exchange declare")
	}
}

// Receive consume rabbitMQ channel
func (r *ChannelRabbitMq) Receive() {
	q, err := r.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel queue declare")
	}

	err = r.Channel.QueueBind(
		q.Name,
		"",
		r.App.Config.ChannelName,
		false,
		nil)
	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel queue bind")
	}

	received, err := r.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		r.App.Logger.LogError(err, "rabbitMQ error channel queue consume")
	}

	go func() {
		for receive := range received {
			r.App.Logger.LogInfo("Receive rabbitmq message: " + string(receive.Body))
			r.App.Job.Run(model.NewReceivedMessage(string(receive.Body)))
		}
	}()

	<-r.App.Channel

	defer r.Channel.Close()
	defer r.Conn.Close()
}
