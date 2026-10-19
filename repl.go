package main

import "strings"

func cleanInput(text string) []string {
	cleanedText := strings.ToLower(text)
	cleanedText = strings.TrimSpace(cleanedText)
	cleanedText = strings.ReplaceAll(cleanedText, ",", " ")
	return strings.Split(cleanedText, " ")
}
