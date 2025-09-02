package httpclient

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"discord-bot/pkg/shared/clock"
)

func NewHttpLoggingTransport(timeSource clock.TimeSource, logger *slog.Logger) http.RoundTripper {
	return &httpLoggingTransport{
		timeSource: timeSource,
		logger:     logger.WithGroup("http_client"),
	}
}

type httpLoggingTransport struct {
	timeSource clock.TimeSource
	logger     *slog.Logger
}

func (h *httpLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	debugEnabled := h.logger.Enabled(req.Context(), slog.LevelDebug)

	if debugEnabled {
		return h.debugLoggingRoundTrip(req)
	} else {
		return h.infoLoggingRoundTrip(req)
	}
}

func (h *httpLoggingTransport) infoLoggingRoundTrip(req *http.Request) (*http.Response, error) {
	start := h.timeSource.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		h.logger.ErrorContext(
			req.Context(),
			"request failed",
			"method", req.Method,
			"url", req.URL.String(),
			"error", err,
		)
		return resp, err
	}

	h.logger.LogAttrs(
		req.Context(),
		slog.LevelInfo,
		"request completed",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("status_code", resp.StatusCode),
		slog.Duration("duration", time.Since(start)),
		slog.String("content_length", resp.Header.Get("Content-Length")),
	)
	return resp, nil
}

func (h *httpLoggingTransport) debugLoggingRoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	if req.Body != nil {
		var err error
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			h.logger.ErrorContext(req.Context(),
				"error reading request body",
				"error", err,
			)
		} else {
			req.Body = io.NopCloser(bytes.NewReader(reqBody))
		}
	}

	start := h.timeSource.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		attrs := []any{
			"method", req.Method,
			"url", req.URL.String(),
			"error", err,
		}
		if len(reqBody) > 0 {
			attrs = append(attrs, "request_body", string(reqBody))
		}
		h.logger.ErrorContext(
			req.Context(),
			"request failed",
			attrs...,
		)
		return resp, err
	}

	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		h.logger.ErrorContext(resp.Request.Context(),
			"error reading response body",
			"error", err,
		)
	} else {
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}

	logAttrs := []slog.Attr{
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("status_code", resp.StatusCode),
		slog.Duration("duration", time.Since(start)),
		slog.String("content_length", resp.Header.Get("Content-Length")),
	}
	if len(reqBody) > 0 {
		logAttrs = append(logAttrs, slog.String("request_body", string(reqBody)))
	}
	if len(respBody) > 0 {
		logAttrs = append(logAttrs, slog.String("response_body", string(respBody)))
	}

	h.logger.LogAttrs(
		req.Context(),
		slog.LevelDebug,
		"request completed",
		logAttrs...,
	)

	return resp, nil
}
