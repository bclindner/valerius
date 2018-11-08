package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
)

// A simple command that just returns "Hello, world!"
// Essentially just an example command for the CommandHandler.
type HelloCommand struct {
	Command
}

// 
func (p HelloCommand) Name() string {
	return "Hello World"
}

func NewHelloCommand() HelloCommand {
	return HelloCommand{}
}

func (h HelloCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!hello"
}

func (h HelloCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	bot.ChannelMessageSend(evt.Message.ChannelID, "Hello, world!")
}
