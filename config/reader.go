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

// MustRead do the same as Read, but panics on error
func (r *Reader) MustRead(key string, dst interface{}) {
	if err := r.provider.Get(key).Populate(dst); err != nil {
		panic("failed to read config: " + err.Error())
	}
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
