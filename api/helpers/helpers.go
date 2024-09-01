package helpers

import (
	"os"
	"strings"
)

func EnforceHTTPS(url string) string {
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
