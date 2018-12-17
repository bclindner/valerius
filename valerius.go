package main

import (
	"encoding/json" // for parsing config file
	"errors"
	"flag"                           // for parsing args at runtime
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"io"                             // for io.MultiWriter (logrus multi-output)
	"io/ioutil"                      // for opening config file
	"os"                             // for opening logging file
	"os/signal"                      // for interrupt signal information
)

// BotConfiguration is the structure for the bot configuration JSON file.
type BotConfiguration struct {
	// Token that the bot logs in with.
	BotToken string `json:"botToken"`
	// Bot status message (when initialized).
	Status string `json:"status"`
	// List of commands to try and create.
	Commands []BaseCommand `json:"commands"`
}

var (
	config     BotConfiguration
	handler    *MessageHandler
	logPath    = flag.String("log", "", "Path to the logfile, if used.")
	configPath = flag.String("conf", "valerius.json", "Path to the config file.")
)

func init() {
	// parse flags
	flag.Parse()
	// setup log
	log.SetFormatter(&log.JSONFormatter{})
	// log to a file as well as stdout if the -log flag was set
	if *logPath != "" {
		logfile, err := os.OpenFile(*logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("Unable to establish logging: ", err)
		}
		// set up the output and formatter
		log.SetOutput(io.MultiWriter(os.Stdout, logfile))
	}
	// have to set err here so config goes in as a global var, 'cuz Go
	// (this is probably the wrong way to do it, though)
	var err error
	config, err = ReadBotConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize the bot.
func initBot() (bot *discordgo.Session, err error) {
	log.Info("Bot initializing")
	// start the bot session
	bot, err = discordgo.New("Bot " + config.BotToken)
	// get the current bot user (to figure out who we are)
	user, err := bot.User("@me")
	if err != nil {
		return
	}
	// log who we are
	log.Info("Bot logged in as ", user.Username, "#", user.Discriminator)
	return
}

// ReadBotConfig reads a config file from a path and parses it into a BotConfiguration.
func ReadBotConfig(path string) (config BotConfiguration, err error) {
	// load bot config file
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return config, errors.New("Unable to read config file: " + err.Error())
	}
	// parse bot config file
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return config, errors.New("Unable to read config file: " + err.Error())
	}
	return config, nil
}

func main() {
	// initialize the bot
	bot, err := initBot()
	if err != nil {
		log.Fatal("Failed to initialize bot: ", err)
	}
	defer bot.Close()
	// instantiate the handler
	handler, err = NewMessageHandler(bot, config.Commands)
	if err != nil {
		log.Fatal(err)
	}
	// open the bot to be used
	bot.Open()
	// set our status
	if len(config.Status) > 0 {
		err = bot.UpdateStatus(0, config.Status)
		if err != nil {
			log.Error("Error setting status:", err)
		}
	}
	// wait for OS interrupt (ctrl-c or a kill or something)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// close the bot websocket and exit the program
	log.Info("Interrupt signal sent, shutting down...")
}
