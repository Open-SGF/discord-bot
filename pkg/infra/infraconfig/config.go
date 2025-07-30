package infraconfig

import (
	"context"
	"discord-bot/pkg/shared/appconfig"
)

const (
	appEnvKey        = "APP_ENV"
	appDomainNameEnv = "APP_DOMAIN_NAME"
)

var configKeys = []string{
	appEnvKey,
	appDomainNameEnv,
}

type Config struct {
	AppEnv string `mapstructure:"app_env"`
}

func NewConfig(ctx context.Context) (*Config, error) {
	var config Config

	err := appconfig.NewParser().
		DefineKeys(configKeys).
		WithEnvFile(".", ".env").
		WithEnvVars().
		Parse(ctx, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
