package clients

//go:generate mockgen -source=http_client.go -destination=../mocks/mock_http_client.go -package=mocks

import (
	"core/log"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
)

type HttpClient interface {
	Get(url string) error
	Post(url string) error
	WithHeader(key string, value string) HttpClient
	WithBody(body interface{}) HttpClient
	WithQueryParameters(parameters map[string]string) HttpClient
	WithQueryParameter(key string, value string) HttpClient
	WithHeaders(parameters map[string]string) HttpClient
	WithResponse(response interface{}) HttpClient
}

var logger = log.NewLogger()

type httpClient struct {
	request  *resty.Request
	response interface{}
}

func (h httpClient) Get(url string) error {
	response, err := h.request.Get(url)
	if err != nil {
		return err
	}

	if response.IsError() {
		err, ok := response.Error().(error)
		if ok {
			return err
		}
		return nil
	}

	return h.bindResponse(err, response)
}

func (h httpClient) Post(url string) error {
	response, err := h.request.Post(url)
	if err != nil {
		return err
	}

	if response.IsError() {
		err, ok := response.Error().(error)
		if ok {
			return err
		}
		return nil
	}

	return h.bindResponse(err, response)
}

func (h httpClient) WithHeader(key string, value string) HttpClient {
	h.request.SetHeader(key, value)

	return h
}

func (h httpClient) WithBody(body interface{}) HttpClient {
	if body == nil {
		return h
	}

	err := validator.New().Struct(body)
	if err != nil {
		logger.Errorf("Error validating request body: %v", err)
		return h
	}

	h.request.SetBody(body)

	return h
}

func (h httpClient) WithQueryParameters(parameters map[string]string) HttpClient {
	h.request.SetQueryParams(parameters)

	return h
}

func (h httpClient) WithQueryParameter(key string, value string) HttpClient {
	h.request.SetQueryParam(key, value)

	return h
}

func (h httpClient) WithHeaders(parameters map[string]string) HttpClient {
	h.request.SetHeaders(parameters)

	return h
}

func (h httpClient) WithResponse(response interface{}) HttpClient {
	return httpClient{
		h.request,
		response,
	}
}

func (h httpClient) bindResponse(err error, response *resty.Response) error {
	if h.response == nil {
		return nil
	}

	err = json.Unmarshal(response.Body(), h.response)
	if err != nil {
		return err
	}

	return nil
}

func NewHttpClient(client *resty.Client) HttpClient {
	return httpClient{
		client.NewRequest(),
		nil,
	}
}
