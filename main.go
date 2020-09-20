package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"rpdSix/cache"
	"rpdSix/commandhandler"
	"rpdSix/commandhandler/commands/bigemojicommand"
	"rpdSix/commandhandler/commands/helpcommand"
	"rpdSix/commandhandler/commands/pingcommand"
	"rpdSix/commandhandler/commands/purgecommand"
	"rpdSix/commandhandler/commands/saycommand"
	"rpdSix/commandhandler/commands/voicechannelmoveallcommand"
	"rpdSix/keepalive"
)

func main() {
	// so that repl.it won't exit after the page is closed
	go keepalive.KeepAlive()

	var bot, newBotErr = discordgo.New("Bot " + os.Getenv("TOKEN"))
	if newBotErr != nil {
		log.Fatalln(newBotErr)
	}

	// register events
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	cache.RegisterEventHandlers(bot)

	// init command map
	commandhandler.InitCommands()
	// register commands
	pingcommand.Initialize()
	saycommand.Initialize()
	helpcommand.Initialize()
	bigemojicommand.Initialize()
	voicechannelmoveallcommand.Initialize()
	purgecommand.Initialize()

	var botOpenErr = bot.Open()
	if botOpenErr != nil {
		log.Fatalln("Error opening Discord session: ", newBotErr)
	}

	defer bot.Close()

	fmt.Println("Bot is now running.")

	// wait forever
	select {}
}

func ready(session *discordgo.Session, event *discordgo.Ready) {
	var sessionUpdateStatusErr = session.UpdateStatus(0, "golang")
	if sessionUpdateStatusErr != nil {
		fmt.Println("Error updating status: ", sessionUpdateStatusErr)
	}
	fmt.Println("Logged in as user " + session.State.User.ID)
}

func messageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	go commandhandler.HandleMessage(session, event)
}
