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
	"github.com/akdilsiz/agente/database/model"
	model2 "github.com/akdilsiz/agente/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

// FileLogController file log api controller
type FileLogController struct {
	Controller
	*API
}

// Index list all file logs
func (c FileLogController) Index(ctx *fasthttp.RequestCtx) {
	paginate, errs, err := c.Paginate(ctx, "id", "node_id", "type", "inserted_at")
	if err != nil {
		c.JSONResponse(ctx, model2.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		}, fasthttp.StatusBadRequest)
		return
	}

	fileID := phi.URLParam(ctx, "fileID")
	log := new(model.FileLog)
	var logs []model.FileLog
	var count int64

	if fileID != "" {
		c.App.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM %s "+
			"WHERE file_id = $1 ORDER BY $2 $3",
			log.TableName()),
			&logs,
			fileID,
			paginate.OrderBy,
			paginate.OrderField)
		c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s WHERE file_id = $1",
			log.TableName()), fileID)
	} else {
		c.App.Database.QueryWithModel(fmt.Sprintf("SELECT * FROM %s "+
			"ORDER BY $1 $2",
			log.TableName()),
			&logs,
			paginate.OrderBy,
			paginate.OrderField)
		c.App.Database.DB.Get(&count, fmt.Sprintf("SELECT count(*) FROM %s", log.TableName()))
	}

	c.JSONResponse(ctx, model2.ResponseSuccess{
		Data:       logs,
		TotalCount: count,
	}, fasthttp.StatusOK)
}
