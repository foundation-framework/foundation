package utils

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRange(t *testing.T) {
	rng := NewRange(100, 200)

	assert.Equal(t, rng.Min, 100)
	assert.Equal(t, rng.Max, 200)

	rng1 := NewRange(200, 100)

	assert.Equal(t, rng1.Min, 100)
	assert.Equal(t, rng1.Max, 200)
}

func TestNewRangeBounds(t *testing.T) {
	rng := NewRangeBounds([2]int{100, 200})

	assert.Equal(t, rng.Min, 100)
	assert.Equal(t, rng.Max, 200)

	rng1 := NewRangeBounds([2]int{200, 100})

	assert.Equal(t, rng1.Min, 100)
	assert.Equal(t, rng1.Max, 200)
}

func TestRange_In(t *testing.T) {
	rng := NewRange(100, 200)

	assert.False(t, rng.In(100))
	assert.False(t, rng.In(200))

	assert.True(t, rng.InBounds(150))
	assert.False(t, rng.InBounds(250))

	rng1 := NewRange(100, 100)

	assert.False(t, rng1.In(100))
}

func TestRange_InBounds(t *testing.T) {
	rng := NewRange(100, 200)

	assert.True(t, rng.InBounds(100))
	assert.True(t, rng.InBounds(200))

	assert.True(t, rng.InBounds(150))
	assert.False(t, rng.InBounds(250))

	rng1 := NewRange(100, 100)

	assert.True(t, rng1.InBounds(100))
}

func TestRange_Slice(t *testing.T) {
	rng := NewRange(100, 200)

	slice := rng.Slice()
	require.Len(t, slice, 101)

	counter := 100
	for _, elem := range slice {
		require.Equal(t, elem, counter)
		counter++
	}
}

func TestRange_Rand(t *testing.T) {
	rng := NewRange(100, 105)

	for i := 0; i < 100; i++ {
		assert.True(t, rng.InBounds(rng.Rand()))
	}
}

func TestRange_UnmarshalYAML_yaml_v2(t *testing.T) {
	type Data struct {
		Range Range `yaml:"range"`
	}

	text := "range: [100, 200]"

	d := &Data{}
	require.NoError(t, yaml.NewDecoder(bytes.NewBufferString(text)).Decode(d))

	assert.Equal(t, 100, d.Range.Min)
	assert.Equal(t, 200, d.Range.Max)

	text1 := "range: [100]"

	d1 := &Data{}
	require.Error(t, yaml.NewDecoder(bytes.NewBufferString(text1)).Decode(d1))
}
