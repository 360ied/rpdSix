package pingcommand

import (
	"github.com/ztrue/tracerr"

	"rpdSix/commandhandler"
)

func Initialize() {
	commandhandler.AddCommand(
		commandhandler.Command{
			Run:   run,
			Names: []string{"ping"},
		},
	)
}

func run(ctx commandhandler.CommandContext) error {
	var _, err = ctx.Message.Reply("Pong!")
	return tracerr.Wrap(err)
}
