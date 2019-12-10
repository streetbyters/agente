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

package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/akdilsiz/release-agent/model/response"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Router struct {
	Api 		*Api
	Server		*fasthttp.Server
	Addr		string
	Handler		*phi.Mux
}

var (
	prefix string
	reqID uint64
	allowHeaders     = "authorization"
	allowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	allowOrigin      = "*"
	allowCredentials = "true"
)

func NewRouter(api *Api) *Router {
	router := &Router{
		Api: api,
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
	r.Use(router.logger)
	r.Use(router.cors)

	r.NotFound(router.notFound)
	r.MethodNotAllowed(router.methodNotAllowed)

	r.Get("/", HomeController{Api: api}.Index)

	router.Server = &fasthttp.Server{
		Handler: 		r.ServeFastHTTP,
		ReadTimeout:	10 * time.Second,
	}
	router.Addr = ":" + strconv.Itoa(api.App.Config.Port)
	router.Handler = r

	return router
}

func (r Router) notFound(ctx *fasthttp.RequestCtx) {
	r.Api.JSONResponse(ctx, response.Error{
		Errors:   nil,
		Detail:   "not found",
	}, http.StatusNotFound)
}

func (r Router) methodNotAllowed(ctx *fasthttp.RequestCtx) {
	r.Api.JSONResponse(ctx, response.Error{
		Errors:   nil,
		Detail:   "method not allowed",
	}, http.StatusMethodNotAllowed)
}

func (r Router) requestID(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		id := atomic.AddUint64(&reqID, 1)
		requestID := fmt.Sprintf("%s-%06d", prefix, id)
		ctx.SetUserValue("requestID", requestID)
		ctx.Response.Header.Add("x-request-id", requestID)
		next(ctx)
	}
}

func (r Router) logger(next phi.HandlerFunc) phi.HandlerFunc {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		defer func() {
			r.Api.App.Logger.LogInfo("Path: " + string(ctx.Path()) +
				" Method: " + string(ctx.Method()) +
				" - " + strconv.Itoa(ctx.Response.StatusCode()))
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

			ctx.SetStatusCode(http.StatusNoContent)
			return
		} else {
			next(ctx)
		}
	}
}

//func (r *Router)
