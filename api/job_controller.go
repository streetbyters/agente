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
	"github.com/fate-lovely/phi"
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
		" ORDER BY $1 $2" +
		" LIMIT $3 OFFSET $4",
		&jobs,
		paginate.OrderField,
		paginate.OrderBy,
		paginate.Limit,
		paginate.Offset)

	var count int64
	err = c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", job.TableName()))

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       jobs,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

func (c JobController) Show(ctx *fasthttp.RequestCtx) {
	var job model2.Job
	 c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", job.TableName()),
		&job,
		phi.URLParam(ctx, "id")).Force()

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: job,
	}, fasthttp.StatusOK)
}

func (c JobController) Create(ctx *fasthttp.RequestCtx) {
	job := model2.NewJob()
	c.JSONBody(ctx, &job)

	job.SourceUserId.SetValid(c.Auth.ID)

	c.App.Database.Insert(new(model2.Job), job, "id", "inserted_at")

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: job,
	}, fasthttp.StatusCreated)
}

func (c JobController) Delete(ctx *fasthttp.RequestCtx) {
	id := phi.URLParam(ctx, "id")

	job := model2.NewJob()
	c.App.Database.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", job.TableName()),
		id).Force()

	c.App.Database.Delete(job.TableName(), "id = $1", id).Force()

	c.JSONResponse(ctx, nil, fasthttp.StatusNoContent)
}
