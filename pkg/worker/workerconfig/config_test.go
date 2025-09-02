package workerconfig

import (
	"context"
	"discord-bot/pkg/shared/appconfig"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	awsConfigManager := appconfig.NewAwsConfigManager()
	ctx := context.Background()

	t.Run("successful load from environment variables", func(t *testing.T) {
		switchToTempTestDir(t)
		t.Setenv(sgfMeetupAPIURL, "test-meetup-url")
		t.Setenv(sgfMeetupAPIClientID, "test-client-id")
		t.Setenv(sgfMeetupAPIClientSecret, "test-client-secret")
		t.Setenv(discordWebhookURL, "test-discord-web")
		cfg, err := NewConfig(ctx, awsConfigManager)
		require.NoError(t, err)

		assert.Equal(t, "test-meetup-url", cfg.SGFMeetupAPIURL)
		assert.Equal(t, "test-client-id", cfg.SGFMeetupAPIClientID)
		assert.Equal(t, "test-client-secret", cfg.SGFMeetupAPIClientSecret)
		assert.Equal(t, "test-discord-web", cfg.DiscordWebhookURL)
	})

	t.Run("successful load from .env file", func(t *testing.T) {
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		envContent := strings.Join([]string{
			sgfMeetupAPIURL + "=test-meetup-url",
			sgfMeetupAPIClientID + "=test-client-id",
			sgfMeetupAPIClientSecret + "=test-client-secret",
			discordWebhookURL + "=test-discord-web",
		}, "\n")

		require.NoError(t, os.WriteFile(envPath, []byte(envContent), 0600))

		origDir, err := os.Getwd()
		require.NoError(t, err)
		t.Cleanup(func() { _ = os.Chdir(origDir) })
		require.NoError(t, os.Chdir(tempDir))

		cfg, err := NewConfig(ctx, awsConfigManager)
		require.NoError(t, err)

		assert.Equal(t, "test-meetup-url", cfg.SGFMeetupAPIURL)
		assert.Equal(t, "test-client-id", cfg.SGFMeetupAPIClientID)
		assert.Equal(t, "test-client-secret", cfg.SGFMeetupAPIClientSecret)
		assert.Equal(t, "test-discord-web", cfg.DiscordWebhookURL)
	})

	t.Run("validation fails with missing fields", func(t *testing.T) {
		switchToTempTestDir(t)

		_, err := NewConfig(ctx, awsConfigManager)
		require.Error(t, err)
		assert.Contains(t, err.Error(), sgfMeetupAPIClientID)
		assert.Contains(t, err.Error(), sgfMeetupAPIClientSecret)
		assert.Contains(t, err.Error(), discordWebhookURL)
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
