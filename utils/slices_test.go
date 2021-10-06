package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	assert.False(t, In(nil, 1))
	assert.False(t, In(map[string]interface{}{}, 1))

	assert.True(t, In([]int{1, 2, 3}, 2))
	assert.False(t, In([]int{1, 2, 3}, 5))

	assert.False(t, In([]int{}, 2))
}

func TestAll(t *testing.T) {
	assert.False(t, All(nil, 1))
	assert.False(t, All(map[string]interface{}{}, 1))

	assert.True(t, All([]int{1, 2, 3, 4}, 1, 2))
	assert.True(t, All([]int{1, 2, 3, 4}))

	assert.False(t, All([]int{1, 2, 3, 4}, 1, 2, 5))
}

func TestAny(t *testing.T) {
	assert.False(t, Any(nil, 1))
	assert.False(t, Any(map[string]interface{}{}, 1))

	assert.True(t, Any([]int{1, 2, 3, 4}, 1, 5, 6))
	assert.False(t, Any([]int{1, 2, 3, 4}))

	assert.False(t, Any([]int{1, 2, 3, 4}, 5, 6, 7))
}
