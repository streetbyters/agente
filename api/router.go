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
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	errors2 "github.com/akdilsiz/agente/errors"
	"github.com/akdilsiz/agente/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Router api router structure
type Router struct {
	API     *API
	Server  *fasthttp.Server
	Addr    string
	Handler *phi.Mux
}

var (
	prefix           string
	reqID            uint64
	allowHeaders     = "authorization"
	allowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	allowOrigin      = "*"
	allowCredentials = "true"
)

// NewRouter building api router
func NewRouter(api *API) *Router {
	router := &Router{
		API: api,
	}

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])

	r := phi.NewRouter()

	r.Use(router.requestID)
	r.Use(router.recover)
	r.Use(router.logger)
	r.Use(router.cors)

	r.NotFound(router.notFound)
	r.MethodNotAllowed(router.methodNotAllowed)

	r.Get("/", HomeController{API: api}.Index)

	r.Route("/api/v1", func(r phi.Router) {
		r.Route("/user", func(r phi.Router) {
			r.Post("/sign_in", LoginController{API: api}.Create)
			r.Post("/token", TokenController{API: api}.Create)
		})

		r.Group(func(r phi.Router) {
			r.Use(api.JWTAuth.Verify)
			// Job Routes
			r.Group(func(r phi.Router) {
				r.Get("/job", JobController{API: api}.Index)
				r.Post("/job", JobController{API: api}.Create)
				r.Route("/job/{jobID}", func(r phi.Router) {
					r.Get("/", JobController{API: api}.Show)
					r.Delete("/", JobController{API: api}.Delete)

					// Detail Routes
					r.Route("/detail", func(r phi.Router) {
						r.Post("/", JobDetailController{API: api}.Create)
					})
				})
			})

			// File routes
			r.Group(func(r phi.Router) {
				r.Get("/file", FileController{API: api}.Index)
				r.Post("/file", FileController{API: api}.Create)
				r.Route("/file/{fileID}", func(r phi.Router) {
					r.Get("/", FileController{API: api}.Show)
					r.Put("/", FileController{API: api}.Update)
					r.Delete("/", FileController{API: api}.Delete)
				})
			})

			r.Get("/upload/dir", UploadController{API: api}.DirIndex)
			r.Post("/upload", UploadController{API: api}.Create)
		})
	})

	router.Server = &fasthttp.Server{
		Handler:            r.ServeFastHTTP,
		ReadTimeout:        10 * time.Second,
		MaxRequestBodySize: 1 * 1024 * 1024 * 1024,
		Logger:             api.App.Logger,
	}
	router.Addr = ":" + strconv.Itoa(api.App.Config.Port)
	router.Handler = r

	return router
}

func (r Router) notFound(ctx *fasthttp.RequestCtx) {
	r.API.JSONResponse(ctx, model.ResponseError{
		Errors: nil,
		Detail: "not found",
	}, http.StatusNotFound)
}

func (r Router) methodNotAllowed(ctx *fasthttp.RequestCtx) {
	r.API.JSONResponse(ctx, model.ResponseError{
		Errors: nil,
		Detail: "method not allowed",
	}, http.StatusMethodNotAllowed)
}

// Reference: https://github.com/go-chi/chi/blob/master/middleware/request_id.go
func (r Router) requestID(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		id := atomic.AddUint64(&reqID, 1)
		requestID := fmt.Sprintf("%s-%06d", prefix, id)
		ctx.SetUserValue("requestID", requestID)
		ctx.Response.Header.Add("x-request-id", requestID)
		next(ctx)
	}
}

// Reference: https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
func (r Router) recover(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				var err error
				switch x := rvr.(type) {
				case *errors2.PluggableError:
					e := rvr.(*errors2.PluggableError)
					r.API.JSONResponse(ctx, model.ResponseError{
						Errors: nil,
						Detail: e.Error(),
					}, e.Status)

					defer func() {
						r.API.App.Logger.LogError(e, "Pluggable error")
					}()
					return
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}

				if r.API.App.Mode == model.Test {
					panic(rvr)
				}

				r.API.App.Logger.LogError(err, "router recover")
				r.API.JSONResponse(ctx, model.ResponseError{
					Errors: nil,
					Detail: http.StatusText(http.StatusInternalServerError),
				}, http.StatusInternalServerError)
				return
			}
		}()

		next(ctx)
	}
}

func (r Router) logger(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		defer func() {
			if r.API.App.Mode != model.Test {
				if r.API.App.Mode == model.Prod {
					r.API.App.Logger.LogInfo("Path: " + string(ctx.Path()) +
						" Method: " + string(ctx.Method()) +
						" - " + strconv.Itoa(ctx.Response.StatusCode()))
				} else {
					r.API.App.Logger.LogDebug("Path: " + string(ctx.Path()) +
						" Method: " + string(ctx.Method()) +
						" - " + strconv.Itoa(ctx.Response.StatusCode()))
				}
			}
		}()
	}
}

func (r Router) cors(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Request.Header.Method()) == "OPTIONS" {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", allowCredentials)
			ctx.Response.Header.Set("Access-Control-Allow-Headers", allowHeaders)
			ctx.Response.Header.Set("Access-Control-Allow-Methods", allowMethods)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", allowOrigin)
			ctx.Response.Header.Set("Accept", "application/json")
			ctx.Response.Header.Set("Accept", "multipart/form-data")

			ctx.SetStatusCode(http.StatusNoContent)
			return
		}
		next(ctx)
	}
}
