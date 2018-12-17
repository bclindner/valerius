package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"  // for running the bot
	log "github.com/sirupsen/logrus" // logging suite
	"github.com/tidwall/gjson"       // for getting items in dot notation
	"io/ioutil"                      // for opening response body
	"net/http"
	"net/url"
	"regexp"
)

// RESTCommand base structure.
type RESTCommand struct {
	BaseCommand
	RESTConfig
	regexp         *regexp.Regexp
	endpointstring string
	endpointgroups []int
	client         http.Client
}

// RESTConfig is the configuration for the RESTCommand.
type RESTConfig struct {
	TriggerRegex string        `json:"triggerregex"`
	Endpoint     []interface{} `json:"endpoint"`
	Method       string        `json:"method"`
	Response     []string      `json:"response"`
	Responses    [][]string    `json:"responses"`
	ErrorMessage string        `json:"errorMessage"`
}

// NewRESTCommand generates a new RESTCommand.
func NewRESTCommand(config BaseCommand) (command RESTCommand, err error) {
	var options RESTConfig
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return command, nil
	}
	// Ensure there is only one of Response or Responses
	if len(options.Response) > 0 && len(options.Responses) > 0 {
		return command, errors.New("Can only have one of 'response' or 'responses'")
	}
	// Ensure the endpoint and response commands are of their correct types.
	endpoint, ok := options.Endpoint[0].(string)
	if !ok {
		return command, errors.New("First of endpoint array should be a string")
	}
	var endpointgroups []int
	for _, item := range options.Endpoint[1:] {
		// it HAS to cast to float64 because of the json package,
		// but this means it allows non-integer numbers without whining which is PURE JANK
		// gfdi
		i, ok := item.(float64)
		if !ok {
			return command, errors.New("All items after string in endpoint must be numbers")
		}
		endpointgroups = append(endpointgroups, int(i))
	}
	// Instantiate the regex.
	rgx, err := regexp.Compile(options.TriggerRegex)
	if err != nil {
		return command, nil
	}
	// Sanity check: is the number of endpoint groups the number of groups in the regex?
	// The command will panic otherwise
	if len(endpointgroups) != rgx.NumSubexp() {
		return command, errors.New("Number of groups in trigger does not match number of groups in regex")
	}
	// generate the command
	command = RESTCommand{
		BaseCommand:    config,
		RESTConfig:     options,
		regexp:         rgx,
		endpointstring: endpoint,
		endpointgroups: endpointgroups,
		client:         http.Client{},
	}
	return command, nil
}

func (r RESTCommand) SendErrorMessage(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	if len(r.ErrorMessage) > 0 {
		bot.ChannelMessageSend(evt.Message.ChannelID, r.ErrorMessage)
	}
}

// Test ensures the compiled regex passes.
func (r RESTCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return r.regexp.MatchString(evt.Message.Content)
}

// Run hits the given REST endpoint, gets a comic, and returns it as an embed.
func (r RESTCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) (err error) {
	// Construct the endpoint
	rgxgroups := r.regexp.FindAllStringSubmatch(evt.Message.Content, -1)[0]
	var reqfmtgroups []interface{}
	for _, i := range r.endpointgroups {
		reqfmtgroups = append(reqfmtgroups, url.QueryEscape(rgxgroups[i]))
	}
	endpoint := fmt.Sprintf(r.endpointstring, reqfmtgroups...)
	// Construct request based on this endpoint
	request, err := http.NewRequest(r.Method, endpoint, nil)
	if err != nil {
		r.SendErrorMessage(bot, evt)
		return err
	}
	// Log that we're about to send the request, in case someone's trying something nasty
	log.WithFields(log.Fields{
		"endpoint": endpoint,
		"method":   r.Method,
	}).Info("Making HTTP request")
	// Send request, ensure nothing failed, get JSON bytes
	resp, err := r.client.Do(request)
	if err != nil {
		r.SendErrorMessage(bot, evt)
		return errors.New("could not make request: " + err.Error())
	}
	// Log some response metadata, again, in case someone's being nasty
	log.WithFields(log.Fields{
		"endpoint": endpoint,
		"response": resp.Status,
	}).Info("HTTP request result")
	if resp.StatusCode >= 400 {
		r.SendErrorMessage(bot, evt)
		return errors.New("request failed with status " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.SendErrorMessage(bot, evt)
		return errors.New("could not parse request body: " + err.Error())
	}
	// For single-response:
	if len(r.Response) > 0 {
		// Get the JSON objects needed to format the response
		var respfmtgroups []interface{}
		items := gjson.GetManyBytes(body, r.Response[1:]...)
		for _, item := range items {
			// Break if a particular lookup is missing
			if !item.Exists() {
				break
			}
			respfmtgroups = append(respfmtgroups, item.Value())
		}
		// Format and send the response
		bot.ChannelMessageSend(evt.Message.ChannelID, fmt.Sprintf(r.Response[0], respfmtgroups...))
		return nil
	}
	// For multi-response:
	if len(r.Responses) > 0 {
	ToNextResponse:
		for _, response := range r.Responses {
			// Get the JSON objects needed to format the response
			var respfmtgroups []interface{}
			items := gjson.GetManyBytes(body, response[1:]...)
			for _, item := range items {
				// Continue into the next response if a particular lookup is missing
				if !item.Exists() {
					continue ToNextResponse
				}
				respfmtgroups = append(respfmtgroups, item.Value())
			}
			// Format and send the response
			bot.ChannelMessageSend(evt.Message.ChannelID, fmt.Sprintf(response[0], respfmtgroups...))
			return nil
		}
	}
	// If this code is reached, no response was valid, which probably shouldn't happen, so we'll throw an error
	r.SendErrorMessage(bot, evt)
	return errors.New("No valid response schema")
}
