package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"strings"
	"text/tabwriter"
)

const (
	BaseURL           = "https://pokeapi.co/api/v2"
	LocationAreasURL  = BaseURL + "/location-area/"
	PokemonURL        = BaseURL + "/pokemon/"
	FirstLocationPage = LocationAreasURL + "?offset=0&limit=20"
)

type Pagination[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type LocationAreas struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Name              string           `json:"name"`
	PokemonEncounters []PokemonVersion `json:"pokemon_encounters"`
}

type PokemonVersion struct {
	Pokemon Pokemon
	Version []string
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
}

func (p Pokemon) String() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "ID:\t%d\n", p.ID)
	fmt.Fprintf(w, "Name:\t%s\n", p.Name)
	fmt.Fprintf(w, "Base Experience:\t%d\n", p.BaseExperience)
	fmt.Fprintf(w, "Height:\t%d\n", p.Height)
	fmt.Fprintf(w, "Is Default:\t%t\n", p.IsDefault)
	fmt.Fprintf(w, "Order:\t%d\n", p.Order)
	fmt.Fprintf(w, "Weight:\t%d\n", p.Weight)

	w.Flush()
	return strings.TrimSuffix(b.String(), "\n")
}

func FetchPokemon(url string, cache *pokecache.Cache) (Pokemon, error) {
	var pokemonData Pokemon

	cachedPage, cachedPageExists := cache.Get(url)
	if cachedPageExists {
		if err := json.Unmarshal(cachedPage, &pokemonData); err != nil {
			return Pokemon{}, err
		}
		return pokemonData, nil
	}

	// No Cached value found so we hit the api
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}

	err = cache.Add(url, body)
	if err != nil {
		return Pokemon{}, err
	}

	if err := json.Unmarshal(body, &pokemonData); err != nil {
		return Pokemon{}, err
	}
	return pokemonData, nil
}

func FetchLocationArea(url string, cache *pokecache.Cache) (LocationArea, error) {
	var locationAreaData LocationArea

	cachedPage, cachedPageExists := cache.Get(url)
	if cachedPageExists {
		if err := json.Unmarshal(cachedPage, &locationAreaData); err != nil {
			return LocationArea{}, err
		}
		return locationAreaData, nil
	}

	// No Cached value found so we hit the api
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}
	if res.StatusCode > 299 {
		return LocationArea{}, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}

	err = cache.Add(url, body)
	if err != nil {
		return LocationArea{}, err
	}

	if err := json.Unmarshal(body, &locationAreaData); err != nil {
		return LocationArea{}, err
	}
	return locationAreaData, nil
}

func FetchLocationAreas(url string, cache *pokecache.Cache) (Pagination[LocationAreas], error) {
	var locationAreaData Pagination[LocationAreas]

	cachedPage, cachedPageExists := cache.Get(url)
	if cachedPageExists {
		if err := json.Unmarshal(cachedPage, &locationAreaData); err != nil {
			return Pagination[LocationAreas]{}, err
		}
		return locationAreaData, nil
	}

	// No Cached value found so we hit the api
	res, err := http.Get(url)
	if err != nil {
		return Pagination[LocationAreas]{}, err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pagination[LocationAreas]{}, err
	}
	if res.StatusCode > 299 {
		return Pagination[LocationAreas]{}, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}

	err = cache.Add(url, body)
	if err != nil {
		return Pagination[LocationAreas]{}, err
	}

	if err := json.Unmarshal(body, &locationAreaData); err != nil {
		return Pagination[LocationAreas]{}, err
	}
	return locationAreaData, nil
}
