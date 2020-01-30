package api

import (
	"fmt"
	"github.com/forgolang/agente/database/model"
	model2 "github.com/forgolang/agente/model"
	"github.com/valyala/fasthttp"
	"strconv"
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
		job.NodeID = s.API.App.Node.ID
		job.SourceUserID.SetValid(s.Auth.User.ID)
		j := new(model.Job)
		err := s.API.App.Database.Insert(j, job, "id")
		s.Nil(err)
		if err != nil {
			break
		}
		jobDetail := model.NewJobDetail()
		jobDetail.JobID = job.ID
		jobDetail.Code = "job" + strconv.Itoa(i)
		jobDetail.Name = "jobName"
		jobDetail.Type = model2.NewRelease
		jobDetail.NodeID = s.API.App.Node.ID
		err = s.API.App.Database.Insert(new(model.JobDetail), jobDetail, "id")
		s.Nil(err)
		if err != nil {
			break
		}
	}

	resp := s.JSON(Get, "/api/v1/job", nil)

	s.Equal(resp.Status, fasthttp.StatusOK)
	s.Greater(resp.Success.TotalCount, int64(49))
	data, _ := resp.Success.Data.([]interface{})

	s.Equal(len(data), 40)
	s.NotNil(data[1].(map[string]interface{})["detail"])

	defaultLogger.LogInfo("List all jobs")
}

func (s JobControllerTest) Test_ShowJobWithGivenIdentifier() {
	job := model.NewJob()
	job.SourceUserID.SetValid(s.Auth.User.ID)
	job.NodeID = s.API.App.Node.ID
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	jobDetail := model.NewJobDetail()
	jobDetail.NodeID = s.API.App.Node.ID
	jobDetail.JobID = job.ID
	jobDetail.Code = "job2"
	jobDetail.Name = "jobName"
	jobDetail.Type = model2.NewRelease
	err = s.API.App.Database.Insert(new(model.JobDetail), jobDetail, "id")

	resp := s.JSON(Get, fmt.Sprintf("/api/v1/job/%d", job.ID), nil)

	s.Equal(resp.Status, fasthttp.StatusOK)

	data := resp.Success.Data.(map[string]interface{})

	s.Equal(data["id"], float64(job.ID))
	s.Equal(data["node_id"], float64(s.API.App.Node.ID))
	s.Equal(data["source_user_id"], float64(s.Auth.User.ID))
	detail := data["detail"].(map[string]interface{})
	s.Equal(detail["id"], float64(jobDetail.ID))
	s.Equal(detail["job_id"], float64(job.ID))
	s.Equal(detail["code"], "job2")
	s.Equal(detail["name"], "jobName")
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
	job.SourceUserID.SetValid(s.Auth.User.ID)
	job.NodeID = s.API.App.Node.ID
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
