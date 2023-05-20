package initializers

import (
	"core/controller"
)

var (
	loginController          controller.UserLoginController
	documentUploadController controller.DocumentUploadController
)

func InitializeControllers() {
	loginController = controller.NewUserLoginController(loginService)
	documentUploadController = controller.NewDocumentUploadController(documentUploadService)
}
