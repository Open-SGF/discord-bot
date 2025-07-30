package workerconfig

import (
	"context"
	"discord-bot/pkg/shared/appconfig"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"strings"
)

const (
	meetupGroupNamesKey         = "MEETUP_GROUP_NAMES"
	proxyFunctionNameKey        = "MEETUP_PROXY_FUNCTION_NAME"
	archivedEventsTableNameKey  = "ARCHIVED_EVENTS_TABLE_NAME"
	eventsTableNameKey          = "EVENTS_TABLE_NAME"
	groupIDDateTimeIndexNameKey = "GROUP_ID_DATE_TIME_INDEX_NAME"
)

var configKeys = []string{
	meetupGroupNamesKey,
	proxyFunctionNameKey,
	archivedEventsTableNameKey,
	eventsTableNameKey,
	groupIDDateTimeIndexNameKey,
}

type Config struct {
	appconfig.Common         `mapstructure:",squash"`
	MeetupGroupNames         []string `mapstructure:"meetup_group_names"`
	ProxyFunctionName        string   `mapstructure:"meetup_proxy_function_name"`
	ArchivedEventsTableName  string   `mapstructure:"archived_events_table_name"`
	EventsTableName          string   `mapstructure:"events_table_name"`
	GroupIDDateTimeIndexName string   `mapstructure:"group_id_date_time_index_name"`
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
	v.SetDefault(strings.ToLower(meetupGroupNamesKey), []string{})
	return nil
}

func (config *Config) validate() error {
	var missing []string

	if config.ProxyFunctionName == "" {
		missing = append(missing, proxyFunctionNameKey)
	}
	if config.EventsTableName == "" {
		missing = append(missing, eventsTableNameKey)
	}
	if config.ArchivedEventsTableName == "" {
		missing = append(missing, archivedEventsTableNameKey)
	}
	if config.GroupIDDateTimeIndexName == "" {
		missing = append(missing, groupIDDateTimeIndexNameKey)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", strings.Join(missing, ", "))
	}

	return nil
}

var ConfigProviders = wire.NewSet(appconfig.ConfigProviders, wire.FieldsOf(new(*Config), "Common"), NewConfig)
