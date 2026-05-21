package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/g-lok/bootdev-pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache) error
}

type config struct {
	Next     string
	Previous string
	Args     []string
	Pokedex  map[string]*Pokemon
}

const (
	cmdExit    = "exit"
	cmdHelp    = "help"
	cmdMap     = "map"
	cmdMapb    = "mapb"
	cmdExplore = "explore"
	cmdCatch   = "catch"
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
		cmdExplore: {
			name:        "explore",
			description: "list possible Pokemon encounters in a location-area (> explore [location-area]",
			callback:    commandExplore,
		},
		cmdCatch: {
			name:        "catch",
			description: "try to catch a pokemon",
			callback:    commandCatch,
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

func commandExplore(cfg *config, cache *pokecache.Cache) error {
	if len(cfg.Args) < 1 || len(cfg.Args) > 1 {
		return fmt.Errorf("invalid number of arguments: %d. 'explore' requires one location-area", len(cfg.Args))
	}

	locationArea := cfg.Args[0]
	res, err := GETPokeLocationAreaDetails(locationArea, cache)
	if err != nil {
		return fmt.Errorf("error getting location-area: %w", err)
	}

	for _, pokemon := range res.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func calculatePokemonDC(baseExp int) float64 {
	floor := 36.0
	peak := 340.0 // Normalizes everything up to Legendaries like Mewtwo

	// 1. Check for the "Crazy High" Outliers first (Critical Tier)
	if baseExp > int(peak) {
		// Forces the player to roll a 0.98 or higher (akin to a Nat 20!)
		return 0.98
	}

	// 2. Normalize the distribution for 99% of normal Pokémon
	dc := (float64(baseExp) - floor) / (peak - floor)

	// Scale the DC slightly so early game isn't a total joke
	// E.g., This maps Mewtwo to a 0.90 DC instead of 1.0, keeping it winnable
	scaledDC := dc * 0.90

	// Ensure the DC never drops below 0.0
	if scaledDC < 0.0 {
		return 0.0
	}

	return scaledDC
}

func commandCatch(cfg *config, cache *pokecache.Cache) error {
	if len(cfg.Args) < 1 || len(cfg.Args) > 1 {
		return fmt.Errorf("invalid number of arguments: %d. 'catch' requires one pokemon name", len(cfg.Args))
	}

	pokemon := cfg.Args[0]
	res, err := GETPokemon(pokemon, cache)
	if err != nil {
		return fmt.Errorf("error getting pokemon: %w", err)
	}

	throwMsg := fmt.Sprintf("Throwing a Pokeball at %s...", pokemon)
	fmt.Println(throwMsg)

	normalizedPokeDC := calculatePokemonDC(res.BaseExperience)
	diceRoll := rand.Float64()

	fmt.Printf("pokeXPRatio: %000f\n", normalizedPokeDC)
	fmt.Printf("diceRoll: %000f\n", diceRoll)

	if diceRoll >= normalizedPokeDC {
		msg := fmt.Sprintf("%s was caught!", pokemon)
		fmt.Println(msg)
		cfg.Pokedex[pokemon] = res
	} else {
		msg := fmt.Sprintf("%s escaped!", pokemon)
		fmt.Println(msg)
	}
	return nil
}
