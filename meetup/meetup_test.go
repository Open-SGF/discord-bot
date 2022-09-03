package meetup

import (
	"discord-bot/config"
	"testing"
)

func TestGetEventFromMeetup(t *testing.T) {
	config.ReadConfig("/home/fred/Projects/events-discord-bot/config.json")
	client := NewClient()
	client.getNextAuthToken()
}
