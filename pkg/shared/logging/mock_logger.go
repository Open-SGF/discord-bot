package logging

import (
	"context"
	"log/slog"
	"sync"
)

type MockHandler struct {
	mu          *sync.Mutex
	groupPrefix string
	attrs       *[]slog.Attr
	entries     *[]MockLogEntry
}

type MockLogEntry struct {
	Level    slog.Level
	Message  string
	Attrs    map[string]any
	Grouping []string
}

func NewMockLogger() *slog.Logger {
	return slog.New(NewMockHandler())
}

func NewMockHandler() *MockHandler {
	entries := make([]MockLogEntry, 0)
	attrs := make([]slog.Attr, 0)
	return &MockHandler{
		mu:      &sync.Mutex{},
		entries: &entries,
		attrs:   &attrs,
	}
}

func (h *MockHandler) appendEntry(level slog.Level, msg string, attrs map[string]any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	*h.entries = append(*h.entries, MockLogEntry{
		Level:   level,
		Message: msg,
		Attrs:   attrs,
	})
}

func (h *MockHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := make(map[string]any)

	for _, attr := range *h.attrs {
		attrs[attr.Key] = attr.Value.Any()
	}

	r.Attrs(func(attr slog.Attr) bool {
		k := h.applyGroups(attr.Key)
		attrs[k] = attr.Value.Any()
		return true
	})

	h.appendEntry(r.Level, r.Message, attrs)
	return nil
}

func (h *MockHandler) Enabled(_ context.Context, l slog.Level) bool {
	return true
}

func (h *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	for _, a := range attrs {
		*h.attrs = append(*h.attrs, slog.Attr{
			Key:   h.applyGroups(a.Key),
			Value: a.Value,
		})
	}

	return &MockHandler{
		mu:          h.mu,
		groupPrefix: h.groupPrefix,
		attrs:       h.attrs,
		entries:     h.entries,
	}
}

func (h *MockHandler) WithGroup(name string) slog.Handler {
	newPrefix := name
	if h.groupPrefix != "" {
		newPrefix = h.groupPrefix + "." + name
	}

	return &MockHandler{
		mu:          h.mu,
		groupPrefix: newPrefix,
		attrs:       h.attrs,
		entries:     h.entries,
	}
}

func (h *MockHandler) applyGroups(key string) string {
	if h.groupPrefix == "" {
		return key
	}
	return h.groupPrefix + "." + key
}

func (h *MockHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	entries := make([]MockLogEntry, 0)
	h.entries = &entries
}

func (h *MockHandler) Entries(level slog.Level) []MockLogEntry {
	h.mu.Lock()
	defer h.mu.Unlock()

	var filtered []MockLogEntry
	for _, e := range *h.entries {
		if e.Level == level {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func (h *MockHandler) AllEntries() []MockLogEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	entries := make([]MockLogEntry, len(*h.entries))
	copy(entries, *h.entries)
	return entries
}
