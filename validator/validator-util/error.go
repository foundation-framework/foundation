package validatorutil

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
	"github.com/pkg/errors"
)

var (
	ErrUnknown = errors.New("unknown validation error")
)

func BeautyError(err error) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrUnknown
	}

	errText := "fields with incorrect format: "
	for _, err := range errs {
		errText += fmt.Sprintf("%s, ", stringy.New(err.Field()).SnakeCase().ToLower())
	}

	return errors.New(strings.TrimRight(errText, ", "))
}
