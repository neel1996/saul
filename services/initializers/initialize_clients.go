package initializers

import (
	"firebase.google.com/go/v4/auth"
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
	firebaseClient              clients.FirebaseClient
)

func InitializeClients(config configuration.Configuration, db *dynamodbPkg.Client, auth *auth.Client) {
	httpClient := clients.NewHttpClient(resty.New())

	huggingFaceDocumentQAClient = hugging_face.NewDocumentQAClient(config, httpClient)
	dynamoDBClient = dynamodb.NewDynamoDBClient(*db)
	firebaseClient = clients.NewFirebaseClient(auth)
}
