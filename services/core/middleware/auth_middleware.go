package middleware

import (
	"core/configuration"
	"core/log"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	Authenticate(ctx *gin.Context)
}

type authMiddleware struct {
	config     configuration.Configuration
	authClient *auth.Client
}

func (middleware authMiddleware) Authenticate(ctx *gin.Context) {
	logger := log.NewLogger(ctx)
	ignoredEndpoints := middleware.config.AuthIgnoreUrls
	path := ctx.Request.URL.Path

	for _, ignoredEndpoint := range ignoredEndpoints {
		if path == ignoredEndpoint {
			logger.Infof("%s is not part of secure list. Ignoring auth", ignoredEndpoint)
			ctx.Next()
			return
		}
	}

	logger.Infof("Authenticating %s", path)

	idToken := ctx.GetHeader("Authorization")
	if idToken == "" || !strings.HasPrefix(idToken, "Bearer ") {
		logger.Errorf("No authorization token found")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	idToken = strings.Replace(idToken, "Bearer ", "", 1)
	token, err := middleware.authClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		logger.Errorf("Error verifying token: %s", err.Error())
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logger.Infof("Successfully authenticated %s", token.UID)
	ctx.Set("authToken", token)
	ctx.Next()
}

func NewAuthMiddleware(config configuration.Configuration, authClient *auth.Client) AuthMiddleware {
	return authMiddleware{
		config,
		authClient,
	}
}
