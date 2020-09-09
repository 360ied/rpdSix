package purgecommandpackage

import (
	"errors"
	"fmt"
	"strconv"

	"rpdSix/commands"
	"rpdSix/commands/checkedrun"
	"rpdSix/helpers/extendeddiscord/extendeddiscordpermissions"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:                         checkedrun.Builder(run, requiredPermissions...),
			Names:                       []string{"purge"},
			ExpectedPositionalArguments: expectedPositionalArguments,
		},
	)
}

var (
	requiredPermissions = []int{extendeddiscordpermissions.MANAGE_MESSAGES}

	expectedPositionalArguments = []string{"1", "2"}
)

func run(ctx commands.CommandContext) error {
	// command overloading
	var arg1, arg1Exists = ctx.Arguments["1"]
	if !arg1Exists {
		return errors.New(fmt.Sprint(commands.MissingArgumentErrorTemplate, "purge requires at least 1 argument"))
	}

	var arg2, arg2Exists = ctx.Arguments["2"]

	var messageDeleteErr = ctx.Message.Delete()
	if messageDeleteErr != nil {
		return messageDeleteErr
	}

	if !arg2Exists {
		return runPurgeUntil(ctx, arg1)
	}

	return runPurgeFromTo(ctx, arg1, arg2)
}

func runPurgeUntil(ctx commands.CommandContext, untilID string) error {
	for {
		var messages, messagesErr = ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 100, "", untilID, "")
		if messagesErr != nil {
			return messagesErr
		}

		if len(messages) == 0 {
			break
		}

		var messageIDs = make([]string, len(messages))
		for index, message := range messages {
			messageIDs[index] = message.ID
		}

		if len(messageIDs) < 100 {
			messageIDs = append(messageIDs, untilID)
		}

		var channelMessagesBulkDeleteErr = ctx.Session.ChannelMessagesBulkDelete(
			ctx.Message.ChannelID, messageIDs)
		if channelMessagesBulkDeleteErr != nil {
			return channelMessagesBulkDeleteErr
		}

		untilID = messageIDs[len(messageIDs)-1]
	}

	return nil
}

func runPurgeFromTo(ctx commands.CommandContext, fromID string, untilID string) error {
	var fromIDuint64, fromIDuint64Err = strconv.ParseUint(fromID, 10, 64)
	if fromIDuint64Err != nil {
		return fromIDuint64Err
	}

	for {
		var messages, messagesErr = ctx.Session.ChannelMessages(
			ctx.Message.ChannelID, 100, "", untilID, "")
		if messagesErr != nil {
			return messagesErr
		}
		if len(messages) == 0 {
			break
		}

		var messageIDs []string
		for _, message := range messages {
			var messageIDuint64, messageIDuint64Err = strconv.ParseUint(message.ID, 10, 64)
			if messageIDuint64Err != nil {
				return messageIDuint64Err
			}

			if messageIDuint64 > fromIDuint64 {
				continue
			}
			messageIDs = append(messageIDs, message.ID)
		}

		if len(messageIDs) == 0 {
			break
		}

		if len(messageIDs) < 100 {
			messageIDs = append(messageIDs, untilID)
		}

		var channelMessagesBulkDeleteErr = ctx.Session.ChannelMessagesBulkDelete(
			ctx.Message.ChannelID, messageIDs)
		if channelMessagesBulkDeleteErr != nil {
			return channelMessagesBulkDeleteErr
		}

		untilID = messageIDs[len(messageIDs)-1]
	}

	return nil
}
