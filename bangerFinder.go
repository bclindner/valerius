package main

import (
	log "github.com/sirupsen/logrus" // logging suite
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
)

type BangerFinderCommand struct {
	Command
	RNG *rand.Rand
	Bangers *[]string
}

func (b BangerFinderCommand) Name() string {
	return "Banger Finder"
}

func NewBangerFinderCommand(bangers *[]string) BangerFinderCommand {
	// initialize RNG if not done already
	rng := rand.New(rand.NewSource(253489732658))
	// Check if bangers are available, if not, panic (for now)
	if len(*bangers) == 0 {
		log.Fatal("No bangers found")
	}
	newCommand := BangerFinderCommand{
			RNG: rng,
			Bangers: bangers,
		}
	return newCommand
}

func (b BangerFinderCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!banger"
}

func (b BangerFinderCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	i := b.RNG.Intn(len(*b.Bangers))
	bot.ChannelMessageSend(evt.Message.ChannelID, (*b.Bangers)[i])
}
