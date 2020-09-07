package cache

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RegisterEventHandlers(bot *discordgo.Session) {
	bot.AddHandler(channelCreateEventHandler)
	bot.AddHandler(channelDeleteEventHandler)
	bot.AddHandler(channelPinsUpdateEventHandler)
	bot.AddHandler(channelUpdateEventHandler)
	bot.AddHandler(connectEventHandler)
	bot.AddHandler(disconnectEventHandler)
	bot.AddHandler(eventEventHandler)
	bot.AddHandler(guildBanAddEventHandler)
	bot.AddHandler(guildBanRemoveEventHandler)
	bot.AddHandler(guildCreateEventHandler)
	bot.AddHandler(guildDeleteEventHandler)
	bot.AddHandler(guildEmojisUpdateEventHandler)
	bot.AddHandler(guildIntegrationsUpdateEventHandler)
	bot.AddHandler(guildMemberAddEventHandler)
	bot.AddHandler(guildMemberRemoveEventHandler)
	bot.AddHandler(guildMemberUpdateEventHandler)
	bot.AddHandler(guildMembersChunkEventHandler)
	bot.AddHandler(guildRoleCreateEventHandler)
	bot.AddHandler(guildRoleDeleteEventHandler)
	bot.AddHandler(guildRoleUpdateEventHandler)
	bot.AddHandler(guildUpdateEventHandler)
	bot.AddHandler(messageAckEventHandler)
	bot.AddHandler(messageCreateEventHandler)
	bot.AddHandler(messageDeleteEventHandler)
	bot.AddHandler(messageDeleteBulkEventHandler)
	bot.AddHandler(messageReactionAddEventHandler)
	bot.AddHandler(messageReactionRemoveEventHandler)
	bot.AddHandler(messageReactionRemoveAllEventHandler)
	bot.AddHandler(messageUpdateEventHandler)
	bot.AddHandler(presenceUpdateEventHandler)
	bot.AddHandler(presencesReplaceEventHandler)
	bot.AddHandler(rateLimitEventHandler)
	bot.AddHandler(readyEventHandler)
	bot.AddHandler(relationshipAddEventHandler)
	bot.AddHandler(relationshipRemoveEventHandler)
	bot.AddHandler(resumedEventHandler)
	bot.AddHandler(typingStartEventHandler)
	bot.AddHandler(userGuildSettingsUpdateEventHandler)
	bot.AddHandler(userNoteUpdateEventHandler)
	bot.AddHandler(userSettingsUpdateEventHandler)
	bot.AddHandler(userUpdateEventHandler)
	bot.AddHandler(voiceServerUpdateEventHandler)
	bot.AddHandler(voiceStateUpdateEventHandler)
	bot.AddHandler(webhooksUpdateEventHandler)
}

func channelCreateEventHandler(session *discordgo.Session, event *discordgo.ChannelCreate) {

}

func channelDeleteEventHandler(session *discordgo.Session, event *discordgo.ChannelDelete) {

}

func channelPinsUpdateEventHandler(session *discordgo.Session, event *discordgo.ChannelPinsUpdate) {

}

func channelUpdateEventHandler(session *discordgo.Session, event *discordgo.ChannelUpdate) {

}

func connectEventHandler(session *discordgo.Session, event *discordgo.Connect) {

}

func disconnectEventHandler(session *discordgo.Session, event *discordgo.Disconnect) {

}

func eventEventHandler(session *discordgo.Session, event *discordgo.Event) {

}

func guildBanAddEventHandler(session *discordgo.Session, event *discordgo.GuildBanAdd) {

}

func guildBanRemoveEventHandler(session *discordgo.Session, event *discordgo.GuildBanRemove) {

}

func guildCreateEventHandler(session *discordgo.Session, event *discordgo.GuildCreate) {
	Cache.GuildsRWMutex.Lock()
	defer Cache.GuildsRWMutex.Unlock()

	Cache.Guilds[event.ID] = event.Guild
}

func guildDeleteEventHandler(session *discordgo.Session, event *discordgo.GuildDelete) {

}

func guildEmojisUpdateEventHandler(session *discordgo.Session, event *discordgo.GuildEmojisUpdate) {

}

func guildIntegrationsUpdateEventHandler(session *discordgo.Session, event *discordgo.GuildIntegrationsUpdate) {

}

func guildMemberAddEventHandler(session *discordgo.Session, event *discordgo.GuildMemberAdd) {

}

func guildMemberRemoveEventHandler(session *discordgo.Session, event *discordgo.GuildMemberRemove) {

}

func guildMemberUpdateEventHandler(session *discordgo.Session, event *discordgo.GuildMemberUpdate) {

}

func guildMembersChunkEventHandler(session *discordgo.Session, event *discordgo.GuildMembersChunk) {

}

func guildRoleCreateEventHandler(session *discordgo.Session, event *discordgo.GuildRoleCreate) {

}

func guildRoleDeleteEventHandler(session *discordgo.Session, event *discordgo.GuildRoleDelete) {

}

func guildRoleUpdateEventHandler(session *discordgo.Session, event *discordgo.GuildRoleUpdate) {

}

func guildUpdateEventHandler(session *discordgo.Session, event *discordgo.GuildUpdate) {

}

func messageAckEventHandler(session *discordgo.Session, event *discordgo.MessageAck) {

}

func messageCreateEventHandler(session *discordgo.Session, event *discordgo.MessageCreate) {

}

func messageDeleteEventHandler(session *discordgo.Session, event *discordgo.MessageDelete) {

}

func messageDeleteBulkEventHandler(session *discordgo.Session, event *discordgo.MessageDeleteBulk) {

}

func messageReactionAddEventHandler(session *discordgo.Session, event *discordgo.MessageReactionAdd) {

}

func messageReactionRemoveEventHandler(session *discordgo.Session, event *discordgo.MessageReactionRemove) {

}

func messageReactionRemoveAllEventHandler(session *discordgo.Session, event *discordgo.MessageReactionRemoveAll) {

}

func messageUpdateEventHandler(session *discordgo.Session, event *discordgo.MessageUpdate) {

}

func presenceUpdateEventHandler(session *discordgo.Session, event *discordgo.PresenceUpdate) {

}

func presencesReplaceEventHandler(session *discordgo.Session, event *discordgo.PresencesReplace) {

}

func rateLimitEventHandler(session *discordgo.Session, event *discordgo.RateLimit) {

}

func readyEventHandler(_ *discordgo.Session, event *discordgo.Ready) {

	fmt.Println("Filling cache...")

	func() {
		Cache.GuildsRWMutex.Lock()
		defer Cache.GuildsRWMutex.Unlock()

		fmt.Println("Filling guilds...")

		for _, guild := range event.Guilds {
			Cache.Guilds[guild.ID] = guild
		}
	}() // Fill guilds

	func() {
		Cache.PrivateChannelsRWMutex.Lock()
		defer Cache.PrivateChannelsRWMutex.Unlock()

		fmt.Println("Filling private channels...")

		for _, privateChannel := range event.PrivateChannels {
			Cache.PrivateChannels[privateChannel.ID] = privateChannel
		}
	}() // Fill private channels

	func() {
		Cache.UserRWMutex.Lock()
		defer Cache.UserRWMutex.Unlock()

		fmt.Println("Filling user...")

		Cache.User = event.User
	}() // Fill user

	func() {
		Cache.NotesRWMutex.Lock()
		defer Cache.NotesRWMutex.Unlock()

		fmt.Println("Filling notes...")

		for key, value := range event.Notes {
			Cache.Notes[key] = value
		}
	}() // Fill notes

	func() {
		Cache.PresencesRWMutex.Lock()
		defer Cache.PresencesRWMutex.Unlock()

		fmt.Println("Filling Presences...")

		for _, presence := range event.Presences {
			Cache.Presences[presence.User.ID] = presence
		}
	}() // Fill presences

	func() {
		Cache.ReadStatesRWMutex.Lock()
		defer Cache.ReadStatesRWMutex.Unlock()

		fmt.Println("Filling ReadStates...")

		for _, readState := range event.ReadState {
			Cache.ReadStates[readState.ID] = readState
		}
	}() // Fill read states

	func() {
		Cache.RelationshipsRWMutex.Lock()
		defer Cache.RelationshipsRWMutex.Unlock()

		fmt.Println("Filling Relationships...")

		for _, relationship := range event.Relationships {
			Cache.Relationships[relationship.ID] = relationship
		}
	}() // Fill relationships

	func() {
		Cache.SessionIDRWMutex.Lock()
		defer Cache.SessionIDRWMutex.Unlock()

		fmt.Println("Filling SessionID...")

		Cache.SessionID = event.SessionID
	}() // Fill session ID

	func() {
		Cache.UserGuildSettingsRWMutex.Lock()
		defer Cache.UserGuildSettingsRWMutex.Unlock()

		fmt.Println("Filling UserGuildSettings...")

		for _, userGuildSettings := range event.UserGuildSettings {
			Cache.UserGuildSettings[userGuildSettings.GuildID] = userGuildSettings
		}
	}() // Fill user guild settings

	func() {
		Cache.SettingsRWMutex.Lock()
		defer Cache.SettingsRWMutex.Unlock()

		fmt.Println("Filling Settings...")

		Cache.Settings = event.Settings
	}() // Fill settings

	func() {
		Cache.VersionRWMutex.Lock()
		defer Cache.VersionRWMutex.Unlock()

		fmt.Println("Filling Version...")

		Cache.Version = event.Version
	}() // Fill version

	fmt.Println("Finished caching Ready event.")
}

func relationshipAddEventHandler(session *discordgo.Session, event *discordgo.RelationshipAdd) {

}

func relationshipRemoveEventHandler(session *discordgo.Session, event *discordgo.RelationshipRemove) {

}

func resumedEventHandler(session *discordgo.Session, event *discordgo.Resumed) {

}

func typingStartEventHandler(session *discordgo.Session, event *discordgo.TypingStart) {

}

func userGuildSettingsUpdateEventHandler(session *discordgo.Session, event *discordgo.UserGuildSettingsUpdate) {

}

func userNoteUpdateEventHandler(session *discordgo.Session, event *discordgo.UserNoteUpdate) {

}

func userSettingsUpdateEventHandler(session *discordgo.Session, event *discordgo.UserSettingsUpdate) {

}

func userUpdateEventHandler(session *discordgo.Session, event *discordgo.UserUpdate) {

}

func voiceServerUpdateEventHandler(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {

}

func voiceStateUpdateEventHandler(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	Cache.GuildsRWMutex.Lock()
	defer Cache.GuildsRWMutex.Unlock()

	var guild = Cache.Guilds[event.GuildID]

	for i, voiceState := range guild.VoiceStates {
		if voiceState.UserID == event.UserID {
			guild.VoiceStates[i] = voiceState
			return
		}
	}

	guild.VoiceStates = append(guild.VoiceStates, event.VoiceState)
}

func webhooksUpdateEventHandler(session *discordgo.Session, event *discordgo.WebhooksUpdate) {

}
