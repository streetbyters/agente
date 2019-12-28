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
	model2 "github.com/akdilsiz/agente/database/model"
	"github.com/akdilsiz/agente/model"
	"github.com/valyala/fasthttp"
)

// FileController job files api controller
type FileController struct {
	Controller
	*API
}

// Index list all job files
func (c FileController) Index(ctx *fasthttp.RequestCtx) {
	paginate, errs, err := c.Paginate(ctx, "id", "inserted_at", "updated_at")
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: err.Error(),
		}, fasthttp.StatusUnprocessableEntity)
	}

	upload := model2.NewFile()
	var uploads []model2.File
	res := c.App.Database.QueryWithModel(fmt.Sprintf("SELECT f.* FROM %s AS f "+
		"ORDER BY f.%s %s LIMIT $1 OFFSET $2",
		upload.TableName(), paginate.OrderField, paginate.OrderBy),
		&uploads,
		paginate.Limit,
		paginate.Offset)

	fmt.Println(res.Error)

	var count int64
	c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", upload.TableName()))

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       uploads,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Create job file
func (c FileController) Create(ctx *fasthttp.RequestCtx) {

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusCreated)
}

// Update job file
func (c FileController) Update(ctx *fasthttp.RequestCtx) {

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusOK)
}

// Delete job file
func (c FileController) Delete(ctx *fasthttp.RequestCtx) {

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusNoContent)
}
