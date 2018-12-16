package main

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
)

// Command is an interface for commands that can be handled by the MessageHandler.
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
	// Returns the whitelists and blacklists for the command.
	GetChannelWhitelist() []string
	GetChannelBlacklist() []string
	GetGuildWhitelist() []string
	GetGuildBlacklist() []string
}

// BaseCommand is the base command structure.
// It also serves as the schema
type BaseCommand struct {
	Command
	// Human-readable name of the command, for logging purposes.
	// In the handler, This is retrieved through GetName().
	Name string `json:"name"`
	// Human-readable type of the command, for logging purposes.
	// In the handler, This is retrieved through GetType().
	Type string `json:"type"`
	// Optional channel whitelist.
	// If set, only channels in this list can use this command.
	ChannelWhitelist []string `json:"channelwhitelist"`
	// Optional channel blacklist.
	// If set, channels in this list cannot use this command.
	ChannelBlacklist []string `json:"channelblacklist"`
	// Optional guild blacklist.
	// If set, guilds in this list cannot use this command.
	GuildWhitelist []string `json:"guildblacklist"`
	// Optional guild blacklist.
	// If set, guilds in this list cannot use this command.
	GuildBlacklist []string `json:"guildwhitelist"`
	// JSON-encoded list of options for the command.
	// This is intended to be parsed and handled by the "NewXCommand" factory function
	// after utilizing this BaseCommand.
	Options json.RawMessage
}

// GetName prints the set name of the BaseCommand.
func (b BaseCommand) GetName() string {
	return b.Name
}

// GetType prints the set type of the BaseCommand.
func (b BaseCommand) GetType() string {
	return b.Type
}

// GetChannelWhitelist returns the list of channels that may use this command.
func (b BaseCommand) GetChannelWhitelist() []string {
	return b.ChannelWhitelist
}

// GetChannelBlacklist returns the list of channels that cannot use this command.
func (b BaseCommand) GetChannelBlacklist() []string {
	return b.ChannelBlacklist
}

// GetGuildWhitelist returns the list of guilds that may use this command.
func (b BaseCommand) GetGuildWhitelist() []string {
	return b.GuildWhitelist
}

// GetGuildBlacklist returns the list of guilds that cannot use this command.
func (b BaseCommand) GetGuildBlacklist() []string {
	return b.GuildBlacklist
}

// MessageHandler handles Discordgo messages, testing them against Valerius-compatible commands.
// The struct itself only contains the list of commands.
type MessageHandler struct {
	Handler
	// List of commands to test.
	commands []Command
}

// Handler is the interface for the bot message handler.
// Has Handle and Add functions that handle commands and add new ones.
type Handler interface {
	// Handle a Discord command.
	Handle(*discordgo.Session, *discordgo.MessageCreate)
}

// NewMessageHandler creates a new handler and binds it to a Session.
func NewMessageHandler(bot *discordgo.Session) *MessageHandler {
	handler := MessageHandler{}
	bot.AddHandler(handler.Handle)
	return &handler
}

// Handle handles a Discord message. This just runs the Test() function of each command,
// and if a command's test passes, the handler calls its Run() function, logging
// the action as well.
func (c *MessageHandler) Handle(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	// Run preliminary tests: is the user sending the message a bot?
	if evt.Message.Author.Bot {
		return
	}
	// Is this message being sent in a guild (i.e. not a PM?)
	if evt.Message.GuildID == "" {
		return
	}
	// For each command:
	for _, cmd := range c.commands {
		// Handle it as a goroutine to speed things up
		go func(cmd Command) {
			// Ensure the command is in guild whitelists, and not in blacklists
			gwhitelist := cmd.GetGuildWhitelist()
			if len(gwhitelist) > 0 {
				guildok := false
				for _, guild := range gwhitelist {
					if guild == evt.Message.GuildID {
						guildok = true
						break
					}
				}
				if !guildok {
					return
				}
			}
			gblacklist := cmd.GetGuildBlacklist()
			if len(gblacklist) > 0 {
				for _, guild := range gblacklist {
					if guild == evt.Message.GuildID {
						return
					}
				}
			}
			// Ensure the command is in channel whitelists, and not in blacklists
			cwhitelist := cmd.GetChannelWhitelist()
			if len(cwhitelist) > 0 {
				chanok := false
				for _, channel := range cwhitelist {
					if channel == evt.Message.ChannelID {
						chanok = true
						break
					}
				}
				if !chanok {
					return
				}
			}
			cblacklist := cmd.GetChannelBlacklist()
			if len(cblacklist) > 0 {
				for _, channel := range cblacklist {
					if channel == evt.Message.ChannelID {
						return
					}
				}
			}
			// Test the command
			if cmd.Test(bot, evt) {
				// If it passed, log it,
				author := *evt.Message.Author
				log.WithFields(log.Fields{
					"text":     evt.Message.Content,
					"command":  cmd.GetName(),
					"type":     cmd.GetType(),
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
						"type":     cmd.GetType(),
						"userID":   author.ID,
						"username": author.Username + "#" + author.Discriminator,
						"error":    err,
					}).Error("Command failed")
				}
			}
		}(cmd)
	}
}

// Add commands to the handler, validating whitelists/blacklists as well.
func (c *MessageHandler) Add(cmd Command) error {
	if len(cmd.GetGuildBlacklist()) > 0 && len(cmd.GetGuildWhitelist()) > 0 {
		return errors.New("can only have one GuildBlacklist or GuildWhitelist")
	}
	if len(cmd.GetChannelBlacklist()) > 0 && len(cmd.GetChannelWhitelist()) > 0 {
		return errors.New("can only have one ChannelBlacklist or ChannelWhitelist")
	}
	c.commands = append(c.commands, cmd)
	return nil
}
