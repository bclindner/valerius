package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"encoding/json"
)

type PingPongCommand struct {
	BaseCommand
	PingPongConfig
}

type PingPongConfig struct {
	Triggers []string `json:"triggers"`
	Response string `json:"response"`
}

func NewPingPongCommand(config CommandConfig) (command PingPongCommand, err error) {
	options := PingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil { return }
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
	for _, trigger := range p.Triggers {
		if evt.Message.Content == trigger {
			return true
		}
	}
	return false
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.Response)
	return
}
