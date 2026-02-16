package ir

import "testing"

func TestNewGraph(t *testing.T) {
	graph := NewGraph()
	if graph.Kind != Flowchart {
		t.Errorf("Kind = %v, want Flowchart", graph.Kind)
	}
	if graph.Direction != TopDown {
		t.Errorf("Direction = %v, want TopDown", graph.Direction)
	}
	if graph.Nodes == nil {
		t.Error("Nodes is nil")
	}
	if graph.Edges != nil {
		t.Error("Edges should be nil (zero-value slice)")
	}
}

func TestEnsureNode(t *testing.T) {
	graph := NewGraph()
	graph.EnsureNode("A", nil, nil)
	if len(graph.Nodes) != 1 {
		t.Fatalf("Nodes = %d, want 1", len(graph.Nodes))
	}
	node := graph.Nodes["A"]
	if node.ID != "A" {
		t.Errorf("ID = %q, want %q", node.ID, "A")
	}
	if node.Label != "A" {
		t.Errorf("Label = %q, want %q", node.Label, "A")
	}
	if node.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", node.Shape)
	}

	// Update with label and shape
	label := "Start"
	shape := Stadium
	graph.EnsureNode("A", &label, &shape)
	node = graph.Nodes["A"]
	if node.Label != "Start" {
		t.Errorf("Label = %q, want %q", node.Label, "Start")
	}
	if node.Shape != Stadium {
		t.Errorf("Shape = %v, want Stadium", node.Shape)
	}
	if len(graph.Nodes) != 1 {
		t.Errorf("Nodes = %d, want 1 (should not duplicate)", len(graph.Nodes))
	}
}

func TestEnsureNodeOrder(t *testing.T) {
	graph := NewGraph()
	graph.EnsureNode("C", nil, nil)
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	if graph.NodeOrder["C"] != 0 {
		t.Errorf("C order = %d, want 0", graph.NodeOrder["C"])
	}
	if graph.NodeOrder["A"] != 1 {
		t.Errorf("A order = %d, want 1", graph.NodeOrder["A"])
	}
	if graph.NodeOrder["B"] != 2 {
		t.Errorf("B order = %d, want 2", graph.NodeOrder["B"])
	}
	// Re-ensure does not change order
	graph.EnsureNode("C", nil, nil)
	if graph.NodeOrder["C"] != 0 {
		t.Errorf("C order = %d after re-ensure, want 0", graph.NodeOrder["C"])
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
