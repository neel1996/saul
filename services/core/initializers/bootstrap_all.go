package initializers

import (
	"context"
	"core/configuration"
	"core/dynamodb/migrations"
)

func Bootstrap(config configuration.Configuration) {
	firebaseAuth := InitializeFirebaseAuth(context.Background())

	// DB setup
	dynamoDb := InitializeDynamoDb()
	migrations.NewMigration(dynamoDb).Setup()

	// Initialize block
	InitializeMiddlewares(config, firebaseAuth)
	InitializeClients(config, dynamoDb, firebaseAuth)
	InitializeRepositories()
	InitializeServices(config)
	InitializeControllers()
}
