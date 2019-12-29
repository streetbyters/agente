package api

import (
	"fmt"
	"github.com/akdilsiz/agente/database/model"
	model2 "github.com/akdilsiz/agente/model"
	"github.com/fate-lovely/phi"
	"github.com/valyala/fasthttp"
)

type FileLogController struct {
	Controller
	*API
}

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
