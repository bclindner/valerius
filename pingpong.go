package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"encoding/json"
)

// Command that sends a response string if any of a set of triggers are sent.
// As simple as it gets.
type PingPongCommand struct {
	BaseCommand
	PingPongConfig
}

// Config for the PingPong command.
type PingPongConfig struct {
	// List of messages that trigger the command.
	Triggers []string `json:"triggers"`
	// Response string to send.
	Response string `json:"response"`
}
// Creates a new PingPongCommand.
func NewPingPongCommand(config CommandConfig) (command PingPongCommand, err error) {
	// Parse config
	options := PingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil { return }
	// Generate command
	command = PingPongCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Type,
		},
		PingPongConfig: options,
	}
	return
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	// Check if any trigger matches the message content
	for _, trigger := range p.Triggers {
		if evt.Message.Content == trigger {
			return true
		}
	}
	return false
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// Send the response
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.Response)
	return
}
