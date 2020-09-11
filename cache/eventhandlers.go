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

func channelCreateEventHandler(_ *discordgo.Session, event *discordgo.ChannelCreate) {
	// a channel is created
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	guildUnion.Channels = append(guildUnion.Channels, event.Channel)
}

func channelDeleteEventHandler(_ *discordgo.Session, event *discordgo.ChannelDelete) {
	// a channel is deleted
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	for channelIndex, channel := range guildUnion.Channels {
		if channel.ID == event.ID {
			// remove the channel
			// swap the member with the last channel
			guildUnion.Channels[channelIndex] = guildUnion.Channels[len(guildUnion.Channels)-1]
			// slice off the last channel
			guildUnion.Channels = guildUnion.Channels[:len(guildUnion.Channels)-1]
			break
		}
	}
}

func channelPinsUpdateEventHandler(_ *discordgo.Session, event *discordgo.ChannelPinsUpdate) {
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	// TODO: I think there's a way to hold the lock for less time, but I'm not so sure about it.
	guildUnion.Lock()
	defer guildUnion.Unlock()

	for _, channel := range guildUnion.Channels {
		if channel.ID == event.ChannelID {
			channel.LastPinTimestamp = discordgo.Timestamp(event.LastPinTimestamp)
			break
		}
	}
}

func channelUpdateEventHandler(_ *discordgo.Session, event *discordgo.ChannelUpdate) {
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	// TODO: I think there's a way to hold the lock for less time, but I'm not so sure about it.
	guildUnion.Lock()
	defer guildUnion.Unlock()

	for channelIndex, channel := range guildUnion.Channels {
		if channel.ID == event.ID {
			guildUnion.Channels[channelIndex] = event.Channel
			break
		}
	}
}

func connectEventHandler(*discordgo.Session, *discordgo.Connect) {
	// synthetic event
}

func disconnectEventHandler(*discordgo.Session, *discordgo.Disconnect) {
	// synthetic event
}

func eventEventHandler(*discordgo.Session, *discordgo.Event) {
	// redundant
}

func guildBanAddEventHandler(*discordgo.Session, *discordgo.GuildBanAdd) {
	// I don't see anywhere I could store this or why I even should store it.
}

func guildBanRemoveEventHandler(*discordgo.Session, *discordgo.GuildBanRemove) {
	// I don't see anywhere I could store this or why I even should store it.
}

func guildCreateEventHandler(_ *discordgo.Session, event *discordgo.GuildCreate) {
	Cache.GuildsRWMutex.Lock()
	defer Cache.GuildsRWMutex.Unlock()

	Cache.Guilds[event.ID] = cacheifyGuild(event.Guild)
}

func guildDeleteEventHandler(_ *discordgo.Session, event *discordgo.GuildDelete) {
	Cache.GuildsRWMutex.Lock()
	defer Cache.GuildsRWMutex.Unlock()

	Cache.Guilds[event.ID] = cacheifyGuild(event.Guild)
}

func guildEmojisUpdateEventHandler(_ *discordgo.Session, event *discordgo.GuildEmojisUpdate) {
	// emojis got updated
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	var emojiMap = map[string]*discordgo.Emoji{}

	for _, emoji := range event.Emojis {
		emojiMap[emoji.ID] = emoji
	}

	// only lock here as locking above is unneeded
	guildUnion.Lock()
	defer guildUnion.Unlock()

	for emojiIndex, emoji := range guildUnion.Emojis {
		if updatedEmoji, in := emojiMap[emoji.ID]; in {
			guildUnion.Emojis[emojiIndex] = updatedEmoji
			delete(emojiMap, emoji.ID)
		}
	}

	for _, emoji := range emojiMap {
		guildUnion.Emojis = append(guildUnion.Emojis, emoji)
	}
}

func guildIntegrationsUpdateEventHandler(*discordgo.Session, *discordgo.GuildIntegrationsUpdate) {
	/*
		type GuildIntegrationsUpdate struct {
			GuildID string `json:"guild_id"`
		}

		...

		What am i supposed to do with this?
	*/
}

func guildMemberAddEventHandler(_ *discordgo.Session, event *discordgo.GuildMemberAdd) {
	// new user joined the guild, so don't need to check whether to update existing user or not
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	guildUnion.Members = append(guildUnion.Members, event.Member)
}

func guildMemberRemoveEventHandler(_ *discordgo.Session, event *discordgo.GuildMemberRemove) {
	// member left guild
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	// TODO: I think there's a way to hold the lock for less time, but I'm not so sure about it.
	guildUnion.Lock()
	defer guildUnion.Unlock()

	for memberIndex, member := range guildUnion.Members {
		if member.User.ID == event.User.ID {
			// remove the member
			// swap the member with the last member
			guildUnion.Members[memberIndex] = guildUnion.Members[len(guildUnion.Members)-1]
			// slice off the last member
			guildUnion.Members = guildUnion.Members[:len(guildUnion.Members)-1]
			break
		}
	}
}

func guildMemberUpdateEventHandler(_ *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	// existing member is updated
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	// TODO: Consider using a different method as this does not scale with large guilds

	for memberIndex, member := range guildUnion.Members {
		if member.User.ID == event.User.ID {
			guildUnion.Members[memberIndex] = event.Member
			break
		}
	}
}

func guildMembersChunkEventHandler(*discordgo.Session, *discordgo.GuildMembersChunk) {
	/*
		From Discord API Documentation:
		Sent in response to Guild Request Members.
		You can use the chunk_index and chunk_count to calculate how many chunks are left for your request.

		...

		What am I supposed to do with this?
	*/
}

func guildRoleCreateEventHandler(_ *discordgo.Session, event *discordgo.GuildRoleCreate) {
	// a new role is created
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	guildUnion.Roles = append(guildUnion.Roles, event.Role)
}

func guildRoleDeleteEventHandler(_ *discordgo.Session, event *discordgo.GuildRoleDelete) {
	// a role is deleted
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	// TODO: I think there's a way to hold the lock for less time, but I'm not so sure about it.
	guildUnion.Lock()
	defer guildUnion.Unlock()

	for roleIndex, role := range guildUnion.Roles {
		if role.ID == event.RoleID {
			// remove the role
			// swap the member with the last role
			guildUnion.Roles[roleIndex] = guildUnion.Roles[len(guildUnion.Roles)-1]
			// slice off the last role
			guildUnion.Roles = guildUnion.Roles[:len(guildUnion.Roles)-1]
			break
		}
	}
}

func guildRoleUpdateEventHandler(_ *discordgo.Session, event *discordgo.GuildRoleUpdate) {
	// existing role is updated
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.GuildID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	// TODO: This might not be very efficient

	for roleIndex, role := range guildUnion.Roles {
		if role.ID == event.Role.ID {
			guildUnion.Roles[roleIndex] = event.Role
			break
		}
	}
}

func guildUpdateEventHandler(_ *discordgo.Session, event *discordgo.GuildUpdate) {
	// entire guild is updated
	// but some elements are not included in the event
	var guildUnion = func() cacheGuild {
		Cache.GuildsRWMutex.RLock()
		defer Cache.GuildsRWMutex.RUnlock()

		return Cache.Guilds[event.ID]
	}()

	guildUnion.Lock()
	defer guildUnion.Unlock()

	// cache elements not included in guild update event
	var joinedAt = guildUnion.JoinedAt
	var large = guildUnion.Large
	var unavailable = guildUnion.Unavailable
	var memberCount = guildUnion.MemberCount
	var voiceStates = guildUnion.VoiceStates
	var members = guildUnion.Members
	var channels = guildUnion.Channels
	var presences = guildUnion.Presences

	// set guild to the updated guild
	guildUnion.Guild = event.Guild

	// restore elements not included in guild update event
	guildUnion.JoinedAt = joinedAt
	guildUnion.Large = large
	guildUnion.Unavailable = unavailable
	guildUnion.MemberCount = memberCount
	guildUnion.VoiceStates = voiceStates
	guildUnion.Members = members
	guildUnion.Channels = channels
	guildUnion.Presences = presences
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
			Cache.Guilds[guild.ID] = cacheifyGuild(guild)
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

func voiceStateUpdateEventHandler(_ *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	Cache.GuildsRWMutex.Lock()
	defer Cache.GuildsRWMutex.Unlock()

	var guild = Cache.Guilds[event.GuildID]

	for i, voiceState := range guild.VoiceStates {
		if voiceState.UserID == event.UserID {
			if event.ChannelID == "" { // member disconnected
				// remove voice state from cache
				// swap voice state with last voice state
				guild.VoiceStates[i] = guild.VoiceStates[len(guild.VoiceStates)-1]
				// slice off the last voice state
				guild.VoiceStates = guild.VoiceStates[:len(guild.VoiceStates)-1]
			} else { // member voice state changed
				guild.VoiceStates[i] = voiceState
			}
			return
		}
	}

	// voice state is new so add it
	guild.VoiceStates = append(guild.VoiceStates, event.VoiceState)
}

func webhooksUpdateEventHandler(session *discordgo.Session, event *discordgo.WebhooksUpdate) {

}
