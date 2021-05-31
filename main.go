package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/caarlos0/env"
	"log"
	"os"
)

// Configure config for aws
type Configure struct {
	AccessKey       string `env:"AWS_ACCESS_KEY" envDefault:""`
	SecretAccessKey string `env:"AWS_SECRET_KEY" envDefault:""`
	Region          string `env:"AWS_REGION" envDefault:"us-west-2"`
	Profile         string `env:"AWS_PROFILE" envDefault:""`
	SecretId        string `env:"AWS_SECRET_ID" envDefault:"LARAVEL-ENV"`
	SecretVersion   string `env:"AWS_SECRET_VERSION" envDefault:"AWSCURRENT"`
	Filepath        string `env:"FILEPATH" envDefault:".env"`
}

var Config Configure

// Setup init aws config
func init() {
	_ = env.Parse(&Config)
}

func main() {
	file, err := os.OpenFile(Config.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		file, err = os.Create(Config.Filepath)
		if err != nil {
			log.Fatalf(err.Error())
		}

		log.Print(fmt.Sprintf("generated [%s]", Config.Filepath))
	}
	defer file.Close()

	secretJson, err := getSecret()
	if err != nil {
		log.Fatalf(err.Error())
	}

	secrets := make(map[string]string)
	json.Unmarshal([]byte(secretJson), &secrets)

	writer := bufio.NewWriter(file)
	for i, v := range secrets {
		fmt.Fprintln(writer, fmt.Sprintf(`%s="%s"`, i, v))
	}
	writer.Flush()

	fmt.Println("done.")
	os.Exit(0)
}

func getSecret() (string, error) {
	//"GP-BILL-ENV"

	var s *secretsmanager.SecretsManager
	if Config.AccessKey != "" && Config.SecretAccessKey != "" {
		s = secretsmanager.New(session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
					AccessKeyID:     Config.AccessKey,
					SecretAccessKey: Config.SecretAccessKey,
				}),
				Region: aws.String(Config.Region),
			},
		})))
	} else {
		//Create a Secrets Manager client
		s = secretsmanager.New(session.New(),
			aws.NewConfig().WithRegion(Config.Region),
		)
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(Config.SecretId),
		VersionStage: aws.String(Config.SecretVersion), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	result, err := s.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	return *result.SecretString, nil
}
