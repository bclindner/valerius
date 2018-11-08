package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"math/rand"
)

// The BangerAlertCommand is designed to post a message when a recognized link
// to a particularly good song is posted. Along with the message, it may also
// post links, specifically of people dancing.
type BangerAlertCommand struct {
	Command
	// Random number generator to use to get a random dance GIF.
	RNG           *rand.Rand
	// Message to display when a banger is posted.
	BangerMessage string
	// List of messages the command will consider bangers.
	Bangers       *[]string
	// Whether or not the dance gifs are enabled.
	// This is set automatically, based on whether or not there are gifs.
	DanceEnabled  bool
	// List of potential gifs to send when a banger is posted.
	DanceGifs     *[]string
}

// Returns the name of the command.
func (b BangerAlertCommand) Name() string {
	return "Banger Alert"
}

// Generates a new BangerAlertCommand.
func NewBangerAlertCommand(bangers *[]string, gifs *[]string) BangerAlertCommand {
	// initialize RNG if not done already
	rng := rand.New(rand.NewSource(253489732658))
	bmessage := "🚨OH🚨SHIT🚨IT'S🚨A🚨BANGER🚨 "
	// Check if bangers are available, if not, panic (for now)
	if len(*bangers) == 0 {
		log.Fatal("No bangers found")
	}
	newCommand := BangerAlertCommand{
		RNG:           rng,
		Bangers:       bangers,
		DanceEnabled:  len(*gifs) != 0,
		DanceGifs:     gifs,
		BangerMessage: bmessage,
	}
	log.Info(len(*newCommand.Bangers), " bangers loaded.")
	log.Info(len(*newCommand.DanceGifs), " dance GIFs loaded.")
	return newCommand
}

// Test for the BangerAlert command - if the message posted is anywhere in
// the Bangers array, it should fire the command.
func (b BangerAlertCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	for _, bangerLink := range *b.Bangers {
		if evt.Message.Content == bangerLink {
			return true
		}
	}
	return false
}

// Posts the BangerMessage, plus a gif link, if the dance is enabled.
func (b BangerAlertCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	if b.DanceEnabled {
		i := b.RNG.Intn(len(*b.DanceGifs))
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage+(*b.DanceGifs)[i])
	} else {
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage)
	}
}
