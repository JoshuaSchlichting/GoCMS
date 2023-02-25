package auth

import (
	"context"
	"log"

	"golang.org/x/oauth2"
)

type Auth struct {
	endpoint oauth2.Endpoint
}

func New(endpoint oauth2.Endpoint) (*Auth, error) {

	return &Auth{
		endpoint: endpoint,
	}, nil
}

func (a *Auth) GetOauthTokenFromEndpoint(authorizationCode string) (*oauth2.Token, error) {

	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUri,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     a.endpoint,
	}
	token, err := config.Exchange(context.Background(), authorizationCode)
	if err != nil {
		log.Printf("Error exchanging code for token: %v\n", err)
		return &oauth2.Token{}, err
	}
	return token, nil
}
