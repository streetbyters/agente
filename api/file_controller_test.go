package api

import (
	"fmt"
	"github.com/akdilsiz/agente/database/model"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type FileControllerTest struct {
	*Suite
}

func (s FileControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s FileControllerTest) Test_ListAllFiles() {
	for i := 0; i < 50; i++ {
		file := model.NewFile()
		file.NodeID = s.API.App.Node.ID
		file.Dir = "/tmp"
		file.File = "file.go"
		err := s.API.App.Database.Insert(new(model.File), file, "id")
		s.Nil(err)
	}

	response := s.JSON(Get, "/api/v1/file", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 40)

	defaultLogger.LogInfo("List all files")
}

func (s FileControllerTest) Test_Should_400Err_ListAllFilesWithInvalidPaginateFields() {
	for i := 0; i < 50; i++ {
		file := model.NewFile()
		file.NodeID = s.API.App.Node.ID
		file.Dir = "/tmp"
		file.File = "file.go"
		err := s.API.App.Database.Insert(new(model.File), file, "id")
		s.Nil(err)
	}

	response := s.JSON(Get, "/api/v1/file?limit=limit&offset=offset", nil)

	s.Equal(response.Status, fasthttp.StatusBadRequest)
	defaultLogger.LogInfo("List all files")
}

func (s FileControllerTest) Test_ListAllFilesWithLimitAndOffsetParams() {
	for i := 0; i < 100; i++ {
		file := model.NewFile()
		file.NodeID = s.API.App.Node.ID
		file.Dir = "/tmp"
		file.File = "file.go"
		err := s.API.App.Database.Insert(new(model.File), file, "id")
		s.Nil(err)
	}

	response := s.JSON(Get, fmt.Sprintf("/api/v1/file?limit=%d&offset=%d", 30, 20), nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	s.Greater(response.Success.TotalCount, int64(49))
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 30)

	defaultLogger.LogInfo("List all files with limit and offset params")
}

func (s FileControllerTest) Test_ShowFileWithGivenIdentifier() {
	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.Dir = "/tmp"
	file.File = "file2.go"
	err := s.API.App.Database.Insert(new(model.File), file, "id", "inserted_at", "updated_at")
	s.Nil(err)

	response := s.JSON(Get, fmt.Sprintf("/api/v1/file/%d", file.ID), nil)

	s.Equal(response.Status, fasthttp.StatusOK)

	data, _ := response.Success.Data.(map[string]interface{})
	s.Equal(data["id"], float64(file.ID))
	s.Equal(data["node_id"], float64(s.API.App.Node.ID))
	s.Equal(data["dir"], "/tmp")
	s.Equal(data["file"], "file2.go")
	s.Equal(data["inserted_at"], file.InsertedAt.Format(time.RFC3339Nano))
	s.Equal(data["updated_at"], file.UpdatedAt.Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a file with given identifier")
}

func (s FileControllerTest) Test_Should_404Err_ShowFileWithGivenIdentifierIfNotExists() {
	response := s.JSON(Get, "/api/v1/file/999999999", nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a file with given identifier" +
		" if does not exists")
}

func (s FileControllerTest) Test_CreateFileWithValidParams() {
	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.Dir = "/tmp"
	file.File = "file_create.go"

	response := s.JSON(Post, "/api/v1/file", file)

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["node_id"], float64(s.API.App.Node.ID))
	s.Equal(data["dir"], "/tmp")
	s.Equal(data["file"], "file_create.go")
	s.NotNil(data["inserted_at"])
	s.NotNil(data["updated_at"])

	defaultLogger.LogInfo("Create a file with valid params")
}

func (s FileControllerTest) Test_CreateFileWithValidParamsAndNodeParams() {
	file := model.NewFile()
	file.Dir = "/tmp"
	file.File = "file_create2.go"

	response := s.JSON(Post, "/api/v1/file", file)

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Greater(data["id"], float64(0))
	s.Equal(data["node_id"], float64(s.API.App.Node.ID))
	s.Equal(data["dir"], "/tmp")
	s.Equal(data["file"], "file_create2.go")
	s.NotNil(data["inserted_at"])
	s.NotNil(data["updated_at"])

	defaultLogger.LogInfo("Create a file with valid params and node param")
}

func (s FileControllerTest) Test_Should_422Err_CreateFileWithInvalidParams() {
	file := model.NewFile()
	file.Dir = "/tmp"

	response := s.JSON(Post, "/api/v1/file", file)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a file with invalid params")
}

func (s FileControllerTest) Test_UpdateFileWithGivenIdentifierAndValidParams() {
	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.Dir = "/tmp"
	file.File = "file_update.go"
	err := s.API.App.Database.Insert(new(model.File), file, "id", "inserted_at", "updated_at")
	s.Nil(err)

	fileRequest := model.NewFile()
	fileRequest.Dir = "/tmp"
	fileRequest.File = "file_update2.go"

	response := s.JSON(Put, fmt.Sprintf("/api/v1/file/%d", file.ID), fileRequest)

	s.Equal(response.Status, fasthttp.StatusOK)

	defaultLogger.LogInfo("Update a file with given identifier and valid params")
}

func (s FileControllerTest) Test_Shoul_404Error_UpdateFileWithGivenIdentifierAndValidParamsIfNotExists() {
	fileRequest := model.NewFile()
	fileRequest.Dir = "/tmp"
	fileRequest.File = "file_update2.go"

	response := s.JSON(Put, fmt.Sprintf("/api/v1/file/fileID"), fileRequest)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error update a file with given identifier " +
		"and valid params if does not exists")
}

func (s FileControllerTest) Test_Should_422Error_UpdateFileWithGivenIdentifierAndInvalidParams() {
	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.Dir = "/tmp"
	file.File = "file_update.go"
	err := s.API.App.Database.Insert(new(model.File), file, "id", "inserted_at", "updated_at")
	s.Nil(err)

	fileRequest := model.NewFile()
	fileRequest.Dir = "/tmp"

	response := s.JSON(Put, fmt.Sprintf("/api/v1/file/%d", file.ID), fileRequest)

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error update a file with given identifier " +
		"and valid params")
}

func (s FileControllerTest) Test_DeleteFileWithGivenIdentifier() {
	file := model.NewFile()
	file.NodeID = s.API.App.Node.ID
	file.Dir = "/tmp"
	file.File = "file_delete.go"
	err := s.API.App.Database.Insert(new(model.File), file, "id", "inserted_at", "updated_at")
	s.Nil(err)

	response := s.JSON(Delete, fmt.Sprintf("/api/v1/file/%d", file.ID), nil)

	s.Equal(response.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete a file with given identifier")

	file2 := model.NewFile()
	result := s.API.App.Database.QueryRowWithModel(fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		file2.TableName()), file2, file.ID)
	s.NotNil(result.Error)

	defaultLogger.LogInfo("Delete a file with given identifier")
}

func (s FileControllerTest) Test_Should_404Err_DeleteFileWithGivenIdentifierIfNotExists() {
	response := s.JSON(Delete, "/api/v1/file/fileID", nil)

	s.Equal(response.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error delete a file with given identifier " +
		"if does not exists")
}

func (s FileControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_FileController(t *testing.T) {
	s := FileControllerTest{NewSuite()}
	Run(t, s)
}
