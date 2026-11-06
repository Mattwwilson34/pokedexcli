package main

import (
	"pokedexcli/internal/pokecache"
	"pokedexcli/internal/pokecatch"
	"strings"
)

type session struct {
	userLocation    userLocation
	pokemonCatchLog *pokecatch.Pokedex
	cache           *pokecache.Cache
}

type userLocation struct {
	previous string
	current  string
	next     string
}

func cleanInput(text string) []string {
	cleanedText := strings.ToLower(text)
	cleanedText = strings.TrimSpace(cleanedText)
	cleanedText = strings.ReplaceAll(cleanedText, ",", " ")
	return strings.Split(cleanedText, " ")
}
