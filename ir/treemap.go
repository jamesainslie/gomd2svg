package ir

// TreemapNode represents a node in the treemap hierarchy.
// Leaf nodes have Value > 0 and no Children.
// Section nodes have Children and their value is the sum of children.
type TreemapNode struct {
	Label    string
	Value    float64 // leaf value (0 for sections)
	Class    string  // CSS class from :::
	Children []*TreemapNode
}

// IsLeaf returns true if this node has no children.
func (n *TreemapNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// TotalValue returns the node's own value if it's a leaf,
// or the recursive sum of children's values if it's a section.
func (n *TreemapNode) TotalValue() float64 {
	if n.IsLeaf() {
		return n.Value
	}
	var sum float64
	for _, c := range n.Children {
		sum += c.TotalValue()
	}
	return sum
}
