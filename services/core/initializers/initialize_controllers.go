package initializers

import (
	"core/controller"
)

var (
	loginController          controller.UserLoginController
	documentUploadController controller.DocumentUploadController
	messageSocketController  controller.MessageSocketController
)

func InitializeControllers() {
	loginController = controller.NewUserLoginController(loginService)
	documentUploadController = controller.NewDocumentUploadController(documentUploadService)
	messageSocketController = controller.NewMessageSocketController(documentAnalyzerService)
}
