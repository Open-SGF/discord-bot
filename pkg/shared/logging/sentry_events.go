package logging

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

// sentryEventHandler keeps error logs visible as Sentry issues. The current
// sentry-go/slog handler sends logs, but no longer creates issue events.
type sentryEventHandler struct {
	attrs  []slog.Attr
	groups []string
}

func (h sentryEventHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level == slog.LevelError
}

func (h sentryEventHandler) Handle(ctx context.Context, record slog.Record) error {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	event := sentry.NewEvent()
	event.Timestamp = record.Time.UTC()
	event.Level = sentry.LevelError
	event.Message = record.Message
	event.Logger = "slog"
	if fn := runtime.FuncForPC(record.PC); fn != nil {
		file, line := fn.FileLine(record.PC)
		event.Tags["source"] = fmt.Sprintf("%s:%d", file, line)
	}

	add := func(attr slog.Attr, groups []string) {
		attr.Value = attr.Value.Resolve()
		if (attr.Key == "error" || attr.Key == "err") && attr.Value.Kind() == slog.KindAny {
			if err, ok := attr.Value.Any().(error); ok {
				event.SetException(err, 100)
				return
			}
		}
		event.Tags[strings.Join(append(groups, attr.Key), ".")] = fmt.Sprint(attr.Value.Any())
	}
	for _, attr := range h.attrs {
		add(attr, nil)
	}
	record.Attrs(func(attr slog.Attr) bool {
		add(attr, h.groups)
		return true
	})
	hub.CaptureEventWithHint(event, &sentry.EventHint{Context: ctx})
	return nil
}

func (h sentryEventHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h
	clone.attrs = append([]slog.Attr(nil), h.attrs...)
	for _, attr := range attrs {
		attr.Key = strings.Join(append(h.groups, attr.Key), ".")
		clone.attrs = append(clone.attrs, attr)
	}
	return clone
}

func (h sentryEventHandler) WithGroup(name string) slog.Handler {
	clone := h
	clone.groups = append(append([]string(nil), h.groups...), name)
	return clone
}
