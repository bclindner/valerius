package main

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"  // for running the bot
	"io/ioutil"                      // for opening response body
	"net/http"
	"regexp"
	"strconv"
)

// The XKCDCommand base structure. Takes a Regexp to test the command.
type XKCDCommand struct {
	BaseCommand
	// Regexp to test with. This is set by the factory function.
	regexp *regexp.Regexp
}

// XKCD API structure, for parsing the XKCD comic API.
type XKCDComic struct {
	// Comic number.
	Number    int    `json:"num"`
	// Comic title.
	Title     string `json:"title"`
	// Safe title (unsure what this does? safe encoding maybe?)
	SafeTitle string `json:"safe_title"`
	// Alt text for the comic.
	Alt       string `json:"alt"`
	// Image URL.
	Image     string `json:"img"`
}

// Generate a new XKCDCommand.
func NewXKCDCommand(config CommandConfig) (command XKCDCommand, err error) {
	// Instantiate the regex.
	rgx, err := regexp.Compile(`^\!xkcd ?([0-9]+)?$`)
	if err != nil { return }
	// generate the command
	command = XKCDCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Name,
		},
		regexp: rgx,
	}
	return
}

func (x XKCDCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return x.regexp.MatchString(evt.Message.Content)
}

func (x XKCDCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// Get the comic number from the regex group matches
	comicNumber := x.regexp.FindStringSubmatch(evt.Message.Content)[1]
	// Get the endpoint necessary
	// The most recent comic has a slightly different endpoint format
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
		return errors.New("Error hitting XKCD API: response code not OK (" + strconv.Itoa(resp.StatusCode) + ")")
	}
	// Read the body to []bytes (to parse as JSON)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil { return }
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
