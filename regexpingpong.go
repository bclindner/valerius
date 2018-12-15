package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo" // for running the bot
	"regexp"
)

// Command which sets a response if the message matches a regular expression.
type RegexPingPongCommand struct {
	BaseCommand
	RegexPingPongConfig
	Regexp *regexp.Regexp
}

type RegexPingPongConfig struct {
	// Regular expression to test messages with.
	Trigger string `json:"trigger"`
	// Response message to send.
	Response string `json:"response"`
}

func NewRegexPingPongCommand(config CommandConfig) (command RegexPingPongCommand, err error) {
	// parse config
	options := RegexPingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return
	}
	// compile regex
	regex, err := regexp.Compile(options.Trigger)
	if err != nil {
		return
	}
	// construct command
	command = RegexPingPongCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Type,
		},
		RegexPingPongConfig: options,
		Regexp:              regex,
	}
	return
}

func (p RegexPingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	// check if regex matches
	return p.Regexp.MatchString(evt.Message.Content)
}
func (p RegexPingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// send the response
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.Response)
	return
}
