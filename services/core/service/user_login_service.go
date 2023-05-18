package service

//go:generate mockgen -destination=../mocks/mock_user_login_service.go -package=mocks -source=user_login_service.go

import (
	"context"
	"core/clients"
	"core/constants"
	"core/log"
	"core/model/request"
	"core/model/response"
	"core/repository"
	"github.com/sirupsen/logrus"
)

type UserLoginService interface {
	Login(ctx context.Context, userRequest request.UserLoginRequest) (response.UserLoginResponse, error)
}

type userValidationService struct {
	repository     repository.UserRepository
	firebaseClient clients.FirebaseClient
}

func (service userValidationService) Login(ctx context.Context, userRequest request.UserLoginRequest) (response.UserLoginResponse, error) {
	logger := log.NewLogger(ctx).WithFields(logrus.Fields{
		"method": "Login",
		"email":  userRequest.Email,
	})
	logger.Info("Logging in user")

	exist, err := service.repository.DoesUserExist(ctx, userRequest.Email)
	if err != nil {
		logger.Errorf("Error occurred while checking if user exists, error: %v", err)
		return response.UserLoginResponse{}, constants.UserLoginError
	}

	if !exist {
		logger.Info("Creating new user")
		err = service.repository.CreateUser(ctx, userRequest)
		if err != nil {
			logger.Errorf("Error occurred while creating new user, error: %v", err)
			return response.UserLoginResponse{}, constants.UserLoginError
		}
	}

	authToken, err := service.getToken(ctx, userRequest)
	if err != nil {
		return response.UserLoginResponse{}, constants.UserLoginError
	}

	return response.UserLoginResponse{
		AuthToken: authToken,
	}, nil
}

func (service userValidationService) getToken(ctx context.Context, userRequest request.UserLoginRequest) (string, error) {
	logger := log.NewLogger(ctx)
	logger.Info("Generating auth token for user")

	authToken, err := service.firebaseClient.GenerateAuthToken(ctx, userRequest.UserId)
	if err != nil {
		logger.Errorf("Error occurred while generating auth token, error: %v", err)
		return "", err
	}

	return authToken, nil
}

func NewUserLoginService(
	repository repository.UserRepository,
	firebaseClient clients.FirebaseClient,
) UserLoginService {
	return userValidationService{
		repository,
		firebaseClient,
	}
}
