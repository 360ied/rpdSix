package helpcommand

import (
	"fmt"
	"strings"

	"github.com/ztrue/tracerr"

	"rpdSix/commandhandler"
)

func Initialize() {
	commandhandler.AddCommand(
		commandhandler.Command{
			Run:                         run,
			Names:                       []string{"help"},
			ExpectedPositionalArguments: []string{commandArg},
		},
	)
}

const (
	commandArg = "command"
)

func run(ctx commandhandler.CommandContext) error {
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		var command, exists_ = commandhandler.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		if exists_ {
			// lmao these variable names

			var formattedCommandNames []string

			for _, value := range command.Names {
				formattedCommandNames = append(
					formattedCommandNames,
					fmt.Sprintf("`%v`", value))
			}

			var formattedExpectedPositionalArguments []string

			for _, value := range command.ExpectedPositionalArguments {
				formattedExpectedPositionalArguments = append(
					formattedExpectedPositionalArguments,
					fmt.Sprintf("`%v`", value))
			}

			var formattedKeywordArgumentAliasesStringArray []string

			for key, value := range command.KeywordArgumentAliases {
				formattedKeywordArgumentAliasesStringArray = append(
					formattedKeywordArgumentAliasesStringArray,
					fmt.Sprintf("`%v`: `%v`", key, value))
			}

			var _, err = ctx.Message.Reply(
				fmt.Sprint(
					strings.Join(formattedCommandNames, ", "),
					":\nExpected Positional Arguments: ",
					strings.Join(formattedExpectedPositionalArguments, ", "),
					"\nKeyword Argument Aliases: ",
					strings.Join(formattedKeywordArgumentAliasesStringArray, ", ")))
			return tracerr.Wrap(err)
		} else {
			// command doesn't exist
			var _, err = ctx.Message.Reply("Command not found!")
			return tracerr.Wrap(err)
		}
	} else {
		// list commands
		var outputStr = "Commands:\n"
		for commandName := range commandhandler.Commands {
			outputStr += fmt.Sprintf("`%v`, ", commandName)
		}
		outputStr = outputStr[:len(outputStr)-len(", ")]
		var _, err = ctx.Message.Reply(outputStr)
		return tracerr.Wrap(err)
	}
}
