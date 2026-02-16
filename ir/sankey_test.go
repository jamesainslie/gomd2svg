package ir

import "testing"

func TestSankeyLink(t *testing.T) {
	link := &SankeyLink{Source: "A", Target: "B", Value: 100.5}
	if link.Source != "A" {
		t.Errorf("Source = %q, want A", link.Source)
	}
	if link.Target != "B" {
		t.Errorf("Target = %q, want B", link.Target)
	}
	if link.Value != 100.5 {
		t.Errorf("Value = %v, want 100.5", link.Value)
	}
}

func TestSankeyGraphFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Sankey
	graph.SankeyLinks = append(graph.SankeyLinks, &SankeyLink{
		Source: "Solar", Target: "Grid", Value: 59.9,
	})
	if len(graph.SankeyLinks) != 1 {
		t.Fatalf("SankeyLinks len = %d, want 1", len(graph.SankeyLinks))
	}
}
