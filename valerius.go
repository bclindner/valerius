package main

import (
	"io/ioutil" // for opening config file
	"encoding/json" // for parsing config file
	"github.com/bwmarrin/discordgo" // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"os"
	"os/signal"
)

type BotConfiguration struct {
	BotToken string `json:"botToken"`
}

func initBot() (bot *discordgo.Session, err error) {
	// load bot config file
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil { return }
	// parse bot config file
	var config BotConfiguration
	err = json.Unmarshal(configFile, &config)
	if err != nil { return }
	// initialize the bot
	bot, err = discordgo.New("Bot "+config.BotToken)
	bot.Open()
	return
}

func main() {
	// initialize the bot
	log.Info("Bot initializing")
	bot, err := initBot()
	if err != nil { log.Fatal("Failed to initialize bot: ",err) }
	// get the current bot user
	user, err := bot.User("@me")
	log.Info("Bot logged in as ",user.Username,"#",user.Discriminator)
	// instantiate and register the handler
	handler := NewMessageHandler(bot, user)
	// add handler commands
	handler.Add(PingPongCommand{
		PingString: "ping",
		PongString: "pong",
	})
	// wait for OS interrupt (ctrl-c or a kill or something)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// close the bot websocket and exit the program
	bot.Close()
}
