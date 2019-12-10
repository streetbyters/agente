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

package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/akdilsiz/release-agent/cmn"
	"github.com/akdilsiz/release-agent/model"
	"github.com/akdilsiz/release-agent/model/response"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"os"
	"path"
	"strings"
	"testing"
)

type Suite struct {
	suite.Suite
	Api *Api
}

type Method string
const (
	Post	Method = "POST"
	Get		Method = "GET"
	Put		Method = "PUT"
	Delete	Method = "DELETE"
)

type ContentType string
const (
	JSON	ContentType = "application/json"
	XML		ContentType = "application/xml"
	HTML	ContentType = "text/html"
)

type TestResponse struct {
	Result	model.Response
	Status	int
	Error	bool
}

func NewSuite() *Suite {
	var mode model.MODE
	var dbPath string

	configFile := "release-agent.env"
	appPath, _ := os.Getwd()
	dirs := strings.SplitAfter(appPath, "release-agent")

	mode = model.Test
	appPath = path.Join(dirs[0])
	dbPath = appPath

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

	newApp := cmn.NewApp(config, logger)
	newApp.Database = database
	newApi := NewApi(newApp)

	return &Suite{Api: newApi}
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

func (s *Suite) JSON(method Method, path string, arg ...interface{}) *TestResponse {
	return s.request(false, "", JSON, method, path, arg...)
}

func (s *Suite) XML(method Method, path string, arg ...interface{}) *TestResponse {
	return s.request(false, "", XML, method, path, arg...)
}

func SetupSuite(s *Suite) {}

func TearDownSuite(s *Suite) {}

func (s *Suite) request(auth bool,
	authToken string,
	contentType ContentType,
	method Method,
	path string,
	body ...interface{}) *TestResponse {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(path)
	req.Header.SetContentType(string(contentType) + "; charset=utf-8")
	if auth {
		req.Header.Set("Authorization", "Bearer " + authToken)
	}
	req.Header.SetMethod(string(method))

	if len(body) > 0 {
		switch contentType {
		case JSON:
				b, err := json.Marshal(body[0])
				if err != nil {
					req.SetBody(b)
				}
			break
		case XML:
			b, err := xml.Marshal(body[0])
			if err != nil {
				req.SetBody(b)
			}
			break
		}
	}

	resp := fasthttp.AcquireResponse()
	err := s.serveApi(s.Api.Router.Handler.ServeFastHTTP, req, resp)
	if err != nil {

	}
	fmt.Print(string(resp.Body()))
	testResponse := &TestResponse{}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		var ts1 response.Success
		err = json.Unmarshal(resp.Body(), &ts1)
		if err == nil {
			testResponse.Result = ts1
		}
		testResponse.Error = false
	} else if resp.StatusCode() >= 400 && resp.StatusCode() < 500 {
		var ts2 response.Error
		err = json.Unmarshal(resp.Body(), &ts2)
		if err == nil {
			testResponse.Result = ts2
		}
		testResponse.Error = true
	} else {
		testResponse.Error = true
	}

	testResponse.Status = resp.StatusCode()
	return testResponse
}

func (s *Suite) serveApi(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(err)
		}
	}()

	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return client.Do(req, res)
}
