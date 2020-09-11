package voicechannelmoveallcommand

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/ztrue/tracerr"

	"rpdSix/cache"
	"rpdSix/commands"
	"rpdSix/commands/checkedrun"
	"rpdSix/helpers/extendeddiscord/extendeddiscordpermissions"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:                         checkedrun.Builder(run, requiredPermissions...),
			Names:                       []string{"voicechannelmoveall", "vcmoveall", "vcmall", "vcma"},
			ExpectedPositionalArguments: []string{destinationChannelArg},
		},
	)
}

const (
	destinationChannelArg = "destinationChannelID"

	memberNotInVoiceChannelErrorTemplate = "VoiceChannelMoveAll memberNotInVoiceChannelError, "
)

var (
	requiredPermissions = []int{extendeddiscordpermissions.MOVE_MEMBERS}
)

func run(ctx commands.CommandContext) error {

	// var messageGuild, messageGuildErr = ctx.Message.Guild()
	// if messageGuildErr != nil {
	// 	return messageGuildErr
	// }

	var authorVoiceState *discordgo.VoiceState

	var cachedGuild, cachedGuildRWMutex = func() (*discordgo.Guild, *sync.RWMutex) {
		cache.Cache.GuildsRWMutex.RLock()
		defer cache.Cache.GuildsRWMutex.RUnlock()

		var get = cache.Cache.Guilds[ctx.Message.GuildID]

		return get.Guild, get.RWMutex
	}()

	cachedGuildRWMutex.RLock()
	defer cachedGuildRWMutex.RUnlock()

	for _, voiceState := range cachedGuild.VoiceStates {
		// fmt.Println(voiceState.UserID)
		if voiceState.UserID == ctx.Message.Author.ID {
			authorVoiceState = voiceState
			goto foundVoiceState
		}
	}

	return tracerr.Wrap(errors.New(fmt.Sprint(memberNotInVoiceChannelErrorTemplate,
		"It does not appear that you are in a voice channel.")))

foundVoiceState:

	var toMove []*discordgo.VoiceState

	for _, voiceState := range cachedGuild.VoiceStates {
		if voiceState.ChannelID == authorVoiceState.ChannelID {
			toMove = append(toMove, voiceState)
		}
	}

	var destinationChannelID, destinationChannelArgExists = ctx.Arguments[destinationChannelArg]
	if !destinationChannelArgExists {
		return tracerr.Wrap(errors.New(fmt.Sprint(commands.MissingArgumentErrorTemplate,
			"Destination Channel ID not found!")))
	}

	for _, memberToMove := range toMove {

		var guildMemberMoveErr = ctx.Session.GuildMemberMove(
			ctx.Message.GuildID,
			memberToMove.UserID,
			&destinationChannelID)

		if guildMemberMoveErr != nil {
			_, _ = ctx.Message.Reply(fmt.Sprintf(
				"An Error occured while moving member ID %v\n%v", memberToMove.UserID, guildMemberMoveErr))
		}
	}

	var _, replyErr = ctx.Message.Reply("Done.")
	return tracerr.Wrap(replyErr)
}
