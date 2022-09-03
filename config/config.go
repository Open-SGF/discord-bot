package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Settings *configStruct
)

type configStruct struct {
	DiscordBotToken                string `json:"discordBotToken"`
	TickRateInSeconds              uint   `json:"tickRateInSeconds"`
	EnableMeetupApi                bool   `json:"enableMeetupApi"`
	EnableCustomMeetupEventMessage bool   `json:"enableCustomMeetupEventMessage"`
	CustomMeetupEventMessage       string `json:"customMeetupEventMessage,omitempty"`
	Meetup                         struct {
		GroupID           string `json:"groupId"`
		UserID            string `json:"userId"`
		JWTSigningString  string `json:"jwtSignedString"`
		JWTPrivateKeyPath string `json:"jwtPrivateKeyPath"`
		OAuthClientKey    string `json:"oauthClientKey"`
	} `json:"meetup,omitempty"`
}

func ReadConfig(configFile string) error {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = json.Unmarshal(file, &Settings)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if Settings.DiscordBotToken[0] == '$' {
		Settings.DiscordBotToken = os.Getenv("OPENSGF_DISCORD_BOT_TOKEN")
	}

	return nil
}
