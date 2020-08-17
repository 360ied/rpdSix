package pingcommand

import (
	"github.com/3sixtyied/rpdSix/commands"
	"github.com/ztrue/tracerr"
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
