package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {

	session := &session{
		cache: pokecache.NewCache(time.Second * 60),
	}
	commands := commands()

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

		command, ok := commands[cleanedInput[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(session, cleanedInput)
		if err != nil {
			fmt.Println(err)
		}

	}
}
