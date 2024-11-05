package urlparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnforceHTTPS(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"happy_path_http":  {"http://example.com", "https://example.com"},
		"happy_path_https": {"https://example.com", "https://example.com"},
		"missing_protocol": {"example.com", "https://example.com"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := EnforceHTTPS(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestRemoveDomainError(t *testing.T) {
	os.Setenv("DOMAIN", "example.com")
	defer os.Unsetenv("DOMAIN")

	tests := map[string]struct {
		input    string
		expected bool
	}{
		"happy_path_http_root_domain":  {"http://example.com", false},
		"happy_path_https_root_domain": {"https://example.com", false},
		"happy_path_http":              {"http://www.example.com", false},
		"happy_path_https":             {"https://www.example.com", false},
		"feedback_loop_http":           {"http://another.com", true},
		"feedback_loop_https":          {"https://another.com", true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := RemoveDomainError(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
