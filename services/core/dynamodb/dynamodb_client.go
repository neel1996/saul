package dynamodb

//go:generate mockgen -source=dynamodb_client.go -destination=../mocks/mock_dynamodb_client.go -package=mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	GetItem(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	TransactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

type dynamoDBClient struct {
	dynamodbClient dynamodb.Client
}

func (c dynamoDBClient) TransactWriteItems(ctx context.Context, input *dynamodb.TransactWriteItemsInput, opts ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	return c.dynamodbClient.TransactWriteItems(ctx, input, opts...)
}

func (c dynamoDBClient) Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return c.dynamodbClient.Query(ctx, input)
}

func (c dynamoDBClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return c.dynamodbClient.GetItem(ctx, input)
}

func NewDynamoDBClient(dynamodbClient dynamodb.Client) DynamoDBClient {
	return dynamoDBClient{
		dynamodbClient: dynamodbClient,
	}
}
