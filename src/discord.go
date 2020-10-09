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
	discordClient.AddHandler(ready)
	discordClient.AddHandler(movieChats)

	if err := discordClient.Open(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Started Bot")
}

func movieChats(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) < 7 {
		return
	}
	if m.Content[0:6] != "_watch" {
		return
	}

	var movie string
	if m.Content[6] == ' ' {
		movie = m.Content[7:][:]
	} else {
		fmt.Println(123)
		movie = m.Content[6:]
	}

	msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Starting Private Party For %v", movie))
	if err := s.MessageReactionAdd(msg.ChannelID, msg.ID, "🙌"); err != nil {
		fmt.Println(err)
	}
	if err := s.MessageReactionAdd(msg.ChannelID, msg.ID, "🏁"); err != nil {
		fmt.Println(err)
	}

	if err := s.MessageReactionAdd(msg.ChannelID, msg.ID, "❌"); err != nil {
		fmt.Println(err)
	}

	for {
		time.Sleep(time.Second)
		//Must repeat every second
		flagUsers, _ := s.MessageReactions(msg.ChannelID, msg.ID, "🏁", 0, "", "")
		if len(flagUsers) > 1 {
			break
		}
	}

	joinMeUsers, _ := s.MessageReactions(msg.ChannelID, msg.ID, "🙌", 0, "", "")
	id := os.Getenv("THEATRUMCHANNELID")
	for _, user := range joinMeUsers {
		must(s.GuildMemberRoleAdd(m.GuildID, user.ID, os.Getenv("WATCHINGID")))
		must(s.GuildMemberMove(m.GuildID, user.ID, &id))
	}

	for {
		time.Sleep(time.Second)
		//Must repeat every second
		flagUsers, _ := s.MessageReactions(msg.ChannelID, msg.ID, "❌", 0, "", "")
		if len(flagUsers) > 1 {
			break
		}
	}

	for _, user := range joinMeUsers {
		must(s.GuildMemberRoleRemove(m.GuildID, user.ID, os.Getenv("WATCHINGID")))
		s.ChannelMessageDelete(m.GuildID, msg.ID)
	}

}

func ready(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	fmt.Println("ready")
}

func startVote(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.ChannelID != os.Getenv("AURELIUSCHANNEL") {
		if m.ChannelID != os.Getenv("AURELIUSTESTINGCHANNELID") {
			return
		}
	}

	if len(m.Message.Content) < 6 {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}
	if m.Message.Content[0:6] == "_start" {

		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v Started The Vote:\n**%v**", m.Message.Author.Mention(), strings.ToUpper(m.Message.Content[6:])))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(message.ChannelID, message.ID) }()
		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "✅"); err != nil {
			fmt.Println(err)
		}

		if err := s.MessageReactionAdd(message.ChannelID, message.ID, "❎"); err != nil {
			fmt.Println(err)
		}

		//Wait Time
		time.Sleep(time.Minute * 5)

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

	if len(checkReactions)+len(crossReactions) < 5 {
		msg, _ := s.ChannelMessageSend(m.ChannelID, "Minimum Of 3 Votes Not Reached.")
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return false
	}

	cmd, _, _ := separateIntoCommand(m.Content)
	msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Results Of**%v** Started By %v:\n✅ **%v**  ❎ **%v**", strings.ToUpper(m.Message.Content[6:]), m.Message.Author.Mention(), len(checkReactions)-1, len(crossReactions)-1))
	go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

	senateAmount, didSenatorVote := countSenators(s, m, checkReactions, crossReactions)
	votes := len(checkReactions) - len(crossReactions)

	switch cmd {
	case "_admin":

		if !didSenatorVote {
			msg, _ = s.ChannelMessageSend(m.ChannelID, "No Senators Voted In Favour")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if didSenatorVote && senateAmount == 0 {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Senate Disagreed")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if votes > 0 {
			return true
		}

		return false

	case "_kick":

		if !didSenatorVote {
			msg, _ = s.ChannelMessageSend(m.ChannelID, "No Senators Voted In Favour")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if didSenatorVote && senateAmount == 0 {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Senate Disagreed")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if votes > 0 {
			return true
		}

		return false

	case "_ban":

		if didSenatorVote && senateAmount == 0 {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Senate Disagreed")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if votes < 0 {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Petition Cancelled")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if votes == 0 {
			msg, _ := s.ChannelMessageSend(m.ChannelID, "Votes Tied")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if !didSenatorVote {
			msg, _ = s.ChannelMessageSend(m.ChannelID, "No Senators Voted In Favour")
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}

		if votes > 0 {
			return true
		}

		return false

	default:

		if len(checkReactions) <= len(crossReactions) {
			msg, err := s.ChannelMessageSend(m.ChannelID, "Petition Cancelled")
			if err != nil {
				fmt.Println(err)
			}
			go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
			return false
		}
	}

	return true

}

func countSenators(s *discordgo.Session, m *discordgo.MessageCreate, checkReactions []*discordgo.User, crossReactions []*discordgo.User) (int, bool) {
	var checkSenateAmount int
	var crossSenateAmount int
	var didSenateVote bool
	for _, user := range checkReactions {
		u, _ := s.GuildMember(m.GuildID, user.ID)
		for _, role := range u.Roles {
			if role == os.Getenv("SENATUSID") {
				checkSenateAmount++
				didSenateVote = true
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
	return checkSenateAmount - crossSenateAmount, didSenateVote
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

		channels, _ := s.GuildChannels(m.GuildID)
		for _, channel := range channels {
			fmt.Printf("Name:%v\nID:%v\n\n\n", channel.Name, channel.ID)
		}

		return

	case "_slave":

		if err := s.GuildMemberRoleAdd(m.GuildID, user, os.Getenv("SERVUSID")); err != nil {
			fmt.Println(err)
		}

		member, _ := s.User(user)
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Slaved %v", member.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

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
			return
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

		member, _ := s.User(user)
		s.GuildMemberDeleteWithReason(m.GuildID, user, "")
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Kicked %v", member.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return

	case "_ban":

		member, _ := s.User(user)
		if err := s.GuildBanCreateWithReason(m.GuildID, user, "", 3); err != nil {
			fmt.Println(err)
		}
		msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Banned %v", member.Mention()))
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()
		return

	case "_unban":

		member, _ := s.User(user)
		must(s.GuildBanDelete(m.GuildID, user))
		msg, someError := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully Unbanned %v", member.Mention()))
		must(someError)
		go func() { time.Sleep(time.Hour * 6); s.ChannelMessageDelete(msg.ChannelID, msg.ID) }()

		return

	}
}

func must(someErr error) {
	if someErr != nil {
		fmt.Println(someErr)
	}
}
