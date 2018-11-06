package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"strings" // for command testing
)

type PingPongCommand struct {
	Command
	PingString string
	PongString string
}

func (p PingPongCommand) Name() string {
	return "PingPong command"
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return strings.Contains(evt.Message.Content, p.PingString)
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	bot.ChannelMessageSend(evt.Message.ChannelID, p.PongString)
}
