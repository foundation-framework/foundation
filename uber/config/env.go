package configext

import (
	"os"
	"strings"

	"go.uber.org/config"
)

type Env struct {
	content string
}

func NewEnv(name string) *Env {
	return &Env{
		content: os.Getenv(name),
	}
}

func (s *Env) Slice(sep string) []string {
	if sep == "" {
		sep = ","
	}

	return strings.Split(s.content, sep)
}

func (s *Env) String() string {
	return s.content
}

func (s *Env) Entries(entries ...Entry) error {
	var provider config.Provider

	var options []config.YAMLOption

	files := s.Slice("")
	for _, path := range files {
		options = append(options, config.File(path))
	}

	provider, err := config.NewYAML(options...)
	if err != nil {
		return err
	}

	return Entries(provider, entries...)
}
