package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	DiscordBotToken   string
	TickRateInSeconds uint
	MeetupGroupID     string
	config            *configStruct
)

type configStruct struct {
	DiscordBotToken   string `json:"discordBotToken"`
	TickRateInSeconds uint   `json:"tickRateInSeconds"`
	MeetupGroupID     string `json:"meetupGroupId"`
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	DiscordBotToken = config.DiscordBotToken
	TickRateInSeconds = config.TickRateInSeconds
	MeetupGroupID = config.MeetupGroupID

	return nil
}
