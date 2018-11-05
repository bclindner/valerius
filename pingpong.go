package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type PingPongCommand struct {
	PingString string
	PongString string
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return strings.Contains(evt.Message.Content, p.PingString)
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	bot.ChannelMessageSend(evt.Message.ChannelID, p.PongString)
}
