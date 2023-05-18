package initializers

import (
	"core/controller"
)

var (
	loginController controller.UserLoginController
)

func InitializeControllers() {
	loginController = controller.NewUserLoginController(loginService)
}
