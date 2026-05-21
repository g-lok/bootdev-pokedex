package main

import (
	"strings"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	split := strings.Fields(lowercase)
	return split
}
