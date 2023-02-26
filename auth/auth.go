package auth

import (
	"errors"
	"log"

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
	return *output.Username, emailAddress, nil

}

func GetAccessJWT(authorizationCode string) (string, error) {
	if authorizationCode == "" {
		return "", errors.New("no authorization code cannot be empty string")
	}
	authClient, _ := New()
	payload, err := authClient.GetOauthTokenFromEndpoint(authorizationCode)
	if err != nil {
		log.Printf("Error getting token endpoint payload: %v\n", err)
		return "", err
	}
	return payload.AccessToken, nil
}
