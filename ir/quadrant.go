package ir

// QuadrantPoint represents a data point in a quadrant chart.
// X and Y are normalized values in the range [0, 1].
type QuadrantPoint struct {
	Label string
	X     float64
	Y     float64
}
