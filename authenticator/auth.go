package authenticator

import (
	"context"
	"errors"
	"sportsify/config"
	"sportsify/user-service/models"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"encoding/json"
	"log"
	"net/http"
	//"github.com/gin-gonic/gin"
	"bytes"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {

    config.LoadSecretsEnv()

	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.Secrets["AUTH0_DOMAIN"]+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     config.Secrets["AUTH0_CLIENT_ID"],
		ClientSecret: config.Secrets["AUTH0_CLIENT_SECRET"],
		RedirectURL:  config.Secrets["AUTH0_CALLBACK_URL"],
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// GetManagementAPIAccessToken returns the Auth0 Management API token
func GetManagementAPIAccessToken() (string, error) {
	config.LoadSecretsEnv()

	auth0Domain := config.Secrets["AUTH0_DOMAIN"]
	clientID := config.Secrets["AUTH0_CLIENT_ID"]
	clientSecret := config.Secrets["AUTH0_CLIENT_SECRET"]
	audience := "https://" + auth0Domain + "/api/v2/"

	tokenReqPayload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
	}

	payloadBytes, _ := json.Marshal(tokenReqPayload)
	tokenReq, _ := http.NewRequest("POST", "https://"+auth0Domain+"/oauth/token", bytes.NewBuffer(payloadBytes))
	tokenReq.Header.Set("Content-Type", "application/json")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return "", err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		log.Printf("Auth0 token request failed with status: %v", tokenResp.StatusCode)
		return "", errors.New("failed to get access token")
	}

	var tokenData models.TokenData
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}