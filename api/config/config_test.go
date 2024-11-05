package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
    t.Run("should load local environment variables", func(t *testing.T) {
        resetState()
        t.Setenv("APP_ENV", "local")

        err := LoadConfig()
        assert.NoError(t, err)
        assert.True(t, IsLocal())
    })

    t.Run("should not load local environment variables if env is not local", func(t *testing.T) {
        resetState()
        t.Setenv("APP_ENV", "production")

        err := LoadConfig()
        assert.NoError(t, err)
        assert.False(t, IsLocal())
    })
}

func resetState() {
	isLocal = false
}
