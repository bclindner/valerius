package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
)

// Interface for commands that can be handled by the MessageHandler.
type Command interface {
	// Returns a human-readable name for the function, for logging purposes.
	Name() string
	// Command test. Whenever a message is sent, this test is run.
	// If it passes, the handler calls the Run() method.
	Test(*discordgo.Session, *discordgo.MessageCreate) bool
	// Runs the function. This can theoretically do anything, but is most
	// commonly used to reply to or otherwise process a message.
	Run(*discordgo.Session, *discordgo.MessageCreate)
}

type BaseCommand struct {
	name string
}

func (b BaseCommand) Name() string {
	return b.name
}

// Struct for the bot message handler. Currently this just contains a list of commands
// Should extend the Handler interface.
type MessageHandler struct {
	Handler
	Commands []Command
}

// Interface for the bot message handler.
// Has Handle and Add functions that handle commands and add new ones.
type Handler interface {
	Handle(*discordgo.Session, *discordgo.MessageCreate)
	Add(...Command)
}

// Create a new handler and bind it to a Session.
func NewMessageHandler(bot *discordgo.Session, user *discordgo.User) *MessageHandler {
	handler := MessageHandler{}
	bot.AddHandler(handler.Handle)
	return &handler
}

// Handle a Discord message. This just runs the Test() function of each command,
// and if a command's test passes, the handler calls its Run() function, logging
// the action as well.
func (c *MessageHandler) Handle(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	// Run preliminary tests: is the user sending the message a bot?
	if !evt.Message.Author.Bot {
		// For each command:
		for _, cmd := range c.Commands {
			// Handle it as a goroutine to speed things up
			go func(cmd Command) {
				// If the test checks out,
				if cmd.Test(bot, evt) {
					// log it,
					author := *evt.Message.Author
					log.WithFields(log.Fields{
						"text":     evt.Message.Content,
						"command":  cmd.Name(),
						"userID":   author.ID,
						"username": author.Username + "#" + author.Discriminator,
					}).Info("Command fired")
					// and run the command
					cmd.Run(bot, evt)
				}
			}(cmd)
		}
	}
}

// Add commands to the handler.
func (c *MessageHandler) Add(cmds ...Command) {
	c.Commands = append(c.Commands, cmds...)
}
