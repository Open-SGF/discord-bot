package meetup

import (
	"discord-bot/util"
	"testing"
	"time"
)

func TestSameWeekWorks(t *testing.T) {
	event := Event{
		DateTime: "2022-10-11T18:00-05:00",
		Timezone: "America/Chicago",
	}
	now := time.Date(2022, 10, 10, 10, 0, 0, 0, util.LocationOrDefault(event.Timezone))

	if !event.SameWeek(now) {
		t.Fatalf("not same week")
	}

}
