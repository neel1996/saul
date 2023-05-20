package kafka

//go:generate mockgen -destination=../mocks/mock_document_details_producer.go -package=mocks -source=document_details_producer.go

import (
	"context"
	"core/configuration"
	"core/log"
	"core/model"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type DocumentDetailsProducer interface {
	ProduceDocumentDetailsEvent(ctx context.Context, documentDetails model.DocumentDetails) error
}

type documentDetailsProducer struct {
	kafkaConnection *kafka.Conn
	config          configuration.Configuration
}

func (producer documentDetailsProducer) ProduceDocumentDetailsEvent(ctx context.Context, documentDetails model.DocumentDetails) error {
	logger := log.NewLogger(ctx)
	logger.Info("Producing document details event to Kafka")

	b, err := json.Marshal(documentDetails)
	if err != nil {
		logger.Errorf("Failed to marshal document details to JSON: %v", err)
		return err
	}

	_, err = producer.kafkaConnection.WriteMessages(kafka.Message{
		Value: b,
		Time:  time.Now(),
	})

	if err != nil {
		logger.Errorf("Failed to write message to Kafka: %v", err)
		return err
	}

	return nil
}

func NewDocumentDetailsProducer(config configuration.Configuration) DocumentDetailsProducer {
	kafkaConnection, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		config.Kafka.BrokerURL,
		config.Kafka.Topics.ProcessDocument.Name,
		0,
	)
	if err != nil {
		panic(err)
	}

	return documentDetailsProducer{
		kafkaConnection,
		config,
	}
}
