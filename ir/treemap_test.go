package ir

import "testing"

func TestTreemapNode(t *testing.T) {
	root := &TreemapNode{
		Label: "Root",
		Children: []*TreemapNode{
			{Label: "Leaf A", Value: 30},
			{Label: "Section B", Children: []*TreemapNode{
				{Label: "Leaf C", Value: 20},
			}},
		},
	}
	if root.Label != "Root" {
		t.Errorf("Label = %q, want %q", root.Label, "Root")
	}
	if len(root.Children) != 2 {
		t.Errorf("children len = %d, want 2", len(root.Children))
	}
	if root.Children[0].Value != 30 {
		t.Errorf("leaf value = %v, want 30", root.Children[0].Value)
	}
}

func TestTreemapNodeIsLeaf(t *testing.T) {
	leaf := &TreemapNode{Label: "Leaf", Value: 10}
	section := &TreemapNode{Label: "Sec", Children: []*TreemapNode{leaf}}
	if !leaf.IsLeaf() {
		t.Error("expected leaf")
	}
	if section.IsLeaf() {
		t.Error("expected non-leaf")
	}
}

func TestTreemapTotalValue(t *testing.T) {
	root := &TreemapNode{
		Label: "Root",
		Children: []*TreemapNode{
			{Label: "A", Value: 30},
			{Label: "B", Children: []*TreemapNode{
				{Label: "C", Value: 20},
				{Label: "D", Value: 10},
			}},
		},
	}
	if got := root.TotalValue(); got != 60 {
		t.Errorf("TotalValue = %v, want 60", got)
	}
}

func TestTreemapGraphField(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Treemap
	graph.TreemapRoot = &TreemapNode{Label: "Root"}
	if graph.TreemapRoot == nil {
		t.Error("TreemapRoot should not be nil")
	}
}
