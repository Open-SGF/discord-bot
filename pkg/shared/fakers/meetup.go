package fakers

import (
	"time"

	"discord-bot/pkg/shared/models"
	"github.com/brianvoe/gofakeit/v7"
)

type MeetupFaker struct {
	faker *gofakeit.Faker
}

func NewMeetupFaker(seed uint64) *MeetupFaker {
	return &MeetupFaker{
		faker: gofakeit.New(seed),
	}
}

func (m *MeetupFaker) CreateEvent(dateTime *time.Time) *models.MeetupEvent {
	var event models.MeetupEvent
	_ = m.faker.Struct(&event)

	event.DateTime = dateTime

	return &event
}

func (m *MeetupFaker) CreateEvents(count int) []*models.MeetupEvent {
	events := make([]*models.MeetupEvent, count)
	m.faker.Slice(&events)
	return events
}
