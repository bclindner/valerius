package main

import (
	"encoding/json"                  // for parsing config file
	"flag"                           // for parsing args at runtime
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"io"                             // for io.MultiWriter (logrus multi-output)
	"io/ioutil"                      // for opening config file
	"os"                             // for opening logging file
	"os/signal"                      // for interrupt signal information
)

// Structure for the bot configuration JSON file.
type BotConfiguration struct {
	BotToken  string   `json:"botToken"`
	Commands []CommandConfig `json:"commands"`
}

type CommandConfig struct {
	Name string
	Type string
	Options json.RawMessage
}

var config BotConfiguration

func init() {
	// set up flags
	logPath := flag.String("log", "", "Path to the logfile, if used.")
	configPath := flag.String("conf", "valerius.json", "Path to the config file.")
	// parse flags
	flag.Parse()
	// log to a file as well as stdout if the -log flag was set
	if *logPath != "" {
		logfile, err := os.OpenFile(*logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("Unable to establish logging: ", err)
		}
		// set up the output and formatter
		log.SetOutput(io.MultiWriter(os.Stdout, logfile))
		log.SetFormatter(&log.JSONFormatter{})
	}

	// setup logrus config
	// load bot config file
	configFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}
	// parse bot config file
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}
}

// Initialize the bot.
func initBot() (bot *discordgo.Session, user *discordgo.User, err error) {
	log.Info("Bot initializing")
	// initialize the bot
	bot, err = discordgo.New("Bot " + config.BotToken)
	// get the current bot user
	user, err = bot.User("@me")
	if err != nil {
		return
	}
	log.Info("Bot logged in as ", user.Username, "#", user.Discriminator)
	return
}

func main() {
	// initialize the bot
	bot, user, err := initBot()
	if err != nil {
		log.Fatal("Failed to initialize bot: ", err)
	}
	// instantiate and register the handler
	handler := NewMessageHandler(bot, user)
	// add handler commands
	for _, config := range config.Commands {
		switch config.Type {
			case "pingpong":
				cmd, err := NewPingPongCommand(config)
				if err != nil { log.Fatal("Error with command "+config.Name+": ",err) }
				handler.Add(cmd)
			case "randompingpong":
				cmd, err := NewRandomPingPongCommand(config)
				if err != nil { log.Fatal("Error with command "+config.Name+": ",err) }
				handler.Add(cmd)
			case "xkcd":
				cmd, err := NewXKCDCommand(config)
				if err != nil { log.Fatal("Error with command "+config.Name+": ",err) }
				handler.Add(cmd)
		}
	}
	// open the bot to be used
	bot.Open()
	// wait for OS interrupt (ctrl-c or a kill or something)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	// close the bot websocket and exit the program
	log.Info("Interrupt signal sent, shutting down...")
	bot.Close()
	log.Info("Bot closed down successfully. Goodbye")
}
