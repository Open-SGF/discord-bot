package appconfig

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAwsConfigManager(t *testing.T) {
	ctx := context.Background()

	t.Run("empty configuration", func(t *testing.T) {
		v := viper.New()

		factory := NewAwsConfigManager()
		assert.Nil(t, factory.Config())

		err := factory.SetConfigFromViper(ctx, v)
		require.NoError(t, err)

		assert.NotNil(t, factory.Config())
	})

	t.Run("region only", func(t *testing.T) {
		v := viper.New()
		v.Set(AWSRegionKey, "us-east-1")

		factory := NewAwsConfigManager()
		assert.Nil(t, factory.Config())

		err := factory.SetConfigFromViper(ctx, v)
		require.NoError(t, err)

		assert.NotNil(t, factory.Config())
		assert.Equal(t, "us-east-1", factory.Config().Region)
	})

	t.Run("access key and secret key", func(t *testing.T) {
		v := viper.New()
		v.Set(AWSAccessKeyKey, "test_access_key")
		v.Set(AWSSecretAccessKeyKey, "test_secret_key")

		factory := NewAwsConfigManager()
		assert.Nil(t, factory.Config())

		err := factory.SetConfigFromViper(ctx, v)
		require.NoError(t, err)

		assert.NotNil(t, factory.Config())
	})

	t.Run("profile", func(t *testing.T) {
		tmpDir := t.TempDir()

		configContent := "[profile test_profile]\nregion = us-east-1"
		configPath := filepath.Join(tmpDir, "config")
		require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o600))

		credsPath := filepath.Join(tmpDir, "credentials")
		require.NoError(t, os.WriteFile(credsPath, []byte{}, 0o600))

		v := viper.New()
		v.Set(AWSProfileKey, "test_profile")
		v.Set(AWSConfigFileKey, configPath)
		v.Set(AWSSharedCredentialsFileKey, credsPath)

		factory := NewAwsConfigManager()
		assert.Nil(t, factory.Config())

		err := factory.SetConfigFromViper(ctx, v)
		require.NoError(t, err)

		assert.NotNil(t, factory.Config())
	})

	t.Run("missing secret key when access key is set", func(t *testing.T) {
		v := viper.New()
		v.Set(AWSAccessKeyKey, "test_access_key")

		factory := NewAwsConfigManager()
		assert.Nil(t, factory.Config())

		err := factory.SetConfigFromViper(ctx, v)
		require.Error(t, err)
	})
}
