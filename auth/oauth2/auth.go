package auth

import (
	"errors"
	"fmt"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/aws"
)

func GetUserInfo(accessToken string) (username, email string, err error) {
	output, err := cognitoProvider.GetUser(&cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		return "", "", err
	}
	var emailAddress string
	for _, attribute := range output.UserAttributes {
		if *attribute.Name == "email" {
			emailAddress = *attribute.Value
			break
		}
	}
	return *output.Username, emailAddress, err

}

func GetAccessJWT(authorizationCode string) (string, error) {
	if authorizationCode == "" {
		return "", errors.New("no authorization code cannot be empty string")
	}
	authClient, _ := New()
	payload, err := authClient.GetOauthTokenFromEndpoint(authorizationCode)
	if err != nil {
		logger.Debug(fmt.Sprintf("error getting token endpoint payload: %v\n", err), "tokenEndpoint", authClient.endpoint.AuthURL)
		return "", err
	}
	logger.Debug("Cognito JWT", "payload", payload.AccessToken)
	return payload.AccessToken, nil
}
