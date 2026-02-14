package ir

import "testing"

func TestStateAnnotationString(t *testing.T) {
	tests := []struct {
		ann  StateAnnotation
		want string
	}{
		{StateChoice, "choice"},
		{StateFork, "fork"},
		{StateJoin, "join"},
		{StateAnnotation(99), ""},
	}
	for _, tt := range tests {
		got := tt.ann.String()
		if got != tt.want {
			t.Errorf("StateAnnotation(%d).String() = %q, want %q", int(tt.ann), got, tt.want)
		}
	}
}

func TestCompositeStateWithInner(t *testing.T) {
	inner := NewGraph()
	inner.EnsureNode("idle", nil, nil)
	inner.EnsureNode("active", nil, nil)
	inner.Edges = append(inner.Edges, &Edge{
		From:     "idle",
		To:       "active",
		Directed: true,
	})

	cs := &CompositeState{
		ID:    "Moving",
		Label: "Moving",
		Inner: inner,
	}

	if cs.ID != "Moving" {
		t.Errorf("ID = %q, want %q", cs.ID, "Moving")
	}
	if cs.Label != "Moving" {
		t.Errorf("Label = %q, want %q", cs.Label, "Moving")
	}
	if cs.Inner == nil {
		t.Fatal("Inner graph is nil")
	}
	if len(cs.Inner.Nodes) != 2 {
		t.Errorf("Inner.Nodes = %d, want 2", len(cs.Inner.Nodes))
	}
	if len(cs.Inner.Edges) != 1 {
		t.Errorf("Inner.Edges = %d, want 1", len(cs.Inner.Edges))
	}
	if cs.Regions != nil {
		t.Error("Regions should be nil when not set")
	}
	if cs.Direction != nil {
		t.Error("Direction should be nil when not set")
	}
}

func TestCompositeStateWithRegions(t *testing.T) {
	region1 := NewGraph()
	region1.EnsureNode("r1_idle", nil, nil)

	region2 := NewGraph()
	region2.EnsureNode("r2_idle", nil, nil)

	dir := LeftRight
	cs := &CompositeState{
		ID:        "Concurrent",
		Label:     "Concurrent State",
		Regions:   []*Graph{region1, region2},
		Direction: &dir,
	}

	if len(cs.Regions) != 2 {
		t.Errorf("Regions = %d, want 2", len(cs.Regions))
	}
	if cs.Direction == nil || *cs.Direction != LeftRight {
		t.Errorf("Direction = %v, want LeftRight", cs.Direction)
	}
}

func TestNewGraphInitializesStateMaps(t *testing.T) {
	g := NewGraph()
	if g.CompositeStates == nil {
		t.Error("CompositeStates is nil")
	}
	if g.StateDescriptions == nil {
		t.Error("StateDescriptions is nil")
	}
	if g.StateAnnotations == nil {
		t.Error("StateAnnotations is nil")
	}
}
