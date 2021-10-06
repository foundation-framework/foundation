package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeCase(t *testing.T) {
	assert.Equal(t, "word", SnakeCase("Word"))
	assert.Equal(t, "word_word_word", SnakeCase("WordWordWord"))

	assert.Equal(t, "word_word", SnakeCase("wordWord"))
	assert.Equal(t, "word_word", SnakeCase("word Word"))

	assert.Equal(t, "word_word_word", SnakeCase("Word wordWord"))
}
