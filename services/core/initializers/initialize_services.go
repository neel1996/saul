package initializers

import (
	"core/configuration"
	"core/kafka"
	"core/service"
)

var (
	loginService            service.UserLoginService
	documentUploadService   service.DocumentUploadService
	inferenceService        service.InferenceService
	documentAnalyzerService service.DocumentAnalyzerService
)

func InitializeServices(config configuration.Configuration) {
	loginService = service.NewUserLoginService(userRepository, firebaseClient)

	documentDetailsProducer := kafka.NewDocumentDetailsProducer(config)
	documentUploadService = service.NewDocumentUploadService(config, minioClient, documentDetailsProducer)
	inferenceService = service.NewInferenceService(config, minioClient, huggingFaceDocumentQAClient)
	documentAnalyzerService = service.NewDocumentAnalyzerService(config, minioClient, inferenceService)
}
