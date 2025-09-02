package main

import (
	"context"
	"log"

	"discord-bot/pkg/infra"
	"discord-bot/pkg/infra/infraconfig"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	config, err := infraconfig.NewConfig(context.Background())
	if err != nil {
		log.Println(err)
	}

	app := awscdk.NewApp(nil)

	infra.NewStack(app, "DiscordBot", &infra.AppStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		AppEnv: config.AppEnv,
	})

	app.Synth(nil)
}

func env() *awscdk.Environment { return nil }
