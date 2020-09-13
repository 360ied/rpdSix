package overload

import (
	"strings"

	"rpdSix/commandhandler"
)

func Builder(
	functions map[int]func(commandhandler.CommandContext, ...string) error,
) func(commandhandler.CommandContext) error {

	return func(ctx commandhandler.CommandContext) error {
		var arguments = strings.Split(ctx.Message.Content, commandhandler.StringSeparator)[1:]

		var numberOfArguments = len(arguments)

		for argumentIndex, argument := range arguments {
			if strings.HasPrefix(argument, commandhandler.KeywordArgumentPrefix) {
				numberOfArguments = argumentIndex
				break
			}
		}

		return functions[numberOfArguments](ctx, arguments[:numberOfArguments]...)
	}
}
