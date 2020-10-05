package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func startDiscordBot() {
	discordClient, err := discordgo.New("Bot " + os.Getenv("CTOKEN"))
	if err != nil {
		fmt.Println(err)
	}

	discordClient.AddHandler(testEndpoint)
	discordClient.AddHandler(extendsHandler)

	if err := discordClient.Open(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Started Bot")
}

func testEndpoint(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.ChannelID != os.Getenv("AURELIUSCHANNEL") {
		return
	}

	if len(m.Message.Content) < 6 {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}
	if m.Message.Content[0:6] == "_start" {

		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v Started The Vote:\n**%v**", m.Message.Author.Mention(), strings.ToUpper(m.Message.Content[6:])))
		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "✅"); err != nil {
			fmt.Println(err)
		}

		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "❎"); err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Minute * 10)

		checkReactionUsers, _ := s.MessageReactions(message.ChannelID, message.ID, "✅", 0, "", "")
		for _, user := range checkReactionUsers {
			fmt.Println(user.Username)
		}
		fmt.Println(len(checkReactionUsers))

		crossReactionUsers, _ := s.MessageReactions(message.ChannelID, message.ID, "❎", 0, "", "")
		for _, user := range crossReactionUsers {
			fmt.Println(user.Username)
		}
		fmt.Println(len(crossReactionUsers))

		s.ChannelMessageDelete(m.ChannelID, m.ID)
		if len(checkReactionUsers) > len(crossReactionUsers) {
			verify(s, m, checkReactionUsers, crossReactionUsers)
			execute(s, m)
		}

	} else if m.Author.Bot && m.Author.Username == "Aurelius" {
		return
	} else {
		time.Sleep(time.Second * 1)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func verify(s *discordgo.Session, m *discordgo.MessageCreate, checkReactions []*discordgo.User, crossReactions []*discordgo.User) bool {
	cmd, userID, _ := separateIntoCommand(m.Content)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Results Of**%v** Started By %v:\n✅ **%v**  ❎ **%v**",
		strings.ToUpper(m.Message.Content[6:]),
		m.Message.Author.Mention(),
		len(checkReactions)-1,
		len(crossReactions)-1))
	switch cmd {
	case "_admin":

		var checkSenateAmount int
		var crossSenateAmount int
		for _, user := range checkReactions {
			u, _ := s.GuildMember(m.GuildID, user.ID)
			for _, role := range u.Roles {
				if role == os.Getenv("SENATUSID") {
					checkSenateAmount++
				}
			}
		}

		for _, user := range crossReactions {
			u, _ := s.GuildMember(m.GuildID, user.ID)
			for _, role := range u.Roles {
				if role == os.Getenv("SENATUSID") {
					crossSenateAmount++
				}
			}
		}

		if checkSenateAmount == 0 {
			s.ChannelMessageSend(m.GuildID, "No Senators Voted In Favour")
			return false
		}

		if checkSenateAmount <= crossSenateAmount {
			s.ChannelMessageSend(m.GuildID, "Senate Disagreed")
			return false
		}

		return true

	case "_kick":

		m, _ := s.GuildMember(m.GuildID, userID)
		for _, role := range m.Roles {
			if role == os.Getenv("SENATUSID") {
				return false
			}
		}

		return true

	case "_ban":

		m, _ := s.GuildMember(m.GuildID, userID)
		for _, role := range m.Roles {
			if role == os.Getenv("SENATUSID") {
				return false
			}
		}

		return true
	}
	return true
}

func execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmd, user, params := separateIntoCommand(m.Content)
	fmt.Printf("CMD:%v\nUSER:%v\nPARAMS:%v", cmd, user, params)

	switch cmd {
	case "_slave":

		s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("SLAVEID"))

		go func() {
			time.Sleep(time.Hour * 6)
			s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("SLAVEID"))
		}()

		return

	case "_free":

		s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("SLAVEID"))
		member, _ := s.GuildMember(m.GuildID, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Freed %v", member.User.Mention()))
		return

	case "_dj":

		s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("DJID"))
		member, _ := s.GuildMember(m.GuildID, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Given DJ To %v", member.User.Mention()))
		return

	case "_undj":

		s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("DJID"))
		member, _ := s.GuildMember(m.GuildID, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Removed DJ From %v", member.User.Mention()))
		return

	case "_admin":

		s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("CAESARID"))
		member, _ := s.GuildMember(m.GuildID, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Given Admin To %v", member.User.Mention()))
		go func() {
			time.Sleep(time.Minute * 30)
			s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("CAESARID"))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Power Has Been Revoked From %v", member.User.Mention()))
		}()

		return
	case "_kick":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildMemberDeleteWithReason(m.GuildID, user, params)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Kicked %v", member.User.Mention()))

		return
	case "_ban":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildBanCreateWithReason(m.GuildID, user, params, 365)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Banned %v", member.User.Mention()))
		return

	case "_unban":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildBanDelete(m.GuildID, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Unbanned %v", member.User.Mention()))

		return
	default:
		fmt.Printf("\nCMD: %v, USERID: %v, PARAMS: %v\n", cmd, user, params)
	}
}
