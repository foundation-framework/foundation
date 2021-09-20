package config

import (
	"go.uber.org/config"
)

// Reader represents YAML configuration reader
type Reader struct {
	provider config.Provider
}

// Read reads configuration section and populates it into a struct using the "yaml:" tag
func (r *Reader) Read(key string, dst interface{}) error {
	return r.provider.Get(key).Populate(dst)
}

// NewReader creates new configuration reader
// Only YAML format supported
func NewReader(files ...string) (*Reader, error) {
	var options []config.YAMLOption

	for _, file := range files {
		options = append(options, config.File(file))
	}

	provider, err := config.NewYAML(options...)
	if err != nil {
		return nil, err
	}

	return &Reader{provider: provider}, nil
}

// NewReaderMust do the same as NewReader, but panics on error
func NewReaderMust(files ...string) *Reader {
	reader, err := NewReader(files...)
	if err != nil {
		panic("failed to create reader: " + err.Error())
	}

	return reader
}
