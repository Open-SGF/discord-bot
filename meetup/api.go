package meetup

import (
	"discord-bot/config"
	"encoding/json"
	"errors"
	"log"
	"time"
)

type Event struct {
	Title            string `json:"title"`
	DateTime         string `json:"dateTime"`
	Timezone         string `json:"timezone"`
	ShortUrl         string `json:"shortUrl"`
	ImageUrl         string `json:"imageUrl"`
	ShortDescription string `json:"shortDescription"`
	Description      string `json:"description"`
	GroupName        string
}

func (e *Event) SameWeek(inTime time.Time) bool {
	eventTime, err := time.Parse("2006-01-02T15:04-07:00", e.DateTime)
	if err != nil {
		return false
	}
	eventYear, eventWeek := eventTime.ISOWeek()
	inYear, inWeek := inTime.ISOWeek()
	return eventYear == inYear && eventWeek == inWeek
}

type Node[T any] struct {
	Node *T `json:"node"`
}

type Edge[T Event] struct {
	Edges []*Node[T] `json:"edges"`
}

type Group struct {
	Id             string       `json:"id"`
	Name           string       `json:"name"`
	UpcomingEvents *Edge[Event] `json:"upcomingEvents"`
}

type GroupEnvelope struct {
	Group *Group `json:"group"`
}

type Envelope[D any] struct {
	Data *D `json:"data"`
}

func GetNextMeetupEvent(client *Client) (*Event, error) {

	query := `query GetUpcomingEventsForGroup ($groupId: ID) {
			group(id: $groupId) {
				id,
				name,
				upcomingEvents (input: {first: 1}) {
					edges {
						node {
							title,
							dateTime,
							timezone,
							shortUrl,
							imageUrl,
							shortDescription,
							description
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"groupId": config.Settings.Meetup.GroupID,
	}

	resp, err := client.MakeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	log.Print(string(resp))

	var out Envelope[GroupEnvelope]
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return nil, err
	}

	if out.Data == nil || out.Data.Group == nil || out.Data.Group.UpcomingEvents == nil {
		return nil, errors.New(string("object has missing data"))
	}

	if len(out.Data.Group.UpcomingEvents.Edges) == 0 || out.Data.Group.UpcomingEvents.Edges[0].Node == nil {
		return nil, errors.New(string("no events"))
	}

	out.Data.Group.UpcomingEvents.Edges[0].Node.GroupName = out.Data.Group.Name

	return out.Data.Group.UpcomingEvents.Edges[0].Node, nil
}
