package config

import (
	"log"
	"os"
)

const (
	local = "local"
)

var (
	isLocal bool = false
)

// LoadConfig loads the environment configuration based on the application environment.
func LoadConfig() error {
	env := os.Getenv("APP_ENV")
	if env == local {
		log.Println("Running application in local configuration")
		isLocal = true
	}
	return nil
}

func IsLocal() bool {
	return isLocal
}
