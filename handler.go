package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
)

// Interface for commands that can be handled by the MessageHandler.
type Command interface {
	// Returns a human-readable name for the function, for logging purposes.
	GetName() string
	GetType() string
	// Command test. Whenever a message is sent, this test is run.
	// If it passes, the handler calls the Run() method.
	Test(*discordgo.Session, *discordgo.MessageCreate) bool
	// Runs the function. This can theoretically do anything, but is most
	// commonly used to reply to or otherwise process a message.
	// Returns an error that the handler can log.
	Run(*discordgo.Session, *discordgo.MessageCreate) error
}

// Base command structure.
type BaseCommand struct {
	Command
	// Human-readable name of the command, for logging purposes.
	// In the handler, This is retrieved through GetName().
	Name string
	// Human-readable type of the command, for logging purposes.
	// In the handler, This is retrieved through GetType().
	Type string
}

// Print the set name of the BaseCommand.
func (b BaseCommand) GetName() string {
	return b.Name
}

// Print the set type of the BaseCommand.
func (b BaseCommand) GetType() string {
	return b.Type
}

// Struct for the bot message handler. Currently this just contains a list of commands.
type MessageHandler struct {
	Handler
	// List of commands to test.
	commands []Command
}

// Interface for the bot message handler.
// Has Handle and Add functions that handle commands and add new ones.
type Handler interface {
	// Handle a Discord command.
	Handle(*discordgo.Session, *discordgo.MessageCreate)
}

// Create a new handler and bind it to a Session.
func NewMessageHandler(bot *discordgo.Session) *MessageHandler {
	handler := MessageHandler{}
	bot.AddHandler(handler.Handle)
	return &handler
}

// Handle a Discord message. This just runs the Test() function of each command,
// and if a command's test passes, the handler calls its Run() function, logging
// the action as well.
func (c *MessageHandler) Handle(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	// Run preliminary tests: is the user sending the message a bot?
	if evt.Message.Author.Bot { return }
	// Is this message being sent in a guild (i.e. not a PM?)
	if evt.Message.GuildID == "" { return }
	// For each command:
	for _, cmd := range c.commands {
		// Handle it as a goroutine to speed things up
		go func(cmd Command) {
			// If the test checks out,
			if cmd.Test(bot, evt) {
				// log it,
				author := *evt.Message.Author
				log.WithFields(log.Fields{
					"text":     evt.Message.Content,
					"command":  cmd.GetName(),
					"type":  cmd.GetType(),
					"userID":   author.ID,
					"username": author.Username + "#" + author.Discriminator,
				}).Info("Command fired")
				// and run the command
				err := cmd.Run(bot, evt)
				if err != nil {
					// Log if it failed, too
					log.WithFields(log.Fields{
						"text":     evt.Message.Content,
						"command":  cmd.GetName(),
						"type":  cmd.GetType(),
						"userID":   author.ID,
						"username": author.Username + "#" + author.Discriminator,
						"error":    err,
					}).Error("Command failed")
				}
			}
		}(cmd)
	}
}

// Add commands to the handler.
func (c *MessageHandler) Add(cmds ...Command) {
	c.commands = append(c.commands, cmds...)
}
