package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// ReloadCommand is a meta-command which reloads commands.
// This should generally only be triggered by admins and bot owners.
// Keep it properly whitelisted.
type ReloadCommand struct {
	BaseCommand
	ReloadConfig
}

// ReloadConfig is the config for the ReloadCommand.
type ReloadConfig struct {
	Trigger string `json:"trigger"`
}

// NewReloadCommand generates a new ReloadCommand.
// Aside from the trigger, no configuration is needed, so this is particularly short.
func NewReloadCommand(config BaseCommand) (cmd ReloadCommand, err error) {
	options := ReloadConfig{}
	err = json.Unmarshal(config.Options, &options)
	if err != nil {
		return cmd, err
	}
	cmd = ReloadCommand{
		BaseCommand:  config,
		ReloadConfig: options,
	}
	return cmd, nil
}

// Test checks if the trigger was sent.
func (c ReloadCommand) Test(bot *discordgo.Session, evt *discordgo.MessageCreate) bool {
	return c.Trigger == evt.Message.Content
}

// Run reloads commands.
func (c ReloadCommand) Run(bot *discordgo.Session, evt *discordgo.MessageCreate) error {
	// Re-read bot config
	var err error
	config, err = ReadBotConfig(*configPath)
	if err != nil {
		bot.ChannelMessageSend(evt.Message.ChannelID, "Failed to reload commands.")
		return err
	}
	// Try to make the new handler
	newhandler, err := NewMessageHandler(bot, config.Commands)
	if err != nil {
		bot.ChannelMessageSend(evt.Message.ChannelID, "Failed to reload commands.")
		return err
	}
	// Destroy the old one and use this new one
	// NOTE Doesn't this mean a small window in which commands could double-trigger?
	handler.DestroySelf()
	handler = newhandler
	// Log the success
	_, err = bot.ChannelMessageSend(evt.Message.ChannelID, fmt.Sprintf("Commands reloaded! Parsed %d commands.", len(handler.commands)))
	if err != nil {
		return err
	}
	return nil
}
