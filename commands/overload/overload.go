package overload

import (
	"strings"

	"rpdSix/commands"
)

func Builder(
	functions map[int]func(commands.CommandContext, ...string) error,
) func(commands.CommandContext) error {
	return func(ctx commands.CommandContext) error {
		var arguments = strings.Split(ctx.Message.Content, commands.StringSeparator)[1:]

		var numberOfArguments = len(arguments)

		for argumentIndex, argument := range arguments {
			if strings.HasPrefix(argument, commands.KeywordArgumentPrefix) {
				numberOfArguments = argumentIndex
				break
			}
		}

		return functions[numberOfArguments](ctx, arguments[:numberOfArguments]...)
	}
}
