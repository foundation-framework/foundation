package utils

import "github.com/intale-llc/foundation/rand"

type IntRange struct {
	Left  int
	Right int
}

func NewIntRange(bounds [2]int) *IntRange {
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

	return &IntRange{
		Left:  left,
		Right: right,
	}
}

func (r *IntRange) Rand() int {
	return rand.Int(r.Left, r.Right)
}

func (r *IntRange) In(value int) bool {
	return r.Left < value && value < r.Right
}

func (r *IntRange) InWithBounds(value int) bool {
	return r.Left <= value && value <= r.Right
}
