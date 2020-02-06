// Copyright 2019 Street Byters Community
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
	"github.com/forgolang/agente/database"
	"github.com/forgolang/agente/database/model"
	model2 "github.com/forgolang/agente/model"
	"github.com/forgolang/agente/utils"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// JobDetailController user defined background job detail api controller
type JobDetailController struct {
	Controller
	*API
}

// Create detail for user defined background job
func (c JobDetailController) Create(ctx *fasthttp.RequestCtx) {
	jobID, notExists := utils.ParseInt(phi.URLParam(ctx, "jobID"), 10, 64)
	if notExists {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: nil,
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	jobDetail := model.NewJobDetail()
	c.API.JSONBody(ctx, &jobDetail)
	jobDetail.JobID = jobID
	if jobDetail.NodeID <= 0 {
		jobDetail.NodeID = c.App.Node.ID
	}
	jobDetail.SourceUserID.SetValid(c.Auth.ID)

	if errs, err := database.ValidateStruct(jobDetail); err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	if res := c.App.Database.QueryRow(fmt.Sprintf("SELECT d.* FROM %s AS d"+
		" LEFT OUTER JOIN %s AS d2 ON d.job_id = d2.job_id AND d.id < d2.id"+
		" WHERE d2.id IS NULL AND d.code = $1", jobDetail.TableName(), jobDetail.TableName()), jobDetail.Code); len(res.Rows) > 0 {
		errs := make(map[string]string)
		errs["code"] = "has been already taken"
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	if err := c.App.Database.Insert(new(model.JobDetail),
		jobDetail,
		"id", "inserted_at"); err != nil {
		errs, err := database.ValidateConstraint(err, jobDetail)
		if err != nil {
			c.JSONResponse(ctx, model2.ResponseError{
				Errors: errs,
				Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
			}, fasthttp.StatusUnprocessableEntity)
			return
		}
	}

	c.JSONResponse(ctx, model2.ResponseSuccessOne{
		Data: jobDetail,
	}, fasthttp.StatusCreated)
}
