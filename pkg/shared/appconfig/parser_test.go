package appconfig

import (
	"context"
	"discord-bot/pkg/shared/logging"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_DefineKeys(t *testing.T) {
	p := NewParser().
		DefineKeys([]string{"TEST_KEY", "TEST_KEY_2"}).
		WithCustomProcessor(func(ctx context.Context, v *viper.Viper) error {
			assert.Len(t, v.AllKeys(), 2)
			return nil
		})

	var cfg struct {
		TestKey  string `mapstructure:"test_key"`
		TestKey2 string `mapstructure:"test_key_2"`
	}

	err := p.Parse(context.Background(), &cfg)
	assert.NoError(t, err)
}

func TestParser_WithEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "test.env")
	_ = os.WriteFile(envPath, []byte("KEY=env_value"), 0644)

	p := NewParser().WithEnvFile(dir, "test")
	var cfg struct {
		Key string `mapstructure:"key"`
	}

	err := p.Parse(context.Background(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "env_value", cfg.Key)
}

func TestParser_WithEnvVars(t *testing.T) {
	t.Setenv("TEST_ENV_VAR", "env_value")

	p := NewParser().
		DefineKeys([]string{"TEST_ENV_VAR"}).
		WithEnvVars()
	var cfg struct {
		TestEnvVar string `mapstructure:"test_env_var"`
	}

	err := p.Parse(context.Background(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "env_value", cfg.TestEnvVar)
}

func TestParser_WithSSMParameters(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"Parameters":[{"Name":"/app/secret","Value":"ssm_value"}]}`
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	t.Setenv("AWS_ENDPOINT_URL_SSM", ts.URL)

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)

	require.NoError(t, err)

	p := NewParser().WithSSMParameters(func(ctx context.Context, v *viper.Viper, opts *SSMParameterOptions) {
		opts.AwsConfig = &awsCfg
		opts.SSMPath = "/app"
	})

	var cfgStruct struct {
		Secret string `mapstructure:"secret"`
	}

	err = p.Parse(ctx, &cfgStruct)
	assert.NoError(t, err)
	assert.Equal(t, "ssm_value", cfgStruct.Secret)
}

func TestParser_ProcessorOrderPrecedence(t *testing.T) {
	t.Setenv("CONFIG_KEY", "env_var_value")
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "test.env"), []byte("CONFIG_KEY=file_value"), 0644)

	p := NewParser().
		DefineKeys([]string{"CONFIG_KEY"}).
		WithEnvFile(dir, "test").
		WithEnvVars().
		WithCustomProcessor(func(ctx context.Context, v *viper.Viper) error {
			v.Set("config_key", "processor_value")
			return nil
		})

	var cfg struct {
		ConfigKey string `mapstructure:"config_key"`
	}

	err := p.Parse(context.Background(), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "processor_value", cfg.ConfigKey)
}

func TestParser_ErrorHandling(t *testing.T) {
	p := NewParser().WithCustomProcessor(func(ctx context.Context, v *viper.Viper) error {
		return fmt.Errorf("processor failed")
	})

	var cfg struct{}
	err := p.Parse(context.Background(), &cfg)
	assert.ErrorContains(t, err, "processor failed")
}

func TestParser_WithCommonConfig(t *testing.T) {
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("LOG_TYPE", "json")
	t.Setenv("AWS_REGION", "test_aws_region")
	t.Setenv("AWS_ACCESS_KEY", "test_aws_access_key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test_aws_secret_access_key")
	t.Setenv("DYNAMODB_ENDPOINT", "test_dynamodb_endpoint")
	t.Setenv("DYNAMODB_AWS_REGION", "test_dynamodb_aws_region")
	t.Setenv("DYNAMODB_AWS_ACCESS_KEY", "test_dynamodb_aws_access_key")
	t.Setenv("DYNAMODB_AWS_SECRET_ACCESS_KEY", "test_dynamodb_aws_secret_access_key")
	t.Setenv("OTHER_KEY", "test_other_key")

	p := NewParser().
		DefineKeys([]string{"OTHER_KEY"}).
		WithCommonConfig().
		WithEnvVars()

	var cfg struct {
		Common   `mapstructure:",squash"`
		OtherKey string `mapstructure:"other_key"`
	}

	err := p.Parse(context.Background(), &cfg)
	require.NoError(t, err)

	assert.Equal(t, slog.LevelWarn, cfg.Logging.Level)
	assert.Equal(t, logging.LogTypeJSON, cfg.Logging.Type)
	assert.Equal(t, "test_aws_region", cfg.Aws.AwsRegion)
	assert.Equal(t, "test_aws_access_key", cfg.Aws.AwsAccessKey)
	assert.Equal(t, "test_aws_secret_access_key", cfg.Aws.AwsSecretAccessKey)
	assert.Equal(t, "test_other_key", cfg.OtherKey)
}

func TestParseFromKey_LogLevel(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		fallbackLvl slog.Level
		expectedLvl slog.Level
	}{
		{"correct string value", "DEBUG", slog.LevelInfo, slog.LevelDebug},
		{"incorrect string value", "invalid", slog.LevelWarn, slog.LevelWarn},
		{"non-string value", 0, slog.LevelError, slog.LevelError},
		{"nil value", nil, slog.LevelInfo, slog.LevelInfo},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			key := "loglevel"
			v := viper.New()
			v.SetDefault(key, tc.input)

			ParseFromKey(v, key, logging.ParseLogLevel, tc.fallbackLvl)

			value := v.Get(key)

			assert.IsType(t, slog.LevelInfo, value)
			assert.Equal(t, value, tc.expectedLvl)
		})
	}
}
