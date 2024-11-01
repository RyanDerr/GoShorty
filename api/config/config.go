package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	local = "local"
)

// LoadConfig loads the environment configuration based on the application environment.
func LoadConfig() error {
	env := os.Getenv("APP_ENV")
	if env == local {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}
	return nil
}
