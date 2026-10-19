# pokedexcli

A REPL-based Pokedex for the terminal, written in Go. It reads commands from
stdin and fetches data from [PokeAPI](https://pokeapi.co/).

## Build and run

```bash
go build -o pokedexcli .
./pokedexcli
```

Or run it directly:

```bash
go run .
```

## Commands

| Command | Description |
| ------- | ----------- |
| `help`  | Display the help message |
| `map`   | Display the next page of location areas |
| `mapb`  | Display the previous page of location areas |
| `exit`  | Close the Pokedex |

## Tests

```bash
go test ./...
```
