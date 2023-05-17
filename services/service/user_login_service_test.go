package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/mocks"
	"github.com/neel1996/saul/model/request"
	"github.com/neel1996/saul/model/response"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserLoginServiceTestSuite struct {
	suite.Suite
	context            context.Context
	mockController     *gomock.Controller
	mockUserRepository *mocks.MockUserRepository
	mockFirebaseClient *mocks.MockFirebaseClient
	service            UserLoginService
}

func TestUserLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginServiceTestSuite))
}

func (suite *UserLoginServiceTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.mockController)
	suite.mockFirebaseClient = mocks.NewMockFirebaseClient(suite.mockController)

	suite.service = NewUserLoginService(suite.mockUserRepository, suite.mockFirebaseClient)
}

func (suite *UserLoginServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserAlreadyExists_ShouldReturnAuthToken() {
	userRequest := request.UserRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "https://test.com",
	}

	expected := response.UserLoginResponse{
		AuthToken: "token",
	}

	suite.mockUserRepository.EXPECT().
		DoesUserExist(suite.context, userRequest.Email).
		Return(true, nil).
		Times(1)

	suite.mockFirebaseClient.EXPECT().
		GenerateAuthToken(suite.context, userRequest.UserId).
		Return("token", nil).
		Times(1)

	loginResponse, err := suite.service.Login(suite.context, userRequest)

	suite.Nil(err)
	suite.Equal(expected, loginResponse)
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserIsNew_ShouldCreatedUserAndReturnAuthToken() {
	userRequest := request.UserRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "https://test.com",
	}

	expected := response.UserLoginResponse{
		AuthToken: "token",
	}

	suite.mockUserRepository.EXPECT().
		DoesUserExist(suite.context, userRequest.Email).
		Return(false, nil).
		Times(1)

	suite.mockUserRepository.EXPECT().
		CreateUser(suite.context, userRequest).
		Return(nil).
		Times(1)

	suite.mockFirebaseClient.EXPECT().
		GenerateAuthToken(suite.context, userRequest.UserId).
		Return("token", nil).
		Times(1)

	loginResponse, err := suite.service.Login(suite.context, userRequest)

	suite.Nil(err)
	suite.Equal(expected, loginResponse)
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserExistenceCheckFails_ShouldReturnError() {
	userRequest := request.UserRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "https://test.com",
	}

	suite.mockUserRepository.EXPECT().
		DoesUserExist(suite.context, userRequest.Email).
		Return(false, errors.New("failed to check user")).
		Times(1)

	_, err := suite.service.Login(suite.context, userRequest)

	suite.NotNil(err)
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserCreationFails_ShouldReturnError() {
	userRequest := request.UserRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "https://test.com",
	}

	suite.mockUserRepository.EXPECT().
		DoesUserExist(suite.context, userRequest.Email).
		Return(false, nil).
		Times(1)

	suite.mockUserRepository.EXPECT().
		CreateUser(suite.context, userRequest).
		Return(errors.New("failed to create user")).
		Times(1)

	_, err := suite.service.Login(suite.context, userRequest)

	suite.NotNil(err)
	suite.Equal(constants.UserLoginError, err)
}
