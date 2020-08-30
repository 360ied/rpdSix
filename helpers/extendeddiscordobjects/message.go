package extendeddiscordobjects

import "github.com/bwmarrin/discordgo"

type ExtendedMessage struct {
	*discordgo.Message
	session *discordgo.Session
}

func ExtendMessage(message *discordgo.Message, session *discordgo.Session) *ExtendedMessage {
	return &ExtendedMessage{message, session}
}

func (message ExtendedMessage) Guild() (*discordgo.Guild, error) {
	return message.session.Guild(message.GuildID)
}

func (message ExtendedMessage) Reply(content string) (*discordgo.Message, error) {
	return message.session.ChannelMessageSend(message.ChannelID, content)
}

func (message ExtendedMessage) Channel() (*discordgo.Channel, error) {
	return message.session.Channel(message.ChannelID)
}

func (message ExtendedMessage) AuthorMember() (*discordgo.Member, error) {
	var messageGuild, err = message.Guild()
	if err != nil {
		return nil, err
	}
	for _, member := range messageGuild.Members {
		if member.User.ID == message.Author.ID {
			return member, nil
		}
	}
	return nil, MemberNotFoundError
}
