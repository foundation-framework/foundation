package configutil

import (
	"go.uber.org/config"
)

// NewProviderFiles creates config.Provider from passed file paths
func NewProviderFiles(paths ...string) (config.Provider, error) {
	var options []config.YAMLOption

	for _, path := range paths {
		options = append(options, config.File(path))
	}

	provider, err := config.NewYAML(options...)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
