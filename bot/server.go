package bot

import (
	"bufio"
	"discord-bot/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	*discordgo.GuildCreate
	PostChannelID string
}

func NewServer(guild *discordgo.GuildCreate) *Server {
	// TODO: Need to get guild permissions (event.Permissions), and then assign those to our cache
	// This is so that we can post what we think we can post and have reasonable fallbacks

	// TODO: We need a way to know which channel to post to per discord server
	// For now we'll use the event.SystemChannelID since that's considered the "announcements" channel
	return &Server{
		guild,
		guild.SystemChannelID,
	}
}

func (s *Server) shouldPost(offset time.Time) bool {
	nextMondayTenAM := util.CalculateNextBlastDate(offset, time.Monday, 10*time.Hour)

	writeDateToFile := func(f *os.File) {
		// Make sure we're at the start of the file
		f.Seek(0, 0)

		_, err := f.WriteString(strconv.FormatInt(nextMondayTenAM.Unix(), 10))
		if err != nil {
			log.Println(err)
		}
	}

	fileName := fmt.Sprintf("./.%s_next_update", s.ID)
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

	if offset.Unix() >= nextRunTime {
		writeDateToFile(file)
		return true
	}

	return false
}
