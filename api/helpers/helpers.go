package helpers

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	local = "local"
)

func EnforceHTTPS(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http://", "https://", 1)
	}

	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}

func RemoveDomainError(url string) bool {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")
}

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
