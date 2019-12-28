package api

import (
	model2 "github.com/akdilsiz/agente/database/model"
	"github.com/akdilsiz/agente/model"
	"github.com/valyala/fasthttp"
	"testing"
)

type LoginControllerTest struct {
	*Suite
}

func (s LoginControllerTest) SetupSuite() {
	SetupSuite(s.Suite)
}

func (s LoginControllerTest) Test_PostLoginWithValidParams() {
	user := model2.NewUser("123456")
	user.Username = "akdilsiz"
	user.Email = "akdilsiz@tecpor.com"
	user.NodeID = s.API.App.Node.ID
	userModel := new(model2.User)

	err := s.API.App.Database.Insert(userModel, user, "id")
	s.Nil(err)

	loginRequest := model.LoginRequest{
		ID:       "akdilsiz",
		Password: "123456",
	}

	resp := s.JSON(Post, "/api/v1/user/sign_in", loginRequest)

	s.Equal(resp.Status, fasthttp.StatusCreated)
	s.Equal(resp.Success.Data.(map[string]interface{})["user_id"], float64(user.ID))
	s.Equal(len(resp.Success.Data.(map[string]interface{})["passphrase"].(string)), 192)

	s.API.App.Logger.LogInfo("Successfully Post login with valid params")
}

func (s LoginControllerTest) Test_Should_422Error_PostLoginWithInvalidParams() {
	loginRequest := model.LoginRequest{
		ID: "akdilsiz",
	}

	resp := s.JSON(Post, "/api/v1/user/sign_in", loginRequest)

	s.Equal(resp.Status, fasthttp.StatusUnprocessableEntity)

	s.API.App.Logger.LogInfo("Should be 422 error post login with invalid params")
}

func (s LoginControllerTest) Test_Should_404Error_PostLoginWithValidParamsIfUserNotExists() {
	loginRequest := model.LoginRequest{
		ID:       "not_found",
		Password: "123456",
	}

	resp := s.JSON(Post, "/api/v1/user/sign_in", loginRequest)

	s.Equal(resp.Status, fasthttp.StatusNotFound)

	s.API.App.Logger.LogInfo("Should be 404 error post login with valid params" +
		"if user does not exists")
}

func (s LoginControllerTest) Test_Should_401Error_PostLoginWithValidParamsIfPasswordNotMatch() {
	user := model2.NewUser("123456789")
	user.Username = "akdilsiz2"
	user.Email = "akdilsiz2@tecpor.com"
	user.NodeID = s.API.App.Node.ID
	userModel := new(model2.User)

	err := s.API.App.Database.Insert(userModel, user, "id")
	s.Nil(err)

	loginRequest := model.LoginRequest{
		ID:       "akdilsiz2",
		Password: "12345",
	}

	resp := s.JSON(Post, "/api/v1/user/sign_in", loginRequest)

	s.Equal(resp.Status, fasthttp.StatusUnauthorized)

	s.API.App.Logger.LogInfo("Should be 404 error post login with valid params" +
		"if password does not match")
}

func (s LoginControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_LoginController(t *testing.T) {
	s := LoginControllerTest{NewSuite()}
	Run(t, s)
}
