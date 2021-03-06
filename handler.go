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
	// Checks if the command can be used on a given guild and channel ID.
	Check(guildID string, channelID string, userID string) bool
}

// Checks if a list contains something.
func listContains(list []string, id string) bool {
	for _, listid := range list {
		if listid == id {
			return true
		}
	}
	return false
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
	// Optional channel whitelist.
	// If set, channels in this list cannot use this command.
	ChannelBlacklist []string `json:"channelblacklist"`
	// Optional guild blacklist.
	// If set, guilds in this list cannot use this command.
	GuildWhitelist []string `json:"guildwhitelist"`
	// Optional guild whitelist.
	// If set, guilds in this list cannot use this command.
	GuildBlacklist []string `json:"guildblacklist"`
	// Optional user whitelist.
	// If set, users in this list cannot use this command.
	UserWhitelist []string `json:"userwhitelist"`
	// Optional user blacklist.
	// If set, users in this list cannot use this command.
	UserBlacklist []string `json:"userblacklist"`
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

// Check ensures the command passes whitelist and blacklist checks.
func (b BaseCommand) Check(guildID, channelID, userID string) bool {
	if len(b.ChannelWhitelist) > 0 && !listContains(b.ChannelWhitelist, channelID) {
		return false
	}
	if len(b.ChannelBlacklist) > 0 && listContains(b.ChannelBlacklist, channelID) {
		return false
	}
	if len(b.GuildWhitelist) > 0 && !listContains(b.GuildWhitelist, guildID) {
		return false
	}
	if len(b.GuildBlacklist) > 0 && listContains(b.GuildBlacklist, guildID) {
		return false
	}
	if len(b.UserWhitelist) > 0 && !listContains(b.UserWhitelist, userID) {
		return false
	}
	if len(b.UserBlacklist) > 0 && listContains(b.UserBlacklist, userID) {
		return false
	}
	return true
}

// The Handler handles Discordgo messages, testing them against Valerius-compatible commands.
// The struct itself only contains the list of commands.
type Handler struct {
	// List of commands to test.
	commands []Command
	// Command to disconnect the handler from the bot.
	DestroySelf func()
}

// NewHandler creates a new handler and binds it to a Session.
func NewHandler(bot *discordgo.Session, commands []BaseCommand) (*Handler, error) {
	handler := Handler{}
	// set variables for use in the loop
	var (
		err error
		cmd Command
	)
	// add handler commands
	for _, config := range commands {
		switch config.Type {
		case "pingpong":
			cmd, err = NewPingPongCommand(config)
		case "iasip":
			cmd, err = NewIASIPCommand(config)
		case "rest":
			cmd, err = NewRESTCommand(config)
		case "reload":
			cmd, err = NewReloadCommand(config)
		default:
			return &handler, errors.New("Command " + config.Name + " is of invalid type (" + config.Type + "). Exiting.")
		}
		// handle any errors
		if err != nil {
			return &handler, errors.New("Error with command " + config.Name + ": " + err.Error())
		}
		// add the command
		handler.Add(cmd)
	}
	// log how many commands we parsed
	log.Info("Parsed ", len(handler.commands), " commands")
	// register self with the bot, and get the function necessary to detach from bot
	handler.DestroySelf = bot.AddHandler(handler.Handle)
	return &handler, nil
}

// Handle handles a Discord message. This just runs the Test() function of each command,
// and if a command's test passes, the handler calls its Run() function, logging
// the action as well.
func (c *Handler) Handle(bot *discordgo.Session, evt *discordgo.MessageCreate) {
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
			// Test the command
			if cmd.Check(evt.Message.GuildID, evt.Message.ChannelID, evt.Message.Author.ID) && cmd.Test(bot, evt) {
				// If it passed, log it,
				author := *evt.Message.Author
				log.WithFields(log.Fields{
					"text":      evt.Message.Content,
					"command":   cmd.GetName(),
					"type":      cmd.GetType(),
					"userID":    author.ID,
					"username":  author.Username + "#" + author.Discriminator,
					"guildID":   evt.Message.GuildID,
					"channelID": evt.Message.ChannelID,
				}).Info("Command fired")
				// and run the command
				err := cmd.Run(bot, evt)
				if err != nil {
					// Log if it failed, too
					log.WithFields(log.Fields{
						"text":      evt.Message.Content,
						"command":   cmd.GetName(),
						"type":      cmd.GetType(),
						"userID":    author.ID,
						"guildID":   evt.Message.GuildID,
						"channelID": evt.Message.ChannelID,
						"username":  author.Username + "#" + author.Discriminator,
						"error":     err,
					}).Error("Command failed")
				}
			}
		}(cmd)
	}
}

// Add commands to the handler, validating whitelists/blacklists as well.
func (c *Handler) Add(cmd Command) {
	c.commands = append(c.commands, cmd)
}
