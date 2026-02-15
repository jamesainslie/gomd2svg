package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseMindmapBasic(t *testing.T) {
	input := `mindmap
    root((Central))
        A[Square]
        B(Rounded)
            C((Circle))`

	out, err := parseMindmap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Mindmap {
		t.Fatalf("Kind = %v, want Mindmap", g.Kind)
	}
	if g.MindmapRoot == nil {
		t.Fatal("MindmapRoot is nil")
	}
	if g.MindmapRoot.Label != "Central" {
		t.Errorf("root label = %q, want %q", g.MindmapRoot.Label, "Central")
	}
	if g.MindmapRoot.Shape != ir.MindmapCircle {
		t.Errorf("root shape = %v, want Circle", g.MindmapRoot.Shape)
	}
	if len(g.MindmapRoot.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(g.MindmapRoot.Children))
	}
	a := g.MindmapRoot.Children[0]
	if a.Label != "Square" || a.Shape != ir.MindmapSquare {
		t.Errorf("child A = %q/%v, want Square/square", a.Label, a.Shape)
	}
	b := g.MindmapRoot.Children[1]
	if len(b.Children) != 1 {
		t.Fatalf("B children = %d, want 1", len(b.Children))
	}
	if b.Children[0].Shape != ir.MindmapCircle {
		t.Errorf("C shape = %v, want Circle", b.Children[0].Shape)
	}
}

func TestParseMindmapAllShapes(t *testing.T) {
	input := `mindmap
    root
        Default
        [Square]
        (Rounded)
        ((Circle))
        ))Bang((
        )Cloud(
        {{Hexagon}}`

	out, err := parseMindmap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.MindmapRoot.Children) != 7 {
		t.Fatalf("children = %d, want 7", len(g.MindmapRoot.Children))
	}
	expected := []ir.MindmapShape{
		ir.MindmapShapeDefault, ir.MindmapSquare, ir.MindmapRounded, ir.MindmapCircle,
		ir.MindmapBang, ir.MindmapCloud, ir.MindmapHexagon,
	}
	for i, want := range expected {
		got := g.MindmapRoot.Children[i].Shape
		if got != want {
			t.Errorf("child[%d] shape = %v, want %v", i, got, want)
		}
	}
}

func TestParseMindmapIconAndClass(t *testing.T) {
	input := `mindmap
    root
        A::icon(fa fa-book)
        B:::urgent`

	out, err := parseMindmap(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.MindmapRoot.Children[0].Icon != "fa fa-book" {
		t.Errorf("icon = %q, want %q", g.MindmapRoot.Children[0].Icon, "fa fa-book")
	}
	if g.MindmapRoot.Children[1].Class != "urgent" {
		t.Errorf("class = %q, want %q", g.MindmapRoot.Children[1].Class, "urgent")
	}
}
