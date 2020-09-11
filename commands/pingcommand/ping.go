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

func run(ctx commands.CommandContext) error {
	var _, err = ctx.Message.Reply("Pong!")
	return tracerr.Wrap(err)
}
