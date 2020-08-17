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
			ExpectedPositionalArguments: []string{"toSay"},
		},
	)
}

func run(ctx commands.CommandContext) {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.Arguments["toSay"])
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
}
