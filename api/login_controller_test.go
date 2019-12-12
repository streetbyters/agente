package api

import (
	"fmt"
	model2 "github.com/akdilsiz/agente/database/model"
	"github.com/akdilsiz/agente/model"
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
	userModel := new(model2.User)

	_, err := s.API.App.Database.Insert(userModel, user)
	s.Nil(err)

	loginRequest := model.LoginRequest{
		ID:       "akdilsiz",
		Password: "123456",
	}

	resp := s.JSON(Post, "/api/v1/user/sign_in", loginRequest)

	fmt.Println(resp.Error)
	fmt.Println(resp.Status)
}

func (s LoginControllerTest) TearDownSuite() {
	TearDownSuite(s.Suite)
}

func Test_LoginController(t *testing.T) {
	s := LoginControllerTest{NewSuite()}
	Run(t, s)
}
