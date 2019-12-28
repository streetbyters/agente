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
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
)

// JobController user defined background job api controller
type JobController struct {
	Controller
	*API
}

// Index list all user defined background jobs
func (c JobController) Index(ctx *fasthttp.RequestCtx) {
	paginate, errs, err := c.Paginate(ctx, "id", "inserted_at")
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: err.Error(),
		}, fasthttp.StatusUnprocessableEntity)
	}

	jobDetail := model2.NewJobDetail()
	job := model2.NewJob()
	var jobs []model2.Job
	res := c.App.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM "+job.TableName()+
		" AS j ORDER BY j.%s %s", paginate.OrderField, paginate.OrderBy)+
		" LIMIT $1 OFFSET $2",
		&jobs,
		paginate.Limit,
		paginate.Offset)

	var jobIDs []int64
	for _, j := range jobs {
		jobIDs = append(jobIDs, j.ID)
	}

	var rJobs []model2.Job
	jobDetails := make([]model2.JobDetail, 0)
	details := make(map[int64]model2.JobDetail)
	if res.Error == nil {
		query, args, _ := sqlx.In(fmt.Sprintf("SELECT d.* FROM %s AS d"+
			" LEFT OUTER JOIN %s AS d2 ON d.job_id = d2.job_id AND d.id < d2.id"+
			" WHERE d2.id IS NULL AND d.job_id IN (?)", jobDetail.TableName(), jobDetail.TableName()),
			jobIDs)
		query = c.App.Database.DB.Rebind(query)
		c.App.Database.QueryWithModel(query, &jobDetails, args...)

		for _, d := range jobDetails {
			details[d.JobID] = d
		}

		for _, j := range jobs {
			if val, ok := details[j.ID]; ok {
				j.Detail = &val
			}
			rJobs = append(rJobs, j)
		}
	}

	var count int64
	err = c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", job.TableName()))

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       rJobs,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Show a user defined background job
func (c JobController) Show(ctx *fasthttp.RequestCtx) {
	var job model2.Job
	c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", job.TableName()),
		&job,
		phi.URLParam(ctx, "jobID")).Force()

	detail := new(model2.JobDetail)
	c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT d.* FROM %s AS d"+
		" LEFT OUTER JOIN %s AS d2 ON d.job_id = d2.job_id AND d.id < d2.id"+
		" WHERE d2.id IS NULL AND d.job_id = $1", detail.TableName(), detail.TableName()),
		detail, job.ID).Force()

	job.Detail = detail

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: job,
	}, fasthttp.StatusOK)
}

// Create user defined background job
func (c JobController) Create(ctx *fasthttp.RequestCtx) {
	job := model2.NewJob()
	c.JSONBody(ctx, &job)

	if job.NodeID <= 0 {
		job.NodeID = c.App.Node.ID
	}

	job.SourceUserId.SetValid(c.Auth.ID)

	c.App.Database.Insert(new(model2.Job), job, "id", "inserted_at")

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: job,
	}, fasthttp.StatusCreated)
}

// Delete user defined background job
func (c JobController) Delete(ctx *fasthttp.RequestCtx) {
	id := phi.URLParam(ctx, "jobID")

	job := model2.NewJob()
	c.App.Database.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", job.TableName()),
		id).Force()

	c.App.Database.Delete(job.TableName(), "id = $1", id).Force()

	c.JSONResponse(ctx, nil, fasthttp.StatusNoContent)
}
