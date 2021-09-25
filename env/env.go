package env

import (
	"os"
	"strings"
)

func Read(name string) string {
	return os.Getenv(name)
}

func ReadSlice(name string) []string {
	return strings.Split(os.Getenv(name), ",")
}
