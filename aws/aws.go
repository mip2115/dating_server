package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

type AWS_Config struct {
	Session *session.Session
}

var AWS_Connection *AWS_Config

func SetAWSConnection() error {

	AWS_Connection = &AWS_Config{}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_KEY"),
			""), // token can be left blank for now
	})
	if err != nil {
		return err
	}
	AWS_Connection.Session = s
	return nil
}

func GetSession() (*session.Session, error) {
	if AWS_Connection == nil {
		return nil, errors.New("AWS Connection is nil")
	}

	if AWS_Connection.Session == nil {
		return nil, errors.New("AWS Connection Session is nil")
	}
	return AWS_Connection.Session, nil
}
