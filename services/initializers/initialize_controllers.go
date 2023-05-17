package initializers

import "github.com/neel1996/saul/controller"

var (
	loginController controller.UserLoginController
)

func InitializeControllers() {
	loginController = controller.NewUserLoginController(loginService)
}
