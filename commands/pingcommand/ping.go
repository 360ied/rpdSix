package pingcommand

import (
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
	var _, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Pong!")
	return err
}
