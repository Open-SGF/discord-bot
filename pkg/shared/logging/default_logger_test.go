package logging

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	t.Run("creates json handler", func(t *testing.T) {
		logger := DefaultLogger(context.Background(), Config{Level: slog.LevelWarn, Type: LogTypeJSON})

		assert.True(t, logger.Handler().Enabled(context.Background(), slog.LevelWarn))
		assert.False(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	})

	t.Run("creates text handler", func(t *testing.T) {
		logger := DefaultLogger(context.Background(), Config{Level: slog.LevelWarn, Type: LogTypeText})

		assert.True(t, logger.Handler().Enabled(context.Background(), slog.LevelWarn))
		assert.False(t, logger.Handler().Enabled(context.Background(), slog.LevelInfo))
	})

	t.Run("adds sentry sink if sentry is enabled", func(t *testing.T) {
		_ = sentry.Init(sentry.ClientOptions{})
		defer sentry.CurrentHub().BindClient(nil)

		// Invalid level to prevent logs from appearing in test output
		level := slog.LevelError + 1
		logger := DefaultLogger(context.Background(), Config{Level: level, Type: LogTypeJSON})

		assert.Empty(t, sentry.CurrentHub().LastEventID())

		logger.Error("test error", "error", errors.New("error"))

		assert.NotEmpty(t, sentry.CurrentHub().LastEventID())
	})

	t.Run("panics for unknown log type", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = DefaultLogger(context.Background(), Config{Level: slog.LevelWarn, Type: LogType(-1)})
		})
	})
}
