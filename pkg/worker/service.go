package worker

import (
	"context"
	"discord-bot/pkg/shared/clock"
	"fmt"
	"log/slog"
	"time"
)

type Service struct {
	meetupEventService     *MeetupEventService
	discordNotifierService *DiscordNotifierService
	timeSource             clock.TimeSource
	logger                 *slog.Logger
}

func NewService(meetupEventService *MeetupEventService, discordNotifierService *DiscordNotifierService, timeSource clock.TimeSource, logger *slog.Logger) *Service {
	return &Service{
		meetupEventService:     meetupEventService,
		discordNotifierService: discordNotifierService,
		timeSource:             timeSource,
		logger:                 logger,
	}
}

func (s *Service) Execute(ctx context.Context) error {
	nextEvent, err := s.meetupEventService.GetNextEvent(ctx)

	if err != nil {
		s.logger.Error("failed to get next event from meetup api", slog.Any("error", err))
		return err
	}

	if nextEvent == nil {
		s.logger.Info("No upcoming events to notify about")
		return nil
	}

	if nextEvent.DateTime == nil {
		s.logger.Error("event DateTime is null")
		return fmt.Errorf("event DateTime is null")
	}

	timeUntil := nextEvent.DateTime.Sub(s.timeSource.Now())

	if timeUntil < 12*time.Hour && timeUntil > 36*time.Hour {
		s.logger.Info("The upcoming event isn't within 1-36 hours from now, skipping notification",
			slog.Any("nextEvent.DateTime", nextEvent.DateTime),
			slog.Any("timeUntil", timeUntil),
		)
		return nil
	}

	if err = s.discordNotifierService.Notify(ctx, nextEvent); err != nil {
		s.logger.Error("failed to get next event from meetup api", slog.Any("error", err))
		return err
	}

	return nil
}
