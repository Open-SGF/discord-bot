package worker

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"discord-bot/pkg/shared/clock"
	"discord-bot/pkg/shared/fakers"
	"discord-bot/pkg/shared/logging"
	"discord-bot/pkg/shared/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockDiscordNotifier struct {
	mock.Mock
}

func (m *MockDiscordNotifier) Notify(ctx context.Context, event *models.MeetupEvent) error {
	m.Called(ctx, event)
	return nil
}

type MockMeetupEventService struct {
	mock.Mock
}

func (m *MockMeetupEventService) GetNextEvent(ctx context.Context) (*models.MeetupEvent, error) {
	args := m.Called(ctx)

	// nil is a valid return value
	var event *models.MeetupEvent
	if v := args.Get(0); v != nil {
		event = v.(*models.MeetupEvent)
	}

	return event, args.Error(1)
}

func TestService_Execute(t *testing.T) {
	ctx := context.Background()
	faker := fakers.NewMeetupFaker(0)
	now := time.Now()
	mockMeetupEventService := new(MockMeetupEventService)
	mockDiscordNotifier := new(MockDiscordNotifier)
	timeSource := clock.NewMockTimeSource(now)
	mockLoggerHandler := logging.NewMockHandler()
	logger := slog.New(mockLoggerHandler)
	service := NewService(mockMeetupEventService, mockDiscordNotifier, timeSource, logger)

	reset := func(t *testing.T) {
		t.Helper()
		timeSource.SetTime(now)
		mockLoggerHandler.Reset()
		mockMeetupEventService.ExpectedCalls = nil
		mockMeetupEventService.Calls = nil
		mockDiscordNotifier.ExpectedCalls = nil
		mockDiscordNotifier.Calls = nil
	}

	t.Run("no event", func(t *testing.T) {
		reset(t)

		mockMeetupEventService.On("GetNextEvent", ctx).Return(nil, nil)

		err := service.Execute(ctx)

		require.NoError(t, err)
		mockDiscordNotifier.AssertNotCalled(t, "Notify")
	})

	t.Run("no datetime", func(t *testing.T) {
		reset(t)

		event := faker.CreateEvent(nil)

		mockMeetupEventService.On("GetNextEvent", ctx).Return(event, nil)

		err := service.Execute(ctx)

		require.Error(t, err)
		mockDiscordNotifier.AssertNotCalled(t, "Notify")
	})

	t.Run("date logic", func(t *testing.T) {
		tests := []struct {
			name            string
			eventTimeOffset time.Duration
			shouldNotify    bool
		}{
			{
				name:            "event too far in future",
				eventTimeOffset: 48 * time.Hour,
				shouldNotify:    false,
			},
			{
				name:            "event too soon",
				eventTimeOffset: 6 * time.Hour,
				shouldNotify:    false,
			},
			{
				name:            "event in past",
				eventTimeOffset: -2 * time.Hour,
				shouldNotify:    false,
			},
			{
				name:            "event at 12 hour boundary",
				eventTimeOffset: 12 * time.Hour,
				shouldNotify:    true,
			},
			{
				name:            "event at 36 hour boundary",
				eventTimeOffset: 36 * time.Hour,
				shouldNotify:    true,
			},
			{
				name:            "event within notification window",
				eventTimeOffset: 24 * time.Hour,
				shouldNotify:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reset(t)

				eventTime := now.Add(tt.eventTimeOffset)
				event := faker.CreateEvent(&eventTime)

				mockMeetupEventService.On("GetNextEvent", ctx).Return(event, nil)
				if tt.shouldNotify {
					mockDiscordNotifier.On("Notify", ctx, event)
				}

				err := service.Execute(ctx)

				require.NoError(t, err)
				if tt.shouldNotify {
					mockDiscordNotifier.AssertCalled(t, "Notify", ctx, event)
				} else {
					mockDiscordNotifier.AssertNotCalled(t, "Notify")
				}
			})
		}
	})
}
