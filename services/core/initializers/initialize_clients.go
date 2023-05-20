package initializers

import (
	"core/clients"
	"core/clients/hugging_face"
	"core/configuration"
	"core/dynamodb"
	"firebase.google.com/go/v4/auth"
	dynamodbPkg "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-resty/resty/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

var (
	huggingFaceDocumentQAClient hugging_face.DocumentQAClient
	dynamoDBClient              dynamodb.DynamoDBClient
	firebaseClient              clients.FirebaseClient
	minioClient                 clients.MinioClient
)

func InitializeClients(config configuration.Configuration, db *dynamodbPkg.Client, auth *auth.Client) {
	httpClient := clients.NewHttpClient(resty.New())

	huggingFaceDocumentQAClient = hugging_face.NewDocumentQAClient(config, httpClient)
	dynamoDBClient = dynamodb.NewDynamoDBClient(*db)
	firebaseClient = clients.NewFirebaseClient(auth)
	minioClient = clients.NewMinioClient(getMinioClient(config))
}

func getMinioClient(config configuration.Configuration) *minio.Client {
	m, err := minio.New(config.Minio.EndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			"",
		),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	return m
}
