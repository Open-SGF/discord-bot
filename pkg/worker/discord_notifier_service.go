package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"discord-bot/pkg/shared/models"
	"discord-bot/pkg/worker/workerconfig"
)

type DiscordNotifierService interface {
	Notify(ctx context.Context, event *models.MeetupEvent) error
}

type HttpDiscordNotifierService struct {
	httpClient *http.Client
	webhookURL string
	logger     *slog.Logger
}

func NewDiscordNotifierService(
	config *workerconfig.Config,
	httpClient *http.Client,
	logger *slog.Logger,
) *HttpDiscordNotifierService {
	return &HttpDiscordNotifierService{
		httpClient: httpClient,
		webhookURL: config.DiscordWebhookURL,
		logger:     logger.WithGroup("HttpDiscordNotifierService"),
	}
}

func (s *HttpDiscordNotifierService) Notify(ctx context.Context, event *models.MeetupEvent) error {
	if event == nil || event.DateTime == nil {
		return ErrMissingDate
	}

	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		s.logger.Error("failed to load CST timezone", slog.Any("error", err))
		return ErrDiscordNotify
	}

	eventTime := event.DateTime.In(loc)

	discordReq := models.DiscordRequest{
		Embeds: []models.DiscordEmbed{
			{
				Title:       event.Title,
				Description: event.Description,
				URL:         event.EventURL,
				Color:       5814783,
				Fields: []models.DiscordEmbedField{
					{Name: "Date", Value: eventTime.Format("January 2 2006"), Inline: true},
					{Name: "Time", Value: eventTime.Format("3:04PM"), Inline: true},
				},
			},
		},
	}

	jsonData, err := json.Marshal(discordReq)
	if err != nil {
		s.logger.Error("failed to marshal json", slog.Any("error", err))
		return ErrDiscordNotify
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.webhookURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		s.logger.Error("failed to create http request", slog.Any("error", err))
		return ErrDiscordNotify
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("failed to execute http request", slog.Any("error", err))
		return ErrDiscordNotify
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		s.logger.Error(
			"unexpected http status code when notifying discord",
			slog.Int("statusCode", resp.StatusCode),
		)
		return ErrDiscordNotify
	}

	return nil
}

var (
	ErrMissingDate   = errors.New("event or event datetime is null")
	ErrDiscordNotify = errors.New("error notifying discord")
)
