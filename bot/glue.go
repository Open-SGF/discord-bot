package bot

import (
	"discord-bot/meetup"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type MeetupEvent struct {
	*meetup.Event
}

func (e *MeetupEvent) toEmbeddedMessage() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Description: e.Description,
		Image: &discordgo.MessageEmbedImage{
			URL: e.ImageUrl,
		},
		Title: fmt.Sprintf("%s: %s", e.GroupName, e.Title),
		URL:   e.ShortUrl,
	}
}
