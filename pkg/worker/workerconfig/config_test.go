package workerconfig

import (
	"context"
	"discord-bot/pkg/shared/appconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	awsConfigManager := appconfig.NewAwsConfigManager()
	ctx := context.Background()

	t.Run("successful load from environment variables", func(t *testing.T) {
		switchToTempTestDir(t)
		t.Setenv(proxyFunctionNameKey, "test-proxy")
		t.Setenv(eventsTableNameKey, "test-events")
		t.Setenv(archivedEventsTableNameKey, "test-archived")
		t.Setenv(groupIDDateTimeIndexNameKey, "test-index")
		t.Setenv(meetupGroupNamesKey, "group1,group2")

		cfg, err := NewConfig(ctx, awsConfigManager)
		require.NoError(t, err)

		assert.Equal(t, "test-proxy", cfg.ProxyFunctionName)
		assert.Equal(t, "test-events", cfg.EventsTableName)
		assert.Equal(t, "test-archived", cfg.ArchivedEventsTableName)
		assert.Equal(t, "test-index", cfg.GroupIDDateTimeIndexName)
		assert.Equal(t, []string{"group1", "group2"}, cfg.MeetupGroupNames)
	})

	t.Run("successful load from .env file", func(t *testing.T) {
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env")

		envContent := strings.Join([]string{
			proxyFunctionNameKey + "=file-proxy",
			eventsTableNameKey + "=file-events",
			archivedEventsTableNameKey + "=file-archived",
			groupIDDateTimeIndexNameKey + "=file-index",
			meetupGroupNamesKey + "=group3,group4",
		}, "\n")

		require.NoError(t, os.WriteFile(envPath, []byte(envContent), 0600))

		origDir, err := os.Getwd()
		require.NoError(t, err)
		t.Cleanup(func() { _ = os.Chdir(origDir) })
		require.NoError(t, os.Chdir(tempDir))

		cfg, err := NewConfig(ctx, awsConfigManager)
		require.NoError(t, err)

		assert.Equal(t, "file-proxy", cfg.ProxyFunctionName)
		assert.Equal(t, "file-events", cfg.EventsTableName)
		assert.Equal(t, "file-archived", cfg.ArchivedEventsTableName)
		assert.Equal(t, "file-index", cfg.GroupIDDateTimeIndexName)
		assert.Equal(t, []string{"group3", "group4"}, cfg.MeetupGroupNames)
	})

	t.Run("sets default values for MeetupGroupNames", func(t *testing.T) {
		switchToTempTestDir(t)
		t.Setenv(proxyFunctionNameKey, "test-proxy")
		t.Setenv(eventsTableNameKey, "test-events")
		t.Setenv(archivedEventsTableNameKey, "test-archived")
		t.Setenv(groupIDDateTimeIndexNameKey, "test-index")

		cfg, err := NewConfig(ctx, awsConfigManager)
		require.NoError(t, err)

		assert.Empty(t, cfg.MeetupGroupNames)
	})

	t.Run("validation fails with missing fields", func(t *testing.T) {
		switchToTempTestDir(t)

		_, err := NewConfig(ctx, awsConfigManager)
		require.Error(t, err)
		assert.Contains(t, err.Error(), proxyFunctionNameKey)
		assert.Contains(t, err.Error(), eventsTableNameKey)
		assert.Contains(t, err.Error(), archivedEventsTableNameKey)
		assert.Contains(t, err.Error(), groupIDDateTimeIndexNameKey)
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
