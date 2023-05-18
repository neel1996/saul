package hugging_face

import (
	"context"
	"core/clients"
	"core/configuration"
	"core/constants"
	"core/mocks"
	"core/model/request"
	"core/model/response"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type DocumentQAClientTestSuite struct {
	suite.Suite
	context        context.Context
	config         configuration.Configuration
	client         DocumentQAClient
	mockController *gomock.Controller
	mockHttpClient *mocks.MockHttpClient
}

func TestDocumentQAClientTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentQAClientTestSuite))
}

func (suite *DocumentQAClientTestSuite) SetupTest() {
	suite.config = configuration.Configuration{
		HuggingFace: configuration.HuggingFace{
			DocumentQA: configuration.DocumentQA{
				Endpoint: "http://localhost:5000",
			},
		},
	}
	suite.context = context.Background()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockHttpClient = mocks.NewMockHttpClient(suite.mockController)

	suite.client = NewDocumentQAClient(suite.config, suite.mockHttpClient)
}

func (suite *DocumentQAClientTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DocumentQAClientTestSuite) TestAnswer_ShouldReturnAnswerWithScore() {
	req := request.DocumentQARequest{
		Inputs: request.Inputs{
			Question: "question",
			Image:    "image",
		},
	}

	_ = os.Setenv("HUGGING_FACE_API_KEY", "TOKEN")
	headers := map[string]string{
		"Authorization": "Bearer TOKEN",
	}

	answers := []response.DocumentQAResponse{
		{
			Answer: "answer",
			Score:  0.9,
		},
	}
	expected := response.DocumentQAResponse{
		Answer: "answer",
		Score:  0.9,
	}

	suite.mockHttpClient.EXPECT().WithBody(req).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithHeaders(headers).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithResponse(gomock.Any()).DoAndReturn(func(response *[]response.DocumentQAResponse) clients.HttpClient {
		*response = answers
		return suite.mockHttpClient
	}).Times(1)
	suite.mockHttpClient.EXPECT().Post(suite.config.HuggingFace.DocumentQA.Endpoint).Return(nil).Times(1)

	answer, err := suite.client.Answer(suite.context, "question", "image")

	suite.Nil(err)
	suite.Equal(expected, answer)
}

func (suite *DocumentQAClientTestSuite) TestAnswer_WhenApiCallFails_ShouldReturnError() {
	req := request.DocumentQARequest{
		Inputs: request.Inputs{
			Question: "question",
			Image:    "image",
		},
	}

	_ = os.Setenv("HUGGING_FACE_API_KEY", "TOKEN")
	headers := map[string]string{
		"Authorization": "Bearer TOKEN",
	}

	suite.mockHttpClient.EXPECT().WithBody(req).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithHeaders(headers).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithResponse(gomock.Any()).DoAndReturn(func(response *[]response.DocumentQAResponse) clients.HttpClient {
		*response = nil
		return suite.mockHttpClient
	})
	suite.mockHttpClient.EXPECT().Post(suite.config.HuggingFace.DocumentQA.Endpoint).Return(errors.New("error")).Times(1)

	_, err := suite.client.Answer(suite.context, "question", "image")

	suite.NotNil(err)
	suite.Equal(constants.ExternalApiError, err)
}

func (suite *DocumentQAClientTestSuite) TestAnswer_WhenNoAnswerFound_ShouldReturnError() {
	req := request.DocumentQARequest{
		Inputs: request.Inputs{
			Question: "question",
			Image:    "image",
		},
	}

	_ = os.Setenv("HUGGING_FACE_API_KEY", "TOKEN")
	headers := map[string]string{
		"Authorization": "Bearer TOKEN",
	}

	suite.mockHttpClient.EXPECT().WithBody(req).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithHeaders(headers).Return(suite.mockHttpClient).Times(1)
	suite.mockHttpClient.EXPECT().WithResponse(gomock.Any()).DoAndReturn(func(response *[]response.DocumentQAResponse) clients.HttpClient {
		*response = nil
		return suite.mockHttpClient
	})
	suite.mockHttpClient.EXPECT().Post(suite.config.HuggingFace.DocumentQA.Endpoint).Return(nil).Times(1)

	_, err := suite.client.Answer(suite.context, "question", "image")

	suite.NotNil(err)
	suite.Equal(constants.DocumentQANoAnswerFoundError, err)
}
