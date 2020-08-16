package pingcommand

import (
	"github.com/ztrue/tracerr"
	"rpdSix/commands"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:   run,
			Names: []string{"ping"},
		},
	)
}

func run(ctx commands.CommandContext) {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Pong!")
	if err != nil {
		tracerr.PrintSourceColor(err)
	}
}
