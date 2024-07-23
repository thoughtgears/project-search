package auth

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
)

func IDTokenTokenSource(ctx context.Context, audience string) (string, error) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
	}
	tokenSource := creds.TokenSource
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Failed to retrieve token: %v", err)
	}

	return token.AccessToken, nil
}
