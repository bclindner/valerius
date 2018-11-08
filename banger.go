package main

import (
	"math/rand"
	"time"
)

// Base command for the Banger commands.
type BangerCommand struct {
	BaseCommand
	// List of bangers. Usually YouTube links or similar.
	Bangers *[]string
	// Random number generator for getting a random banger from the list.
	RNG *rand.Rand
}

// Gets a random banger from the list.
func (b BangerCommand) GetRandomBanger() string {
	i := b.RNG.Intn(len(*b.Bangers))
	return (*b.Bangers)[i]
}

// Initializes a new BangerCommand.
func NewBangerCommand(name string, bangers *[]string) BangerCommand {
	// generate the Rand object to use in the object
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return BangerCommand{
		BaseCommand: BaseCommand{
			name: name,
		},
		Bangers: bangers,
		RNG:     rng,
	}
}
