package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("should load local environment variables", func(t *testing.T) {
		resetState()
		t.Setenv("APP_ENV", local)
		LoadConfig()

		assert.True(t, IsLocal())
	})

	t.Run("should not load local environment variables", func(t *testing.T) {
		resetState()
		t.Setenv("APP_ENV", "dev")
		LoadConfig()
		assert.False(t, IsLocal())
	})
}

func resetState() {
	isLocal = false
}
