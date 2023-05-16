package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/mocks"
	"github.com/neel1996/saul/model/request"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserLoginServiceTestSuite struct {
	suite.Suite
	context            context.Context
	mockController     *gomock.Controller
	mockUserRepository *mocks.MockUserRepository
	service            UserLoginService
}

func TestUserLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserLoginServiceTestSuite))
}

func (suite *UserLoginServiceTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.mockController)

	suite.service = NewUserLoginService(suite.mockUserRepository)
}

func (suite *UserLoginServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserAlreadyExists_ShouldReturnExistingUserStatus() {
	userRequest := request.UserRequest{
		UserId: "123",
		Email:  "test@test.com",
		Name:   "Test",
		Avatar: "https://test.com",
	}

	suite.mockUserRepository.EXPECT().
		DoesUserExist(suite.context, userRequest.Email).
		Return(true, nil).
		Times(1)

	status, err := suite.service.Login(suite.context, userRequest)

	suite.Nil(err)
	suite.Equal(constants.ExistingUser, status)
}

func (suite *UserLoginServiceTestSuite) TestLogin_WhenUserIsNew_ShouldCreatedUserAndReturnNewUserStatus() {
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
		Return(nil).
		Times(1)

	status, err := suite.service.Login(suite.context, userRequest)

	suite.Nil(err)
	suite.Equal(constants.NewUser, status)
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

	status, err := suite.service.Login(suite.context, userRequest)

	suite.NotNil(err)
	suite.Equal(constants.NoStatus, status)
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

	status, err := suite.service.Login(suite.context, userRequest)

	suite.NotNil(err)
	suite.Equal(constants.UserLoginError, err)
	suite.Equal(constants.NoStatus, status)
}
