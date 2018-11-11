package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
	"time"
	"encoding/json"
)

type RandomPingPongConfig struct {
	Triggers []string `json:"triggers"`
	Responses []string `json:"responses"`
	ResponsePrefix string `json:"responsePrefix"`
	ResponseSuffix string `json:"responseSuffix"`
}

type RandomPingPongCommand struct {
	BaseCommand
	RandomPingPongConfig
	RNG *rand.Rand
}

func NewRandomPingPongCommand(config CommandConfig) (command RandomPingPongCommand, err error) {
	options := RandomPingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	command = RandomPingPongCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Type,
		},
		RandomPingPongConfig: options,
		RNG: rng,
	}
	return
}

func (p RandomPingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	for _, trigger := range p.Triggers {
		if evt.Message.Content == trigger {
			return true
		}
	}
	return false
}

func (p RandomPingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	i := p.RNG.Intn(len(p.Responses))
	msg := p.ResponsePrefix + p.Responses[i] + p.ResponseSuffix
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, msg)
	return
}

