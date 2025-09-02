package worker

import (
	"context"
	"discord-bot/pkg/shared/fakers"
	"discord-bot/pkg/shared/logging"
	"discord-bot/pkg/worker/workerconfig"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeetupEventService_GetNextEvent_Success(t *testing.T) {
	faker := fakers.NewMeetupFaker(0)
	eventTime := time.Now().Add(24 * time.Hour)
	expectedEvent := faker.CreateEvent(&eventTime)
	authToken := "test-auth-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth":
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			authResponse := meetupAuthResponse{AccessToken: authToken}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(authResponse)
		case "/v1/groups/open-sgf/events/next":
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer "+authToken, r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(expectedEvent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	event, err := service.GetNextEvent(context.Background())

	assert.NoError(t, err)
	require.NotNil(t, event)
	assert.Equal(t, expectedEvent.ID, event.ID)
	assert.Equal(t, expectedEvent.Title, event.Title)
}

func TestMeetupEventService_GetNextEvent_AuthFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/auth" {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "invalid-client-id",
		SGFMeetupAPIClientSecret: "invalid-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	event, err := service.GetNextEvent(context.Background())

	assert.Nil(t, event)
	require.Error(t, err)
	assert.Equal(t, ErrMeetupAuth, err)
}

func TestMeetupEventService_GetNextEvent_EventFetchFailure(t *testing.T) {
	authToken := "test-auth-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth":
			authResponse := meetupAuthResponse{AccessToken: authToken}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(authResponse)
		case "/v1/groups/open-sgf/events/next":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	event, err := service.GetNextEvent(context.Background())

	assert.Nil(t, event)
	require.Error(t, err)
	assert.Equal(t, ErrMeetupEventFetch, err)
}

func TestMeetupEventService_GetNextEvent_InvalidEventJSON(t *testing.T) {
	authToken := "test-auth-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth":
			authResponse := meetupAuthResponse{AccessToken: authToken}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(authResponse)
		case "/v1/groups/open-sgf/events/next":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid json`))
		}
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	event, err := service.GetNextEvent(context.Background())

	assert.Nil(t, event)
	require.Error(t, err)
	assert.Equal(t, ErrMeetupEventFetch, err)
}

func TestMeetupEventService_getAuthToken_Success(t *testing.T) {
	expectedToken := "test-auth-token-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/v1/auth", r.URL.Path)

		var authReq meetupAuthRequest
		err := json.NewDecoder(r.Body).Decode(&authReq)
		assert.NoError(t, err)
		assert.Equal(t, "test-client-id", authReq.ClientID)
		assert.Equal(t, "test-client-secret", authReq.ClientSecret)

		authResponse := meetupAuthResponse{AccessToken: expectedToken}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(authResponse)
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	token, err := service.getAuthToken(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestMeetupEventService_getAuthToken_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	token, err := service.getAuthToken(context.Background())

	assert.Empty(t, token)
	require.Error(t, err)
	assert.Equal(t, ErrMeetupAuth, err)
}

func TestMeetupEventService_getAuthToken_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	config := &workerconfig.Config{
		SGFMeetupAPIURL:          server.URL,
		SGFMeetupAPIClientID:     "test-client-id",
		SGFMeetupAPIClientSecret: "test-client-secret",
	}
	logger := logging.NewMockLogger()
	service := NewMeetupEventService(config, server.Client(), logger)

	token, err := service.getAuthToken(context.Background())

	assert.Empty(t, token)
	require.Error(t, err)
	assert.Equal(t, ErrMeetupAuth, err)
}
