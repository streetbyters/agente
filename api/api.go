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
	"github.com/akdilsiz/release-agent/cmn"
	"github.com/akdilsiz/release-agent/model"
	"github.com/valyala/fasthttp"
)

type Api struct {
	App 	*cmn.App
	Router	*Router
}

func NewApi(app *cmn.App) *Api {
	api := &Api{App: app}
	api.Router = NewRouter(api)

	return api
}

func (a *Api) JSONResponse(ctx *fasthttp.RequestCtx, response model.Response, status int) {
	ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	ctx.SetBody([]byte(response.ToJson()))
	ctx.SetStatusCode(status)
}
