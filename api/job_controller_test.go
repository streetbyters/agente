package api

import (
	"fmt"
	"github.com/akdilsiz/agente/database/model"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type JobControllerTest struct {
	*Suite
}

func (s JobControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s JobControllerTest) Test_ListAllJobs() {
	for i := 0; i < 50; i++ {
		job := model.NewJob()
		job.SourceUserId.SetValid(s.Auth.User.ID)
		j := new(model.Job)
		err := s.API.App.Database.Insert(j, job, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	resp := s.JSON(Get, "/api/v1/job", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(50))

	defaultLogger.LogInfo("List all jobs")
}

func (s JobControllerTest) Test_ShowJobWithGivenIdentifier() {
	job := model.NewJob()
	job.SourceUserId.SetValid(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/job/%d", job.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(job.ID))
	s.Equal(data["source_user_id"], float64(s.Auth.User.ID))
	s.Equal(data["inserted_at"], job.InsertedAt.Format(time.RFC3339Nano))

	defaultLogger.LogInfo("Show a job with given identifier")
}

func (s JobControllerTest) Test_Should_404Error_ShowJobWithGivenIdentifierIfDoesNotExists() {
	resp := s.JSON(Get, "/api/v1/job/999999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error show a job with given identifier " +
		"if does not exists")
}

func (s JobControllerTest) Test_CreateJobWithValidParams() {
	resp := s.JSON(Post, "/api/v1/job", nil)

	s.Equal(resp.Status, fasthttp.StatusCreated)

	defaultLogger.LogInfo("Create a job with valid params")
}

func (s JobControllerTest) Test_DeleteJobWithGivenIdentifier() {
	job := model.NewJob()
	job.SourceUserId.SetValid(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	resp := s.JSON(Delete, fmt.Sprintf("/api/v1/job/%d", job.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusNoContent)

	defaultLogger.LogInfo("Delete a job with given identifier")
}

func (s JobControllerTest) Test_Should_404Error_DeleteJobWithGivenIdentifier() {
	resp := s.JSON(Delete, "/api/v1/job/999999", nil)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	defaultLogger.LogInfo("Should be 404 error delete a job with given identifier " +
		"if does not exists")
}

func (s JobControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_JobController(t *testing.T) {
	s := JobControllerTest{Suite: NewSuite()}
	Run(t, s)
}
