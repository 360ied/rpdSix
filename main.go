package main

import (
	"fmt"
	"github.com/3sixtyied/rpdSix/commands"
	"github.com/3sixtyied/rpdSix/commands/pingcommand"
	"github.com/3sixtyied/rpdSix/commands/saycommand"
	"github.com/3sixtyied/rpdSix/keepalive"
	"github.com/bwmarrin/discordgo"
	"os"
)

func main() {
	// so that repl.it won't exit after the page is closed
	go keepalive.KeepAlive()

	bot, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		panic(err)
	}

	// register events
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	// init command map
	commands.InitCommands()
	// register commands
	pingcommand.Initialize()
	saycommand.Initialize()

	err = bot.Open()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is now running.")

	// wait forever
	select {}

	// unreachable code
	// bot.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateStatus(0, "golang")
	if err != nil {
		fmt.Println("Error updating status: ", err)
	}
	fmt.Println("Logged in as user " + s.State.User.ID)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go commands.HandleMessage(s, m)
}
