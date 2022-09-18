package meetup

import (
	"discord-bot/config"
	"log"
	"testing"
)

func TestApi_GetNextEvent(t *testing.T) {
	config.ReadConfig("/home/fred/Projects/events-discord-bot/config.json")
	client := NewClient()

	data, err := GetNextMeetupEvent(client)
	if err != nil {
		t.Fatal(err)
	}

	log.Print(data)
}
