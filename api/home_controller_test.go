package api

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type HomeControllerTest struct {
	*Suite
}

func (s *HomeControllerTest) Test_GetHome() {
	resp := s.JSON(Get, "/")

	s.Equal(resp.Status, 200)
	s.Equal(resp.Success.Data, "Agente")

	s.API.App.Logger.LogInfo("Success get home")
}

func Test_HomeController(t *testing.T) {
	s := &HomeControllerTest{NewSuite()}
	suite.Run(t, s)
}
