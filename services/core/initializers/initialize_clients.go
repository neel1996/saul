package initializers

import (
	clients2 "core/clients"
	"core/clients/hugging_face"
	"core/configuration"
	"core/dynamodb"
	"firebase.google.com/go/v4/auth"
	dynamodbPkg "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-resty/resty/v2"
)

var (
	huggingFaceDocumentQAClient hugging_face.DocumentQAClient
	dynamoDBClient              dynamodb.DynamoDBClient
	firebaseClient              clients2.FirebaseClient
)

func InitializeClients(config configuration.Configuration, db *dynamodbPkg.Client, auth *auth.Client) {
	httpClient := clients2.NewHttpClient(resty.New())

	huggingFaceDocumentQAClient = hugging_face.NewDocumentQAClient(config, httpClient)
	dynamoDBClient = dynamodb.NewDynamoDBClient(*db)
	firebaseClient = clients2.NewFirebaseClient(auth)
}
