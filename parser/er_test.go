package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseERSimple(t *testing.T) {
	input := `erDiagram
CUSTOMER ||--o{ ORDER : places`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Er {
		t.Errorf("Kind = %v, want Er", graph.Kind)
	}
	if len(graph.Edges) != 1 {
		t.Errorf("Edges = %d, want 1", len(graph.Edges))
	}
	if _, ok := graph.Nodes["CUSTOMER"]; !ok {
		t.Error("expected node CUSTOMER")
	}
	if _, ok := graph.Nodes["ORDER"]; !ok {
		t.Error("expected node ORDER")
	}
}

func TestParseERAttributes(t *testing.T) {
	input := `erDiagram
PRODUCT {
  int id PK
  string name UK
  float price "retail price"
}`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	graph := out.Graph
	ent, ok := graph.Entities["PRODUCT"]
	if !ok {
		t.Fatalf("Entities missing PRODUCT")
	}
	if len(ent.Attributes) != 3 {
		t.Fatalf("Attributes = %d, want 3", len(ent.Attributes))
	}
	// Check keys
	if len(ent.Attributes[0].Keys) != 1 || ent.Attributes[0].Keys[0] != ir.KeyPrimary {
		t.Errorf("attr[0] keys = %v, want [KeyPrimary]", ent.Attributes[0].Keys)
	}
	if len(ent.Attributes[1].Keys) != 1 || ent.Attributes[1].Keys[0] != ir.KeyUnique {
		t.Errorf("attr[1] keys = %v, want [KeyUnique]", ent.Attributes[1].Keys)
	}
	if ent.Attributes[2].Comment != "retail price" {
		t.Errorf("attr[2] comment = %q, want %q", ent.Attributes[2].Comment, "retail price")
	}
}

func TestParseERCardinality(t *testing.T) {
	tests := []struct {
		name      string
		card      string
		wantLeft  ir.EdgeDecoration
		wantRight ir.EdgeDecoration
	}{
		{"one-to-zero-many", "||--o{", ir.DecCrowsFootOne, ir.DecCrowsFootZeroMany},
		{"one-to-one", "||--||", ir.DecCrowsFootOne, ir.DecCrowsFootOne},
		{"zero-one-to-zero-many", "|o--o{", ir.DecCrowsFootZeroOne, ir.DecCrowsFootZeroMany},
		{"many-to-many", "}|--|{", ir.DecCrowsFootMany, ir.DecCrowsFootMany},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "erDiagram\nA " + tt.card + " B : rel"
			out, err := parseER(input)
			if err != nil {
				t.Fatalf("parseER() error: %v", err)
			}
			if len(out.Graph.Edges) != 1 {
				t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
			}
			edge := out.Graph.Edges[0]
			if edge.StartDecoration == nil {
				t.Fatal("StartDecoration is nil")
			}
			if *edge.StartDecoration != tt.wantLeft {
				t.Errorf("StartDecoration = %v, want %v", *edge.StartDecoration, tt.wantLeft)
			}
			if edge.EndDecoration == nil {
				t.Fatal("EndDecoration is nil")
			}
			if *edge.EndDecoration != tt.wantRight {
				t.Errorf("EndDecoration = %v, want %v", *edge.EndDecoration, tt.wantRight)
			}
		})
	}
}

func TestParseERLineStyle(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		style ir.EdgeStyle
	}{
		{"solid", "A ||--o{ B : rel", ir.Solid},
		{"dotted", "A ||..o{ B : rel", ir.Dotted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "erDiagram\n" + tt.line
			out, err := parseER(input)
			if err != nil {
				t.Fatalf("parseER() error: %v", err)
			}
			if len(out.Graph.Edges) != 1 {
				t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
			}
			if out.Graph.Edges[0].Style != tt.style {
				t.Errorf("Style = %v, want %v", out.Graph.Edges[0].Style, tt.style)
			}
		})
	}
}

func TestParseERLabel(t *testing.T) {
	input := `erDiagram
CUSTOMER ||--o{ ORDER : places`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
	}
	e := out.Graph.Edges[0]
	if e.Label == nil {
		t.Fatal("Label is nil")
	}
	if *e.Label != "places" {
		t.Errorf("Label = %q, want %q", *e.Label, "places")
	}
}

func TestParseERAlias(t *testing.T) {
	input := `erDiagram
p["Person"] {
  string firstName
}`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	ent, ok := out.Graph.Entities["p"]
	if !ok {
		t.Fatalf("Entities missing p")
	}
	if ent.Label != "Person" {
		t.Errorf("Label = %q, want %q", ent.Label, "Person")
	}
	if ent.ID != "p" {
		t.Errorf("ID = %q, want %q", ent.ID, "p")
	}
}

func TestParseERCompositeKey(t *testing.T) {
	input := `erDiagram
ENROLLMENT {
  int id PK,FK
}`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	ent, ok := out.Graph.Entities["ENROLLMENT"]
	if !ok {
		t.Fatalf("Entities missing ENROLLMENT")
	}
	if len(ent.Attributes) != 1 {
		t.Fatalf("Attributes = %d, want 1", len(ent.Attributes))
	}
	attr := ent.Attributes[0]
	if len(attr.Keys) != 2 {
		t.Fatalf("Keys = %d, want 2", len(attr.Keys))
	}
	if attr.Keys[0] != ir.KeyPrimary {
		t.Errorf("Keys[0] = %v, want KeyPrimary", attr.Keys[0])
	}
	if attr.Keys[1] != ir.KeyForeign {
		t.Errorf("Keys[1] = %v, want KeyForeign", attr.Keys[1])
	}
}

func TestParseERDirection(t *testing.T) {
	input := `erDiagram
direction LR`
	out, err := parseER(input)
	if err != nil {
		t.Fatalf("parseER() error: %v", err)
	}
	if out.Graph.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", out.Graph.Direction)
	}
}
