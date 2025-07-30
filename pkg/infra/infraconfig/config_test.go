package infraconfig

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("successful load from environment variables", func(t *testing.T) {
		switchToTempTestDir(t)
		t.Setenv(appEnvKey, "staging")

		cfg, err := NewConfig(ctx)
		require.NoError(t, err)

		assert.Equal(t, "staging", cfg.AppEnv)
	})

	t.Run("successful load from .env file", func(t *testing.T) {
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		envContent := strings.Join([]string{
			appEnvKey + "=staging",
		}, "\n")

		require.NoError(t, os.WriteFile(envPath, []byte(envContent), 0600))

		origDir, err := os.Getwd()
		require.NoError(t, err)
		t.Cleanup(func() { _ = os.Chdir(origDir) })
		require.NoError(t, os.Chdir(tempDir))

		cfg, err := NewConfig(ctx)
		require.NoError(t, err)

		assert.Equal(t, "staging", cfg.AppEnv)
	})
}

func switchToTempTestDir(t *testing.T) {
	t.Helper()

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	tempDir := t.TempDir()
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	require.NoError(t, os.Chdir(tempDir))
}
