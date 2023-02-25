package auth

import (
	"context"
	"log"
	"os"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/oauth2"
)

var poolId string
var region string
var clientId string
var clientSecret string

var awsAccessKeyId string
var awsSecretAccessKey string

const redirectUri string = "http://localhost:8000/getjwtandlogin"

var cognitoProvider cognito.CognitoIdentityProvider

func init() {
	poolId = os.Getenv("POOL_ID")
	region = os.Getenv("REGION")
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyId, awsSecretAccessKey, ""),
	})
	if err != nil {
		log.Printf("Error creating session for cognito: %v\n", err)
	}
	cognitoProvider = *cognito.New(session, aws.NewConfig().WithRegion(region))
}

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
