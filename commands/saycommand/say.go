package saycommand

import (
	"github.com/ztrue/tracerr"
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

func run(ctx commands.CommandContext) {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.Arguments[toSayArg])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
}
