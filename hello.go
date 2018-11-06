package main

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus" // logging suite
)

type HelloCommand struct {
	Name string
}

func (h HelloCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!hello"
}

func (h HelloCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	author := *evt.Message.Author
	log.WithFields(log.Fields{
		"command": h.Name,
		"userID": author.ID,
		"username": author.Username + "#" + author.Discriminator,
	}).Info("sending hello")
	bot.ChannelMessageSend(evt.Message.ChannelID, "Hello, world!")
}
