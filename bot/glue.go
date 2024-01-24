package bot

import (
	"discord-bot/meetup"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type MeetupEvent struct {
	*meetup.Event
}

func (e *MeetupEvent) toEmbeddedMessage() *discordgo.MessageSend {
	var embed = discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Description: e.Description,
		Image: &discordgo.MessageEmbedImage{
			URL: e.ImageUrl,
		},
		Title: fmt.Sprintf("%s: %s", e.GroupName, e.Title),
		URL:   e.ShortUrl,
	}

	var dateTimeFmt = "Mon Jan _2 @ 3:4pm MST"
	var str = "Hey everyone! We have a new event!"
	eventTime, err := time.Parse("2006-01-02T15:04-07:00", e.DateTime)
	if err == nil {
		str += "\n**Event Day & Time**: " + eventTime.Format(dateTimeFmt)
	}

	return &discordgo.MessageSend{
		Content: str,
		Embed:   &embed,
	}
}
