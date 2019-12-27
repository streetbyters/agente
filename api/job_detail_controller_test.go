package api

import (
	"fmt"
	"github.com/akdilsiz/agente/database/model"
	model2 "github.com/akdilsiz/agente/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type JobDetailControllerTest struct {
	*Suite
}

func (s JobDetailControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
	UserAuth(s.Suite)
}

func (s JobDetailControllerTest) Test_CreateJobDetailWithValidParams() {
	job := model.NewJob()
	job.SourceUserId.SetValid(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	detail := model.NewJobDetail()
	detail.Type = model2.NewRelease
	detail.Code = "job3"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	resp := s.JSON(Post, fmt.Sprintf("/api/v1/job/%d/detail", job.ID), detail)

	s.Equal(resp.Status, fasthttp.StatusCreated)

	data, _ := resp.Success.Data.(map[string]interface{})

	s.Greater(data["id"], float64(0))
	s.Equal(data["job_id"], float64(job.ID))
	s.Equal(data["type"], "new_release")
	s.Equal(data["code"], "job3")
	s.Equal(data["name"], "Test Job")
	s.Equal(data["detail"], "Test Job Detail")
	s.Equal(data["before"], true)
	s.Equal(data["before_jobs"], "job4,job5")
	s.Equal(data["after"], true)
	s.Equal(data["after_jobs"], "job1,job2")
	s.Equal(data["script"], `
		#!/bin/sh
		echo "Test Job running"
	`)
	s.NotNil(data["inserted_at"])
	defaultLogger.LogInfo("Create a job detail with valid params")
}

func (s JobDetailControllerTest) Test_Should_422Error_CreateJobDetailWithInvalidParams() {
	job := model.NewJob()
	job.SourceUserId.SetValid(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	detail := model.NewJobDetail()
	detail.Type = model2.NewRelease
	detail.Code = "jo"
	detail.Name = "Te"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	resp := s.JSON(Post, fmt.Sprintf("/api/v1/job/%d/detail", job.ID), detail)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a job detail with invalid params")
}

func (s JobDetailControllerTest) Test_Should_422Error_CreateJobDetailWithValidParamsIfCodeNotUnique() {
	job := model.NewJob()
	job.SourceUserId.SetValid(s.Auth.User.ID)
	err := s.API.App.Database.Insert(new(model.Job), job, "id", "inserted_at")
	s.Nil(err)

	detail := model.NewJobDetail()
	detail.JobID = job.ID
	detail.Type = model2.NewRelease
	detail.Code = "job4"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	err = s.API.App.Database.Insert(new(model.JobDetail), detail, "id")
	s.Nil(err)

	detail = model.NewJobDetail()
	detail.Type = model2.NewRelease
	detail.Code = "job4"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	resp := s.JSON(Post, fmt.Sprintf("/api/v1/job/%d/detail", job.ID), detail)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	defaultLogger.LogInfo("Should be 422 error create a job detail with valid params " +
		"if code has already been taken")
}

func (s JobDetailControllerTest) Test_Should_422Error_CreateJobDetailWithValidParamsIfRelationalError() {
	detail := model.NewJobDetail()
	detail.Type = model2.NewRelease
	detail.Code = "job3"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	resp := s.JSON(Post, fmt.Sprintf("/api/v1/job/999999999/detail"), detail)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)
	defaultLogger.LogInfo("Should be 422 error create a job detail with valid params" +
		"if relational error")
}

func (s JobDetailControllerTest) Test_Should_400Error_CreateJobDetailWithValidParams() {
	detail := model.NewJobDetail()
	detail.Type = model2.NewRelease
	detail.Code = "job3"
	detail.Name = "Test Job"
	detail.Detail.SetValid("Test Job Detail")
	detail.Before = true
	detail.BeforeJobs.SetValid("job4,job5")
	detail.After = true
	detail.AfterJobs.SetValid("job1,job2")
	detail.Script.SetValid(`
		#!/bin/sh
		echo "Test Job running"
	`)

	resp := s.JSON(Post, fmt.Sprintf("/api/v1/job/jobID/detail"), detail)

	s.Equal(resp.Status, fasthttp.StatusBadRequest)
	defaultLogger.LogInfo("Should be 400 error create a job detail with valid params")
}

func (s JobDetailControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_JobDetailController(t *testing.T) {
	s := JobDetailControllerTest{Suite: NewSuite()}
	Run(t, s)
}
