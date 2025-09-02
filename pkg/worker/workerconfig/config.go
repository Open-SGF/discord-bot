package workerconfig

import (
	"context"
	"discord-bot/pkg/shared/appconfig"
	"fmt"
	"strings"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

const (
	sgfMeetupAPIURL          = "SGF_MEETUP_API_URL"
	sgfMeetupAPIClientID     = "SGF_MEETUP_API_CLIENT_ID"
	sgfMeetupAPIClientSecret = "SGF_MEETUP_API_CLIENT_SECRET"
	discordWebhookURL        = "DISCORD_WEBHOOK_URL"
)

var configKeys = []string{
	sgfMeetupAPIURL,
	sgfMeetupAPIClientID,
	sgfMeetupAPIClientSecret,
	discordWebhookURL,
}

type Config struct {
	appconfig.Common         `mapstructure:",squash"`
	SGFMeetupAPIURL          string `mapstructure:"sgf_meetup_api_url"`
	SGFMeetupAPIClientID     string `mapstructure:"sgf_meetup_api_client_id"`
	SGFMeetupAPIClientSecret string `mapstructure:"sgf_meetup_api_client_secret"`
	DiscordWebhookURL        string `mapstructure:"discord_webhook_url"`
}

func NewConfig(ctx context.Context, awsConfigFactory appconfig.AwsConfigManager) (*Config, error) {
	var config Config

	err := appconfig.NewParser().
		WithCommonConfig().
		DefineKeys(configKeys).
		WithEnvFile(".", ".env").
		WithEnvVars().
		WithCustomProcessor(awsConfigFactory.SetConfigFromViper).
		WithSSMParameters(func(ctx context.Context, v *viper.Viper, opts *appconfig.SSMParameterOptions) {
			opts.AwsConfig = awsConfigFactory.Config()
			opts.SSMPath = v.GetString(appconfig.SSMPathKey)
		}).
		WithCustomProcessor(setDefaults).
		Parse(ctx, &config)

	if err != nil {
		return nil, err
	}

	if err = config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults(_ context.Context, v *viper.Viper) error {
	v.SetDefault(strings.ToLower(sgfMeetupAPIURL), "https://sgf-meetup-api.opensgf.org")
	return nil
}

func (config *Config) validate() error {
	var missing []string

	if config.SGFMeetupAPIURL == "" {
		missing = append(missing, sgfMeetupAPIURL)
	}
	if config.SGFMeetupAPIClientID == "" {
		missing = append(missing, sgfMeetupAPIClientID)
	}
	if config.SGFMeetupAPIClientSecret == "" {
		missing = append(missing, sgfMeetupAPIClientSecret)
	}
	if config.DiscordWebhookURL == "" {
		missing = append(missing, discordWebhookURL)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", strings.Join(missing, ", "))
	}

	return nil
}

var ConfigProviders = wire.NewSet(appconfig.ConfigProviders, wire.FieldsOf(new(*Config), "Common"), NewConfig)
