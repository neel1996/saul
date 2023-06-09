package consumers

import (
	"context"
	"core/configuration"
	"core/log"
	"core/model"
	"core/service"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type DocumentStatusConsumer interface {
	ConsumeDocumentStatus()
}

type documentStatusConsumer struct {
	kafkaReader      *kafka.Reader
	config           configuration.Configuration
	documentAnalyzer service.DocumentAnalyzerService
}

func (consumer documentStatusConsumer) ConsumeDocumentStatus() {
	logger := log.NewLogger()
	logger.Info("Consuming document status from Kafka")

	for {
		m, err := consumer.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			logger.Errorf("Failed to read message from Kafka: %v", err)
			return
		}

		logger.Infof("Received message: %v", string(m.Value))

		value := m.Value
		var status model.DocumentStatus
		err = json.Unmarshal(value, &status)
		if err != nil {
			logger.Errorf("Failed to unmarshal document status from JSON: %v", err)
			return
		}

		if status.Status == model.DocumentStatusError {
			logger.Infof("Skipping document status with status %v", status.Status)
			continue
		}
	}
}

func NewDocumentStatusConsumer(config configuration.Configuration, documentAnalyzer service.DocumentAnalyzerService) DocumentStatusConsumer {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.Kafka.BrokerURL},
		GroupID:  config.Kafka.Topics.ProcessDocumentStatus.GroupId,
		Topic:    config.Kafka.Topics.ProcessDocumentStatus.Name,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})

	return documentStatusConsumer{
		kafkaReader,
		config,
		documentAnalyzer,
	}
}
