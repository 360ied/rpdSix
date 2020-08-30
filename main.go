package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"rpdSix/commands"
	"rpdSix/commands/bigemojicommand"
	"rpdSix/commands/helpcommand"
	"rpdSix/commands/pingcommand"
	"rpdSix/commands/saycommand"
	"rpdSix/keepalive"
)

func main() {
	// so that repl.it won't exit after the page is closed
	go keepalive.KeepAlive()

	var bot, err = discordgo.New("Bot " + os.Getenv("TOKEN"))

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
	helpcommand.Initialize()
	bigemojicommand.Initialize()

	err = bot.Open()

	if err != nil {
		log.Fatalln("Error opening Discord session: ", err)
	}

	fmt.Println("Bot is now running.")

	// wait forever
	select {}

	// unreachable code
	// bot.Close()
}

func ready(session *discordgo.Session, event *discordgo.Ready) {
	var err = session.UpdateStatus(0, "golang")
	if err != nil {
		fmt.Println("Error updating status: ", err)
	}
	fmt.Println("Logged in as user " + session.State.User.ID)
}

func messageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	go commands.HandleMessage(session, event)
}
