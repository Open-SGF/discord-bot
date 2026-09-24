package logging

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSentryEventHandler(t *testing.T) {
	transport := &sentry.MockTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: transport,
	})
	require.NoError(t, err)
	ctx := sentry.SetHubOnContext(context.Background(), sentry.NewHub(client, sentry.NewScope()))
	logger := slog.New(sentryEventHandler{})

	logger.InfoContext(ctx, "ignored", "error", errors.New("ignored"))
	logger.ErrorContext(ctx, "direct", "error", errors.New("direct failure"))
	logger.With("error", errors.New("bound failure")).ErrorContext(ctx, "bound")
	grouped := logger.WithGroup("request").With("err", errors.New("grouped failure"))
	grouped.ErrorContext(ctx, "grouped")
	logger.ErrorContext(ctx, "message only")

	events := transport.Events()
	require.Len(t, events, 4)
	for i, message := range []string{"direct", "bound", "grouped", "message only"} {
		assert.Equal(t, message, events[i].Message)
		assert.Equal(t, sentry.LevelError, events[i].Level)
	}
	for i, message := range []string{"direct failure", "bound failure", "grouped failure"} {
		require.NotEmpty(t, events[i].Exception)
		assert.Equal(t, message, events[i].Exception[0].Value)
	}
	assert.Empty(t, events[3].Exception)
}
