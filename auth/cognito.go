package auth

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"golang.org/x/oauth2"
)

var poolId string
var region string
var clientId string
var clientSecret string

var awsAccessKeyId string
var awsSecretAccessKey string

const redirectUri string = "http://localhost:8000/getjwtandlogin"

func init() {
	poolId = os.Getenv("POOL_ID")
	region = os.Getenv("REGION")
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

type CognitoPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
}

type Cognito struct {
	session session.Session
	client  cognito.CognitoIdentityProvider
}

func New() (*Cognito, error) {

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyId, awsSecretAccessKey, ""),
	})
	if err != nil {
		return &Cognito{}, err
	}
	var cognitoClient cognito.CognitoIdentityProvider = *cognito.New(session, aws.NewConfig().WithRegion(region))

	return &Cognito{
		session: *session,
		client:  cognitoClient,
	}, nil
}

func (c *Cognito) GetCognitoTokenEndpointPayload(authorizationCode string) (CognitoPayload, error) {
	poolDesc, err := c.client.DescribeUserPool(&cognito.DescribeUserPoolInput{UserPoolId: aws.String(poolId)})
	if err != nil {
		log.Printf("Error describing user pool: %v\n", err)
		return CognitoPayload{}, err
	}
	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUri,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + *poolDesc.UserPool.Domain + ".auth." + region + ".amazoncognito.com/oauth2/authorize",
			TokenURL: "https://" + *poolDesc.UserPool.Domain + ".auth." + region + ".amazoncognito.com/oauth2/token",
		},
	}
	token, err := config.Exchange(context.Background(), authorizationCode)
	if err != nil {
		log.Printf("Error exchanging code for token: %v\n", err)
		return CognitoPayload{}, err
	}
	cognitoPayload := CognitoPayload{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int(token.Expiry.Unix()),
		IdToken:      token.Extra("id_token").(string),
	}
	return cognitoPayload, nil
}

func (c *Cognito) GetUserInfo(accessToken string) (username, email string, err error) {
	// use access token to get user info from cognito
	// return user info
	output, err := c.client.GetUser(&cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		return "", "", err
	}
	// search *output.UserAttributes for email
	var emailAddress string
	for _, attribute := range output.UserAttributes {
		if *attribute.Name == "email" {
			emailAddress = *attribute.Value
			break
		}
	}
	return *output.Username, emailAddress, nil

}

func GetAccessToken(authorizationCode string) (string, error) {
	if authorizationCode == "" {
		return "", errors.New("no authorization code found")
	}
	cognitoClient, _ := New()
	payload, err := cognitoClient.GetCognitoTokenEndpointPayload(authorizationCode)
	if err != nil {
		log.Printf("Error getting token endpoint payload: %v\n", err)
		return "", err
	}
	return payload.AccessToken, nil
}
