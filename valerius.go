package main

import (
	"io"
	"io/ioutil" // for opening config file
	"encoding/json" // for parsing config file
	"github.com/bwmarrin/discordgo" // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"os"
	"os/signal"
)

func init() {
	// setup logging file
	logfile, err := os.OpenFile("valerius.log", os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { log.Fatal("Unable to establish logging: ",err) }
	// set up the output and formatter
	log.SetOutput(io.MultiWriter(os.Stdout, logfile))
	log.SetFormatter(&log.JSONFormatter{})
}

// Structure for the bot configuration JSON file.
type BotConfiguration struct {
	BotToken string `json:"botToken"`
}

// Initialize the bot.
func initBot() (bot *discordgo.Session, user *discordgo.User, err error) {
	log.Info("Bot initializing")
	// setup logrus config
	// load bot config file
	configFile, err := ioutil.ReadFile("valerius.json")
	if err != nil { return }
	// parse bot config file
	var config BotConfiguration
	err = json.Unmarshal(configFile, &config)
	if err != nil { return }
	// initialize the bot
	bot, err = discordgo.New("Bot "+config.BotToken)
	bot.Open()
	// get the current bot user
	user, err = bot.User("@me")
	if err != nil { return }
	log.Info("Bot logged in as ",user.Username,"#",user.Discriminator)
	return
}

func main() {
	// initialize the bot
	bot, user, err := initBot()
	if err != nil { log.Fatal("Failed to initialize bot: ",err) }
	// instantiate and register the handler
	handler := NewMessageHandler(bot, user)
	// add handler commands
	handler.Add(PingPongCommand{
		Name: "example PingPongCommand",
		PingString: "ping",
		PongString: "pong",
	})
	// wait for OS interrupt (ctrl-c or a kill or something)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// close the bot websocket and exit the program
	log.Info("Interrupt signal sent, shutting down...")
	bot.Close()
}
