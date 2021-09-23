package configutil

import (
	"os"
	"strings"
)

func ReadEnv(name string) string {
	return os.Getenv(name)
}

func ReadEnvArray(name string) []string {
	return strings.Split(os.Getenv(name), ",")
}
