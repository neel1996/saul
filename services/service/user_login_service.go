package service

import (
	"context"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/log"
	"github.com/neel1996/saul/model/request"
	"github.com/neel1996/saul/repository"
)

type UserLoginService interface {
	Login(ctx context.Context, userRequest request.UserRequest) (constants.UserStatus, error)
}

type userValidationService struct {
	repository repository.UserRepository
}

func (service userValidationService) Login(ctx context.Context, userRequest request.UserRequest) (constants.UserStatus, error) {
	logger := log.NewLogger(ctx)
	logger.Infof("Validating user with email: %s", userRequest.Email)

	exist, err := service.repository.DoesUserExist(ctx, userRequest.Email)
	if err != nil {
		logger.Errorf("Error occurred while checking if user with email: %s exists, error: %v", userRequest.Email, err)
		return constants.NoStatus, constants.UserLoginError
	}

	if exist {
		logger.Infof("User with email: %s already exists", userRequest.Email)
		return constants.ExistingUser, nil
	}

	logger.Infof("Creating new user with email: %s", userRequest.Email)
	err = service.repository.CreateUser(ctx, userRequest)
	if err != nil {
		logger.Errorf("Error occurred while creating new user with email: %s, error: %v", userRequest.Email, err)
		return constants.NoStatus, constants.UserLoginError
	}

	return constants.NewUser, nil
}

func NewUserLoginService(repository repository.UserRepository) UserLoginService {
	return userValidationService{
		repository,
	}
}
