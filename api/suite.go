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
	"github.com/akdilsiz/agente/cmn"
	"github.com/akdilsiz/agente/model"
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

// Suite application test structure
type Suite struct {
	suite.Suite
	API *API
}

// Method request method for test api request
type Method string
const (
	// Post method for api request
	Post	Method = "POST"
	// Get method for api request
	Get		Method = "GET"
	// Put method for api request
	Put		Method = "PUT"
	// Delete method for api request
	Delete	Method = "DELETE"
)

// ContentType request content type for test api request
type ContentType string
const (
	// JSON Content type for api request
	JSON	ContentType = "application/json"
	// XML Content type for api request
	XML		ContentType = "application/xml"
	// HTML Content type for api request
	HTML	ContentType = "text/html"
)

// TestResponse response model for test api request
type TestResponse struct {
	Success model.ResponseSuccess
	Error   model.ResponseError
	Status  int
}

// NewSuite build test application
func NewSuite() *Suite {
	var mode model.MODE
	var dbPath string

	configFile := "agente.test.env"
	appPath, _ := os.Getwd()
	dirs := strings.SplitAfter(appPath, "agente")

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
		RedisDB:		viper.GetInt("REDIS_DB"),
		Versioning:		viper.GetBool("VERSIONING"),
	}

	database, err := cmn.NewDB(config, logger)
	if err != nil {
		logger.Panic().Err(err)
	}

	newApp := cmn.NewApp(config, logger)
	newApp.Database = database
	newApp.Mode = model.Test

	newAPI := NewAPI(newApp)

	return &Suite{API: newAPI}
}

// Run run test suites
func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

// JSON api json request
func (s *Suite) JSON(method Method, path string, arg ...interface{}) *TestResponse {
	return s.request(false, "", JSON, method, path, arg...)
}

// XML api xml request
func (s *Suite) XML(method Method, path string, arg ...interface{}) *TestResponse {
	return s.request(false, "", XML, method, path, arg...)
}

// SetupSuite before suite processes
func SetupSuite(s *Suite) {}

// TearDownSuite after suite processes
func TearDownSuite(s *Suite) {}

// request test request for api
func (s *Suite) request(auth bool, authToken string, contentType ContentType, method Method, path string, body ...interface{}) *TestResponse {
	req := fasthttp.AcquireRequest()
	req.Header.SetHost(s.API.Router.Addr)
	req.Header.SetRequestURI(path)
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
	err := s.serveAPI(s.API.Router.Handler.ServeFastHTTP, req, resp)
	if err != nil {
		//fmt.Println(err)
	}

	testResponse := &TestResponse{}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		var ts1 model.ResponseSuccess
		err = json.Unmarshal(resp.Body(), &ts1)
		if err == nil {
			testResponse.Success = ts1
		}
	} else if resp.StatusCode() >= 400 && resp.StatusCode() < 500 {
		var ts2 model.ResponseError
		err = json.Unmarshal(resp.Body(), &ts2)
		if err == nil {
			testResponse.Error = ts2
		}
	}

	testResponse.Status = resp.StatusCode()
	return testResponse
}

// serveAPI
func (s *Suite) serveAPI(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
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
