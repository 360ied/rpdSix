package extendeddiscordobjects

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"rpdSix/helpers"
)

type ExtendedMessage struct {
	*discordgo.Message
	session *discordgo.Session
}

func ExtendMessage(message *discordgo.Message, session *discordgo.Session) *ExtendedMessage {
	return &ExtendedMessage{
		Message: message,
		session: session,
	}
}

// short form for message.session.Guild(message.GuildID)
func (message *ExtendedMessage) Guild() (*discordgo.Guild, error) {
	return message.session.Guild(message.GuildID)
}

// short form for message.session.ChannelMessageSend(message.ChannelID, content)
func (message *ExtendedMessage) Reply(content string) (*discordgo.Message, error) {
	return message.session.ChannelMessageSend(message.ChannelID, content)
}

// short form for message.session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{...})
func (message *ExtendedMessage) ComplexReply(send *discordgo.MessageSend) (*discordgo.Message, error) {
	return message.session.ChannelMessageSendComplex(message.ChannelID, send)
}

// short form for message.session.Channel(message.ChannelID)
func (message *ExtendedMessage) Channel() (*discordgo.Channel, error) {
	return message.session.Channel(message.ChannelID)
}

func (message *ExtendedMessage) AuthorMember() (*discordgo.Member, error) {
	var messageGuild, err = message.Guild()
	if err != nil {
		return nil, err
	}
	for _, member := range messageGuild.Members {
		if member.User.ID == message.Author.ID {
			return member, nil
		}
	}
	return nil, errors.New(fmt.Sprint(helpers.MemberNotFoundErrorTemplate,
		"member ", message.Author.ID, " not found in guild ", message.GuildID))
}
