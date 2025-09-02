package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"discord-bot/pkg/shared/fakers"
	"discord-bot/pkg/shared/logging"
	"discord-bot/pkg/worker/workerconfig"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscordNotifierService_Notify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	config := &workerconfig.Config{DiscordWebhookURL: server.URL}
	logger := logging.NewMockLogger()
	service := NewDiscordNotifierService(config, server.Client(), logger)

	faker := fakers.NewMeetupFaker(0)
	eventTime := time.Now().Add(24 * time.Hour)
	event := faker.CreateEvent(&eventTime)

	err := service.Notify(context.Background(), event)

	assert.NoError(t, err)
}

func TestDiscordNotifierService_Notify_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	config := &workerconfig.Config{DiscordWebhookURL: server.URL}
	logger := logging.NewMockLogger()
	service := NewDiscordNotifierService(config, server.Client(), logger)

	faker := fakers.NewMeetupFaker(0)
	eventTime := time.Now().Add(24 * time.Hour)
	event := faker.CreateEvent(&eventTime)

	err := service.Notify(context.Background(), event)

	require.Error(t, err)
	assert.Equal(t, ErrDiscordNotify, err)
}

func TestDiscordNotifierService_Notify_NilDateTime(t *testing.T) {
	config := &workerconfig.Config{DiscordWebhookURL: "http://test.com"}
	logger := logging.NewMockLogger()
	service := NewDiscordNotifierService(config, &http.Client{}, logger)

	faker := fakers.NewMeetupFaker(0)
	event := faker.CreateEvent(nil)
	event.DateTime = nil

	err := service.Notify(context.Background(), event)

	require.Error(t, err)
	assert.Equal(t, ErrMissingDate, err)
}
