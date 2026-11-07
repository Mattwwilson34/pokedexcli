package main

import (
	"cmp"
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecatch"
)

type cliCommand struct {
	name        string
	description string
	callback    func(s *session, args []string) error
}

func commands() map[string]cliCommand {
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
		"mapb": {
			name:        "map back",
			description: "Displays previous location areas",
			callback:    commandMapBack,
		},
		"mapc": {
			name:        "map current",
			description: "Displays current location areas",
			callback:    commandMapCurrent,
		},
		"mapf": {
			name:        "map forward",
			description: "Displays next location areas",
			callback:    commandMapForward,
		},
		"explore": {
			name:        "show location pokemon encounters",
			description: "Displays the pokemon encounters for a location",
			callback:    commandExplore,
		},
		"inspect": {
			name:        "show pokedex information on caught pokemon",
			description: "Displays the pokedex informaiton for a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "show information on a pokemon",
			description: "Displays the information for a pokemon",
			callback:    commandPokedex,
		},
		"catch": {
			name:        "attempt to catch a pokemon",
			description: "Attempts to catch a pokemon",
			callback:    commandCatchPokemon,
		},
		"cache": {
			name:        "print cache stats",
			description: "Displays current cache state values",
			callback:    commandPrintCacheStats,
		},
	}

}

func commandExit(_ *session, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *session, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands() {
		fmt.Printf("%s: %s", cmd.name, cmd.description)
		fmt.Println()
	}
	fmt.Println()
	return nil
}

func commandMapBack(s *session, _ []string) error {
	if s.userLocation.previous == "" {
		return fmt.Errorf("you are on the first page")
	}

	locationAreasData, err := pokeapi.FetchLocationAreas(s.userLocation.previous, s.cache)
	if err != nil {
		return err
	}

	s.userLocation.previous = locationAreasData.Previous
	s.userLocation.current = s.userLocation.previous
	s.userLocation.next = locationAreasData.Next

	printLocationAreas(locationAreasData, s.userLocation.current)
	return nil
}

func commandMapCurrent(s *session, _ []string) error {
	if s.userLocation.current == "" {
		return fmt.Errorf("you have not navigated to any location yet")
	}

	locationAreasData, err := pokeapi.FetchLocationAreas(s.userLocation.current, s.cache)
	if err != nil {
		return err
	}

	printLocationAreas(locationAreasData, s.userLocation.current)
	return nil
}

func commandMapForward(s *session, _ []string) error {
	nextUrl := cmp.Or(s.userLocation.next, pokeapi.FirstLocationPage)
	locationAreasData, err := pokeapi.FetchLocationAreas(nextUrl, s.cache)
	if err != nil {
		return err
	}

	s.userLocation.previous = locationAreasData.Previous
	s.userLocation.current = nextUrl
	s.userLocation.next = locationAreasData.Next

	printLocationAreas(locationAreasData, s.userLocation.current)
	return nil
}

func commandExplore(s *session, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("must specify a location explore <location>")
	}

	currentLocationUrl := pokeapi.LocationAreasURL + args[1]
	currentLocation, err := pokeapi.FetchLocationArea(currentLocationUrl, s.cache)
	if err != nil {
		return err
	}

	fmt.Println("\nPOKEMON")
	fmt.Println("======================")
	for _, pokemonVersion := range currentLocation.PokemonEncounters {
		fmt.Println(pokemonVersion.Pokemon.Name)
	}
	fmt.Println()

	return nil
}

func commandPokedex(s *session, _ []string) error {
	if len(s.pokedex.Entries) == 0 {
		fmt.Println("you have not caught any pokemon")
	}
	for _, p := range s.pokedex.Entries {
		fmt.Printf("- %s", p.Name)
	}
	return nil
}

func commandInspect(s *session, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("must specify a pokemon")
	}

	pokemonName := args[1]
	pokedexEntry, ok := s.pokedex.Entries[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Println("\nPOKEMON")
		fmt.Println("======================")
		fmt.Println(pokedexEntry.String())
		fmt.Println()
	}

	return nil
}

func commandCatchPokemon(s *session, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("must specify a pokemon")
	}

	defaultCatchRate := float32(0.25)
	pokemonName := args[1]

	pokedexEntry, ok := s.pokedex.Entries[pokemonName]
	fmt.Printf("Throwing a Pokeball at %s...", pokemonName)
	fmt.Println()

	// No log entry for pokemon so we fetch data from api
	if !ok {
		pokemonURL := pokeapi.PokemonURL + args[1]
		pokemon, err := pokeapi.FetchPokemon(pokemonURL, s.cache)
		if err != nil {
			return err
		}

		difficultyFactor := float32(pokemon.BaseExperience) / float32(100)
		modifiedCatchRate := float32(defaultCatchRate / difficultyFactor)

		catchRoll := rand.Float32()
		pokemonCaught := catchRoll < modifiedCatchRate
		if pokemonCaught {
			// Pokemon caught update log with catch rate set to default
			fmt.Printf("%s was caught!", pokemonName)
			fmt.Println()
			newPokedexEntry := pokecatch.PokedexEntry{
				Name:                  pokemonName,
				TotalCaught:           1,
				BaseExperience:        pokemon.BaseExperience,
				CatchDifficultyFactor: difficultyFactor,
				CatchRate:             defaultCatchRate,
				Stats:                 pokemon.Stats,
				Types:                 pokemon.Types,
			}
			s.pokedex.Entries[pokemonName] = newPokedexEntry
			fmt.Println("#######################")
			fmt.Println("   POKEDEX UPDATED")
			fmt.Println("#######################")
			fmt.Println(newPokedexEntry.String())
		} else {
			// Pokemon not caught increase likelyhood of catch for next attempt
			fmt.Printf("%s escaped!", pokemonName)
			fmt.Println()
			updatedPokedexEntry := pokecatch.PokedexEntry{
				Name:                  pokemonName,
				TotalCaught:           0,
				BaseExperience:        pokemon.BaseExperience,
				CatchDifficultyFactor: difficultyFactor,
				CatchRate:             modifiedCatchRate * defaultCatchRate,
				Stats:                 pokemon.Stats,
				Types:                 pokemon.Types,
			}
			s.pokedex.Entries[pokemonName] = updatedPokedexEntry
		}
	} else {
		// Pokemon log entry exists skip fetch and attempt catch
		catchRoll := rand.Float32()
		pokemonCaught := catchRoll < pokedexEntry.CatchRate
		if pokemonCaught {
			// Pokemon caught update log with catch rate set to default
			fmt.Printf("%s was caught!", pokemonName)
			fmt.Println()
			pokedexEntry.TotalCaught++
			pokedexEntry.CatchRate = defaultCatchRate
			s.pokedex.Entries[pokemonName] = pokedexEntry
			fmt.Println("#######################")
			fmt.Println("   POKEDEX UPDATED")
			fmt.Println("#######################")
			fmt.Println(pokedexEntry.String())
		} else {
			// Pokemon not caught increase likelyhood of catch for next attempt
			fmt.Printf("%s escaped!", pokemonName)
			fmt.Println()
			pokedexEntry.CatchRate += defaultCatchRate
			s.pokedex.Entries[pokemonName] = pokedexEntry
		}
	}

	return nil
}

func commandPrintCacheStats(s *session, _ []string) error {
	err := s.cache.PrintCacheStats()
	if err != nil {
		return err
	}
	return nil
}

func printLocationAreas(
	data pokeapi.Pagination[pokeapi.LocationAreas],
	currentLocation string,
) {
	fmt.Println()
	fmt.Printf("MAP: %s\n", currentLocation)
	fmt.Println("================================")
	for _, locationArea := range data.Results {
		fmt.Println(locationArea.Name)
	}
	fmt.Println()
}
