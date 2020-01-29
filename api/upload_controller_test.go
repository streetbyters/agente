package api

import (
	"github.com/valyala/fasthttp"
	"path/filepath"
	"testing"
)

type UploadControllerTest struct {
	*Suite
}

func (s UploadControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s UploadControllerTest) Test_ListLibDirs() {
	response := s.JSON(Get, "/api/v1/upload/dir", nil)

	s.Equal(response.Status, fasthttp.StatusOK)
	data, _ := response.Success.Data.([]interface{})
	s.Equal(len(data), 3)

	defaultLogger.LogInfo("List all lib dirs")
}

func (s UploadControllerTest) Test_PostUploadFile() {
	file1 := filepath.Join(s.API.App.Config.Path, "files", "tests", "agente.png")

	body := make(map[string]interface{})
	body["file"] = file1
	body["dir"] = filepath.Join(s.API.App.Config.Path, "files", "tests", "upload")

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusCreated)
	data, _ := response.Success.Data.(map[string]interface{})
	s.Equal(data["dir"], body["dir"].(string))
	s.Equal(data["filename"], "agente.png")

	defaultLogger.LogInfo("Post upload file")
}

func (s UploadControllerTest) Test_Should_400Err_PostUploadFileIfFileNotValid() {
	file1 := filepath.Join(s.API.App.Config.Path, "files", "tests", "notfound.png")

	body := make(map[string]interface{})
	body["file"] = file1
	body["dir"] = filepath.Join(s.API.App.Config.Path, "files", "tests", "upload")

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusBadRequest)

	defaultLogger.LogInfo("Should 400 error post upload file if file is not valid")
}

func (s UploadControllerTest) Test_Should_422Err_PostUploadFileIfDirIsNil() {
	file1 := filepath.Join(s.API.App.Config.Path, "files", "tests", "agente.png")

	body := make(map[string]interface{})
	body["file"] = file1

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error post upload file if dir is nil")
}

func (s UploadControllerTest) Test_Should_500Error_PostUploadFileIfDirPermissionError() {
	file1 := filepath.Join(s.API.App.Config.Path, "files", "tests", "agente.png")

	body := make(map[string]interface{})
	body["file"] = file1
	body["dir"] = "/root"

	response := s.File(Post, "/api/v1/upload", body, "file")

	s.Equal(response.Status, fasthttp.StatusInternalServerError)

	defaultLogger.LogInfo("Post upload file")
}

func (s UploadControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_UploadController(t *testing.T) {
	s := UploadControllerTest{NewSuite()}
	Run(t, s)
}
