package initializers

import (
	"context"
	"github.com/neel1996/saul/configuration"
)

func InitializeObjects(config configuration.Configuration) {
	firebaseAuth := InitializeFirebaseAuth(context.Background())
	InitializeMiddlewares(config, firebaseAuth)
}
