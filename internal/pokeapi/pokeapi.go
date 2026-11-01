package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

const (
	BaseUrl           = "https://pokeapi.co/api/v2"
	LocationAreasUrl  = BaseUrl + "/location-area/"
	FirstLocationPage = LocationAreasUrl + "?offset=0&limit=20"
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
	Name string `json:"name"`
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
