package main

import (
	log "github.com/sirupsen/logrus" // logging suite
	"github.com/bwmarrin/discordgo" // for running the bot
	"encoding/json"
	"net/http"
	"io/ioutil" // for opening response body
	"regexp"
)

type XKCDCommand struct {
	Command
	Regexp *regexp.Regexp
}

type XKCDComic struct {
	Number int `json:"num"`
	Title string `json:"title"`
	SafeTitle string `json:"safe_title"`
	Alt string `json:"alt"`
	Image string `json:"img"`
}

func (p XKCDCommand) Name() string {
	return "XKCD Viewer"
}

func NewXKCDCommand() XKCDCommand {
	rgx , err := regexp.Compile(`^\!xkcd ?([0-9]+)?$`)
	if err != nil { log.Fatal(err) }
	return XKCDCommand{
		Regexp: rgx,
	}
}

func (x XKCDCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return x.Regexp.MatchString(evt.Message.Content)
}

func (x XKCDCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	comicNumber := x.Regexp.FindStringSubmatch(evt.Message.Content)[1]
	var endpoint string
	if comicNumber != "" {
		endpoint = "https://xkcd.com/"+comicNumber+"/info.0.json"
	} else {
		endpoint = "https://xkcd.com/info.0.json"
	}
	resp, err := http.Get(endpoint)
	if err != nil || resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"name": x.Name(),
			"statusCode": resp.StatusCode,
			"error": err,
		}).Error("Command failed")
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"name": x.Name(),
			"error": err,
		}).Error("Command failed")
		return
	}
	comic := XKCDComic{}
	err = json.Unmarshal(data, &comic)
	if err != nil {
		log.WithFields(log.Fields{
			"name": x.Name(),
			"error": err,
		}).Error("Command failed")
		return
	}
	bot.ChannelMessageSendEmbed(evt.Message.ChannelID, &discordgo.MessageEmbed{
		URL: "https://xkcd.com/"+comicNumber,
		Title: comic.SafeTitle,
		Description: comic.Alt,
		Image: &discordgo.MessageEmbedImage{
			URL: comic.Image,
		},
	})
}
