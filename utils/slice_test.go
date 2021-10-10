package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Has(t *testing.T) {
	assert.False(t, NewSlice(nil).Has(1))
	assert.False(t, NewSlice(map[string]interface{}{}).Has(1))

	assert.True(t, NewSlice([]int{1, 2, 3}).Has(2))
	assert.False(t, NewSlice([]int{1, 2, 3}).Has(5))

	assert.False(t, NewSlice([]int{}).Has(2))
}

func TestSlice_HasAll(t *testing.T) {
	assert.False(t, NewSlice(nil).HasAll(1))
	assert.False(t, NewSlice(map[string]interface{}{}).HasAll(1))

	assert.True(t, NewSlice([]int{1, 2, 3, 4}).HasAll(1, 2))
	assert.True(t, NewSlice([]int{1, 2, 3, 4}).HasAll())

	assert.False(t, NewSlice([]int{1, 2, 3, 4}).HasAll(1, 2, 5))
}

func TestSlice_HasAny(t *testing.T) {
	assert.False(t, NewSlice(nil).HasAny(1))
	assert.False(t, NewSlice(map[string]interface{}{}).HasAny(1))

	assert.True(t, NewSlice([]int{1, 2, 3, 4}).HasAny(1, 5, 6))
	assert.False(t, NewSlice([]int{1, 2, 3, 4}).HasAny())

	assert.False(t, NewSlice([]int{1, 2, 3, 4}).HasAny(5, 6, 7))
}
