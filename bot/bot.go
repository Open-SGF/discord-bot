package bot

import (
	"bufio"
	"discord-bot/config"
	"discord-bot/meetup"
	"discord-bot/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botID string

// TODO: Turn this into a concurrent set to handle safe
// discord server removals
var subscribedServers map[string]string

func Run() {
	subscribedServers = make(map[string]string)

	session, err := discordgo.New("Bot " + config.Settings.DiscordBotToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err := session.User("@me")
	if err != nil {
		log.Fatal(err)
		return
	}

	botID = user.ID
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
				for guildID, _ := range subscribedServers {
					// TODO: Run these in parallel?
					postNextEvent(session, guildID)
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
	// TODO: Need to get guild permissions (event.Permissions), and then assign those to our cache
	// This is so that we can post what we think we can post and have reasonable fallbacks

	// TODO: We need a way to know which channel to post to per discord server
	// For now we'll use the event.SystemChannelID since that's considered the "announcements" channel
	subscribedServers[event.ID] = event.SystemChannelID

	postNextEvent(session, event.ID)
}

// TODO: We need a handler for when the bot leaves a guild so we can remove it from
// the guild from the subscribedGuilds list. There doesn't seem to be a GuildRemove event
// on first glance. Will need to setup a second Discord server to test the multi-server setup.

func postNextEvent(session *discordgo.Session, guildID string) {

	if !shouldGetNextMeetupEvent(guildID) {
		return
	}

	if !config.Settings.EnableMeetupApi && !config.Settings.EnableCustomMeetupEventMessage {
		return
	}

	nextEventMessage := ""
	if config.Settings.EnableMeetupApi {
		client := meetup.NewClient()
		meetupEvent, err := meetup.GetNextMeetupEvent(client)
		if err != nil {
			log.Print(err)

			// Go ahead and default to custom message if this fails
			config.Settings.EnableCustomMeetupEventMessage = true
		} else {
			// If the upcoming event is not this week, then we know we should skip posting this event
			if meetupEvent.SameWeek(util.TimeNow("America/Chicago")) {
				nextEventMessage = constructMessage(meetupEvent)
			}
		}
	}

	// TODO: When config.Settings.EnableCustomMeetupEventMessage is set, check config if there's a holiday so we don't post anything

	if config.Settings.EnableCustomMeetupEventMessage {
		nextEventMessage = config.Settings.CustomMeetupEventMessage
	}

	if len(nextEventMessage) <= 0 {
		return
	}

	// Use: session.GuildScheduledEvents()
	// To grab current events in discord, and see if our next fetched event is
	// already there.
	// If not, create the event in Discord and post the meetup event in channel
	// If yes, don't post to channel -- this is our final check to not blast
	// servers upon joining them (helps for the case when the bot goes down in
	// and boots back up in the middle of a cycle)

	channelID := subscribedServers[guildID]
	fmt.Printf("[%s] Posting next event to %s: %s\n", guildID, channelID, nextEventMessage)
	session.ChannelMessageSend(channelID, nextEventMessage)
}

func shouldGetNextMeetupEvent(guildID string) bool {
	now := util.TimeNow("America/Chicago")
	nextMondayTenAM := util.CalculateNextBlastDate(now, time.Monday, 10*time.Hour)

	// Use of goto is a bit funky with GO's variable
	// definitions, so opt for lambda instead
	writeDateToFile := func(f *os.File) {
		// Make sure we're at the start of the file
		f.Seek(0, 0)

		_, err := f.WriteString(strconv.FormatInt(nextMondayTenAM.Unix(), 10))
		if err != nil {
			log.Println(err)
		}
	}

	fileName := fmt.Sprintf("./.%s_next_update", guildID)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Only read the first line because that's all we care about
	scanner.Scan()
	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		log.Println(err)
		log.Println("writing new time to file anyway")
		writeDateToFile(file)
		return false
	}

	if len(line) <= 0 {
		writeDateToFile(file)
		return true
	}

	// This is a tricky case because if someone mutates this file with garbage,
	// do we want to post anyway, or don't post anything? This will either
	// cause a double post in the former case, and no post in the later.
	// We should opt for double post because the chances of this file
	// being touched by something else should be slim, and we'd rather
	// have _some_ post than none.
	nextRunTime, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		log.Println(err)
		log.Println("writing new time to file anyway")
		writeDateToFile(file)
		return true
	}

	if now.Unix() >= nextRunTime {
		writeDateToFile(file)
		return true
	}

	return false
}

func constructMessage(event *meetup.Event) string {
	return ""
}
