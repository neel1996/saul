package service

import (
	"context"
	"core/clients"
	"core/configuration"
	"core/constants"
	"core/log"
	"core/model/response"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"sort"
)

type DocumentAnalyzerService interface {
	AnalyzeDocument(ctx context.Context, checksum string, question string) (response.LayoutLMAnswer, error)
}

type documentAnalyzerService struct {
	config           configuration.Configuration
	minioClient      clients.MinioClient
	inferenceService InferenceService
}

func (service documentAnalyzerService) AnalyzeDocument(ctx context.Context, checksum string, question string) (response.LayoutLMAnswer, error) {
	var answers []response.LayoutLMAnswer

	logger := log.NewLogger().WithFields(logrus.Fields{"checksum": checksum, "question": question})
	logger.Info("Analyzing document")

	imagePaths, err := service.imagePaths(ctx, checksum)
	if err != nil {
		return response.LayoutLMAnswer{}, err
	}

	channel := make(chan response.LayoutLMAnswer, len(imagePaths))
	for _, imagePath := range imagePaths {
		go service.inferenceService.GetAnswersFromInferenceAPI(ctx, imagePath, question, channel)
	}

	for i := 0; i < len(imagePaths); i++ {
		answer := <-channel
		if answer.Err != nil {
			logger.WithError(answer.Err).Error("Error getting answer from inference API")
			continue
		}

		if answer.Answer != "" {
			answers = append(answers, answer)
		}
	}

	if len(answers) == 0 {
		logger.Info("No answer found")
		return response.LayoutLMAnswer{}, constants.DocumentQANoAnswerFoundError
	}

	return service.accurateAnswer(answers), nil
}

func (service documentAnalyzerService) imagePaths(ctx context.Context, checksum string) ([]string, error) {
	logger := log.NewLogger()
	logger.Info("Fetching image list from Minio")

	images := service.minioClient.ListObjects(ctx, service.config.Minio.Bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    fmt.Sprintf("%s/%s", checksum, constants.ImagePrefix),
	})

	var imagePaths []string
	for image := range images {
		if image.Err != nil {
			logger.WithError(image.Err).Error("Error listing images")
			return nil, image.Err
		}

		imagePaths = append(imagePaths, image.Key)
	}

	logger.Infof("Found %d images", len(imagePaths))
	return imagePaths, nil
}

func (service documentAnalyzerService) accurateAnswer(answers []response.LayoutLMAnswer) response.LayoutLMAnswer {
	logger := log.NewLogger()
	logger.Info("Picking most accurate answer")

	// sort answers in descending order by score
	sort.Slice(answers, func(i, j int) bool {
		return answers[i].Score > answers[j].Score
	})

	pickedAnswer := answers[0]

	logger.Infof("Found answer: %s", pickedAnswer.Answer)
	return pickedAnswer
}

func NewDocumentAnalyzerService(
	config configuration.Configuration,
	minioClient clients.MinioClient,
	inferenceService InferenceService,
) DocumentAnalyzerService {
	return documentAnalyzerService{
		config:           config,
		minioClient:      minioClient,
		inferenceService: inferenceService,
	}
}
