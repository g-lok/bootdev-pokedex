package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	config      config
	callback    func(*config) error
}

type config struct {
	Next     string
	Previous string
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "list next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "list previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
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

func commandMap(cfg *config) error {
	url := cfg.Next
	if cfg.Previous != "" && cfg.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	res, err := GETPokeNamedAPIResourceList("location-area", url)
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

func commandMapb(cfg *config) error {
	url := cfg.Previous
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := GETPokeNamedAPIResourceList("location-area", url)
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
