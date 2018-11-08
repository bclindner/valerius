package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"io/ioutil"                      // for opening response body
	"net/http"
	"regexp"
	"strconv"
	"errors"
)

// The XKCDCommand base structure. Takes a Regexp to test the command.
type XKCDCommand struct {
	Command
	BaseCommand
	name string
	Regexp *regexp.Regexp
}

// XKCD API structure, for parsing the XKCD comic API.
type XKCDComic struct {
	Number    int    `json:"num"`
	Title     string `json:"title"`
	SafeTitle string `json:"safe_title"`
	Alt       string `json:"alt"`
	Image     string `json:"img"`
}

func (x XKCDCommand) Name() string {
	return x.name
}

func NewXKCDCommand(name string) XKCDCommand {
	// Instantiate the regex.
	rgx, err := regexp.Compile(`^\!xkcd ?([0-9]+)?$`)
	if err != nil {
		log.Fatal(err)
	}
	return XKCDCommand{
		BaseCommand: BaseCommand{
			name: name,
		},
		name: name,
		Regexp: rgx,
	}
}

func (x XKCDCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return x.Regexp.MatchString(evt.Message.Content)
}

func (x XKCDCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// Get the comic number from the regex group matches
	comicNumber := x.Regexp.FindStringSubmatch(evt.Message.Content)[1]
	// Get the endpoint necessary
	var endpoint string
	if comicNumber != "" {
		endpoint = "https://xkcd.com/" + comicNumber + "/info.0.json"
	} else {
		endpoint = "https://xkcd.com/info.0.json"
	}
	// Run a GET request to the endpoint
	resp, err := http.Get(endpoint)
	// Error out if it failed or did not return 200
	if err != nil { return }
	if resp.StatusCode != 200 {
		return errors.New("Error hitting XKCD API: response code not OK ("+strconv.Itoa(resp.StatusCode)+")")
	}
	// Read the body to []bytes
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"name":  x.Name(),
			"error": err,
		}).Error("Command failed")
		return
	}
	// Close the body (it is no longer necessary)
	resp.Body.Close()
	// Instantiate an XKCDComic object and map the JSON to the object
	comic := XKCDComic{}
	err = json.Unmarshal(data, &comic)
	if err != nil { return }
	// Send the message as an embed
	_, err = bot.ChannelMessageSendEmbed(evt.Message.ChannelID, &discordgo.MessageEmbed{
		URL:         "https://xkcd.com/" + comicNumber,
		Title:       "XKCD " + strconv.Itoa(comic.Number) + ": " + comic.SafeTitle,
		Description: comic.Alt,
		Image: &discordgo.MessageEmbedImage{
			URL: comic.Image,
		},
	})
	return
}
