package purgecommand

import (
	"strconv"

	"github.com/ztrue/tracerr"

	"rpdSix/commandhandler"
	"rpdSix/commandhandler/helpers/checkedrun"
	"rpdSix/commandhandler/helpers/overload"
	"rpdSix/helpers/extendeddiscord/extendeddiscordpermissions"
)

func Initialize() {
	commandhandler.AddCommand(
		commandhandler.Command{
			Run: checkedrun.Builder(
				overload.Builder(map[int]func(commandhandler.CommandContext, ...string) error{
					1: runPurgeUntil,
					2: runPurgeFromTo,
				}), requiredPermissions...),
			Names: []string{"purge"},
		},
	)
}

var (
	requiredPermissions = []int{extendeddiscordpermissions.MANAGE_MESSAGES}
)

func runPurgeUntil(ctx commandhandler.CommandContext, args ...string) error {
	var untilID = args[0]
	for {
		var messages, messagesErr = ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 100, "", untilID, "")
		if messagesErr != nil {
			return tracerr.Wrap(messagesErr)
		}

		if len(messages) == 0 {
			break
		}

		var messageIDs = make([]string, len(messages))
		for index, message := range messages {
			messageIDs[index] = message.ID
		}

		var channelMessagesBulkDeleteErr = ctx.Session.ChannelMessagesBulkDelete(
			ctx.Message.ChannelID, messageIDs)
		if channelMessagesBulkDeleteErr != nil {
			return tracerr.Wrap(channelMessagesBulkDeleteErr)
		}

		untilID = messageIDs[len(messageIDs)-1]
	}

	_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, untilID)
	return nil
}

func runPurgeFromTo(ctx commandhandler.CommandContext, args ...string) error {
	var fromID = args[0]
	var untilID = args[1]

	var fromIDuint64, fromIDuint64Err = strconv.ParseUint(fromID, 10, 64)
	if fromIDuint64Err != nil {
		return tracerr.Wrap(fromIDuint64Err)
	}

	for {
		var messages, messagesErr = ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 100, "", untilID, "")
		if messagesErr != nil {
			return tracerr.Wrap(messagesErr)
		}

		if len(messages) == 0 {
			break
		}

		var messageIDs []string
		for _, message := range messages {
			var messageIDuint64, messageIDuint64Err = strconv.ParseUint(message.ID, 10, 64)
			if messageIDuint64Err != nil {
				return tracerr.Wrap(messageIDuint64Err)
			}

			if messageIDuint64 > fromIDuint64 {
				continue
			}

			messageIDs = append(messageIDs, message.ID)
		}

		if len(messageIDs) == 0 {
			break
		}

		var channelMessagesBulkDeleteErr = ctx.Session.ChannelMessagesBulkDelete(
			ctx.Message.ChannelID, messageIDs)
		if channelMessagesBulkDeleteErr != nil {
			return tracerr.Wrap(channelMessagesBulkDeleteErr)
		}

		untilID = messageIDs[len(messageIDs)-1]
	}

	_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, untilID)
	return nil
}
