package pokecatch

import (
	"fmt"
	"pokedexcli/internal/pokeapi"
	"strings"
	"text/tabwriter"
)

type Pokedex struct {
	Entries map[string]PokedexEntry
}

type PokedexEntry struct {
	Name                  string
	TotalCaught           int
	BaseExperience        int
	CatchDifficultyFactor float32
	CatchRate             float32
	Stats                 []pokeapi.PokemonStat
	Types                 []pokeapi.PokemonType
}

func NewPokedex() *Pokedex {
	return &Pokedex{Entries: make(map[string]PokedexEntry)}
}

func (p PokedexEntry) String() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Name:\t%s\n", p.Name)
	fmt.Fprintf(w, "Caught:\t%d\n", p.TotalCaught)
	fmt.Fprintf(w, "Base Experience:\t%d\n", p.BaseExperience)
	fmt.Fprintf(w, "Catch Difficulty Factor:\t%f\n", p.CatchDifficultyFactor)
	fmt.Fprintf(w, "Catch Rate:\t%f\n", p.CatchRate)
	fmt.Fprintf(w, "Stats\n")
	for _, stat := range p.Stats {
		fmt.Fprintf(w, "    - %s:\t%d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Fprintf(w, "Types\n")
	for _, t := range p.Types {
		fmt.Fprintf(w, "    - Type:\t%s\n", t.Type.Name)
	}

	w.Flush()
	return strings.TrimSuffix(b.String(), "\n")
}
