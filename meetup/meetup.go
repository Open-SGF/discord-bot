package meetup

import (
	"discord-bot/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetNextMeetupEvent() string {

	// TODO: Need to figure out how to use the Meetup API
	// getNextMeetupEvent()

	query := `{
		"query": "query GetUpcomingEventsForGroup ($groupId: ID) {
			group(id: $groupId) {
				id,
				name,
				upcomingEvents (input: {first: 1}) {
					edges {
						node {
							dateTime,
							timezone,
							shortUrl,
							tickets {
								edges {
									node {
										user {
											name
										}
									}
								}
							}
						}
					}
				}
			}
		}",
		"variables": {
			"groupId": ` + config.MeetupGroupID + `
		}
	}`

	reader := strings.NewReader(query)
	// Need to add custom headers
	// req = http.NewRequest()
	// http.Do(req)
	resp, err := http.Post("https://api.meetup.com/gql", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// fmt.Printf("%s", body)
	return string(body)
}
