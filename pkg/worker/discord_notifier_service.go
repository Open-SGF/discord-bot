package worker

import (
	"context"
	"discord-bot/pkg/worker/workerconfig"
	"net/http"
)

type DiscordNotifierService struct {
	client     *http.Client
	webhookURL string
}

func NewDiscordNotifierService(config *workerconfig.Config, client *http.Client) *DiscordNotifierService {
	return &DiscordNotifierService{
		client:     client,
		webhookURL: config.DiscordWebhookURL,
	}
}

func (s *DiscordNotifierService) Notify(ctx context.Context, event *MeetupEvent) error {

}
