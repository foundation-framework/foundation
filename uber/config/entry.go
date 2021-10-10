package configext

import (
	"go.uber.org/config"
	"go.uber.org/multierr"

	"github.com/intale-llc/foundation/errors"
)

type Entry interface {
	// Name returns config entry name
	Name() string

	// Process starts entry processing
	// Any check or initialization logic should be there
	Process() error
}

// Entries reads & processes passed entries with config.Provider
//
// Error returned from function can be
// formatted by '%+v' for pretty output
func Entries(provider config.Provider, entries ...Entry) error {
	var entryErrs error
	for _, entry := range entries {
		value := provider.Get(entry.Name())

		if err := value.Populate(entry); err != nil {
			multierr.AppendInto(
				&entryErrs,
				errors.Wrapf(err, "failed to read \"%s\" entry", entry.Name()),
			)
			continue
		}

		multierr.AppendInto(
			&entryErrs,
			errors.Wrapf(entry.Process(), "failed to process \"%s\" entry", entry.Name()),
		)
	}

	return entryErrs
}
