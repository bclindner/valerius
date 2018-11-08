package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"strings"                       // for command testing
)

type PingPongCommand struct {
	BaseCommand
	PingString string
	PongString string
}

func (p PingPongCommand) Name() string {
	return p.name
}

func NewPingPongCommand(name string, ping string, pong string) PingPongCommand {
	return PingPongCommand{
		BaseCommand: BaseCommand{
			name: name,
		},
		PingString: ping,
		PongString: pong,
	}
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return strings.Contains(evt.Message.Content, p.PingString)
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.PongString)
	return
}
