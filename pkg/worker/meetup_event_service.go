package worker

import (
	"bytes"
	"context"
	"discord-bot/pkg/shared/models"
	"discord-bot/pkg/worker/workerconfig"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

const meetupGroup = "open-sgf"

type MeetupEventService interface {
	GetNextEvent(ctx context.Context) (*models.MeetupEvent, error)
}

type SgfMeetupApiEventService struct {
	httpClient   *http.Client
	baseURL      string
	clientID     string
	clientSecret string
	logger       *slog.Logger
}

func NewMeetupEventService(config *workerconfig.Config, httpClient *http.Client, logger *slog.Logger) *SgfMeetupApiEventService {
	return &SgfMeetupApiEventService{
		httpClient:   httpClient,
		logger:       logger.WithGroup("SgfMeetupApiEventService"),
		baseURL:      config.SGFMeetupAPIURL,
		clientID:     config.SGFMeetupAPIClientID,
		clientSecret: config.SGFMeetupAPIClientSecret,
	}
}

func (s *SgfMeetupApiEventService) GetNextEvent(ctx context.Context) (*models.MeetupEvent, error) {
	authToken, err := s.getAuthToken(ctx)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.baseURL+"/v1/groups/"+meetupGroup+"/events/next",
		nil,
	)

	if err != nil {
		s.logger.Error("failed to create http request", slog.Any("error", err))
		return nil, ErrMeetupEventFetch
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := s.httpClient.Do(req)

	if err != nil {
		s.logger.Error("failed to execute http request", slog.Any("error", err))
		return nil, ErrMeetupEventFetch
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected http status code when fetching event", slog.Int("statusCode", resp.StatusCode))
		return nil, ErrMeetupEventFetch
	}

	var event models.MeetupEvent
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		s.logger.Error("failed to parse meetup event", slog.Any("error", err))
		return nil, ErrMeetupEventFetch
	}

	return &event, nil
}

type meetupAuthRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type meetupAuthResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *SgfMeetupApiEventService) getAuthToken(ctx context.Context) (string, error) {
	authRequest := meetupAuthRequest{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
	}

	jsonData, err := json.Marshal(authRequest)

	if err != nil {
		s.logger.Error("failed to marshal json", slog.Any("error", err))
		return "", ErrMeetupAuth
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/v1/auth",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		s.logger.Error("failed to create http request", slog.Any("error", err))
		return "", ErrMeetupAuth
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := s.httpClient.Do(req)

	if err != nil {
		s.logger.Error("failed to execute http request", slog.Any("error", err))
		return "", ErrMeetupAuth
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected http status code when fetching auth token", slog.Int("statusCode", resp.StatusCode))
		return "", ErrMeetupAuth
	}

	var authResponse meetupAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		s.logger.Error("failed to parse auth response", slog.Any("error", err))
		return "", ErrMeetupAuth
	}

	return authResponse.AccessToken, nil
}

var ErrMeetupAuth = errors.New("meetup auth failed")
var ErrMeetupEventFetch = errors.New("failed to get meetup event")
