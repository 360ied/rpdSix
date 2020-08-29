package helpcommand

import (
	"fmt"
	"github.com/ztrue/tracerr"
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

func run(ctx commands.CommandContext) {
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		command, exists_ := commands.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		if exists_ {
			// lmao
			var keywordArgumentAliasesStringArray []string

			for key, value := range command.KeywordArgumentAliases {
				keywordArgumentAliasesStringArray = append(
					keywordArgumentAliasesStringArray,
					fmt.Sprint(key, ": ", value))
			}

			_, err := ctx.Session.ChannelMessageSend(
				ctx.Message.ChannelID,
				fmt.Sprint(
					strings.Join(command.Names, ", "),
					"\nExpected Positional Arguments: ",
					strings.Join(command.ExpectedPositionalArguments, ", "),
					"\nKeyword Argument Aliases: ",
					strings.Join(keywordArgumentAliasesStringArray, ", ")))
			if err != nil {
				tracerr.PrintSourceColor(err)
			}
		} else {
			// command doesn't exist
			_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Command not found!")
			if err != nil {
				tracerr.PrintSourceColor(err)
			}
		}
	} else {
		// list commands
		var outputStr = "Commands:\n"
		for commandName := range commands.Commands {
			outputStr += fmt.Sprint("`", commandName, "`, ")
		}
		outputStr = outputStr[:len(outputStr)-len(", ")]
		_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, outputStr)
		if err != nil {
			tracerr.PrintSourceColor(err)
		}
	}
}
