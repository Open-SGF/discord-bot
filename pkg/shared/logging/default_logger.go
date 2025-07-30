package logging

import (
	"context"
	"github.com/getsentry/sentry-go"
	slogSentry "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
)

type Config struct {
	Level slog.Level `mapstructure:"log_level"`
	Type  LogType    `mapstructure:"log_type"`
}

func DefaultLogger(context context.Context, config Config) *slog.Logger {

	handlers := make([]slog.Handler, 0, 2)

	opts := &slog.HandlerOptions{
		Level: config.Level,
	}

	switch config.Type {
	case LogTypeText:
		handlers = append(handlers, slog.NewTextHandler(os.Stdout, opts))
	case LogTypeJSON:
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, opts))
	default:
		panic("unknown log type")
	}

	if sentry.CurrentHub().Client() != nil {
		handlers = append(handlers, slogSentry.Option{
			EventLevel: []slog.Level{slog.LevelError},
			AddSource:  true,
		}.NewSentryHandler(context))
	}

	return slog.New(slogmulti.Fanout(handlers...))
}
