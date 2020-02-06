// Copyright 2019 StreetByters Community
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
	"github.com/fate-lovely/phi"
	"github.com/streetbyters/agente/database"
	model2 "github.com/streetbyters/agente/database/model"
	"github.com/streetbyters/agente/model"
	"github.com/valyala/fasthttp"
)

// FileController job files api controller
type FileController struct {
	Controller
	*API
}

// Index list all job files
func (c FileController) Index(ctx *fasthttp.RequestCtx) {
	paginate, _, _ := c.Paginate(ctx, "id", "inserted_at", "updated_at")

	upload := model2.NewFile()
	var uploads []model2.File
	c.App.Database.QueryWithModel(fmt.Sprintf("SELECT f.* FROM %s AS f "+
		"ORDER BY f.%s %s LIMIT $1 OFFSET $2",
		upload.TableName(), paginate.OrderField, paginate.OrderBy),
		&uploads,
		paginate.Limit,
		paginate.Offset)

	var count int64
	c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", upload.TableName()))

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       uploads,
		TotalCount: count,
	}, fasthttp.StatusOK)
}

// Show job file
func (c FileController) Show(ctx *fasthttp.RequestCtx) {
	var file model2.File
	c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", file.TableName()),
		&file,
		phi.URLParam(ctx, "fileID")).Force()

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: file,
	}, fasthttp.StatusOK)
}

// Create job file
func (c FileController) Create(ctx *fasthttp.RequestCtx) {
	file := model2.NewFile()
	c.JSONBody(ctx, &file)

	if file.NodeID <= 0 {
		file.NodeID = c.App.Node.ID
	}

	if errs, err := database.ValidateStruct(file); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	// TODO: file distribute all worker nodes

	var err error

	db := c.App.Database.Transaction(func(tx *database.Tx) error {
		if err := c.App.Database.Insert(new(model2.File), file, "id", "inserted_at", "updated_at"); err != nil {
			return err
		}

		log := model2.NewFileLog(file.ID)
		log.NodeID = file.NodeID
		log.Type = model.Insert
		log.Data = model2.FileLogData{
			Dir:  file.Dir,
			File: file.File,
			Type: file.Type,
		}

		c.App.Database.Insert(new(model2.FileLog), log, "id")
		return nil
	})

	err = db.Error

	if err != nil {
		if errs, err := database.ValidateConstraint(err, file); err != nil {
			c.JSONResponse(ctx, model.ResponseError{
				Errors: errs,
				Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
			}, fasthttp.StatusUnprocessableEntity)
		}
		return
	}

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: file,
	}, fasthttp.StatusCreated)
}

// Update job file
func (c FileController) Update(ctx *fasthttp.RequestCtx) {
	var file model2.File
	c.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", file.TableName()),
		&file,
		phi.URLParam(ctx, "fileID")).Force()

	fileRequest := model2.NewFile()
	c.JSONBody(ctx, &fileRequest)
	if fileRequest.NodeID <= 0 {
		fileRequest.NodeID = c.App.Node.ID
	}

	if errs, err := database.ValidateStruct(fileRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	var err error

	db := c.App.Database.Transaction(func(tx *database.Tx) error {
		if err := c.App.Database.Update(&file, fileRequest, nil, "id", "updated_at"); err != nil {
			return err
		}

		log := model2.NewFileLog(file.ID)
		log.NodeID = fileRequest.NodeID
		log.Type = model.Update
		log.Data = model2.FileLogData{
			Dir:  file.Dir,
			File: file.File,
			Type: file.Type,
		}

		c.App.Database.Insert(new(model2.FileLog), log, "id")
		return nil
	})
	err = db.Error

	if errs, err := database.ValidateConstraint(err, fileRequest); err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusOK)
}

// Delete job file
func (c FileController) Delete(ctx *fasthttp.RequestCtx) {
	file := model2.NewFile()
	c.App.Database.Delete(file.TableName(), "id = $1", phi.URLParam(ctx, "fileID")).Force()

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: nil,
	}, fasthttp.StatusNoContent)
}
