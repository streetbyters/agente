package api

import (
	"fmt"
	"github.com/forgolang/agente/database/model"
	model2 "github.com/forgolang/agente/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type FileLogControllerTest struct {
	*Suite
}

func (s FileLogControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s FileLogControllerTest) Test_ListAllFileLogs() {
	for i := 0; i < 500; i++ {
		file := model.NewFile()
		file.NodeID = s.API.App.Node.ID
		file.File = "file.go"
		file.Dir = "/tmp"
		err := s.API.App.Database.Insert(new(model.File), file, "id")
		s.Nil(err)

		log := model.NewFileLog(0)
		log.FileID = file.ID
		log.Type = model2.Insert
		log.NodeID = s.API.App.Node.ID
		log.SourceUserID.SetValid(s.Auth.User.ID)
		err = s.API.App.Database.Insert(new(model.FileLog), log, "id")
		s.Nil(err)
	}

	response := s.JSON(Get, "/api/v1/file/log", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(499))
	defaultLogger.LogInfo("List all file logs")
}

func (s FileLogControllerTest) Test_ListAllFileLogsWithFileIDParam() {
	for i := 0; i < 100; i++ {
		file := model.NewFile()
		file.NodeID = s.API.App.Node.ID
		file.File = "file.go"
		file.Dir = "/tmp"
		err := s.API.App.Database.Insert(new(model.File), file, "id")
		s.Nil(err)

		log := model.NewFileLog(0)
		log.FileID = file.ID
		log.Type = model2.Insert
		log.NodeID = s.API.App.Node.ID
		log.SourceUserID.SetValid(s.Auth.User.ID)
		err = s.API.App.Database.Insert(new(model.FileLog), log, "id")
		s.Nil(err)
	}

	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.File = "file.go"
	file.Dir = "/tmp"
	err := s.API.App.Database.Insert(new(model.File), file, "id")
	s.Nil(err)

	for i := 0; i < 500; i++ {
		log := model.NewFileLog(0)
		log.FileID = file.ID
		log.Type = model2.Insert
		log.NodeID = s.API.App.Node.ID
		log.SourceUserID.SetValid(s.Auth.User.ID)
		err = s.API.App.Database.Insert(new(model.FileLog), log, "id")
		s.Nil(err)
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/file/%d/log", file.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Equal(response.Success.TotalCount, int64(500))
	defaultLogger.LogInfo("List all file logs with file identifier")
}

func (s FileLogControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_FileLogController(t *testing.T) {
	s := FileLogControllerTest{NewSuite()}
	Run(t, s)
}
