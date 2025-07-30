package worker

import "time"

type MeetupEvent struct {
	ID          string        `json:"id"`
	Group       MeetupGroup   `json:"group"`
	Title       string        `json:"title"`
	EventURL    string        `json:"eventUrl"`
	Description string        `json:"description"`
	DateTime    *time.Time    `json:"dateTime"`
	Duration    string        `json:"duration"`
	Venue       *MeetupVenue  `json:"venue"`
	Host        *MeetupHost   `json:"host"`
	Images      []MeetupImage `json:"images"`
}

type MeetupGroup struct {
	URLName string `json:"urlname"`
	Name    string `json:"name"`
}

type MeetupVenue struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
}

type MeetupHost struct {
	Name string `json:"name"`
}

type MeetupImage struct {
	BaseURL string `json:"baseUrl"`
	Preview string `json:"preview"`
}

type DiscordRequest struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	URL         string              `json:"url"`
	Timestamp   string              `json:"timestamp"`
	Color       int                 `json:"color"`
	Fields      []DiscordEmbedField `json:"fields"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
