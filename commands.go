package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ardamertdedeoglu/pokedexcli/internal/pokecache"
)

var cache *pokecache.Cache

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
	config      config
}

type config struct {
	Next     string
	Previous string
}

var commands = make(map[string]cliCommand)
var inventory = make(map[string]Pokemon)

func init() {
	cache = pokecache.NewCache(5 * time.Second)
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Maps 20 next locations",
		callback:    commandMap,
		config: config{
			Next:     "https://pokeapi.co/api/v2/location-area/",
			Previous: "",
		},
	}
	commands["mapb"] = cliCommand{
		name:        "mapback",
		description: "Maps 20 previous locations",
		callback:    commandMapBack,
		config: config{
			Next:     "https://pokeapi.co/api/v2/location-area/",
			Previous: "",
		},
	}
	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explores a locations pokemons.",
		callback:    commandExplore,
	}

	commands["catch"] = cliCommand{
		name:        "catch",
		description: "Tries to catch a pokemon.",
		callback:    commandCatch,
	}

	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Inspect an already caught pokemon",
		callback:    commandInspect,
	}

	commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List the caught pokemons.",
		callback:    commandPokedex,
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}

	return nil
}

type location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonEncounter struct {
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

func commandMap(cfg *config, args ...string) error {
	var result location
	baseURL := cfg.Next

	if res, ok := cache.Get(baseURL); ok {
		err := json.Unmarshal(res, &result)
		fmt.Println("command from cache")
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(baseURL)
		fmt.Println("cache no")
		if err != nil {
			fmt.Printf("error getting locations")
			return err
		}
		bodyData, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bodyData, &result)
		if err != nil {
			fmt.Printf("error decoding locations %s", err)
			return err
		}
		cache.Add(baseURL, bodyData)
	}

	next := result.Next
	previous := result.Previous
	locations := result.Results
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	if previous == nil {
		previous = ""
	}
	newConfig := config{
		Next:     next,
		Previous: cfg.Next,
	}
	tmp := commands["map"]
	tmp.config = newConfig
	commands["map"] = tmp

	return nil
}

func commandMapBack(cfg *config, args ...string) error {
	var result location
	baseURL := cfg.Previous

	if commands["map"].config.Next == "https://pokeapi.co/api/v2/location-area/" {
		fmt.Println("you're on the first page")
		return nil
	}

	if res, ok := cache.Get(baseURL); ok {
		fmt.Println("mapb cache")
		err := json.Unmarshal(res, &result)
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(baseURL)
		if err != nil {
			fmt.Printf("error getting locations")
			return err
		}
		bodyData, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bodyData, &result)
		if err != nil {
			fmt.Printf("error decoding locations %s", err)
			return err
		}
		cache.Add(baseURL, bodyData)
	}

	next := result.Next
	previous := result.Previous
	locations := result.Results
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	if previous == nil {
		previous = ""
	}
	newConfig := config{
		Next:     next,
		Previous: previous.(string),
	}
	tmp := commands["map"]
	tmp.config = newConfig
	commands["map"] = tmp

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	cityName := args[0]
	fmt.Printf("Exploring %s...\n", cityName)
	baseURL := "https://pokeapi.co/api/v2/location-area/" + cityName + "/"

	var result PokemonEncounter

	if res, ok := cache.Get(baseURL); ok {
		err := json.Unmarshal(res, &result)
		fmt.Println("command from cache")
		if err != nil {
			return err
		}
	} else {
		res, err := http.Get(baseURL)
		fmt.Println("cache no")
		if err != nil {
			fmt.Printf("error getting locations")
			return err
		}
		bodyData, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bodyData, &result)
		if err != nil {
			fmt.Printf("error decoding locations %s", err)
			return err
		}
		cache.Add(baseURL, bodyData)

	}

	if len(result.EncounterMethodRates) == 0 {
		return fmt.Errorf("Found no pokemon.")
	}

	for _, val := range result.PokemonEncounters {
		fmt.Println("- " + val.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	pokemonName := args[0]
	baseURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName + "/"
	var result Pokemon

	res, err := http.Get(baseURL)
	if err != nil {
		fmt.Printf("error getting locations")
		return err
	}
	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyData, &result)
	if err != nil {
		fmt.Printf("error decoding locations %s", err)
		return err
	}

	if _, ok := inventory[pokemonName]; ok {
		return fmt.Errorf("You already have this pokemon.")
	} else {
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
		catchChance := rand.Intn(1000)
		if result.BaseExperience < catchChance {
			fmt.Printf("%s was caught!\n", pokemonName)
			inventory[pokemonName] = result
			fmt.Println("You may now inspect it with the inspect command.")
		} else {
			fmt.Printf("%s escaped!\n", pokemonName)
		}
	}
	return nil

}

func commandInspect(cfg *config, args ...string) error {
	pokemonName := args[0]
	if res, ok := inventory[pokemonName]; !ok {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", res.Name)
		fmt.Printf("Height: %v\n", res.Height)
		fmt.Printf("Weight: %v\n", res.Weight)
		fmt.Printf("Stats:\n")
		for _, val := range res.Stats {
			fmt.Printf("  -%s: %v\n", val.Stat.Name, val.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, val := range res.Types {
			fmt.Printf("  - %s\n", val.Type.Name)
		}
	}
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(inventory) == 0 {
		return fmt.Errorf("You have no pokemons.")
	}
	for _, val := range inventory {
		fmt.Printf(" - %s\n", val.Name)
	}
	return nil
}
