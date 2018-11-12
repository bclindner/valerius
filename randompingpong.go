package main

import (
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
	"time"
	"encoding/json"
)

// Command which sends a random response if any of a set of triggers are sent.
type RandomPingPongCommand struct {
	BaseCommand
	RandomPingPongConfig
	RNG *rand.Rand
}

// Config for the PingPong command.
type RandomPingPongConfig struct {
	// List of messages that trigger the command.
	Triggers []string `json:"triggers"`
	// List of responses to pick randomly from.
	Responses []string `json:"responses"`
	// Prefix of the random response.
	ResponsePrefix string `json:"responsePrefix"`
	// Suffix of the random response.
	ResponseSuffix string `json:"responseSuffix"`
}

// Creates a new RandomPingPongCommand.
func NewRandomPingPongCommand(config CommandConfig) (command RandomPingPongCommand, err error) {
	// parse config
	options := RandomPingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil { return }
	// generate new RNG object
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// generate the command
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
	// Check if any trigger matches the message content
	for _, trigger := range p.Triggers {
		if evt.Message.Content == trigger {
			return true
		}
	}
	return false
}

func (p RandomPingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// get an index between 0 and the length of the message
	// we use this to get a random entry from the response list
	i := p.RNG.Intn(len(p.Responses))
	// construct the message
	msg := p.ResponsePrefix + p.Responses[i] + p.ResponseSuffix
	// send it
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, msg)
	return
}

