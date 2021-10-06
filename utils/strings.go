package utils

import (
	"strings"
	"unicode"
)

func SnakeCase(value string) string {
	words := strings.Split(value, " ")

	for i := range words {
		words[i] = snakeCaseWord(words[i])
	}

	return strings.Join(words, "_")
}

func snakeCaseWord(word string) string {
	for i, letter := range word {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			continue
		}

		lower := string(unicode.ToLower(letter))
		if i != 0 {
			lower = "_" + lower
		}

		word = strings.Replace(word, string(letter), lower, 1)
	}

	return word
}
