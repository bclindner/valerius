package main

import (
	"bytes"
	"encoding/json"
	"github.com/bclindner/iasipgenerator/iasipgen"
	"github.com/bwmarrin/discordgo"
	"image/jpeg"
	"regexp"
)

type IASIPCommand struct {
	BaseCommand
	IASIPConfig
}

type IASIPConfig struct {
	Prefix       string `json:"prefix"`
	FontPath     string `json:"fontpath"`
	ImageQuality int    `json:"Quality"`
	TriggerRegex *regexp.Regexp
}

func NewIASIPCommand(config CommandConfig) (cmd IASIPCommand, err error) {
	options := IASIPConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return cmd, err
	}
	err = iasipgen.LoadFont("textile.ttf")
	if err != nil {
		return cmd, err
	}
	regex, err := regexp.Compile("^" + options.Prefix + " (.*)$")
	if err != nil {
		return cmd, err
	}
	options.TriggerRegex = regex
	cmd = IASIPCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Type,
		},
		IASIPConfig: options,
	}
	return cmd, nil
}

func (i IASIPCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return i.TriggerRegex.MatchString(evt.Message.Content)
}

func (i IASIPCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	msgstring := i.TriggerRegex.FindStringSubmatch(evt.Message.Content)[1]
	img, err := iasipgen.Generate(msgstring)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, &jpeg.Options{
		Quality: i.ImageQuality,
	})
	if err != nil {
		return err
	}
	bot.ChannelFileSend(evt.Message.ChannelID, "iasip.jpg", buf)
	return nil
}
