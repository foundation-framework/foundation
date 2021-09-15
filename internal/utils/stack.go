package utils

import (
	"runtime/debug"
)

func Stack() []byte {
	return debug.Stack()
}

func StackString() string {
	return string(debug.Stack())
}
