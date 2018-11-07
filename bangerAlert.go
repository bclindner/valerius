package main

import (
	log "github.com/sirupsen/logrus" // logging suite
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
)

type BangerAlertCommand struct {
	Command
	RNG *rand.Rand
	BangerMessage string
	Bangers *[]string
	DanceEnabled bool
	DanceGifs *[]string
}

func (b BangerAlertCommand) Name() string {
	return "Banger Alert"
}

func NewBangerAlertCommand(bangers *[]string, gifs *[]string) BangerAlertCommand {
	// initialize RNG if not done already
	rng := rand.New(rand.NewSource(253489732658))
	bmessage := "🚨OH🚨SHIT🚨IT'S🚨A🚨BANGER🚨 "
	// Check if bangers are available, if not, panic (for now)
	if len(*bangers) == 0 {
		log.Fatal("No bangers found")
	}
	newCommand := BangerAlertCommand{
			RNG: rng,
			Bangers: bangers,
			DanceEnabled: len(*gifs) != 0,
			DanceGifs: gifs,
			BangerMessage: bmessage,
		}
	log.Info(len(*newCommand.Bangers)," bangers loaded.")
	log.Info(len(*newCommand.DanceGifs)," dance GIFs loaded.")
return newCommand
}

func (b BangerAlertCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	for _, bangerLink := range *b.Bangers {
		if evt.Message.Content == bangerLink {
			return true
		}
	}
	return false
}

func (b BangerAlertCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	if b.DanceEnabled {
		i := b.RNG.Intn(len(*b.DanceGifs))
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage+(*b.DanceGifs)[i])
	} else {
		bot.ChannelMessageSend(evt.Message.ChannelID, b.BangerMessage)
	}
}
