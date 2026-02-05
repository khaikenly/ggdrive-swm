package auth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleapi "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	ScopeDriveReadonly = "https://www.googleapis.com/auth/drive.readonly"
)

func NewOAuth2Config(clientID, clientSecret, backendURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  backendURL + "/api/auth/callback",
		Scopes:       []string{ScopeDriveReadonly},
		Endpoint:     google.Endpoint,
	}
}

func DriveServiceFromToken(ctx context.Context, token *oauth2.Token) (*googleapi.Service, error) {
	config := &oauth2.Config{Endpoint: google.Endpoint}
	client := config.Client(ctx, token)

	svc, err := googleapi.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}

	return svc, nil
}
