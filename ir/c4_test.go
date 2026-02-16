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
	g := NewGraph()
	g.Kind = C4
	g.C4SubKind = C4Container
	g.C4Elements = append(g.C4Elements, &C4Element{
		ID:    "user",
		Label: "User",
		Type:  C4Person,
	})
	g.C4Boundaries = append(g.C4Boundaries, &C4Boundary{
		ID:       "system",
		Label:    "My System",
		Type:     "Software System",
		Children: []string{"webapp"},
	})
	g.C4Rels = append(g.C4Rels, &C4Rel{
		From:  "user",
		To:    "webapp",
		Label: "Uses",
	})
	if g.C4SubKind != C4Container {
		t.Errorf("C4SubKind = %v", g.C4SubKind)
	}
	if len(g.C4Elements) != 1 {
		t.Errorf("C4Elements = %d", len(g.C4Elements))
	}
	if len(g.C4Boundaries) != 1 {
		t.Errorf("C4Boundaries = %d", len(g.C4Boundaries))
	}
	if len(g.C4Rels) != 1 {
		t.Errorf("C4Rels = %d", len(g.C4Rels))
	}
}
