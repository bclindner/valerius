package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"math/rand"
)

// The BangerFinderCommand is designed to pull a link from a list of
// predefined links and post it. It's mostly made for posting bangers
// from the list already defined for BangerAlertCommand.
type BangerFinderCommand struct {
	Command
	RNG     *rand.Rand
	Bangers *[]string
}

// Returns the name of the command.
func (b BangerFinderCommand) Name() string {
	return "Banger Finder"
}

// Generates a new BangerFinderCommand.
func NewBangerFinderCommand(bangers *[]string) BangerFinderCommand {
	// initialize RNG if not done already
	rng := rand.New(rand.NewSource(253489732658))
	// Check if bangers are available, if not, panic (for now)
	if len(*bangers) == 0 {
		log.Fatal("No bangers found")
	}
	newCommand := BangerFinderCommand{
		RNG:     rng,
		Bangers: bangers,
	}
	return newCommand
}

// Test for the BangerFinder command - really only checks if !banger is the message.
func (b BangerFinderCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!banger"
}

// Posts a random banger from the Bangers list.
func (b BangerFinderCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	i := b.RNG.Intn(len(*b.Bangers))
	bot.ChannelMessageSend(evt.Message.ChannelID, (*b.Bangers)[i])
}
