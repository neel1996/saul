package initializers

import (
	"context"
	"core/log"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func InitializeFirebaseAuth(ctx context.Context) *auth.Client {
	logger := log.NewLogger(ctx)

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		logger.Panic(err)
	}

	authInstance, err := app.Auth(ctx)
	if err != nil {
		logger.Panic(err)
	}

	return authInstance
}
