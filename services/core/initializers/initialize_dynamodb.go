package initializers

import (
	"context"
	"core/log"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func InitializeDynamoDb() *dynamodb.Client {
	logger := log.NewLogger()

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Panic(err)
	}

	return dynamodb.NewFromConfig(awsConfig)
}
