package initializers

import "github.com/neel1996/saul/service"

var (
	loginService service.UserLoginService
)

func InitializeServices() {
	loginService = service.NewUserLoginService(userRepository, firebaseClient)
}
