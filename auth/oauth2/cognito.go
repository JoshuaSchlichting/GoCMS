package auth

import (
	"os"

	"log/slog"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"golang.org/x/oauth2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var poolId string
var region string
var clientId string
var clientSecret string

var awsAccessKeyId string
var awsSecretAccessKey string

const redirectUri string = "http://localhost:8000/getjwtandlogin"

var cognitoProvider cognito.CognitoIdentityProvider
var cognitoEndpoint oauth2.Endpoint

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
	logger.Debug("", "poolId", poolId)
	logger.Debug("", "region", region)
	logger.Debug("", "clientId", clientId)
}
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
		logger.Error("error creating session for cognito:", err)
	}
	cognitoProvider = *cognito.New(session, aws.NewConfig().WithRegion(region))

	poolDesc, err := cognitoProvider.DescribeUserPool(&cognito.DescribeUserPoolInput{UserPoolId: aws.String(poolId)})
	if err != nil {
		logger.Error("error describing user pool:", err)
	}
	cognitoEndpoint = oauth2.Endpoint{
		AuthURL:  "https://" + *poolDesc.UserPool.Domain + ".auth." + region + ".amazoncognito.com/oauth2/authorize",
		TokenURL: "https://" + *poolDesc.UserPool.Domain + ".auth." + region + ".amazoncognito.com/oauth2/token",
	}

}
