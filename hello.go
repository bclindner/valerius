package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
)

type HelloCommand struct {
	Command
}

func (p HelloCommand) Name() string {
	return "Hello World"
}
func (h HelloCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!hello"
}

func (h HelloCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	bot.ChannelMessageSend(evt.Message.ChannelID, "Hello, world!")
}
