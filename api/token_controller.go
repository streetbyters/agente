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
	"github.com/forgolang/agente/database"
	model2 "github.com/forgolang/agente/database/model"
	"github.com/forgolang/agente/model"
	"github.com/valyala/fasthttp"
	"net/http"
)

// TokenController user authentication token controller
type TokenController struct {
	Controller
	*API
}

// Create generate user jwt method
func (c TokenController) Create(ctx *fasthttp.RequestCtx) {
	tokenRequest := new(model.TokenRequest)

	c.JSONBody(ctx, &tokenRequest)
	if errs, err := database.ValidateStruct(tokenRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: http.StatusText(http.StatusUnprocessableEntity),
		}, http.StatusUnprocessableEntity)
		return
	}

	passphrase := new(model2.UserPassphrase)
	result := c.App.Database.QueryRowWithModel(passphrase.PassphraseQuery(c.App.Database),
		passphrase,
		tokenRequest.Passphrase)
	if result.Error != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}, fasthttp.StatusNotFound)
		return
	}

	user := new(model2.User)
	result = c.App.Database.QueryRowWithModel("SELECT u.* FROM "+user.TableName()+" AS u "+
		"WHERE u.id = $1 AND u.is_active = true",
		user,
		passphrase.UserID)
	if result.Error != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Detail: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}, fasthttp.StatusNotFound)
		return
	}

	jwt, _ := c.API.JWTAuth.Generate(user.ID)

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: model.ResponseToken{
			JWT:    jwt,
			UserID: user.ID,
		},
	}, http.StatusCreated)
}
