package logging

import (
	"context"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

// sentryEventHandler keeps error logs visible as Sentry issues. The current
// sentry-go/slog handler sends logs, but no longer creates issue events.
type sentryEventHandler struct {
	attrs []slog.Attr
}

func (h sentryEventHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level == slog.LevelError
}

func (h sentryEventHandler) Handle(ctx context.Context, record slog.Record) error {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	client := hub.Client()
	if client == nil {
		return nil
	}

	var logErr error
	findError := func(attr slog.Attr) bool {
		if attr.Key == "error" || attr.Key == "err" {
			if err, ok := attr.Value.Resolve().Any().(error); ok && err != nil {
				logErr = err
			}
		}
		return true
	}
	for _, attr := range h.attrs {
		findError(attr)
	}
	record.Attrs(findError)

	var event *sentry.Event
	if logErr != nil {
		event = client.EventFromException(logErr, sentry.LevelError)
		event.Message = record.Message
	} else {
		event = client.EventFromMessage(record.Message, sentry.LevelError)
	}
	event.Timestamp = record.Time.UTC()
	hub.CaptureEventWithHint(event, &sentry.EventHint{Context: ctx, OriginalException: logErr})
	return nil
}

func (h sentryEventHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return h
}

// Only error values are used; group names are not needed for issue metadata.
func (h sentryEventHandler) WithGroup(_ string) slog.Handler {
	return h
}
