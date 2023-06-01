package service

import (
	"context"
	"core/configuration"
	"core/constants"
	"core/mocks"
	"core/model/response"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DocumentAnalyzerServiceTestSuite struct {
	suite.Suite
	config               configuration.Configuration
	context              context.Context
	mockController       *gomock.Controller
	mockMinioClient      *mocks.MockMinioClient
	mockInferenceService *mocks.MockInferenceService
	service              DocumentAnalyzerService
}

func TestDocumentAnalyzerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentAnalyzerServiceTestSuite))
}

func (suite *DocumentAnalyzerServiceTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.config = configuration.Configuration{
		Minio: configuration.Minio{
			Bucket: "documents",
		},
	}
	suite.mockController = gomock.NewController(suite.T())
	suite.mockMinioClient = mocks.NewMockMinioClient(suite.mockController)
	suite.mockInferenceService = mocks.NewMockInferenceService(suite.mockController)

	suite.service = NewDocumentAnalyzerService(suite.config, suite.mockMinioClient, suite.mockInferenceService)
}

func (suite *DocumentAnalyzerServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DocumentAnalyzerServiceTestSuite) TestAnalyzeDocument_ShouldReturnMaximumScoredAnswer() {
	expectedAnswer := response.LayoutLMAnswer{
		Answer: "answer2",
		Score:  0.9,
	}

	suite.mockMinioClient.EXPECT().
		ListObjects(suite.context, "documents", minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    "checksum/images",
		}).
		DoAndReturn(func(ctx context.Context, bucket string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
			objectInfos := make(chan minio.ObjectInfo, 2)
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image1.jpg",
				Err: nil,
			}
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image2.jpg",
				Err: nil,
			}
			close(objectInfos)

			return objectInfos
		}).
		Times(1)

	suite.mockInferenceService.EXPECT().
		GetAnswersFromInferenceAPI(suite.context, "checksum/images/image2.jpg", "question", gomock.Any()).
		DoAndReturn(func(ctx context.Context, imagePath string, question string, responseChan chan response.LayoutLMAnswer) {
			responseChan <- response.LayoutLMAnswer{
				Answer: "answer1",
				Score:  0.8,
			}
		}).
		Times(1)
	suite.mockInferenceService.EXPECT().
		GetAnswersFromInferenceAPI(suite.context, "checksum/images/image1.jpg", "question", gomock.Any()).
		DoAndReturn(func(ctx context.Context, imagePath string, question string, responseChan chan response.LayoutLMAnswer) {
			responseChan <- response.LayoutLMAnswer{
				Answer: "answer2",
				Score:  0.9,
			}
		}).
		Times(1)

	answer, err := suite.service.AnalyzeDocument(suite.context, "checksum", "question")

	suite.Nil(err)
	suite.Equal(expectedAnswer, answer)
}

func (suite *DocumentAnalyzerServiceTestSuite) TestAnalyzeDocument_WhenListingObjectsFail_ShouldReturnError() {
	suite.mockMinioClient.EXPECT().
		ListObjects(suite.context, "documents", minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    "checksum/images",
		}).
		DoAndReturn(func(ctx context.Context, bucket string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
			objectInfos := make(chan minio.ObjectInfo, 2)
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image1.jpg",
				Err: errors.New("error"),
			}
			close(objectInfos)

			return objectInfos
		}).
		Times(1)

	_, err := suite.service.AnalyzeDocument(suite.context, "checksum", "question")

	suite.NotNil(err)
}

func (suite *DocumentAnalyzerServiceTestSuite) TestAnalyzeDocument_WhenInferenceApiInvocationFails_ShouldOmitTheFailedAnswer() {
	expectedAnswer := response.LayoutLMAnswer{
		Answer: "answer1",
		Score:  0.8,
	}

	suite.mockMinioClient.EXPECT().
		ListObjects(suite.context, "documents", minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    "checksum/images",
		}).
		DoAndReturn(func(ctx context.Context, bucket string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
			objectInfos := make(chan minio.ObjectInfo, 2)
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image1.jpg",
				Err: nil,
			}
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image2.jpg",
				Err: nil,
			}
			close(objectInfos)

			return objectInfos
		}).
		Times(1)

	suite.mockInferenceService.EXPECT().
		GetAnswersFromInferenceAPI(suite.context, "checksum/images/image2.jpg", "question", gomock.Any()).
		DoAndReturn(func(ctx context.Context, imagePath string, question string, responseChan chan response.LayoutLMAnswer) {
			responseChan <- response.LayoutLMAnswer{
				Answer: "answer1",
				Score:  0.8,
			}
		}).
		Times(1)
	suite.mockInferenceService.EXPECT().
		GetAnswersFromInferenceAPI(suite.context, "checksum/images/image1.jpg", "question", gomock.Any()).
		DoAndReturn(func(ctx context.Context, imagePath string, question string, responseChan chan response.LayoutLMAnswer) {
			responseChan <- response.LayoutLMAnswer{
				Answer: "answer2",
				Score:  0.9,
				Err:    errors.New("error"),
			}
		}).
		Times(1)

	answer, err := suite.service.AnalyzeDocument(suite.context, "checksum", "question")

	suite.Nil(err)
	suite.Equal(expectedAnswer, answer)
}

func (suite *DocumentAnalyzerServiceTestSuite) TestAnalyzeDocument_WhenNoAnswersQualify_ShouldReturnError() {
	suite.mockMinioClient.EXPECT().
		ListObjects(suite.context, "documents", minio.ListObjectsOptions{
			Recursive: true,
			Prefix:    "checksum/images",
		}).
		DoAndReturn(func(ctx context.Context, bucket string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
			objectInfos := make(chan minio.ObjectInfo, 2)
			objectInfos <- minio.ObjectInfo{
				Key: "checksum/images/image1.jpg",
				Err: nil,
			}
			close(objectInfos)

			return objectInfos
		}).
		Times(1)

	suite.mockInferenceService.EXPECT().
		GetAnswersFromInferenceAPI(suite.context, "checksum/images/image1.jpg", "question", gomock.Any()).
		DoAndReturn(func(ctx context.Context, imagePath string, question string, responseChan chan response.LayoutLMAnswer) {
			responseChan <- response.LayoutLMAnswer{
				Answer: "answer1",
				Score:  0.8,
				Err:    errors.New("error"),
			}
		}).
		Times(1)

	_, err := suite.service.AnalyzeDocument(suite.context, "checksum", "question")

	suite.NotNil(err)
	suite.Equal(constants.DocumentQANoAnswerFoundError, err)
}
