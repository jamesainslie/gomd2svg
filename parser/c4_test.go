package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseC4Context(t *testing.T) {
	input := `C4Context
Person(user, "User", "A user of the system")
System(webapp, "Web App", "The main web application")
Rel(user, webapp, "Uses", "HTTPS")`

	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.C4 {
		t.Fatalf("Kind = %v, want C4", g.Kind)
	}
	if g.C4SubKind != ir.C4Context {
		t.Fatalf("C4SubKind = %v, want C4Context", g.C4SubKind)
	}
	if len(g.C4Elements) != 2 {
		t.Fatalf("Elements = %d, want 2", len(g.C4Elements))
	}
	if g.C4Elements[0].Type != ir.C4Person {
		t.Errorf("elem[0] type = %v, want Person", g.C4Elements[0].Type)
	}
	if g.C4Elements[0].Description != "A user of the system" {
		t.Errorf("elem[0] desc = %q", g.C4Elements[0].Description)
	}
	if len(g.C4Rels) != 1 {
		t.Fatalf("Rels = %d, want 1", len(g.C4Rels))
	}
	if g.C4Rels[0].Technology != "HTTPS" {
		t.Errorf("rel tech = %q", g.C4Rels[0].Technology)
	}
}

func TestParseC4Container(t *testing.T) {
	input := `C4Container
Person(user, "User", "End user")
Container_Boundary(system, "My System") {
Container(api, "API", "Go", "REST API")
ContainerDb(db, "Database", "PostgreSQL", "Stores data")
}
Rel(user, api, "Calls")
Rel(api, db, "Reads/Writes")`

	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.C4SubKind != ir.C4Container {
		t.Fatalf("C4SubKind = %v, want C4Container", g.C4SubKind)
	}
	if len(g.C4Boundaries) != 1 {
		t.Fatalf("Boundaries = %d, want 1", len(g.C4Boundaries))
	}
	b := g.C4Boundaries[0]
	if b.Label != "My System" {
		t.Errorf("boundary label = %q", b.Label)
	}
	if len(b.Children) != 2 {
		t.Fatalf("boundary children = %d, want 2", len(b.Children))
	}
	var api *ir.C4Element
	for _, e := range g.C4Elements {
		if e.ID == "api" {
			api = e
		}
	}
	if api == nil {
		t.Fatal("api element not found")
	}
	if api.Technology != "Go" {
		t.Errorf("api tech = %q, want Go", api.Technology)
	}
	if api.BoundaryID != "system" {
		t.Errorf("api boundary = %q, want system", api.BoundaryID)
	}
}

func TestParseC4Empty(t *testing.T) {
	input := `C4Context`
	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.C4Elements) != 0 {
		t.Errorf("Elements = %d, want 0", len(out.Graph.C4Elements))
	}
}

func TestParseC4Args(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{`user, "User", "Description"`, 3},
		{`api, "API", "Go", "REST API"`, 4},
		{`a, b`, 2},
	}
	for _, tt := range tests {
		args := parseC4Args(tt.input)
		if len(args) != tt.want {
			t.Errorf("parseC4Args(%q) = %d args, want %d", tt.input, len(args), tt.want)
		}
	}
}
