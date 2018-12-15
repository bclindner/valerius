package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo" // for running the bot
	"regexp"
)

type RegexPingPongCommand struct {
	BaseCommand
	RegexPingPongConfig
	Regexp *regexp.Regexp
}

type RegexPingPongConfig struct {
	Trigger  string `json:"trigger"`
	Response string `json:"response"`
}

func NewRegexPingPongCommand(config CommandConfig) (command RegexPingPongCommand, err error) {
	options := RegexPingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return
	}
	regex, err := regexp.Compile(options.Trigger)
	if err != nil {
		return
	}
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
	return p.Regexp.MatchString(evt.Message.Content)
}
func (p RegexPingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.Response)
	return
}
