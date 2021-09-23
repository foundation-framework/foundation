package validatorutil

import (
	"github.com/go-playground/validator/v10"
)

var singleton = validator.New()

func Struct(value interface{}) error {
	return singleton.Struct(value)
}

func Var(field interface{}, tag string) error {
	return singleton.Var(field, tag)
}
