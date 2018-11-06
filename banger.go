package main

import (
	log "github.com/sirupsen/logrus" // logging suite
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
)

type BangerAlertCommand struct {
	Command
	BangerMessage string
	Bangers []string
	DanceEnabled bool
	DanceGifs []string
}

func NewBangerAlertCommand(bangers []string, gifs []string) BangerAlertCommand {
	// initialize RNG if not done already
	rand.Seed(253489732658)
	// Check if bangers are available, if not, panic (for now)
	if len(bangers) == 0 {
		log.Fatal("No bangers specified for BangerAlertCommand!")
	}
	// Check if dance gifs should be enabled
	if len(gifs) > 0 {
		// return the object DanceGifs and enable them
		return BangerAlertCommand{
			Bangers: bangers,
			DanceEnabled: true,
			DanceGifs: gifs,
			BangerMessage: "🚨OH🚨SHIT🚨IT'S🚨A🚨BANGER🚨 ",
		}
	} else {
		return BangerAlertCommand{
			Bangers: bangers,
			DanceEnabled: false,
			BangerMessage: "🚨OH🚨SHIT🚨IT'S🚨A🚨BANGER🚨 ",
		}
	}
}

func (b BangerAlertCommand) Name() string {
	return "Banger Alert"
}
func (b BangerAlertCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	for _, bangerLink := range b.Bangers {
		if evt.Message.Content == bangerLink {
			return true
		}
	}
	return false
}

func (b BangerAlertCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	if b.DanceEnabled {
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage+b.DanceGifs[rand.Intn(len(b.DanceGifs))])
	} else {
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage)

	}
}
