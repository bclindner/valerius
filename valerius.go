package main

import (
	"io" // for io.MultiWriter (logrus multi-output)
	"io/ioutil" // for opening config file
	"encoding/json" // for parsing config file
	"github.com/bwmarrin/discordgo" // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"os" // for opening logging file
	"os/signal" // for interrupt signal information
	"flag" // for parsing args at runtime
)

// Structure for the bot configuration JSON file.
type BotConfiguration struct {
	BotToken string `json:"botToken"`
	Bangers []string `json:"bangers"`
	DanceGifs []string `json:"danceGifs"`
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
		logfile, err := os.OpenFile(*logPath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		if err != nil { log.Fatal("Unable to establish logging: ",err) }
		// set up the output and formatter
		log.SetOutput(io.MultiWriter(os.Stdout, logfile))
	}
	log.SetFormatter(&log.JSONFormatter{})
	// setup logrus config
	// load bot config file
	configFile, err := ioutil.ReadFile(*configPath)
	if err != nil { log.Fatal("Unable to read config file: ",err) }
	// parse bot config file
	err = json.Unmarshal(configFile, &config)
	if err != nil { log.Fatal("Unable to read config file: ",err) }
}

// Initialize the bot.
func initBot() (bot *discordgo.Session, user *discordgo.User, err error) {
	log.Info("Bot initializing")
	// initialize the bot
	bot, err = discordgo.New("Bot "+config.BotToken)
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
	handler.Add(HelloCommand{})
	handler.Add(PingPongCommand{
		PingString: "ping",
		PongString: "pong",
	})
	handler.Add(NewXKCDCommand())
	if len(config.Bangers) == 0 {
		log.Warn("No bangers in config file; not adding banger alert command")
	} else {
		handler.Add(NewBangerAlertCommand(config.Bangers, config.DanceGifs))
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
