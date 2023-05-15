package initializers

import (
	dynamodbPkg "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-resty/resty/v2"
	"github.com/neel1996/saul/clients"
	"github.com/neel1996/saul/clients/hugging_face"
	"github.com/neel1996/saul/configuration"
	"github.com/neel1996/saul/dynamodb"
)

var (
	huggingFaceDocumentQAClient hugging_face.DocumentQAClient
	dynamoDBClient              dynamodb.DynamoDBClient
)

func InitializeClients(config configuration.Configuration, db *dynamodbPkg.Client) {
	httpClient := clients.NewHttpClient(resty.New())

	huggingFaceDocumentQAClient = hugging_face.NewDocumentQAClient(config, httpClient)
	dynamoDBClient = dynamodb.NewDynamoDBClient(*db)
}
