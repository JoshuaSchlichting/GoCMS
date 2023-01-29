package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var poolId string
var region string
var clientId string
var clientSecret string

var awsAccessKeyId string
var awsSecretAccessKey string

const redirectUri string = "http://localhost:8000/secure"

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

func NewCognito() (*Cognito, error) {

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
	var tokenEndpoint string
	if poolDesc.UserPool.CustomDomain != nil {
		fmt.Println("Found custom domain to build endpoint")
		tokenEndpoint = "https://" + *poolDesc.UserPool.CustomDomain + "/oauth2/token"
	} else if poolDesc.UserPool.Domain != nil {
		fmt.Println("Found cognito domain to build endpoint")
		tokenEndpoint = "https://" + *poolDesc.UserPool.Domain + ".auth." + region + ".amazoncognito.com/oauth2/token"
	} else {
		log.Printf("No domain present for user pool %v, unable to build token endpoint.\n", poolId)
		return CognitoPayload{}, err
	}
	log.Printf("Token endpoint: %v\n", tokenEndpoint)

	var authHeader []byte = []byte("Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	headerMap := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": string(authHeader),
	}
	body := []byte(`grant_type=authorization_code&client_id=` + clientId + `&code=` + authorizationCode + `&redirect_uri=` + redirectUri)
	statusCode, responseBody, err := postRequest(tokenEndpoint, headerMap, body)
	if err != nil {
		log.Printf("Error calling token endpoint: %v\n", err)
		return CognitoPayload{}, err
	}
	if statusCode != 200 {
		return CognitoPayload{}, errors.New("error calling token endpoint: " + string(responseBody) + ": status code: " + strconv.Itoa(statusCode))
	}
	cognitoPayload := CognitoPayload{}
	json.Unmarshal(responseBody, &cognitoPayload)
	return cognitoPayload, nil
}

func postRequest(url string, headersMap map[string]string, body []byte) (statusCode int, responseBody []byte, err error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	for key, value := range headersMap {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()
	body, _ = io.ReadAll(response.Body)
	return response.StatusCode, body, err
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
