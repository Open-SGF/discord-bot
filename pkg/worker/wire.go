//go:build wireinject
// +build wireinject

package worker

import (
	"context"
	"discord-bot/pkg/shared/clock"
	"discord-bot/pkg/shared/httpclient"
	"discord-bot/pkg/shared/logging"
	"discord-bot/pkg/worker/workerconfig"
	"github.com/google/wire"
)

var CommonProviders = wire.NewSet(workerconfig.ConfigProviders, clock.RealClockProvider, logging.DefaultLogger, httpclient.DefaultClient)

func InitService(ctx context.Context) (*Service, error) {
	panic(wire.Build(
		CommonProviders,
		NewMeetupEventService,
		NewDiscordNotifierService,
		NewService,
	))
}
