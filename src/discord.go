package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var isCaesar bool

func startDiscordBot() {
	discordClient, err := discordgo.New("Bot " + os.Getenv("CTOKEN"))
	if err != nil {
		fmt.Println(err)
	}

	discordClient.AddHandler(startVote)
	discordClient.AddHandler(extendsHandler)

	if err := discordClient.Open(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Started Bot")
}

func startVote(s *discordgo.Session, m *discordgo.MessageCreate) {

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
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(m.ChannelID, m.ID) }()
		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "✅"); err != nil {
			fmt.Println(err)
		}

		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "❎"); err != nil {
			fmt.Println(err)
		}

		//Wait Time
		time.Sleep(time.Second * 10)

		checkReactionUsers, _ := s.MessageReactions(message.ChannelID, message.ID, "✅", 0, "", "")
		for _, user := range checkReactionUsers {
			fmt.Println(user.Username)
		}

		crossReactionUsers, _ := s.MessageReactions(message.ChannelID, message.ID, "❎", 0, "", "")
		for _, user := range crossReactionUsers {
			fmt.Println(user.Username)
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		if verify(s, m, checkReactionUsers, crossReactionUsers) {
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
	msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Results Of**%v** Started By %v:\n✅ **%v**  ❎ **%v**",
		strings.ToUpper(m.Message.Content[6:]),
		m.Message.Author.Mention(),
		len(checkReactions)-1,
		len(crossReactions)-1))
	go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

	var checkSenateAmount int
	var crossSenateAmount int

	switch cmd {
	case "_admin":

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
			msg, _ = s.ChannelMessageSend(m.ChannelID, "No Senators Voted In Favour")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if checkSenateAmount <= crossSenateAmount {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Senate Disagreed")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		return true

	case "_kick":

		member, _ := s.GuildMember(m.GuildID, userID)
		for _, role := range member.Roles {
			if role == os.Getenv("SENATUSID") {
				return true
			}
		}

		msg, _ := s.ChannelMessageSend(m.ChannelID, "No Senators Voted")
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return false

	case "_ban":

		member, _ := s.GuildMember(m.GuildID, userID)
		for _, role := range member.Roles {
			if role == os.Getenv("SENATUSID") {
				return true
			}
		}

		msg, _ := s.ChannelMessageSend(m.ChannelID, "No Senators Voted")
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return false

	default:

		if len(checkReactions) <= len(crossReactions) {
			msg, err := s.ChannelMessageSend(m.ChannelID, "Same Number Of Votes. Petition Cancelled.")
			if err != nil {
				fmt.Println(err)
			}
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}
	}

	return true

}

func execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmd, user, params := separateIntoCommand(m.Content)
	fmt.Printf("CMD:%v\nUSER:%v\nPARAMS:%v\n", cmd, user, params)

	switch cmd {

	case "_id":

		roles, _ := s.GuildRoles(m.GuildID)
		for _, role := range roles {
			fmt.Printf("Role:%v, ID:%v\n", role.Name, role.ID)
		}

		return

	case "_slave":

		if err := s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("SERVUSID")); err != nil {
			fmt.Println(err)
		}

		go func() {
			time.Sleep(time.Hour * 6)
			s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("SERVUSID"))
		}()

		return

	case "_free":

		s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("SERVUSID"))
		member, _ := s.GuildMember(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Freed %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_dj":

		s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("DJID"))
		member, _ := s.GuildMember(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Given DJ To %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_undj":

		s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("DJID"))
		member, _ := s.GuildMember(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Removed DJ From %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_admin":

		if isCaesar {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "There's An Existing Caesar Already")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		}

		s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("CAESARID"))
		member, _ := s.GuildMember(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Given Admin To %v", member.User.Mention()))
		isCaesar = true
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		go func() {
			time.Sleep(time.Minute * 30)
			s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("CAESARID"))
			msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Power Has Been Revoked From %v", member.User.Mention()))
			isCaesar = false
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		}()

		return

	case "_unadmin":

		s.GuildMemberRoleRemove(m.GuildID, user, os.Getenv("CAESARID"))
		member, _ := s.GuildMember(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Removed Caesar From %v", member.User.Mention()))
		isCaesar = false
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_kick":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildMemberDeleteWithReason(m.GuildID, user, params)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Kicked %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return
	case "_ban":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildBanCreateWithReason(m.GuildID, user, params, 365)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Banned %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_unban":

		member, _ := s.GuildMember(m.GuildID, user)
		s.GuildBanDelete(m.GuildID, user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Unbanned %v", member.User.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return

	}
}
