package ir

import "testing"

func TestC4Kind(t *testing.T) {
	tests := []struct {
		kind C4Kind
		str  string
	}{
		{C4Context, "C4Context"},
		{C4Container, "C4Container"},
		{C4Component, "C4Component"},
		{C4Dynamic, "C4Dynamic"},
		{C4Deployment, "C4Deployment"},
	}
	for _, tt := range tests {
		if tt.kind.String() != tt.str {
			t.Errorf("C4Kind(%d).String() = %q, want %q", tt.kind, tt.kind.String(), tt.str)
		}
	}
}

func TestC4ElementType(t *testing.T) {
	if C4Person.String() != "Person" {
		t.Errorf("C4Person = %q", C4Person.String())
	}
	if C4ContainerPlain.String() != "Container" {
		t.Errorf("C4ContainerPlain = %q", C4ContainerPlain.String())
	}
	if C4ExternalSystem.String() != "System_Ext" {
		t.Errorf("C4ExternalSystem = %q", C4ExternalSystem.String())
	}
}

func TestC4ElementTypePredicates(t *testing.T) {
	if !C4ExternalSystem.IsExternal() {
		t.Error("C4ExternalSystem should be external")
	}
	if C4System.IsExternal() {
		t.Error("C4System should not be external")
	}
	if !C4Person.IsPerson() {
		t.Error("C4Person should be person")
	}
	if !C4SystemDb.IsDatabase() {
		t.Error("C4SystemDb should be database")
	}
	if !C4ContainerQueue.IsQueue() {
		t.Error("C4ContainerQueue should be queue")
	}
}

func TestC4GraphFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = C4
	graph.C4SubKind = C4Container
	graph.C4Elements = append(graph.C4Elements, &C4Element{
		ID:    "user",
		Label: "User",
		Type:  C4Person,
	})
	graph.C4Boundaries = append(graph.C4Boundaries, &C4Boundary{
		ID:       "system",
		Label:    "My System",
		Type:     "Software System",
		Children: []string{"webapp"},
	})
	graph.C4Rels = append(graph.C4Rels, &C4Rel{
		From:  "user",
		To:    "webapp",
		Label: "Uses",
	})
	if graph.C4SubKind != C4Container {
		t.Errorf("C4SubKind = %v", graph.C4SubKind)
	}
	if len(graph.C4Elements) != 1 {
		t.Errorf("C4Elements = %d", len(graph.C4Elements))
	}
	if len(graph.C4Boundaries) != 1 {
		t.Errorf("C4Boundaries = %d", len(graph.C4Boundaries))
	}
	if len(graph.C4Rels) != 1 {
		t.Errorf("C4Rels = %d", len(graph.C4Rels))
	}
}
