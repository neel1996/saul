package initializers

import "core/service"

var (
	loginService service.UserLoginService
)

func InitializeServices() {
	loginService = service.NewUserLoginService(userRepository, firebaseClient)
}
