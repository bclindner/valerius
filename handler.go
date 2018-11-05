package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
)

type Command interface {
	Test(*discordgo.Session, *discordgo.MessageCreate) bool
	Run(*discordgo.Session, *discordgo.MessageCreate)
}

type MessageHandler struct {
	Commands []Command
	UserID string
}

// Create a new handler and bind it to a Session.
func NewMessageHandler(bot *discordgo.Session, user *discordgo.User) *MessageHandler {
	handler := MessageHandler{
		UserID: user.ID,
	}
	bot.AddHandler(handler.Handle)
	return &handler
}

// Handle a Discord message.
func (c *MessageHandler) Handle(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	// Run preliminary tests: is the user sending the message a bot?
	if !evt.Message.Author.Bot {
		// For each command:
		for _, cmd := range c.Commands {
			// If the test checks out,
			if cmd.Test(bot, evt) {
				// run the command
				cmd.Run(bot, evt)
			}
		}
	}
}

// Add commands to the handler.
func (c *MessageHandler) Add(cmds ...Command) {
	c.Commands = append(c.Commands, cmds...)
}
