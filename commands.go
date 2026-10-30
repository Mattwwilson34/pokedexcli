package main

import (
	"cmp"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(s *session) error
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
		"cache": {
			name:        "print cache stats",
			description: "Displays current cache state values",
			callback:    commandPrintCacheStats,
		},
	}

}

func commandExit(_ *session) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *session) error {
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

func commandMapBack(s *session) error {
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

func commandMapCurrent(s *session) error {
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

func commandMapForward(s *session) error {
	nextUrl := cmp.Or(s.userLocation.next, pokeapi.LocationAreasUrl)
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

func commandPrintCacheStats(s *session) error {
	err := s.cache.PrintCacheStats()
	if err != nil {
		return err
	}
	return nil
}

func printLocationAreas(
	data pokeapi.Pagination[pokeapi.LocationArea],
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
