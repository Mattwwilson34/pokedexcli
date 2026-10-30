package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

const (
	BaseUrl          = "https://pokeapi.co/api/v2"
	LocationAreasUrl = BaseUrl + "/location-area/?offset=0&limit=20"
)

type Pagination[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func FetchLocationAreas(url string, cache *pokecache.Cache) (Pagination[LocationArea], error) {
	var locationAreaData Pagination[LocationArea]

	cachedPage, cachedPageExists := cache.Get(url)
	if cachedPageExists {
		if err := json.Unmarshal(cachedPage, &locationAreaData); err != nil {
			return Pagination[LocationArea]{}, err
		}
		return locationAreaData, nil
	}

	// No Cached value found so we hit the api
	res, err := http.Get(url)
	if err != nil {
		return Pagination[LocationArea]{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pagination[LocationArea]{}, err
	}
	if res.StatusCode > 299 {
		return Pagination[LocationArea]{}, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}

	err = cache.Add(url, body)
	if err != nil {
		return Pagination[LocationArea]{}, err
	}

	if err := json.Unmarshal(body, &locationAreaData); err != nil {
		return Pagination[LocationArea]{}, err
	}
	return locationAreaData, nil
}
