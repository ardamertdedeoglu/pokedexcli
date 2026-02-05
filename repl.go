package main

import (
	"strings"
)

func cleanInput(text string) []string {
	lowered_string := strings.ToLower(text)
	word_list := strings.Fields(lowered_string)
	return word_list
}