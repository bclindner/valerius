package main

import (
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
)

// The BangerFinderCommand is designed to pull a link from a list of
// predefined links and post it. It's mostly made for posting bangers
// from the list already defined for BangerAlertCommand.
type BangerFinderCommand struct {
	BangerCommand
}

// Generates a new BangerFinderCommand.
func NewBangerFinderCommand(name string, bangers *[]string) BangerFinderCommand {
	// initialize RNG if not done already
	// Check if bangers are available, if not, panic (for now)
	if len(*bangers) == 0 {
		log.Fatal("No bangers found")
	}
	newCommand := BangerFinderCommand{
		BangerCommand: NewBangerCommand(name, bangers),
	}
	return newCommand
}

func (b BangerFinderCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return evt.Message.Content == "!banger"
}

func (b BangerFinderCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, b.GetRandomBanger())
	return
}
