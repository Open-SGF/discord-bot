package main

import (
	"discord-bot/bot"
	"discord-bot/config"
	"log"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Run()
	<-make(chan struct{})
}
