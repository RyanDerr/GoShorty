package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	local = "local"
)

var (
	isLocal bool = false
	// TODO: Should this be read in from a yml file?
	awsSecrets = []string{"REDIS_ADDRESS"}
	awsRegion  = "us-east-1"
)

// LoadConfig loads the environment configuration based on the application environment.
func LoadConfig() error {
	env := os.Getenv("APP_ENV")
	if env == local {
		log.Println("Running application in local configuration")
		isLocal = true
	} else {
		log.Println("Running application in AWS configuration")

		if os.Getenv("AWS_REGION") != "" {
			awsRegion = os.Getenv("AWS_REGION")
		}

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))

		if err != nil {
			return err
		}

		svc := secretsmanager.NewFromConfig(config)

		for _, secret := range awsSecrets {
			input := &secretsmanager.GetSecretValueInput{
				SecretId: &secret,
			}

			res, err := svc.GetSecretValue(context.TODO(), input)
			if err != nil {
				return err
			}
			os.Setenv(secret, *res.SecretString)
		}
	}
	return nil
}

func IsLocal() bool {
	return isLocal
}

func GetRegion() string {
	return awsRegion
}
