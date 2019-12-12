package api

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type HomeControllerTest struct {
	*Suite
}

func (s *HomeControllerTest) Test_GetHome() {
	resp := s.JSON(Get, "/", nil)

	s.Equal(resp.Status, 200)
	s.Equal(resp.Success.Data, "Agente")

	s.API.App.Logger.LogInfo("Success get home")
}

func (s *HomeControllerTest) Test_GetHomeWithAuthRequset() {
	s.Auth.Token = "token"
	resp := s.JSON(Get, "/", nil)
	s.Auth.Token = ""

	s.Equal(resp.Status, 200)
	s.Equal(resp.Success.Data, "Agente")

	s.API.App.Logger.LogInfo("Success get home")
}

func (s *HomeControllerTest) Test_OptionsHome() {
	resp := s.JSON(Options, "/", nil)

	s.Equal(resp.Status, 204)
	s.API.App.Logger.LogInfo("Success options home")
}

func (s *HomeControllerTest) Test_404Home() {
	resp := s.JSON(Get, "/home", nil)

	s.Equal(resp.Status, 404)

	s.API.App.Logger.LogInfo("Should be 404 get home request if not found page")
}

func Test_HomeController(t *testing.T) {
	s := &HomeControllerTest{NewSuite()}
	suite.Run(t, s)
}
