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

		logger.Info("ignored")
		logger.Error("test error", "error", errors.New("error"))
		assert.True(t, sentry.Flush(time.Second))

		events := transport.Events()
		assert.Len(t, events, 2)
		var gotLog, gotIssue bool
		for _, event := range events {
			if event.Message == "test error" {
				gotIssue = true
				assert.Equal(t, sentry.LevelError, event.Level)
				assert.NotEmpty(t, event.Exception)
			}
			if len(event.Logs) == 1 {
				gotLog = true
				assert.Equal(t, "test error", event.Logs[0].Body)
				assert.Equal(t, sentry.LogLevelError, event.Logs[0].Level)
			}
		}
		assert.True(t, gotLog)
		assert.True(t, gotIssue)
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
