package main

import (
	"bytes"
	"encoding/json"
	"github.com/bclindner/iasipgenerator/iasipgen"
	"github.com/bwmarrin/discordgo"
	"image/jpeg"
	"regexp"
)

// IASIPCommand generates title cards from It's Always Sunny in Philadelphia.
// This is based on a title card generator I wrote some time ago.
type IASIPCommand struct {
	BaseCommand
	IASIPConfig
}

// IASIPConfig is the config for the IASIPCommand.
type IASIPConfig struct {
	// Prefix is the command prefix, any text after which will be made into a title card.
	Prefix string `json:"prefix"`
	// FontPath is the path to the Textile font, in TTF format.
	// Without this, the generator will not function.
	// The generator makes no attempt to ensure this is the right font,
	// so make sure you're loading the right one!
	FontPath string `json:"fontpath"`
	// ImageQuality is the quality of the JPEG to render.
	// If you need to lower bandwidth usage, you may consider lowering this.
	// That said, even at 100 quality, the JPEG file is still under 50kB.
	ImageQuality int `json:"Quality"`
	TriggerRegex *regexp.Regexp
}

// NewIASIPCommand generates a new IASIPCommand.
func NewIASIPCommand(config BaseCommand) (cmd IASIPCommand, err error) {
	options := IASIPConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return cmd, err
	}
	err = iasipgen.LoadFont(options.FontPath)
	if err != nil {
		return cmd, err
	}
	regex, err := regexp.Compile(`^` + options.Prefix + ` ([\S\s]*)$`)
	if err != nil {
		return cmd, err
	}
	options.TriggerRegex = regex
	cmd = IASIPCommand{
		BaseCommand: config,
		IASIPConfig: options,
	}
	return cmd, nil
}

// Test checks if the compiled regex matches the sent string.
func (i IASIPCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return i.TriggerRegex.MatchString(evt.Message.Content)
}

// Run generates an IASIP title card and sends it as a file to the channel.
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
