package service

//go:generate mockgen -destination=../mocks/mock_inference_service.go -source=inference_service.go -package=mocks

import (
	"context"
	"core/clients"
	"core/clients/hugging_face"
	"core/configuration"
	"core/log"
	"core/model/response"
	"encoding/base64"
)

type InferenceService interface {
	GetAnswersFromInferenceAPI(ctx context.Context, imagePath string, question string, channel chan response.LayoutLMAnswer)
}

type inferenceService struct {
	config            configuration.Configuration
	minioClient       clients.MinioClient
	huggingFaceClient hugging_face.DocumentQAClient
}

func (service inferenceService) GetAnswersFromInferenceAPI(ctx context.Context, imagePath string, question string, channel chan response.LayoutLMAnswer) {
	logger := log.NewLogger(ctx)
	logger.Info("Getting answer from inference API")

	object, err := service.minioClient.GetObject(ctx, service.config.Minio.Bucket, imagePath)
	if err != nil {
		logger.Errorf("Error getting image from Minio: %v", err)
		channel <- response.LayoutLMAnswer{
			Err: err,
		}
		return
	}

	logger.Infof("Getting answer from inference API for image %s", imagePath)
	answer, err := service.huggingFaceClient.Answer(ctx, question, base64.StdEncoding.EncodeToString(object))
	if err != nil {
		logger.Errorf("Error getting answer from inference API for image %s: %v", imagePath, err)
		channel <- response.LayoutLMAnswer{
			Err: err,
		}
		return
	}

	logger.Infof("Answer for image %s: %s", imagePath, answer.Answer)
	channel <- response.LayoutLMAnswer{
		Score:  answer.Score,
		Answer: answer.Answer,
	}
}

func NewInferenceService(
	config configuration.Configuration,
	minioClient clients.MinioClient,
	huggingFaceClient hugging_face.DocumentQAClient,
) InferenceService {
	return inferenceService{
		config:            config,
		minioClient:       minioClient,
		huggingFaceClient: huggingFaceClient,
	}
}
