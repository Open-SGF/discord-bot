package httpclient

import (
	"discord-bot/pkg/shared/clock"
	"discord-bot/pkg/shared/logging"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpLoggingTransport_SuccessfulRequest(t *testing.T) {
	mockHandler := logging.NewMockHandler()
	transport := NewHttpLoggingTransport(clock.NewRealTimeSource(), slog.New(mockHandler))
	client := &http.Client{Transport: transport}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1234")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	resp, err := client.Do(req)

	require.NoError(t, err)

	defer func() { _ = resp.Body.Close() }()

	infoEntries := mockHandler.Entries(slog.LevelDebug)

	require.Len(t, infoEntries, 1)

	entry := infoEntries[0]

	assert.Equal(t, "request completed", entry.Message)

	expectedAttrs := map[string]any{
		"http_client.method":         "GET",
		"http_client.url":            ts.URL,
		"http_client.status_code":    200,
		"http_client.content_length": "1234",
	}

	for k, v := range expectedAttrs {
		assert.Contains(t, entry.Attrs, k)
		assert.EqualValues(t, v, entry.Attrs[k])
	}

	assert.Contains(t, entry.Attrs, "http_client.duration")
	assert.IsType(t, time.Duration(0), entry.Attrs["http_client.duration"])
}

func TestHttpLoggingTransport_FailedRequest(t *testing.T) {
	mockHandler := logging.NewMockHandler()
	transport := NewHttpLoggingTransport(clock.NewRealTimeSource(), slog.New(mockHandler))
	client := &http.Client{Transport: transport}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	_, err := client.Do(req)

	require.Error(t, err)

	errorEntries := mockHandler.Entries(slog.LevelError)

	require.Len(t, errorEntries, 1)

	entry := errorEntries[0]

	assert.Equal(t, "request failed", entry.Message)
	assert.Contains(t, entry.Attrs, "http_client.method")
	assert.IsType(t, "", entry.Attrs["http_client.method"])
	assert.Contains(t, entry.Attrs, "http_client.url")
	assert.IsType(t, "", entry.Attrs["http_client.url"])

	assert.Contains(t, entry.Attrs, "http_client.error")
	if _, ok := entry.Attrs["http_client.error"].(error); !ok {
		assert.Fail(t, "Missing error in error log")
	}
}
