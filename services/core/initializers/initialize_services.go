package initializers

import (
	"core/configuration"
	"core/kafka"
	"core/service"
)

var (
	loginService          service.UserLoginService
	documentUploadService service.DocumentUploadService
)

func InitializeServices(config configuration.Configuration) {
	loginService = service.NewUserLoginService(userRepository, firebaseClient)

	documentDetailsProducer := kafka.NewDocumentDetailsProducer(config)
	documentUploadService = service.NewDocumentUploadService(config, minioClient, documentDetailsProducer)
}
