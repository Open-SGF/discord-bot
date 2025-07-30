package httpclient

import (
	"discord-bot/pkg/shared/clock"
	"log/slog"
	"net/http"
)

func DefaultClient(timeSource clock.TimeSource, logger *slog.Logger) *http.Client {
	return &http.Client{Transport: NewHttpLoggingTransport(timeSource, logger)}
}
