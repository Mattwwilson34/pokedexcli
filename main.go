package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/internal"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type locationState struct {
	next     *string
	previous *string
}

type REPL struct {
	commands      map[string]cliCommand
	locationState locationState
	cache         *internal.Cache
}

func newREPL() *REPL {
	r := &REPL{}

	r.cache = internal.NewCache(time.Second * 60)

	r.commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    r.commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    r.commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays location areas",
			callback:    r.commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Displays previous location areas",
			callback:    r.commandMapBack,
		},
	}

	return r
}

func main() {

	repl := newREPL()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput) == 0 {
			continue
		}

		command, ok := repl.commands[cleanedInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback()
		repl.cache.PrintCacheStats()
		if err != nil {
			fmt.Println(err)
		}

	}
}

func (r *REPL) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (r *REPL) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range r.commands {
		fmt.Printf("%s: %s", cmd.name, cmd.description)
		fmt.Println()
	}
	fmt.Println()
	return nil
}

type Pagination[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r *REPL) commandMap() error {
	return r.showPage(r.locationState.next)
}

func (r *REPL) commandMapBack() error {
	if r.locationState.previous == nil {
		return fmt.Errorf("you are on the first page")
	}
	return r.showPage(r.locationState.previous)
}

func fetchLocations(url string, cache *internal.Cache) (Pagination[Location], error) {

	res, err := http.Get(url)
	if err != nil {
		return Pagination[Location]{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pagination[Location]{}, err
	}
	if res.StatusCode > 299 {
		return Pagination[Location]{}, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}

	err = cache.Add(url, body)
	if err != nil {
		return Pagination[Location]{}, err
	}

	var page Pagination[Location]
	if err := json.Unmarshal(body, &page); err != nil {
		return Pagination[Location]{}, err
	}
	return page, nil

}

func (r *REPL) showPage(pageURL *string) error {
	url := baseUrl
	if pageURL != nil {
		url = *pageURL
	}

	var page Pagination[Location]
	cachedPage, cachedPageExists := r.cache.Get(url)

	if cachedPageExists {
		if err := json.Unmarshal(cachedPage, &page); err != nil {
			return err
		}
	}

	if !cachedPageExists {
		nonCachedPage, err := fetchLocations(url, r.cache)
		if err != nil {
			return err
		}
		page = nonCachedPage
	}

	r.locationState.next = page.Next
	r.locationState.previous = page.Previous

	for _, location := range page.Results {
		fmt.Println(location.Name)
	}
	return nil
}
