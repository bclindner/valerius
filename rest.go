package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo" // for running the bot
	"github.com/gregjones/httpcache"
	log "github.com/sirupsen/logrus" // logging suite
	"io/ioutil"                      // for opening response body
	"net/http"
	"net/url"
	"regexp"
	"text/template"
)

// RESTCommand base structure.
type RESTCommand struct {
	BaseCommand
	RESTConfig
	regexp         *regexp.Regexp
	endpointstring string
	endpointgroups []int
	template       *template.Template
	client         http.Client
}

// RESTConfig is the configuration for the RESTCommand.
type RESTConfig struct {
	TriggerRegex     string            `json:"triggerregex"`
	Endpoint         []interface{}     `json:"endpoint"`
	Method           string            `json:"method"`
	Response         string            `json:"response"`
	ResponseFilepath string            `json:"responseFile"`
	ErrorMessage     string            `json:"errorMessage"`
	Headers          map[string]string `json:"headers"`
	DisableCache     bool              `json:"disablecache"`
}

// NewRESTCommand generates a new RESTCommand.
func NewRESTCommand(config BaseCommand) (command RESTCommand, err error) {
	var options RESTConfig
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return command, nil
	}
	// Ensure only one of Response and ResponseFilepath is set
	if len(options.Response) > 0 && len(options.ResponseFilepath) > 0 {
		return command, errors.New("Can only have one of response and responseFile")
	}
	// Get template text to use
	var tmplstr string
	if len(options.Response) > 0 {
		tmplstr = options.Response
	} else {
		tmplbytes, err := ioutil.ReadFile(options.ResponseFilepath)
		if err != nil {
			return command, errors.New("Error reading response file: " + err.Error())
		}
		tmplstr = string(tmplbytes)
	}
	// Compile the template
	tmpl, err := template.New(config.Name).Parse(tmplstr)
	if err != nil {
		return command, errors.New("Failed to compile template: " + err.Error())
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
		template:       tmpl,
	}
	// set the client based on if this restcommand is cached
	if options.DisableCache {
		command.client = http.Client{}
	} else {
		// use a caching transport to stop the bot from flooding servers with identical requests, if the config allows
		command.client = http.Client{
			Transport: httpcache.NewMemoryCacheTransport(),
		}
		command.client = http.Client{}
	}
	return command, nil
}

func (r RESTCommand) sendErrorMessage(bot *discordgo.Session, evt *discordgo.MessageCreate) {
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
		r.sendErrorMessage(bot, evt)
		return err
	}
	// Set headers
	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}
	// Log that we're about to send the request, in case someone's trying something nasty
	log.WithFields(log.Fields{
		"endpoint": endpoint,
		"method":   r.Method,
	}).Info("Making HTTP request")
	// Send request, ensure nothing failed, get JSON bytes
	resp, err := r.client.Do(request)
	if err != nil {
		r.sendErrorMessage(bot, evt)
		return errors.New("could not make request: " + err.Error())
	}
	// Log some response metadata, again, in case someone's being nasty
	log.WithFields(log.Fields{
		"endpoint": endpoint,
		"response": resp.Status,
	}).Info("HTTP request result")
	if resp.StatusCode >= 400 {
		r.sendErrorMessage(bot, evt)
		return errors.New("request failed with status " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.sendErrorMessage(bot, evt)
		return errors.New("could not read request body: " + err.Error())
	}
	// Unmarshal the response JSON into an interface{}
	var bodyjson interface{}
	err = json.Unmarshal(body, &bodyjson)
	if err != nil {
		r.sendErrorMessage(bot, evt)
		return errors.New("could not unmarshal request body: " + err.Error())
	}
	msgbuf := new(bytes.Buffer)
	err = r.template.Execute(msgbuf, bodyjson)
	if err != nil {
		r.sendErrorMessage(bot, evt)
		return errors.New("could not execute template: " + err.Error())
	}
	bot.ChannelMessageSend(evt.Message.ChannelID, msgbuf.String())
	return nil
}
