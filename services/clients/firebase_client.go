package clients

//go:generate mockgen -destination=../mocks/mock_firebase_client.go -package=mocks -source=firebase_client.go

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/neel1996/saul/log"
)

type FirebaseClient interface {
	GenerateAuthToken(ctx context.Context, userId string) (string, error)
}

type firebaseClient struct {
	authClient *auth.Client
}

func (client firebaseClient) GenerateAuthToken(ctx context.Context, userId string) (string, error) {
	logger := log.NewLogger(ctx)
	logger.Infof("Generating new auth token for user: %s", userId)

	err := client.authClient.RevokeRefreshTokens(ctx, userId)
	if err != nil {
		logger.Errorf("Error occurred while revoking refresh tokens for user: %s, error: %v", userId, err)
		return "", err
	}

	token, err := client.authClient.CustomToken(ctx, userId)
	if err != nil {
		logger.Errorf("Error occurred while generating custom token for user: %s, error: %v", userId, err)
		return "", err
	}

	return token, nil
}

func NewFirebaseClient(authClient *auth.Client) FirebaseClient {
	return firebaseClient{
		authClient,
	}
}
