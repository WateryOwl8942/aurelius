package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func extendsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	m.Content = strings.ToLower(m.Content)

	var isRythmChan bool
	var rythmChanID string

	for _, guild := range s.State.Guilds {

		// Get channels for this guild
		channels, _ := s.GuildChannels(guild.ID)

		for _, c := range channels {
			// Check if channel is a guild text channel and not a voice or DM channel
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			if c.Name == "rythm" {
				rythmChanID = c.ID
			}

		}
	}

	rex := regexp.MustCompile("(^\\!([a-z]|[A-Z])*)|(\\.\\.([a-z])*)")
	paramRex := regexp.MustCompile("([a-z]|[0-9])*$")
	timeCmdRex := regexp.MustCompile("^(\\.\\.t) ([0-9]*)?")
	numberRex := regexp.MustCompile("[0-9]+")
	ms := string(rex.Find([]byte(m.Content))[:])

	switch ms {
	case "!fs":

		deleteMsg(s, m, m.ChannelID)

		return

	case "!shuffle":

		deleteMsg(s, m, m.ChannelID)

		return

	case "!q":

		deleteMsg(s, m, m.ChannelID)

		return

	case "!play":

		deleteMsg(s, m, m.ChannelID)

		return

	case "!p":

		deleteMsg(s, m, m.ChannelID)

		return

	case ">fs":

		deleteMsg(s, m, m.ChannelID)

		return

	case ">shuffle":

		deleteMsg(s, m, m.ChannelID)

		return

	case ">q":

		deleteMsg(s, m, m.ChannelID)

		return

	case ">play":

		deleteMsg(s, m, m.ChannelID)

		return

	case ">p":

		deleteMsg(s, m, m.ChannelID)

		return

	case "..clean":

		num, err := strconv.ParseInt(string(paramRex.Find([]byte(m.Content)[:])), 10, 64)

		if err != nil {

			if err.Error() == "strconv.ParseInt: parsing \"clean\": invalid syntax" {

				num = 6
			} else {
				fmt.Println(err)
			}

		}

		deleteAllMsgs(s, m, num+1)

		return

	case "..c":

		if m.ChannelID == os.Getenv("AURELIUSCHANNEL") {
			return
		}

		num, err := strconv.ParseInt(string(paramRex.Find([]byte(m.Content)[:])), 10, 64)

		if err != nil {

			if err.Error() == "strconv.ParseInt: parsing \"c\": invalid syntax" {

				num = 6
			} else {
				fmt.Println(err)
			}

		}

		deleteAllMsgs(s, m, num+1)

		//s.ChannelMessageSend(m.ChannelID, "Use Format Is\n```..clean amount```\nor\n```..cl amount```")

		return

	case "..t":

		num, err := strconv.ParseInt(string(numberRex.Find([]byte(timeCmdRex.Find([]byte(m.Content)[:]))[:])), 10, 64)

		if err != nil {

			if err.Error() == "strconv.ParseInt: parsing \"\": invalid syntax" {

				num = 5
			} else {
				fmt.Println(err)
			}

		}

		time.Sleep(time.Second * time.Duration(num))
		deleteMsg(s, m, m.ChannelID)

		//s.ChannelMessageSend(m.ChannelID, "Use Format Is\n```..clean amount```\nor\n```..cl amount```")

		return

	}

	if m.Author.Username == "Rythm" {

		time.Sleep(time.Second * 30)

		deleteMsg(s, m, m.ChannelID)

		return
	}

	if m.Author.Username == "Rythm 2" {

		time.Sleep(time.Second * 30)

		deleteMsg(s, m, m.ChannelID)

		return
	}

	isRythmChan = m.ChannelID == rythmChanID

	//All Code Must Be After This.
	if !isRythmChan {
		return
	}

	if m.Author.Username != "Rythm" {
		time.Sleep(time.Millisecond * 1500)
		deleteMsg(s, m, m.ChannelID)
		return
	}

}

func deleteMsg(s *discordgo.Session, m *discordgo.MessageCreate, rythmChanID string) string {
	s.ChannelMessageDelete(rythmChanID, m.ID)
	return m.Content
}

func deleteAllMsgs(s *discordgo.Session, m *discordgo.MessageCreate, num int64) {

	msgs, err := s.ChannelMessages(m.ChannelID, int(num), "", "", "")

	if err == nil {
		fmt.Println(err)
	}

	var msgLs []string

	for _, d := range msgs {
		msgLs = append(msgLs, d.ID)
	}

	if err := s.ChannelMessagesBulkDelete(m.ChannelID, msgLs); err != nil {
		fmt.Println(err)
	}
}
