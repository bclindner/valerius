package main

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo" // for running the bot
	"math/rand"
	"regexp"
	"time"
)

// PingPongCommand is a generic command that sends responses based on triggers.
type PingPongCommand struct {
	BaseCommand
	PingPongConfig
	RNG          *rand.Rand
	Regexp       *regexp.Regexp
	TriggerType  int
	ResponseType int
}

// Trigger types.
const (
	// Set if "trigger" is set in the config.
	// Only triggers if one exact string is matched.
	triggerSingle int = iota
	// Set if "triggers" is set in the config.
	// Triggers if any one defined string is exactly matched.
	triggerMultiple
	// Set if "triggerregex" is set in the config.
	// Triggers if the message matches the given regular expression.
	triggerRegex
)
const (
	// Set if "response" is set in the config.
	// Sends a single, static response.
	responseSingle = iota
	// Set if "responses" is set in the config.
	// Sends one response from a list pseudo-randomly.
	responseMultiple
)

// PingPongConfig is the configurator for the PingPong command.
type PingPongConfig struct {
	// Regular expression to trigger the command.
	TriggerRegex string `json:"triggerregex"`
	// List of messages that may trigger the command.
	Triggers []string `json:"triggers"`
	// Message that may trigger the command.
	Trigger string `json:"trigger"`
	// Response to send if the command is triggered.
	Response string `json:"response"`
	// List of responses to randomly send from if the command is triggered.
	Responses []string `json:"responses"`
	// Prefix to put before each response.
	ResponsePrefix string `json:"responseprefix"`
	// Suffix to put after each response.
	ResponseSuffix string `json:"responsesuffix"`
}

// NewPingPongCommand creates a new PingPongCommand.
func NewPingPongCommand(config CommandConfig) (command PingPongCommand, err error) {
	// Parse config
	options := PingPongConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return command, err
	}
	// Sanity check: cannot have more than one of Trigger, Triggers, or TriggerRegex
	// Also determine trigger type here (to keep from doing a bunch of extra stuff on Test()
	actives := 0
	ttype := -1
	if len(options.Trigger) > 0 {
		actives++
		ttype = triggerSingle
	}
	if len(options.Triggers) > 0 {
		actives++
		ttype = triggerMultiple
	}
	if len(options.TriggerRegex) > 0 {
		actives++
		ttype = triggerRegex
	}
	if actives > 1 {
		return command, errors.New("Cannot have more than one of 'trigger', 'triggers', or 'triggerregex' in the same PingPongCommand")
	}
	// Sanity check: cannot have Response and Responses in the same command
	if len(options.Response) > 0 && len(options.Responses) > 0 {
		return command, errors.New("Cannot have 'response' and 'responses' in the same PingPongCommand")
	}
	// Determine response type (to keep from doing a bunch of extra stuff on Test())
	rtype := -1
	if len(options.Response) > 0 {
		rtype = responseSingle
	}
	if len(options.Responses) > 0 {
		rtype = responseMultiple
	}
	// Initialize command
	command = PingPongCommand{
		BaseCommand: BaseCommand{
			Name: config.Name,
			Type: config.Type,
		},
		PingPongConfig: options,
		TriggerType:    ttype,
		ResponseType:   rtype,
	}
	// Initialize regex, if necessary
	if len(options.TriggerRegex) > 0 {
		command.Regexp, err = regexp.Compile(options.TriggerRegex)
	}
	// Initialize RNG, if necessary
	if len(options.Responses) > 1 {
		command.RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return command, nil
}

func (p PingPongCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	// Run the necessary test based on set trigger type
	switch p.TriggerType {
	case triggerSingle:
		if len(p.Trigger) > 0 {
			if evt.Message.Content == p.Trigger {
				return true
			}
		}
	case triggerMultiple:
		if len(p.Triggers) > 0 {
			for _, trigger := range p.Triggers {
				if evt.Message.Content == trigger {
					return true
				}
			}
		}
	case triggerRegex:
		if len(p.TriggerRegex) > 0 {
			if p.Regexp.MatchString(evt.Message.Content) {
				return true
			}
		}
	default: //uhhhHHHH
		panic("No test for " + p.GetName())
	}
	return false
}
func (p PingPongCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// If Response is used, send the response
	if len(p.Response) > 0 {
		// Send the response
		_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.ResponsePrefix+p.Response+p.ResponseSuffix)

	}
	// If Responses are used, send a random response from the list
	if len(p.Responses) > 0 {
		i := p.RNG.Intn(len(p.Responses))
		// Send the response
		_, err = bot.ChannelMessageSend(evt.Message.ChannelID, p.ResponsePrefix+p.Responses[i]+p.ResponseSuffix)
	}
	return
}
