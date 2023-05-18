package controller

import (
	"bytes"
	"core/constants"
	"core/mocks"
	"core/model/request"
	"core/model/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type UserLoginControllerTestSuite struct {
	suite.Suite
	context              *gin.Context
	recorder             *httptest.ResponseRecorder
	mockController       *gomock.Controller
	mockUserLoginService *mocks.MockUserLoginService
	controller           UserLoginController
}

func TestUserLoginControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginControllerTestSuite))
}

func (suite *UserLoginControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockUserLoginService = mocks.NewMockUserLoginService(suite.mockController)

	suite.controller = NewUserLoginController(suite.mockUserLoginService)
}

func (suite *UserLoginControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *UserLoginControllerTestSuite) TestLogin_ShouldLoginUserAndReturnToken() {
	loginRequest := request.UserLoginRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "test",
	}
	requestBytes, err := json.Marshal(loginRequest)
	suite.Nil(err)

	suite.context.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBytes))

	expectedResponse := response.UserLoginResponse{AuthToken: "token"}

	suite.mockUserLoginService.EXPECT().
		Login(suite.context, loginRequest).
		Return(expectedResponse, nil).
		Times(1)

	suite.controller.Login(suite.context)

	var actualResponse response.UserLoginResponse
	err = json.Unmarshal(suite.recorder.Body.Bytes(), &actualResponse)
	suite.Nil(err)

	suite.Equal(200, suite.recorder.Code)
	suite.Equal(expectedResponse, actualResponse)
}

func (suite *UserLoginControllerTestSuite) TestLogin_WhenRequestBodyIsInvalid_ShouldReturnBadRequest() {
	suite.context.Request = httptest.NewRequest("POST", "/login", nil)

	suite.controller.Login(suite.context)

	suite.Equal(400, suite.recorder.Code)
}

func (suite *UserLoginControllerTestSuite) TestLogin_WhenServiceInvocationFails_ShouldReturnError() {
	loginRequest := request.UserLoginRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "test",
	}
	requestBytes, err := json.Marshal(loginRequest)
	suite.Nil(err)

	suite.context.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(requestBytes))

	suite.mockUserLoginService.EXPECT().
		Login(suite.context, loginRequest).
		Return(response.UserLoginResponse{}, constants.UserLoginError).
		Times(1)

	suite.controller.Login(suite.context)

	suite.Equal(500, suite.recorder.Code)
	suite.Equal(constants.UserLoginError.String(), suite.recorder.Body.String())
}
