package bot

import (
	"discord-bot/config"
	"discord-bot/meetup"
	"discord-bot/util"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TODO: Turn this into a concurrent set to handle safe
// discord server removals
var subscribedServers map[string]*Server

func Run() {
	subscribedServers = make(map[string]*Server)

	session, err := discordgo.New("Bot " + config.Settings.DiscordBotToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	session.AddHandler(onJoinGuild)

	err = session.Open()
	if err != nil {
		return
	}
	defer session.Close()

	// This Discord client doesn't have a way to attach event handlers to the
	// heartbeat, therefore we need our own timer to post to all the servers.
	ticker := time.NewTicker(time.Duration(config.Settings.TickRateInSeconds) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// Loop over Discord servers and do a thing
				for _, server := range subscribedServers {
					// TODO: Run these in parallel?
					postNextEvent(session, server)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// In order to create/post events per discord server, we need to know which
// server we're currently subscribed to.
func onJoinGuild(session *discordgo.Session, event *discordgo.GuildCreate) {
	server := NewServer(event)
	subscribedServers[server.ID] = server
	postNextEvent(session, server)
}

// TODO: We need a handler for when the bot leaves a guild so we can remove it from
// the guild from the subscribedGuilds list. There doesn't seem to be a GuildRemove event
// on first glance. Will need to setup a second Discord server to test the multi-server setup.

func postEmbeddedMessage(session *discordgo.Session, server *Server, event *MeetupEvent) {
	embed, err := session.ChannelMessageSendComplex(server.PostChannelID, event.toEmbeddedMessage())
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("[%s] Posting next event to %s: %s\n", server.ID, server.PostChannelID, embed.Content)
}

func postCustomEventMessage(session *discordgo.Session, server *Server) {
	nextEventMessage := config.Settings.CustomMeetupEventMessage
	send, err := session.ChannelMessageSend(server.PostChannelID, nextEventMessage)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("[%s] Posting next event to %s: %s\n", server.ID, server.PostChannelID, send.Content)
}

func postNextEvent(session *discordgo.Session, server *Server) {
	now := util.TimeNow("America/Chicago")

	if !server.shouldPost(now) {
		return
	}

	if !config.Settings.EnableMeetupApi && !config.Settings.EnableCustomMeetupEventMessage {
		log.Println("Config does not allow sending of any messages!")
		return
	}

	if config.Settings.EnableMeetupApi {
		client := meetup.NewClient()
		meetupEvent, err := meetup.GetNextMeetupEvent(client)
		if err != nil {
			log.Print(err)
		} else {
			// If the upcoming event is not this week, then we know we should skip posting this event
			if meetupEvent.SameWeek(now) {
				postEmbeddedMessage(session, server, &MeetupEvent{meetupEvent})
			}
		}
	}

	// TODO: When config.Settings.EnableCustomMeetupEventMessage is set, check config if there's a holiday so we don't post anything

	if config.Settings.EnableCustomMeetupEventMessage {
		postCustomEventMessage(session, server)
	}

	// Use: session.GuildScheduledEvents()
	// To grab current events in discord, and see if our next fetched event is
	// already there.
	// If not, create the event in Discord and post the meetup event in channel
	// If yes, don't post to channel -- this is our final check to not blast
	// servers upon joining them (helps for the case when the bot goes down in
	// and boots back up in the middle of a cycle)
}
