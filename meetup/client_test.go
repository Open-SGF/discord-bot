package meetup

import (
	"discord-bot/config"
	"strings"
	"testing"
)

func TestClient_GetNextAuthToken(t *testing.T) {
	config.ReadConfig("/home/fred/Projects/events-discord-bot/config.json")
	client := NewClient()
	expected, err := client.GetNextAuthToken()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := client.GetNextAuthToken()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(expected, actual) != 0 {
		t.Fatalf("expecting %s got %s", expected, actual)
	}

	refreshToken, err := client.refreshToken()
	if err != nil {
		t.Fatal(err)
	}

	if len(refreshToken) <= 0 {
		t.Fatalf("Refresh token is empty")
	}
}
