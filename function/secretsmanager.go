package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/tidwall/gjson"
)

// SecretsManagerClient SecretsManager client interface
type SecretsManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
	secretName  string
	secretValue *secretsmanager.GetSecretValueOutput
}

// New instanciates a SecretsManagerClient struct
//
// Minimum IAM permission:
//
// * secretsmanager:GetSecretValue
//
// * kms:Decrypt
//
//
// It returns a SecretsManagerClient or any error encountered
func New() (*SecretsManagerClient, error) {
	// Instanciate a new aws session
	awssession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Instanciate a new SecretsManager client with an aws session
	svc := secretsmanager.New(awssession)

	// Prepare input to retrieve secret's value string
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsSecretsManagerName),
	}

	// Retrieve secret's value
	output, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	return &SecretsManagerClient{svc, awsSecretsManagerName, output}, nil
}

// GetSlackToken will lookup the AWS Secrets Manager json value to find the Slack token
// It returns the Slack token as a string or any error encountered
func (s *SecretsManagerClient) GetSlackToken() (string, error) {
	// Lookup if the key already exists.
	// gjson.Get() return the value if the key exists
	jsonValue := gjson.Get(*s.secretValue.SecretString, jsonAWSSecretsManagerNameKeyName)
	if jsonValue.String() == "" {
		return "", errors.New("Error finding the Slack token value in AWS Secrets Manager! Key " + jsonAWSSecretsManagerNameKeyName + " doesn't exists! Exiting")
	}
	return jsonValue.String(), nil
}
