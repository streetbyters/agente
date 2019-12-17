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
	model2 "github.com/akdilsiz/agente/database/model"
	"github.com/akdilsiz/agente/model"
	"github.com/valyala/fasthttp"
)

type JobController struct {
	Controller
	*API
}

func (c JobController) Index(ctx *fasthttp.RequestCtx) {
	paginate, errs, err := c.Paginate(ctx, "id", "inserted_at")
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: err.Error(),
		}, fasthttp.StatusUnprocessableEntity)
	}

	job := model2.NewJob()
	var jobs []model2.Job
	c.App.Database.QueryWithModel("SELECT * FROM "+job.TableName()+" AS j"+
		" ORDER BY "+paginate.OrderField+" "+paginate.OrderBy+
		" LIMIT $3 OFFSET $4",
		&jobs,
		paginate.Limit,
		paginate.Offset)

	var count int64
	err = c.App.Database.DB.Get(&count, "SELECT count(*) FROM "+job.TableName())

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       jobs,
		TotalCount: count,
	}, fasthttp.StatusOK)
}
