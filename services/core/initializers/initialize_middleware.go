package initializers

import (
	"core/configuration"
	"core/middleware"
	"firebase.google.com/go/v4/auth"
)

var (
	authMiddleware middleware.AuthMiddleware
)

func InitializeMiddlewares(config configuration.Configuration, firebaseAuth *auth.Client) {
	authMiddleware = middleware.NewAuthMiddleware(config, firebaseAuth)
}
