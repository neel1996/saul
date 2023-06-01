package service

import (
	"context"
	"core/configuration"
	"core/mocks"
	"core/model/response"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InferenceServiceTestSuite struct {
	suite.Suite
	context              context.Context
	config               configuration.Configuration
	mockController       *gomock.Controller
	mockMinioClient      *mocks.MockMinioClient
	mockDocumentQAClient *mocks.MockDocumentQAClient
	service              InferenceService
}

func TestInferenceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InferenceServiceTestSuite))
}

func (suite *InferenceServiceTestSuite) SetupTest() {
	suite.context = context.Background()
	suite.config = configuration.Configuration{
		Minio: configuration.Minio{
			Bucket: "documents",
		},
	}
	suite.mockController = gomock.NewController(suite.T())
	suite.mockMinioClient = mocks.NewMockMinioClient(suite.mockController)
	suite.mockDocumentQAClient = mocks.NewMockDocumentQAClient(suite.mockController)

	suite.service = NewInferenceService(suite.config, suite.mockMinioClient, suite.mockDocumentQAClient)
}

func (suite *InferenceServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *InferenceServiceTestSuite) TestGetAnswersFromInferenceAPI_ShouldReturnAnswersForTheDocumentAndQuestion() {
	testChan := make(chan response.LayoutLMAnswer, 1)
	b := []byte("test")

	suite.mockMinioClient.EXPECT().
		GetObject(suite.context, "documents", "image1.jpg").
		Return(b, nil).
		Times(1)

	suite.mockDocumentQAClient.EXPECT().
		Answer(suite.context, "question1", "dGVzdA==").
		Return(response.DocumentQAResponse{
			Score:  0.9,
			Answer: "new answer",
		}, nil).
		Times(1)

	suite.service.GetAnswersFromInferenceAPI(suite.context, "image1.jpg", "question1", testChan)
	close(testChan)

	for answer := range testChan {
		suite.Equal("new answer", answer.Answer)
		suite.Equal(0.9, answer.Score)
		suite.Nil(answer.Err)
	}
}

func (suite *InferenceServiceTestSuite) TestGetAnswersFromInferenceAPI_WhenGettingObjectFromMinioFails_ShouldReturnError() {
	testChan := make(chan response.LayoutLMAnswer, 1)

	suite.mockMinioClient.EXPECT().
		GetObject(suite.context, "documents", "image1.jpg").
		Return(nil, errors.New("error")).
		Times(1)

	suite.service.GetAnswersFromInferenceAPI(suite.context, "image1.jpg", "question1", testChan)
	close(testChan)

	for answer := range testChan {
		suite.NotNil(answer.Err)
	}
}

func (suite *InferenceServiceTestSuite) TestGetAnswersFromInferenceAPI_WhenInvokingHuggingFaceApiFails_ShouldReturnError() {
	testChan := make(chan response.LayoutLMAnswer, 1)
	b := []byte("test")

	suite.mockMinioClient.EXPECT().
		GetObject(suite.context, "documents", "image1.jpg").
		Return(b, nil).
		Times(1)

	suite.mockDocumentQAClient.EXPECT().
		Answer(suite.context, "question1", "dGVzdA==").
		Return(response.DocumentQAResponse{}, errors.New("error")).
		Times(1)

	suite.service.GetAnswersFromInferenceAPI(suite.context, "image1.jpg", "question1", testChan)
	close(testChan)

	for answer := range testChan {
		suite.NotNil(answer.Err)
	}
}
