//-build !test
// Copyright 2019 Forgolang Community
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forgolang/agente/cmn"
	"github.com/forgolang/agente/database"
	model2 "github.com/forgolang/agente/database/model"
	"github.com/forgolang/agente/model"
	"github.com/forgolang/agente/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"io/ioutil"
	"mime/multipart"
	"net"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"
)

var defaultLogger *utils.Logger

// Suite application test structure
type Suite struct {
	suite.Suite
	API  *API
	Auth struct {
		User  *model2.User
		Token string
	}
}

// Method request method for test api request
type Method string

const (
	// Options method for api request
	Options Method = "OPTIONS"
	// Post method for api request
	Post Method = "POST"
	// Get method for api request
	Get Method = "GET"
	// Put method for api request
	Put Method = "PUT"
	// Delete method for api request
	Delete Method = "DELETE"
)

// ContentType request content type for test api request
type ContentType string

const (
	// JSON Content type for api request
	JSON ContentType = "application/json"
	// XML Content type for api request
	XML ContentType = "application/xml"
	// HTML Content type for api request
	HTML ContentType = "text/html"
	// FormData Multipart/form-data Content Type for api request
	FormData ContentType = "multipart/form-data"
)

// TestResponse response model for test api request
type TestResponse struct {
	RequestError error
	Success      model.ResponseSuccess
	Error        model.ResponseError
	Other        interface{}
	Status       int
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

	logger := utils.NewLogger(string(mode))
	defaultLogger = logger
	viper.SetConfigName(configFile)
	viper.AddConfigPath(appPath)
	err := viper.ReadInConfig()
	cmn.FailOnError(logger, err)

	config := &model.Config{
		NodeType:     model.Node(viper.GetString("TYPE")),
		Path:         appPath,
		LibPath:      path.Join(appPath, "files"),
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

	db, err := database.NewDB(config)
	cmn.FailOnError(logger, err)
	db.Logger = logger
	db.Reset = true
	database.InstallDB(db)

	newApp := cmn.NewApp(config, logger)
	newApp.Database = db
	newApp.Mode = model.Test

	ch := make(chan bool)
	go genNode(ch, newApp)
	<-ch

	newAPI := NewAPI(newApp)

	return &Suite{API: newAPI}
}

func genNode(ch chan bool, app *cmn.App) {
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
			panic(errors.New("node information could not be created on the database, " + err.Error()))
		}
	}

	app.Node = node

	app.Logger.LogInfo("Node information was created")

	time.AfterFunc(time.Millisecond * 100, func() {
		ch<-true
	})
}

// Run run test suites
func Run(t *testing.T, s suite.TestingSuite) {
	defaultLogger.LogInfo(fmt.Sprintf("Started %s tests", reflect.TypeOf(s).Name()))
	suite.Run(t, s)
	defaultLogger.LogInfo(fmt.Sprintf("Finish %s tests", reflect.TypeOf(s).Name()))
}

// JSON api json request
func (s *Suite) JSON(method Method, path string, arg interface{}) *TestResponse {
	return s.request(JSON, method, path, arg)
}

// File api form-data request
func (s *Suite) File(method Method, path string, arg interface{}, fileParam ...string) *TestResponse {
	return s.request(FormData, method, path, arg, fileParam...)
}

//// XML api xml request
//func (s *Suite) XML(method Method, path string, arg ...interface{}) *TestResponse {
//	return s.request(false, "", XML, method, path, arg...)
//}

// SetupSuite before suite processes
func SetupSuite(s *Suite) {}

// TearDownSuite after suite processes
func TearDownSuite(s *Suite) {}

// UserAuth generate test request auth provides
func UserAuth(s *Suite) {
	user := model2.NewUser("1234")
	user.Username = "testUser"
	user.Email = "testUser@tecpor.com"
	user.IsActive = true
	user.NodeID = s.API.App.Node.ID
	userModel := new(model2.User)

	err := s.API.App.Database.Insert(userModel, user,
		"id", "inserted_at")
	if err != nil {
		panic(err)
	}

	token, _ := s.API.JWTAuth.Generate(user.ID)
	s.Auth.Token = token
	s.Auth.User = user
}

// request test request for api
func (s *Suite) request(contentType ContentType, method Method, path string, body interface{}, fileParam ...string) *TestResponse {
	testResponse := &TestResponse{}
	var err error
	req := fasthttp.AcquireRequest()
	req.Header.SetHost(s.API.Router.Addr)
	req.Header.SetRequestURI(path)
	if s.Auth.Token != "" {
		req.Header.Set("Authorization", "Bearer "+s.Auth.Token)
	}
	req.Header.SetMethod(string(method))
	resp := fasthttp.AcquireResponse()

	if body != nil {
		switch contentType {
		case JSON:
			req.Header.SetContentType(string(contentType) + "; charset=utf-8")
			b, err := json.Marshal(body)
			if err == nil {
				req.SetBody(b)
			}
			break
		case FormData:
			b, _ := body.(map[string]interface{})
			body2 := new(bytes.Buffer)
			writer := multipart.NewWriter(body2)

			for _, f := range fileParam {
				file, err := os.Open(b[f].(string))
				if err == nil {
					fileContents, err := ioutil.ReadAll(file)
					if err == nil {
						fi, err := file.Stat()
						if err == nil {
							file.Close()
							part, err := writer.CreateFormFile(f, fi.Name())
							if err == nil {
								part.Write(fileContents)
							}
						}
					}
				}
			}

			for key, val := range b {
				if ok, _ := utils.InArray(key, fileParam); !ok {
					_ = writer.WriteField(key, val.(string))
				}
			}
			err = writer.Close()
			if err != nil {
				testResponse.RequestError = err
				return testResponse
			}
			req.SetBodyStream(body2, body2.Len())
			req.Header.SetContentType(writer.FormDataContentType())
			break
		}
	}

	err = s.serveAPI(s.API.Router.Handler.ServeFastHTTP, req, resp)
	if err != nil {
		testResponse.RequestError = err
		return testResponse
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		var ts1 model.ResponseSuccess
		err = json.Unmarshal(resp.Body(), &ts1)
		if err == nil {
			testResponse.Success = ts1
		}
	} else if resp.StatusCode() >= 400 && resp.StatusCode() < 501 {
		var ts2 model.ResponseError
		err = json.Unmarshal(resp.Body(), &ts2)
		if err == nil {
			testResponse.Error = ts2
		}
	} else {
		testResponse.Other = resp.Body()
	}

	testResponse.Status = resp.StatusCode()
	return testResponse
}

func (s *Suite) serveAPI(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		cmn.FailOnError(s.API.App.Logger, err)
	}()

	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return client.Do(req, res)
}
