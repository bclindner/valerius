package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
)

type Command interface {
	Name() string
	Test(*discordgo.Session, *discordgo.MessageCreate) bool
	Run(*discordgo.Session, *discordgo.MessageCreate)
}

type MessageHandler struct {
	Commands []Command
	UserID   string
}

type Handler interface {
	Handle(bot *discordgo.Session, evt *discordgo.MessageCreate)
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
