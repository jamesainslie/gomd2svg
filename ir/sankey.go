package ir

// SankeyLink represents a flow from source to target with a value.
type SankeyLink struct {
	Source string
	Target string
	Value  float64
}
