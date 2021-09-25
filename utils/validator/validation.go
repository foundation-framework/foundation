package validatorutil

import (
	"github.com/go-playground/validator/v10"
)

var Singleton = validator.New()

func Struct(value interface{}) error {
	return Singleton.Struct(value)
}

func Var(field interface{}, tag string) error {
	return Singleton.Var(field, tag)
}
