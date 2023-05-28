package request

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type DocumentQaSocketMessage struct {
	DocumentId string `json:"documentId" required:"true"`
	Question   string `json:"question" required:"true"`
}

func (documentQaSocketMessage DocumentQaSocketMessage) Validate() error {
	err := validator.New().Struct(documentQaSocketMessage)
	if err != nil {
		return err
	}

	return nil
}

func (documentQaSocketMessage DocumentQaSocketMessage) UnmarshallJSON(data []byte) error {
	err := json.Unmarshal(data, &documentQaSocketMessage)
	if err != nil {
		return err
	}

	return nil
}
