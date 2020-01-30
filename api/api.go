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
	"github.com/forgolang/agente/cmn"
	pluggableError "github.com/forgolang/agente/errors"
	"github.com/forgolang/agente/model"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
)

// API rest api structure
type API struct {
	App     *cmn.App
	Router  *Router
	JWTAuth *JWTAuth
	Auth    struct {
		ID int64
	}
}

// NewAPI building api
func NewAPI(app *cmn.App) *API {
	api := &API{App: app}
	api.JWTAuth = NewJWTAuth(api)
	api.Router = NewRouter(api)

	return api
}

// ParseQuery parse url query string
func (a *API) ParseQuery(ctx *fasthttp.RequestCtx) map[string]string {
	qs, _ := url.ParseQuery(string(ctx.URI().QueryString()))
	values := make(map[string]string)
	for key, val := range qs {
		values[key] = val[0]
	}

	return values
}

// Paginate request paginate build
func (a *API) Paginate(ctx *fasthttp.RequestCtx, orderFields ...string) (model.Pagination, map[string]string, error) {
	var err error
	errs := make(map[string]string)
	pagination := model.NewPagination()
	queryParams := a.ParseQuery(ctx)

	if val, ok := queryParams["limit"]; ok {
		pagination.Limit, err = strconv.Atoi(val)
		if err != nil {
			errs["limit"] = "is not valid"
		}
	}

	if val, ok := queryParams["offset"]; ok {
		pagination.Offset, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			errs["offset"] = "is not valid"
		}
	}

	if val, ok := queryParams["order_by"]; ok {
		pagination.OrderBy = val
	}

	if val, ok := queryParams["order_field"]; ok {
		pagination.OrderField = val
	}

	if err != nil {
		panic(pluggableError.New(fasthttp.StatusMessage(fasthttp.StatusBadRequest),
			fasthttp.StatusBadRequest,
			"paginate params not valid",
			errs))
	}

	errs, err = pagination.Validate(orderFields...)

	if err != nil {
		panic(pluggableError.New(fasthttp.StatusMessage(fasthttp.StatusBadRequest),
			fasthttp.StatusBadRequest,
			"paginate params not valid",
			errs))
	}

	return pagination, errs, err
}

// JSONBody parse given model request body
func (a *API) JSONBody(ctx *fasthttp.RequestCtx, model interface{}) {
	r := bytes.NewReader(ctx.PostBody())
	json.NewDecoder(r).Decode(&model)
}

// JSONResponse building json response
func (a *API) JSONResponse(ctx *fasthttp.RequestCtx, response model.ResponseInterface, status int) {
	ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	if response != nil {
		ctx.SetBody([]byte(response.ToJSON()))
	}
	ctx.SetStatusCode(status)
}
