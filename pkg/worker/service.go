package worker

import (
	"context"
	"discord-bot/pkg/shared/clock"
	"fmt"
	"log/slog"
)

type Service struct {
	meetupEventService     *MeetupEventService
	discordNotifierService *DiscordNotifierService
	timeSource             clock.TimeSource
	logger                 *slog.Logger
}

func NewService(meetupEventService *MeetupEventService, discordNotifierService *DiscordNotifierService, timeSource clock.TimeSource) *Service {
	return &Service{
		meetupEventService:     meetupEventService,
		discordNotifierService: discordNotifierService,
		timeSource:             timeSource,
	}
}

func (s *Service) Execute(ctx context.Context) error {
	nextEvent, err := s.meetupEventService.GetNextEvent(ctx)

	if err != nil {
		return fmt.Errorf("failed to get next event from meetup api %w", err)
	}

	if err = s.discordNotifierService.Notify(ctx, nextEvent); err != nil {
		return fmt.Errorf("failed to get next event from meetup api %w", err)
	}

	return nil
}
