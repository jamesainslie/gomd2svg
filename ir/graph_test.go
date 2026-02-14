package ir

import "testing"

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g.Kind != Flowchart {
		t.Errorf("Kind = %v, want Flowchart", g.Kind)
	}
	if g.Direction != TopDown {
		t.Errorf("Direction = %v, want TopDown", g.Direction)
	}
	if g.Nodes == nil {
		t.Error("Nodes is nil")
	}
	if g.Edges != nil {
		t.Error("Edges should be nil (zero-value slice)")
	}
}

func TestEnsureNode(t *testing.T) {
	g := NewGraph()
	g.EnsureNode("A", nil, nil)
	if len(g.Nodes) != 1 {
		t.Fatalf("Nodes = %d, want 1", len(g.Nodes))
	}
	n := g.Nodes["A"]
	if n.ID != "A" {
		t.Errorf("ID = %q, want %q", n.ID, "A")
	}
	if n.Label != "A" {
		t.Errorf("Label = %q, want %q", n.Label, "A")
	}
	if n.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", n.Shape)
	}

	// Update with label and shape
	label := "Start"
	shape := Stadium
	g.EnsureNode("A", &label, &shape)
	n = g.Nodes["A"]
	if n.Label != "Start" {
		t.Errorf("Label = %q, want %q", n.Label, "Start")
	}
	if n.Shape != Stadium {
		t.Errorf("Shape = %v, want Stadium", n.Shape)
	}
	if len(g.Nodes) != 1 {
		t.Errorf("Nodes = %d, want 1 (should not duplicate)", len(g.Nodes))
	}
}

func TestEnsureNodeOrder(t *testing.T) {
	g := NewGraph()
	g.EnsureNode("C", nil, nil)
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	if g.NodeOrder["C"] != 0 {
		t.Errorf("C order = %d, want 0", g.NodeOrder["C"])
	}
	if g.NodeOrder["A"] != 1 {
		t.Errorf("A order = %d, want 1", g.NodeOrder["A"])
	}
	if g.NodeOrder["B"] != 2 {
		t.Errorf("B order = %d, want 2", g.NodeOrder["B"])
	}
	// Re-ensure does not change order
	g.EnsureNode("C", nil, nil)
	if g.NodeOrder["C"] != 0 {
		t.Errorf("C order = %d after re-ensure, want 0", g.NodeOrder["C"])
	}
}

func TestEdgeArrowheadValues(t *testing.T) {
	heads := []EdgeArrowhead{
		OpenTriangle,
		ClassDependency,
		ClosedTriangle,
		FilledDiamond,
		OpenDiamond,
		Lollipop,
	}
	seen := make(map[EdgeArrowhead]bool)
	for _, h := range heads {
		if seen[h] {
			t.Errorf("duplicate arrowhead value: %d", h)
		}
		seen[h] = true
	}
}
