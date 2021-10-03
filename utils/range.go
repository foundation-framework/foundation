package utils

import "github.com/intale-llc/foundation/rand"

type Range struct {
	Left  int
	Right int
}

func NewRange(bounds [2]int) *Range {
	var (
		left  int
		right int
	)

	if bounds[0] <= bounds[1] {
		left = bounds[0]
		right = bounds[1]
	} else {
		left = bounds[1]
		right = bounds[0]
	}

	return &Range{
		Left:  left,
		Right: right,
	}
}

func (r *Range) Rand() int {
	return rand.Int(r.Left, r.Right)
}

func (r *Range) In(value int) bool {
	return r.Left < value && value < r.Right
}

func (r *Range) InBounds(value int) bool {
	return r.Left <= value && value <= r.Right
}
