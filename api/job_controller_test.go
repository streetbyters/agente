package api

import (
	"github.com/akdilsiz/agente/database/model"
	"github.com/valyala/fasthttp"
	"testing"
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
	s.Equal(resp.Success.TotalCount, int64(50))
}

func (s JobControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_JobController(t *testing.T) {
	s := JobControllerTest{Suite: NewSuite()}
	Run(t, s)
}
