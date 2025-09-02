package models

import (
	"time"
)

type MeetupEvent struct {
	ID          string        `json:"id" fake:"{uuid}"`
	Group       MeetupGroup   `json:"group"`
	Title       string        `json:"title" fake:"{sentence:5}"`
	EventURL    string        `json:"eventUrl" fake:"https://www.meetup.com/{username}/events/{digit:10}"`
	Description string        `json:"description" fake:"{paragraph:3,5,10, }"`
	DateTime    *time.Time    `json:"dateTime"`
	Duration    string        `json:"duration" fake:"{number:1,4} hours"`
	Venue       *MeetupVenue  `json:"venue"`
	Host        *MeetupHost   `json:"host"`
	Images      []MeetupImage `json:"images" fakesize:"3"`
}

type MeetupGroup struct {
	URLName string `json:"urlname" fake:"{username}"`
	Name    string `json:"name" fake:"{company} {noun}"`
}

type MeetupVenue struct {
	Name       string `json:"name" fake:"{company} {noun}"`
	Address    string `json:"address" fake:"{street}"`
	City       string `json:"city" fake:"{city}"`
	State      string `json:"state" fake:"{stateabr}"`
	PostalCode string `json:"postalCode" fake:"{zip}"`
}

type MeetupHost struct {
	Name string `json:"name" fake:"{name}"`
}

type MeetupImage struct {
	BaseURL string `json:"baseUrl" fake:"{url}"`
	Preview string `json:"preview" fake:"{url}/preview"`
}
