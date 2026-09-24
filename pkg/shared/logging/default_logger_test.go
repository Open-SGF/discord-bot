package logging

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	t.Run("creates json handler", func(t *testing.T) {
		logger := DefaultLogger(
			context.Background(),
			Config{Level: slog.LevelWarn, Type: LogTypeJSON},
		)

		assert.True(t, logger.Handler().Enabled(context.Background(), slog.LevelWarn))
		assert.False(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	})

	t.Run("creates text handler", func(t *testing.T) {
		logger := DefaultLogger(
			context.Background(),
			Config{Level: slog.LevelWarn, Type: LogTypeText},
		)

		assert.True(t, logger.Handler().Enabled(context.Background(), slog.LevelWarn))
		assert.False(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	})

	t.Run("adds sentry sink if sentry is enabled", func(t *testing.T) {
		transport := &sentry.MockTransport{}
		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn:       "https://public@example.com/1",
			Transport: transport,
		})
		assert.NoError(t, err)
		sentry.CurrentHub().BindClient(client)
		defer sentry.CurrentHub().BindClient(nil)

		// Invalid level to prevent logs from appearing in test output
		level := slog.LevelError + 1
		logger := DefaultLogger(context.Background(), Config{Level: level, Type: LogTypeJSON})

		logger.Error("test error", "error", errors.New("error"))
		assert.True(t, sentry.Flush(time.Second))

		events := transport.Events()
		if assert.Len(t, events, 1) && assert.Len(t, events[0].Logs, 1) {
			assert.Equal(t, "test error", events[0].Logs[0].Body)
			assert.Equal(t, sentry.LogLevelError, events[0].Logs[0].Level)
		}
	})

	t.Run("panics for unknown log type", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = DefaultLogger(
				context.Background(),
				Config{Level: slog.LevelWarn, Type: LogType(-1)},
			)
		})
	})
}
