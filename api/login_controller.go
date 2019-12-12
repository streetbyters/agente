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
	"fmt"
	"github.com/akdilsiz/agente/database"
	model2 "github.com/akdilsiz/agente/database/model"
	"github.com/akdilsiz/agente/model"
	"github.com/valyala/fasthttp"
	"net/http"
)

// LoginController user authentication controller
type LoginController struct {
	Controller
	*API
}

// Create user sign in method
func (c LoginController) Create(ctx *fasthttp.RequestCtx) {
	var loginRequest model.LoginRequest

	c.JSONBody(ctx, &loginRequest)
	fmt.Println(loginRequest)
	if errs, err := database.ValidateStruct(loginRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: http.StatusText(http.StatusUnprocessableEntity),
		}, http.StatusUnprocessableEntity)
		return
	}

	var userModel model2.User

	result := c.App.Database.QueryRow("SELECT * FROM " + userModel.TableName() + " AS u" +
		"WHERE u.username = $1 OR u.email = $1")
	if result.Error != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: nil,
			Detail: http.StatusText(http.StatusNotFound),
		}, http.StatusNotFound)
		return
	}

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: "OK",
	}, http.StatusCreated)
}

