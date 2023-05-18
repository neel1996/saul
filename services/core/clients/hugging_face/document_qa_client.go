package hugging_face

import (
	"context"
	"core/clients"
	"core/configuration"
	"core/constants"
	"core/log"
	"core/model/request"
	"core/model/response"
	"os"
)

type DocumentQAClient interface {
	Answer(ctx context.Context, question string, imageBase64 string) (response.DocumentQAResponse, error)
}

type documentQAClient struct {
	client clients.HttpClient
	config configuration.Configuration
}

func (d documentQAClient) Answer(ctx context.Context, question string, imageBase64 string) (response.DocumentQAResponse, error) {
	logger := log.NewLogger(ctx)
	req := request.DocumentQARequest{
		Inputs: request.Inputs{
			Question: question,
			Image:    imageBase64,
		},
	}
	logger.Info("Invoking DocumentQA API")

	apiKey := os.Getenv("HUGGING_FACE_API_KEY")
	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
	}

	var answers []response.DocumentQAResponse
	err := d.client.
		WithHeaders(headers).
		WithBody(req).
		WithResponse(&answers).
		Post(d.config.HuggingFace.DocumentQA.Endpoint)

	if err != nil {
		logger.Errorf("Error invoking DocumentQA API: %v", err)
		return response.DocumentQAResponse{}, constants.ExternalApiError
	}

	logger.Info("Successfully invoked DocumentQA API")

	if len(answers) == 0 {
		logger.Info("No answer found")
		return response.DocumentQAResponse{}, constants.DocumentQANoAnswerFoundError
	}

	logger.Infof("Returning answer received from DocumentQA API")

	return answers[0], nil
}

func NewDocumentQAClient(config configuration.Configuration, client clients.HttpClient) DocumentQAClient {
	return documentQAClient{
		client,
		config,
	}
}
