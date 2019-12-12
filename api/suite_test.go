package api

import "testing"

func Test_Run(t *testing.T) {
	suite := NewSuite()

	Run(t, suite)
	suite.API.App.Logger.LogInfo("Suite Run")
}

func Test_SetupSuite(t *testing.T) {
	suite := NewSuite()

	SetupSuite(suite)
	suite.API.App.Logger.LogInfo("Run SetupSuite")
}

func Test_TearDownSuite(t *testing.T) {
	suite := NewSuite()

	TearDownSuite(suite)
	suite.API.App.Logger.LogInfo("Run TearDownSuite")
}
