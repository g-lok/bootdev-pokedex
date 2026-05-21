package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/g-lok/bootdev-pokedex/internal/pokecache"
)

type PokeNamedAPIResourceList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var pokeAPIs = map[string]string{
	"ability":                   "https://pokeapi.co/api/v2/ability/",
	"berry":                     "https://pokeapi.co/api/v2/berry/",
	"berry-firmness":            "https://pokeapi.co/api/v2/berry-firmness/",
	"berry-flavor":              "https://pokeapi.co/api/v2/berry-flavor/",
	"characteristic":            "https://pokeapi.co/api/v2/characteristic/",
	"contest-effect":            "https://pokeapi.co/api/v2/contest-effect/",
	"contest-type":              "https://pokeapi.co/api/v2/contest-type/",
	"egg-group":                 "https://pokeapi.co/api/v2/egg-group/",
	"encounter-condition":       "https://pokeapi.co/api/v2/encounter-condition/",
	"encounter-condition-value": "https://pokeapi.co/api/v2/encounter-condition-value/",
	"encounter-method":          "https://pokeapi.co/api/v2/encounter-method/",
	"evolution-chain":           "https://pokeapi.co/api/v2/evolution-chain/",
	"evolution-trigger":         "https://pokeapi.co/api/v2/evolution-trigger/",
	"gender":                    "https://pokeapi.co/api/v2/gender/",
	"generation":                "https://pokeapi.co/api/v2/generation/",
	"growth-rate":               "https://pokeapi.co/api/v2/growth-rate/",
	"item":                      "https://pokeapi.co/api/v2/item/",
	"item-attribute":            "https://pokeapi.co/api/v2/item-attribute/",
	"item-category":             "https://pokeapi.co/api/v2/item-category/",
	"item-fling-effect":         "https://pokeapi.co/api/v2/item-fling-effect/",
	"item-pocket":               "https://pokeapi.co/api/v2/item-pocket/",
	"language":                  "https://pokeapi.co/api/v2/language/",
	"location":                  "https://pokeapi.co/api/v2/location/",
	"location-area":             "https://pokeapi.co/api/v2/location-area/",
	"machine":                   "https://pokeapi.co/api/v2/machine/",
	"meta":                      "https://pokeapi.co/api/v2/meta/",
	"move":                      "https://pokeapi.co/api/v2/move/",
	"move-ailment":              "https://pokeapi.co/api/v2/move-ailment/",
	"move-battle-style":         "https://pokeapi.co/api/v2/move-battle-style/",
	"move-category":             "https://pokeapi.co/api/v2/move-category/",
	"move-damage-class":         "https://pokeapi.co/api/v2/move-damage-class/",
	"move-learn-method":         "https://pokeapi.co/api/v2/move-learn-method/",
	"move-target":               "https://pokeapi.co/api/v2/move-target/",
	"nature":                    "https://pokeapi.co/api/v2/nature/",
	"pal-park-area":             "https://pokeapi.co/api/v2/pal-park-area/",
	"pokeathlon-stat":           "https://pokeapi.co/api/v2/pokeathlon-stat/",
	"pokedex":                   "https://pokeapi.co/api/v2/pokedex/",
	"pokemon":                   "https://pokeapi.co/api/v2/pokemon/",
	"pokemon-color":             "https://pokeapi.co/api/v2/pokemon-color/",
	"pokemon-form":              "https://pokeapi.co/api/v2/pokemon-form/",
	"pokemon-habitat":           "https://pokeapi.co/api/v2/pokemon-habitat/",
	"pokemon-shape":             "https://pokeapi.co/api/v2/pokemon-shape/",
	"pokemon-species":           "https://pokeapi.co/api/v2/pokemon-species/",
	"region":                    "https://pokeapi.co/api/v2/region/",
	"stat":                      "https://pokeapi.co/api/v2/stat/",
	"super-contest-effect":      "https://pokeapi.co/api/v2/super-contest-effect/",
	"type":                      "https://pokeapi.co/api/v2/type/",
	"version":                   "https://pokeapi.co/api/v2/version/",
	"version-group":             "https://pokeapi.co/api/v2/version-group/",
}

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

func getValidURL(api string, url string) (string, error) {
	prefix, ok := pokeAPIs[api]
	if !ok {
		return "", fmt.Errorf("invalid poke API: %s", api)
	}

	var validURL string
	if url != "" {
		normalizedURL := url
		if after, found := strings.CutPrefix(url, "http://"); found {
			normalizedURL = "https://" + after
		}
		if strings.HasPrefix(normalizedURL, prefix) {
			validURL = url
		} else {
			return "", fmt.Errorf("invalid url(%s) for pokeAPI '%s'", url, api)
		}
	} else {
		validURL = prefix
	}
	return validURL, nil
}

func GETPokeNamedAPIResourceList(api string, url string, cache *pokecache.Cache) (*PokeNamedAPIResourceList, error) {
	validURL, err := getValidURL(api, url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate url: %w", err)
	}

	cacheVal, ok := cache.Get(validURL)
	var data []byte
	if ok {
		data = cacheVal
	} else {
		// Backoff/retry failed calls
		var resp *http.Response
		for _, backoff := range backoffSchedule {
			resp, err = http.Get(validURL)
			if err == nil {
				break
			}
			time.Sleep(backoff)
		}
		if err != nil {
			return nil, fmt.Errorf("http.Get(%s) failed: %w", validURL, err)
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			resp.Body.Close()
			return nil, fmt.Errorf("bad status code: %d from url: %s", resp.StatusCode, validURL)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll() failed to read resp.Body: %w", err)
		}
		err = cache.Add(validURL, data)
		if err != nil {
			return nil, fmt.Errorf("error write to cache: %w", err)
		}
	}

	// Unmarshal data
	var result PokeNamedAPIResourceList
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding pokeapi response: %w", err)
	}

	return &result, nil
}
