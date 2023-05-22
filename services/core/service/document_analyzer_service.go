package service

import (
	"context"
	"core/clients"
	"core/clients/hugging_face"
	"core/configuration"
	"core/constants"
	"core/log"
	"core/model/response"
	"encoding/base64"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
)

type DocumentAnalyzerService interface {
	AnalyzeDocument(ctx context.Context, checksum string, question string) (response.LayoutLMAnswer, error)
}

type documentAnalyzerService struct {
	config            configuration.Configuration
	minioClient       clients.MinioClient
	huggingFaceClient hugging_face.DocumentQAClient
}

func (service documentAnalyzerService) AnalyzeDocument(ctx context.Context, checksum string, question string) (response.LayoutLMAnswer, error) {
	logger := log.NewLogger().WithFields(logrus.Fields{"checksum": checksum, "question": question})
	logger.Info("Analyzing document")
	wg := new(sync.WaitGroup)

	images := service.minioClient.ListObjects(ctx, service.config.Minio.Bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    fmt.Sprintf("%s/%s", checksum, constants.ImagePrefix),
	})

	var imagePaths []string
	for image := range images {
		if image.Err != nil {
			logger.WithError(image.Err).Error("Error listing images")
			return response.LayoutLMAnswer{}, image.Err
		}

		imagePaths = append(imagePaths, image.Key)
	}

	var answers []response.LayoutLMAnswer
	wg.Add(len(imagePaths))
	for _, imagePath := range imagePaths {
		go service.getAnswerFromInferenceAPI(ctx, answers, imagePath, question, wg)
	}
	wg.Wait()

	if len(answers) > 0 {
		return answers[0], nil
	}

	logger.Info("No answer found")
	return response.LayoutLMAnswer{}, constants.DocumentQANoAnswerFoundError
}

func (service documentAnalyzerService) getAnswerFromInferenceAPI(ctx context.Context, answers []response.LayoutLMAnswer, imagePath string, question string, wg *sync.WaitGroup) {
	logger := log.NewLogger()
	defer wg.Done()

	object, err := service.minioClient.GetObject(ctx, service.config.Minio.Bucket, imagePath)
	if err != nil {
		logger.Errorf("Error getting image from Minio: %v", err)
		return
	}
	defer object.Close()

	objectInfo, err := object.Stat()
	if err != nil || objectInfo.Size == 0 {
		logger.Errorf("Error getting image info from Minio: %v", err)
		return
	}

	logger.Infof("Reading %d bytes from Minio", objectInfo.Size)
	b := make([]byte, objectInfo.Size)
	for {
		_, err = object.Read(b)
		if err != nil {
			break
		}
	}

	answer, err := service.huggingFaceClient.Answer(ctx, question, base64.StdEncoding.EncodeToString(b))
	if err != nil {
		logger.WithError(err).Error("Error getting answer from inference API")
		return
	}

	answers = append(answers, response.LayoutLMAnswer{
		Score:  answer.Score,
		Answer: answer.Answer,
	})

	// sort answers in descending order by score
	sort.Slice(answers, func(i, j int) bool {
		return answers[i].Score > answers[j].Score
	})
}

func NewDocumentAnalyzerService(
	config configuration.Configuration,
	minioClient clients.MinioClient,
	huggingFaceClient hugging_face.DocumentQAClient,
) DocumentAnalyzerService {
	return documentAnalyzerService{
		config:            config,
		minioClient:       minioClient,
		huggingFaceClient: huggingFaceClient,
	}
}
