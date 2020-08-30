package saycommand

import (
	"rpdSix/commands"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:                         run,
			Names:                       []string{"say"},
			ExpectedPositionalArguments: []string{toSayArg},
		},
	)
}

const (
	toSayArg = "toSay"
)

func run(ctx commands.CommandContext) error {
	var _, err = ctx.Message.Reply(ctx.Arguments[toSayArg])
	return err
}
