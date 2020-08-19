package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandContext struct {
	Session   *discordgo.Session
	Message   *discordgo.MessageCreate
	Arguments map[string]string
}

type Command struct {
	Run                         func(ctx CommandContext)
	Names                       []string
	ExpectedPositionalArguments []string
	KeywordArgumentAliases      map[string]string
}

var Commands map[string]Command

// Initialize the Commands map
func InitCommands() {
	Commands = make(map[string]Command)
}

// Adds a command to the Commands map
func AddCommand(command Command) {
	for _, name := range command.Names {
		Commands[name] = command
	}
}

const prefix = "'"

// Handle a message creation event
func HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if !strings.HasPrefix(message.Content, prefix) {
		return
	}

	var endIndex = strings.Index(message.Content, stringSeparator)

	var commandName string

	if endIndex == -1 {
		commandName = message.Content[len(prefix):]
	} else {
		commandName = message.Content[
			len(prefix):strings.Index(message.Content, stringSeparator)]
	}

	// fmt.Println(commandName) // debug

	command, exists := Commands[commandName]

	if !exists {
		return
	}

	var context = CommandContext{
		Session: session,
		Message: message,
		Arguments: parseArguments(
			message.Content,
			command.ExpectedPositionalArguments,
			command.KeywordArgumentAliases),
	}

	go command.Run(context)

}

const keywordArgumentPrefix = "--"

const stringSeparator = " "

// Parse command arguments
func parseArguments(
	content string,
	expectedPositionalArguments []string,
	keywordArgumentAliases map[string]string) map[string]string {

	var separated = strings.Split(content, stringSeparator)

	// do not process the command name and prefix
	separated = separated[1:]

	var returnArguments = make(map[string]string)

	//if len(separated) == 0 {
	//	return returnArguments
	//}

	var currentPosition = 0

	for len(separated) > 0 {
		// emulate pop
		var currentItem = separated[0]
		separated = separated[1:]

		var currentArgumentValue []string

		if strings.HasPrefix(currentItem, keywordArgumentPrefix) {
			for len(separated) > 0 && !strings.HasPrefix(separated[0], keywordArgumentPrefix) {
				_ = append(currentArgumentValue, separated[0])
				separated = separated[1:]
				currentPosition++
			}
		} else {
			if currentPosition >= len(expectedPositionalArguments) {
				_, exists := returnArguments[
					expectedPositionalArguments[
						len(expectedPositionalArguments)-1]]
				if exists {
					// The length checks should prevent the value from being nil
					//goland:noinspection GoNilness
					returnArguments[
						expectedPositionalArguments[
							len(expectedPositionalArguments)-1]] += stringSeparator + currentItem
				} else {
					//goland:noinspection GoNilness
					returnArguments[
						expectedPositionalArguments[
							len(expectedPositionalArguments)-1]] = currentItem
				}
			} else {
				//goland:noinspection GoNilness
				returnArguments[
					expectedPositionalArguments[
						currentPosition]] = currentItem
			}
		}

		currentPosition++
	}

	//goland:noinspection GoNilness
	for key, value := range returnArguments {
		key = strings.ToLower(key)
		_, exists := keywordArgumentAliases[key]
		if exists {
			returnArguments[keywordArgumentAliases[key]] = value
		}
	}

	return returnArguments

}
