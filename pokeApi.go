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

type PokeLocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastStats []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		} `json:"stats"`
	} `json:"past_stats"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationIx struct {
				ScarletViolet struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				BrilliantDiamondShiningPearl struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"brilliant-diamond-shining-pearl"`
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var PokeAPIs = map[string]string{
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

// T is a generic placeholder for whatever struct type being Unmarshaled into
func fetchAndCache[T any](cache *pokecache.Cache, url string) (*T, error) {
	// 1. Check Cache
	cacheVal, ok := cache.Get(url)
	if ok {
		var result T
		if err := json.Unmarshal(cacheVal, &result); err == nil {
			return &result, nil
		}
	}

	// Backoff/retry failed calls
	var resp *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(backoff)
	}
	if err != nil {
		return nil, fmt.Errorf("http.Get(%s) failed: %w", url, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status code: %d from url: %s", resp.StatusCode, url)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = cache.Add(url, data)
	if err != nil {
		return nil, err
	}

	var result T
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getValidURL(api string, url string) (string, error) {
	prefix, ok := PokeAPIs[api]
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
		return nil, err
	}

	return fetchAndCache[PokeNamedAPIResourceList](cache, validURL)
}

func GETPokeLocationAreaDetails(locationArea string, cache *pokecache.Cache) (*PokeLocationArea, error) {
	url := PokeAPIs["location-area"] + locationArea + "/"

	return fetchAndCache[PokeLocationArea](cache, url)
}

func GETPokemon(pokemon string, cache *pokecache.Cache) (*Pokemon, error) {
	url := PokeAPIs["pokemon"] + pokemon + "/"

	return fetchAndCache[Pokemon](cache, url)
}
