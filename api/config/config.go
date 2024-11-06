package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gopkg.in/yaml.v3"
)

const (
	local      = "local"
	configFile = "../config.yml"
)

var (
	isLocal   bool = false
	awsConfig *Config
)

type Config struct {
	AwsSecrets []string `yaml:"awsSecrets"`
	AwsRegion  string   `yaml:"awsRegion"`
}

// LoadConfig loads the environment configuration based on the application environment.
func LoadConfig() error {
	env := os.Getenv("APP_ENV")
	if env == local {
		log.Println("Running application in local configuration")
		isLocal = true
	} else {
		log.Println("Running application in AWS configuration")
		err := loadYamlConfig(configFile)
		if err != nil {
			return err
		}

		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsConfig.AwsRegion))

		if err != nil {
			return err
		}

		svc := secretsmanager.NewFromConfig(config)
		err = loadSecrets(svc)
		if err != nil {
			return err
		}

	}
	return nil
}

func loadSecrets(svc *secretsmanager.Client) error {
	for _, secret := range awsConfig.AwsSecrets {
		input := &secretsmanager.GetSecretValueInput{
			SecretId: &secret,
		}

		res, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			return err
		}

		var secretMap map[string]string
		err = json.Unmarshal([]byte(*res.SecretString), &secretMap)
		if err != nil {
			return err
		}

		for key, value := range secretMap {
			os.Setenv(key, value)
		}
	}
	return nil
}

func loadYamlConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, awsConfig)
	if err != nil {
		return err
	}

	return nil
}

func IsLocal() bool {
	return isLocal
}

func GetRegion() string {
	return awsConfig.AwsRegion
}
