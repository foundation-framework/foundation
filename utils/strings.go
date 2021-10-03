package utils

import (
	"strings"
	"unicode"
)

func SnakeCase(value string) string {
	for i, letter := range value {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			continue
		}

		lower := string(unicode.ToLower(letter))
		if i != 0 {
			lower = "_" + lower
		}

		value = strings.ReplaceAll(value, string(letter), lower)
	}

	return value
}
