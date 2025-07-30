package appconfig

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
)

type AwsConfigManager interface {
	SetConfigFromViper(ctx context.Context, v *viper.Viper) error
	Config() *aws.Config
}

type AwsConfigManagerImpl struct {
	config *aws.Config
}

func NewAwsConfigManager() *AwsConfigManagerImpl {
	return &AwsConfigManagerImpl{}
}

func (f *AwsConfigManagerImpl) SetConfigFromViper(ctx context.Context, v *viper.Viper) error {
	var cfgOpts []func(*config.LoadOptions) error

	if region := v.GetString(AWSRegionKey); region != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(region))
	}

	if accessKey := v.GetString(AWSAccessKeyKey); accessKey != "" {
		secretKey := v.GetString(AWSSecretAccessKeyKey)

		if secretKey == "" {
			return fmt.Errorf("missing %s", AWSSecretAccessKeyKey)
		}

		cfgOpts = append(cfgOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				v.GetString(AWSSessionTokenKey),
			),
		))
	} else if profile := v.GetString(AWSProfileKey); profile != "" {
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))

		if configFile := v.GetString(AWSConfigFileKey); configFile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedConfigFiles([]string{configFile}))
		}

		if credentialsFile := v.GetString(AWSSharedCredentialsFileKey); credentialsFile != "" {
			cfgOpts = append(cfgOpts, config.WithSharedCredentialsFiles([]string{credentialsFile}))
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return err
	}

	f.config = &cfg

	return nil
}

func (f *AwsConfigManagerImpl) Config() *aws.Config {
	return f.config
}

func AwsConfigProvider(manager AwsConfigManager) *aws.Config {
	return manager.Config()
}
