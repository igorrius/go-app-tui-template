package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadWithInterpolation(t *testing.T) {
	// Create a temporary .env file
	err := os.WriteFile(".env", []byte("TEST_VAR=interpolated_value\nTEST_VAR_DEFAULT=from_env"), 0644)
	require.NoError(t, err)
	defer os.Remove(".env") //nolint:errcheck

	// The godotenv.Load() in config.Load() will pick up the .env we just created.
	// However, since godotenv.Load() is called in Load, we don't need to load it here
	// but we should be aware that it might affect other tests if not cleaned up.

	cfg, err := Load("app-config.test.yaml")
	require.NoError(t, err)

	// Test case 1: Existing env var
	val, err := Get[string](cfg, "test.env")
	assert.NoError(t, err)
	assert.Equal(t, "interpolated_value", val)

	// Test case 2: Missing env var with default
	val, err = Get[string](cfg, "test.default_missing")
	assert.NoError(t, err)
	assert.Equal(t, "default_val", val)

	// Test case 3: Existing env var with default
	val, err = Get[string](cfg, "test.default_exists")
	assert.NoError(t, err)
	assert.Equal(t, "from_env", val)

	// Test case 4: Nested value
	val, err = Get[string](cfg, "test.nested.val")
	assert.NoError(t, err)
	assert.Equal(t, "interpolated_value", val)

	// Test case 5: List value
	list, err := Get[[]any](cfg, "test.list")
	assert.NoError(t, err)
	assert.Equal(t, "interpolated_value", list[0])
	assert.Equal(t, "static", list[1])
}
