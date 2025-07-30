package worker

import (
	"bytes"
	"context"
	"discord-bot/pkg/worker/workerconfig"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type DiscordNotifierService struct {
	httpClient *http.Client
	webhookURL string
	logger     *slog.Logger
}

func NewDiscordNotifierService(config *workerconfig.Config, httpClient *http.Client, logger *slog.Logger) *DiscordNotifierService {
	return &DiscordNotifierService{
		httpClient: httpClient,
		webhookURL: config.DiscordWebhookURL,
		logger:     logger.WithGroup("DiscordNotifierService"),
	}
}

func (s *DiscordNotifierService) Notify(ctx context.Context, event *MeetupEvent) error {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		s.logger.Error("failed to load CST timezone", slog.Any("error", err))
		return ErrDiscordNotify
	}

	eventTime := event.DateTime.In(loc)

	discordReq := DiscordRequest{
		Embeds: []DiscordEmbed{
			{
				Title:       event.Title,
				Description: event.Description,
				URL:         event.EventURL,
				Timestamp:   eventTime.Format(time.RFC3339),
				Color:       5814783,
				Fields: []DiscordEmbedField{
					{Name: "Date", Value: eventTime.Format("Januay 2 2006"), Inline: true},
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
		s.logger.Error("unexpected http status code when notifying discord", slog.Int("statusCode", resp.StatusCode))
		return ErrDiscordNotify
	}

	return nil
}

var ErrDiscordNotify = errors.New("error notifying discord")
