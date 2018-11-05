package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	log "github.com/sirupsen/logrus" // logging suite
)

type PingPongCommand struct {
	Name string
	PingString string
	PongString string
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return strings.Contains(evt.Message.Content, p.PingString)
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	author := *evt.Message.Author
	log.WithFields(log.Fields{
		"command": p.Name,
		"userID": author.ID,
		"username": author.Username + "#" + author.Discriminator,
	}).Info("sending pong")
	bot.ChannelMessageSend(evt.Message.ChannelID, p.PongString)
}
