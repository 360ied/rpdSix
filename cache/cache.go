package cache

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	Cache = cache{
		GuildsRWMutex:            &sync.RWMutex{},
		Guilds:                   make(map[string]*cacheGuild),
		PrivateChannelsRWMutex:   &sync.RWMutex{},
		PrivateChannels:          make(map[string]*discordgo.Channel),
		UserRWMutex:              &sync.RWMutex{},
		User:                     nil,
		NotesRWMutex:             &sync.RWMutex{},
		Notes:                    make(map[string]string),
		PresencesRWMutex:         &sync.RWMutex{},
		Presences:                make(map[string]*discordgo.Presence),
		ReadStatesRWMutex:        &sync.RWMutex{},
		ReadStates:               make(map[string]*discordgo.ReadState),
		RelationshipsRWMutex:     &sync.RWMutex{},
		Relationships:            make(map[string]*discordgo.Relationship),
		SessionIDRWMutex:         &sync.RWMutex{},
		SessionID:                "",
		UserGuildSettingsRWMutex: &sync.RWMutex{},
		UserGuildSettings:        make(map[string]*discordgo.UserGuildSettings),
		SettingsRWMutex:          &sync.RWMutex{},
		Settings:                 nil,
		VersionRWMutex:           &sync.RWMutex{},
		Version:                  -1,
	}
)

type cache struct {
	// Session *discordgo.Session

	GuildsRWMutex *sync.RWMutex
	Guilds        map[string]*cacheGuild // key is guild ID

	PrivateChannelsRWMutex *sync.RWMutex
	PrivateChannels        map[string]*discordgo.Channel // key is channel ID

	UserRWMutex *sync.RWMutex
	User        *discordgo.User

	NotesRWMutex *sync.RWMutex
	Notes        map[string]string // key is prob user ID

	PresencesRWMutex *sync.RWMutex
	Presences        map[string]*discordgo.Presence // User ID as key

	ReadStatesRWMutex *sync.RWMutex
	ReadStates        map[string]*discordgo.ReadState // the ID value of the value

	RelationshipsRWMutex *sync.RWMutex
	Relationships        map[string]*discordgo.Relationship // the ID value of the value

	SessionIDRWMutex *sync.RWMutex
	SessionID        string // the session ID

	UserGuildSettingsRWMutex *sync.RWMutex
	UserGuildSettings        map[string]*discordgo.UserGuildSettings // key is guild ID

	SettingsRWMutex *sync.RWMutex
	Settings        *discordgo.Settings // A Settings stores data for a specific users Discord client settings.

	VersionRWMutex *sync.RWMutex
	Version        int // field from ready event
}

type cacheGuild struct {
	*discordgo.Guild
	*sync.RWMutex
}

func cacheifyGuild(guild *discordgo.Guild) *cacheGuild {
	return &cacheGuild{
		Guild:   guild,
		RWMutex: &sync.RWMutex{},
	}
}
