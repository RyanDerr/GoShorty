package urlparser

import (
	"os"
	"strings"
)

// EnforceHTTPS ensures the URL uses HTTPS.
func EnforceHTTPS(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http://", "https://", 1)
	}

	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}

// RemoveDomainError checks if the URL domain matches the configured domain.
func RemoveDomainError(url string) bool {
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")
}
