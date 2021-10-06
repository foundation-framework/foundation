package utils

import "github.com/intale-llc/foundation/rand"

type Range struct {
	Min, Max int
}

func NewRange(min, max int) *Range {
	if min > max {
		min, max = max, min
	}

	return &Range{
		Min: min,
		Max: max,
	}
}

func NewRangeBounds(bounds [2]int) *Range {
	return NewRange(bounds[0], bounds[1])
}

func (r *Range) Rand() int {
	return rand.Int(r.Min, r.Max)
}

func (r *Range) In(value int) bool {
	return r.Min < value && value < r.Max
}

func (r *Range) InBounds(value int) bool {
	return r.Min <= value && value <= r.Max
}

func (r *Range) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var bounds [2]int

	if err := unmarshal(&bounds); err != nil {
		return err
	}

	*r = *NewRangeBounds(bounds)
	return nil
}
