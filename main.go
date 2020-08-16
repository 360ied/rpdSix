package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"rpdSix/commands"
	"rpdSix/commands/pingcommand"
	"rpdSix/commands/saycommand"
)

func main() {
	// so that repl.it won't exit after the page is closed
	go KeepAlive()

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
	err := s.UpdateStatus(0, "FUCK PYTHON I AM NOW A GO CODER")
	if err != nil {
		fmt.Println("Error updating status: ", err)
	}
	fmt.Println("Logged in as user " + s.State.User.ID)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go commands.HandleMessage(s, m)
}
