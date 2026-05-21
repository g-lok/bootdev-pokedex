package main

import (
	"fmt"
	"os"

	"github.com/g-lok/bootdev-pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	config      config
	callback    func(*config, *pokecache.Cache) error
}

type config struct {
	Next     string
	Previous string
}

const (
	cmdExit = "exit"
	cmdHelp = "help"
	cmdMap  = "map"
	cmdMapb = "mapb"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		cmdExit: {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		cmdHelp: {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		cmdMap: {
			name:        "map",
			description: "list next 20 location areas",
			callback:    commandMap,
		},
		cmdMapb: {
			name:        "map",
			description: "list previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	pokeCommands := getCommands()
	for k, v := range pokeCommands {
		helpStr := fmt.Sprintf("%s: %s", k, v.description)
		fmt.Println(helpStr)
	}
	return nil
}

func commandMap(cfg *config, cache *pokecache.Cache) error {
	url := cfg.Next
	if cfg.Previous != "" && cfg.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	res, err := GETPokeNamedAPIResourceList("location-area", url, cache)
	if err != nil {
		return fmt.Errorf("error getting location-area: %w", err)
	}
	for _, location := range res.Results {
		fmt.Println(location.Name)
	}
	cfg.Next = res.Next
	cfg.Previous = res.Previous
	return nil
}

func commandMapb(cfg *config, cache *pokecache.Cache) error {
	url := cfg.Previous
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := GETPokeNamedAPIResourceList("location-area", url, cache)
	if err != nil {
		return fmt.Errorf("error getting location-area: %w", err)
	}
	for _, location := range res.Results {
		fmt.Println(location.Name)
	}
	cfg.Next = res.Next
	cfg.Previous = res.Previous
	return nil
}
