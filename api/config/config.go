package config

import (
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
		isLocal = true
	}
	return nil
}

func IsLocal() bool {
	return isLocal
}
