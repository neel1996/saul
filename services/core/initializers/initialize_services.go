package initializers

import (
	"core/configuration"
	"core/kafka"
	"core/service"
)

var (
	loginService            service.UserLoginService
	documentUploadService   service.DocumentUploadService
	documentAnalyzerService service.DocumentAnalyzerService
)

func InitializeServices(config configuration.Configuration) {
	loginService = service.NewUserLoginService(userRepository, firebaseClient)

	documentDetailsProducer := kafka.NewDocumentDetailsProducer(config)
	documentUploadService = service.NewDocumentUploadService(config, minioClient, documentDetailsProducer)
	documentAnalyzerService = service.NewDocumentAnalyzerService(config, minioClient, huggingFaceDocumentQAClient)
}
