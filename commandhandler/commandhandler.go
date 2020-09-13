package commandhandler

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ztrue/tracerr"

	"rpdSix/helpers/extendeddiscord/extendeddiscordobjects"
)

const (
	CommandPrefix         = "'"
	KeywordArgumentPrefix = "--"
	StringSeparator       = " "

	IgnoredPositionalArgumentName = "_"
)

var (
	Commands map[string]Command
)

type CommandContext struct {
	Session   *discordgo.Session
	Message   *extendeddiscordobjects.ExtendedMessage
	Arguments map[string]string
}

type Command struct {
	Run                         func(ctx CommandContext) error
	Names                       []string
	ExpectedPositionalArguments []string
	KeywordArgumentAliases      map[string]string
}

// Initialize the Commands map
func InitCommands() {
	Commands = make(map[string]Command)
}

// Adds a command to the Commands map
func AddCommand(command Command) {
	// Ignore positional arguments instead of panicking
	// when command does not accept any positional arguments but positional arguments are given
	if len(command.ExpectedPositionalArguments) == 0 {
		command.ExpectedPositionalArguments = []string{IgnoredPositionalArgumentName}
	}
	for _, name := range command.Names {
		Commands[name] = command
	}
}

// Handle a message creation event
func HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if !strings.HasPrefix(message.Content, CommandPrefix) {
		return
	}

	var endIndex = strings.Index(message.Content, StringSeparator)

	var commandName string

	if endIndex == -1 {
		commandName = message.Content[len(CommandPrefix):]
	} else {
		commandName = message.Content[len(CommandPrefix):endIndex]
	}

	var command, exists = Commands[commandName]

	if !exists {
		return
	}

	var context = CommandContext{
		Session: session,
		Message: extendeddiscordobjects.ExtendMessage(message.Message, session),
		Arguments: parseArguments(
			message.Content,
			command.ExpectedPositionalArguments,
			command.KeywordArgumentAliases),
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			var panicMessage = fmt.Sprintf("Panic while executing %v!\nMessage: %v\nArguments: %v\nError: %v",
				commandName, message.Content, context.Arguments, panicErr)
			fmt.Println(panicMessage)
			_, _ = context.Message.Reply(panicMessage)
		}
	}()

	var err = command.Run(context)
	if err != nil {
		tracerr.PrintSourceColor(err)
		_, err = session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprint("An Error occurred while executing the command!\n", err))
		if err != nil {
			tracerr.PrintSourceColor(err)
			_, err = session.ChannelMessageSend(
				message.ChannelID,
				"An Error occurred while executing the command and sending the error message!")
			if err != nil {
				tracerr.PrintSourceColor(err)
			}
		}
	}
}

// Parse command arguments
func parseArguments(
	content string,
	expectedPositionalArguments []string,
	keywordArgumentAliases map[string]string) map[string]string {

	var separated = strings.Split(content, StringSeparator)

	// do not process the command name and prefix
	separated = separated[1:]

	var returnArguments = make(map[string]string)

	// Keep track of current position
	var currentPosition = 0
	for len(separated) > 0 {
		// Set current item
		var currentItem = separated[0]
		// Remove current item from processing queue
		separated = separated[1:]
		// Parse --arg arguments
		if strings.HasPrefix(currentItem, KeywordArgumentPrefix) {
			var currentArgumentValue []string
			// Do not take the next keyword argument as part of the current value
			// Iterate through the processing queue
			// Note: For presence arguments, check for presence in the map, not for its boolean value
			for len(separated) > 0 && !strings.HasPrefix(separated[0], KeywordArgumentPrefix) {
				// Append value and Remove processed value from the processing queue
				currentArgumentValue = append(currentArgumentValue, separated[0])
				separated = separated[1:]
				currentPosition++
			}
			// Set the current value
			returnArguments[currentItem[len(KeywordArgumentPrefix):]] = strings.Join(
				currentArgumentValue, StringSeparator)
		} else {
			// Set by positional argument
			// Allow the last positional argument to have the separator in between
			if currentPosition >= len(expectedPositionalArguments) {
				var _, in = returnArguments[
					expectedPositionalArguments[
						len(expectedPositionalArguments)-1]]
				if in {
					// The length checks should prevent the value from being nil
					// goland:noinspection GoNilness
					returnArguments[
						expectedPositionalArguments[
							len(expectedPositionalArguments)-1]] += StringSeparator + currentItem
				} else {
					// goland:noinspection GoNilness
					returnArguments[
						expectedPositionalArguments[
							len(expectedPositionalArguments)-1]] = currentItem
				}
			} else {
				// goland:noinspection GoNilness
				returnArguments[
					expectedPositionalArguments[
						currentPosition]] = currentItem
			}
		}

		currentPosition++
	}

	// goland:noinspection GoNilness
	for key, value := range returnArguments {
		key = strings.ToLower(key)
		var _, in = keywordArgumentAliases[key]
		if in {
			returnArguments[keywordArgumentAliases[key]] = value
		}
	}

	return returnArguments
}
