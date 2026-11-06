package pokecatch

import "fmt"

type Pokedex struct {
	Entries map[string]PokedexEntry
}

type PokedexEntry struct {
	Name                  string
	TotalCaught           int
	BaseExperience        int
	CatchDifficultyFactor float32
	CatchRate             float32
}

func NewPokemonCatchLog() *Pokedex {
	return &Pokedex{Entries: make(map[string]PokedexEntry)}
}

func (p *PokedexEntry) PrintLogEntry() {
	fmt.Println("____POKEDEX____")
	fmt.Println("Name:", p.Name)
	fmt.Println("Caught: ", p.TotalCaught)
	fmt.Println("Base Experience: ", p.BaseExperience)
	fmt.Println("Catch Difficulty Factor:", p.CatchDifficultyFactor)
	fmt.Println("Catch Rate: ", p.CatchRate)
	fmt.Println()
}
