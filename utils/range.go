package utils

type Range struct {
	left  int
	right int
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
		left:  left,
		right: right,
	}
}

func (r *Range) In(value int) bool {
	return r.left < value && value < r.right
}

func (r *Range) InWithBounds(value int) bool {
	return r.left <= value && value <= r.right
}
