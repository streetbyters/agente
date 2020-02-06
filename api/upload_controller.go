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
	"github.com/forgolang/agente/model"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
)

// UploadController file upload api controller
type UploadController struct {
	Controller
	*API
}

// DirIndex directory list in lib path
func (c UploadController) DirIndex(ctx *fasthttp.RequestCtx) {
	var dirs []model.Dir
	i := 0
	filepath.Walk(c.App.Config.LibPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && i != 0 {
			dirs = append(dirs, model.Dir{
				Path:    path,
				Name:    info.Name(),
				Mode:    info.Mode(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}
		i++
		return nil
	})

	c.JSONResponse(ctx, model.ResponseSuccess{
		Data:       dirs,
		TotalCount: int64(len(dirs)),
	}, fasthttp.StatusOK)
}

// Create file upload method
func (c UploadController) Create(ctx *fasthttp.RequestCtx) {
	errs := make(map[string]string)
	file, err := ctx.FormFile("file")
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: nil,
			Detail: err.Error(),
		}, fasthttp.StatusBadRequest)
		return
	}

	dir := ctx.FormValue("dir")
	if string(dir) == "" {
		errs["dir"] = "is not nil"
		c.JSONResponse(ctx, model.ResponseError{
			Errors: errs,
			Detail: fasthttp.StatusMessage(fasthttp.StatusUnprocessableEntity),
		}, fasthttp.StatusUnprocessableEntity)
		return
	}

	err = fasthttp.SaveMultipartFile(file, filepath.Join(string(dir), file.Filename))
	if err != nil {
		c.JSONResponse(ctx, model.ResponseError{
			Errors: nil,
			Detail: err.Error(),
		}, fasthttp.StatusInternalServerError)
		return
	}

	// TODO: If the uploaded file is not used within 5 minutes, it needs to be deleted.

	resp := make(map[string]string)
	resp["dir"] = string(dir)
	resp["filename"] = file.Filename

	c.JSONResponse(ctx, model.ResponseSuccessOne{
		Data: resp,
	}, fasthttp.StatusCreated)
}
