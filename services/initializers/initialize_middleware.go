package initializers

import (
	"firebase.google.com/go/v4/auth"
	"github.com/neel1996/saul/configuration"
	"github.com/neel1996/saul/middleware"
)

var (
	authMiddleware middleware.AuthMiddleware
)

func InitializeMiddlewares(config configuration.Configuration, firebaseAuth *auth.Client) {
	authMiddleware = middleware.NewAuthMiddleware(config, firebaseAuth)
}
