package initializers

import (
	"context"
	"github.com/neel1996/saul/configuration"
	"github.com/neel1996/saul/dynamodb/migrations"
)

func Bootstrap(config configuration.Configuration) {
	firebaseAuth := InitializeFirebaseAuth(context.Background())
	InitializeMiddlewares(config, firebaseAuth)

	// DB setup
	dynamoDb := InitializeDynamoDb()
	migrations.NewMigration(dynamoDb).Setup()
}
