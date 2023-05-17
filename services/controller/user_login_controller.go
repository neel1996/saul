package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/neel1996/saul/constants"
	"github.com/neel1996/saul/log"
	"github.com/neel1996/saul/model/request"
	"github.com/neel1996/saul/service"
)

type UserLoginController interface {
	Login(ctx *gin.Context)
}

type userLoginController struct {
	service service.UserLoginService
}

func (controller userLoginController) Login(ctx *gin.Context) {
	logger := log.NewLogger(ctx)
	logger.Info("Logging in user")

	var userRequest request.UserLoginRequest

	err := ctx.ShouldBindJSON(&userRequest)
	if err != nil {
		logger.Errorf("Error occurred while binding request to struct, error: %v", err)
		ctx.JSON(constants.RequestValidationError.GetGinResponse())
		return
	}

	userResponse, err := controller.service.Login(ctx, userRequest)
	if err != nil {
		logger.Errorf("Error occurred while logging in user, error: %v", err)
		e, ok := err.(constants.Error)
		if !ok {
			ctx.JSON(500, err)
			return
		}
		ctx.JSON(e.GetGinResponse())
		return
	}

	logger.Info("Successfully logged in user")
	ctx.JSON(200, userResponse)
}

func NewUserLoginController(service service.UserLoginService) UserLoginController {
	return userLoginController{
		service,
	}
}
