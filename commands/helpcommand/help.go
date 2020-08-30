package helpcommand

import (
	"fmt"
	"rpdSix/commands"
	"strings"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:                         run,
			Names:                       []string{"help"},
			ExpectedPositionalArguments: []string{commandArg},
		},
	)
}

const (
	commandArg = "command"
)

func run(ctx commands.CommandContext) error {
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		var command, exists_ = commands.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		if exists_ {
			// lmao these variable names

			var formattedCommandNames []string

			for _, value := range command.Names {
				formattedCommandNames = append(
					formattedCommandNames,
					fmt.Sprint("`", value, "`"))
			}

			var formattedExpectedPositionalArguments []string

			for _, value := range command.ExpectedPositionalArguments {
				formattedExpectedPositionalArguments = append(
					formattedExpectedPositionalArguments,
					fmt.Sprint("`", value, "`"))
			}

			var formattedKeywordArgumentAliasesStringArray []string

			for key, value := range command.KeywordArgumentAliases {
				formattedKeywordArgumentAliasesStringArray = append(
					formattedKeywordArgumentAliasesStringArray,
					fmt.Sprint("`", key, "`: `", value, "`"))
			}

			var _, err = ctx.Session.ChannelMessageSend(
				ctx.Message.ChannelID,
				fmt.Sprint(
					strings.Join(formattedCommandNames, ", "),
					":\nExpected Positional Arguments: ",
					strings.Join(formattedExpectedPositionalArguments, ", "),
					"\nKeyword Argument Aliases: ",
					strings.Join(formattedKeywordArgumentAliasesStringArray, ", ")))
			return err
		} else {
			// command doesn't exist
			var _, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Command not found!")
			return err
		}
	} else {
		// list commands
		var outputStr = "Commands:\n"
		for commandName := range commands.Commands {
			outputStr += fmt.Sprint("`", commandName, "`, ")
		}
		outputStr = outputStr[:len(outputStr)-len(", ")]
		var _, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, outputStr)
		return err
	}
}
