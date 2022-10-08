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
		GroupID              string `json:"groupId"`
		UserID               string `json:"userId"`
		JWTSigningString     string `json:"jwtSignedString"`
		JWTPrivateKeyPath    string `json:"jwtPrivateKeyPath"`
		OAuthClientKey       string `json:"oauthClientKey"`
		OAuthClientSecretKey string `json:"OAuthClientSecretKey"`
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

	if Settings.Meetup.GroupID[0] == '$' {
		Settings.Meetup.GroupID = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_GROUP_ID")
	}

	if Settings.Meetup.UserID[0] == '$' {
		Settings.Meetup.UserID = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_USER_ID")
	}

	if Settings.Meetup.JWTSigningString[0] == '$' {
		Settings.Meetup.JWTSigningString = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_JWT_PUBLIC_KEY")
	}

	if Settings.Meetup.JWTPrivateKeyPath[0] == '$' {
		Settings.Meetup.JWTPrivateKeyPath = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_JWT_PRIVATE_KEY_PATH")
	}

	if Settings.Meetup.OAuthClientKey[0] == '$' {
		Settings.Meetup.OAuthClientKey = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_OAUTH_CUSTOMER_KEY")
	}

	if Settings.Meetup.OAuthClientSecretKey[0] == '$' {
		Settings.Meetup.OAuthClientSecretKey = os.Getenv("OPENSGF_DISCORD_BOT_MEETUP_OAUTH_CLIENT_SECRET_KEY")
	}

	return nil
}
