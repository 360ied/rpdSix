package saycommand

import (
	"github.com/ztrue/tracerr"

	"rpdSix/commandhandler"
)

func Initialize() {
	commandhandler.AddCommand(
		commandhandler.Command{
			Run:                         run,
			Names:                       []string{"say"},
			ExpectedPositionalArguments: []string{toSayArg},
		},
	)
}

const (
	toSayArg = "toSay"
)

func run(ctx commandhandler.CommandContext) error {
	var _, err = ctx.Message.Reply(ctx.Arguments[toSayArg])
	return tracerr.Wrap(err)
}
