package logging

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"sync"
	"testing"
	"time"
)

func TestMockLogger(t *testing.T) {
	t.Run("basic logging", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler)

		logger.Info("test message", "user", "john", "age", 30)
		logger.Warn("warning message", "error", "something wrong")

		entries := handler.AllEntries()

		require.Len(t, entries, 2)

		infoEntry := entries[0]
		assert.Equal(t, slog.LevelInfo, infoEntry.Level)
		assert.Equal(t, "test message", infoEntry.Message)

		assert.Equal(t, "john", infoEntry.Attrs["user"])
		assert.Equal(t, int64(30), infoEntry.Attrs["age"])

		warnEntry := entries[1]
		assert.Equal(t, slog.LevelWarn, warnEntry.Level)
		assert.Equal(t, "warning message", warnEntry.Message)
	})

	t.Run("with attributes", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler).With("service", "auth", "version", 1.2)

		logger.Error("failed request", "path", "/login", "status", 500)

		entries := handler.AllEntries()

		require.Len(t, entries, 1)

		verifyAttributes(t, entries[0].Attrs, map[string]any{
			"service": "auth",
			"version": 1.2,
			"path":    "/login",
			"status":  500,
		})
	})

	t.Run("with group and attributes", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler).WithGroup("request").With("method", "GET")

		logger.Debug("processing request", "path", "/api", "duration", 150*time.Millisecond)

		entries := handler.AllEntries()
		verifyAttributes(t, entries[0].Attrs, map[string]any{
			"request.method":   "GET",
			"request.path":     "/api",
			"request.duration": 150 * time.Millisecond,
		})
	})

	t.Run("with nested group and attributes", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler).WithGroup("service").With("name", "serviceName")

		logger.Debug("starting service", "someAttr", "someValue")
		logger.Debug("starting service", "otherAttr", "otherValue")

		nestedLogger := logger.WithGroup("repository").With("name", "repositoryName")

		nestedLogger.Debug("starting repository", "someAttr", "someValue")

		entries := handler.AllEntries()
		verifyAttributes(t, entries[0].Attrs, map[string]any{
			"service.name":     "serviceName",
			"service.someAttr": "someValue",
		})

		verifyAttributes(t, entries[1].Attrs, map[string]any{
			"service.name":      "serviceName",
			"service.otherAttr": "otherValue",
		})

		verifyAttributes(t, entries[2].Attrs, map[string]any{
			"service.name":                "serviceName",
			"service.repository.name":     "repositoryName",
			"service.repository.someAttr": "someValue",
		})
	})

	t.Run("concurrent logging", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler)
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				logger.Info("concurrent message", "worker", n)
			}(i)
		}

		wg.Wait()
		require.Len(t, handler.AllEntries(), 100)
	})

	t.Run("reset", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler)

		logger.Info("first message")
		handler.Reset()
		logger.Info("second message")

		entries := handler.AllEntries()

		require.Len(t, entries, 1)
		assert.Equal(t, "second message", entries[0].Message)
	})

	t.Run("filter entries by level", func(t *testing.T) {
		handler := NewMockHandler()
		logger := slog.New(handler)

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warning message")

		infoEntries := handler.Entries(slog.LevelInfo)

		require.Len(t, infoEntries, 1)
		assert.Equal(t, "info message", infoEntries[0].Message)

		warnEntries := handler.Entries(slog.LevelWarn)

		require.Len(t, warnEntries, 1)
		assert.Equal(t, "warning message", warnEntries[0].Message)
	})
}

func verifyAttributes(t *testing.T, actual map[string]any, expected map[string]any) {
	t.Helper()

	require.Equal(t, len(expected), len(actual), "attribute count mismatch")

	for k, expectedVal := range expected {
		assert.Contains(t, actual, k, "missing attribute %q", k)
		if actualVal, ok := actual[k]; ok {
			assert.EqualValues(t, expectedVal, actualVal, "attribute %q value mismatch", k)
		}
	}
}
