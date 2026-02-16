package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
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
	graph := out.Graph
	if graph.Kind != ir.Mindmap {
		t.Fatalf("Kind = %v, want Mindmap", graph.Kind)
	}
	if graph.MindmapRoot == nil {
		t.Fatal("MindmapRoot is nil")
	}
	if graph.MindmapRoot.Label != "Central" {
		t.Errorf("root label = %q, want %q", graph.MindmapRoot.Label, "Central")
	}
	if graph.MindmapRoot.Shape != ir.MindmapCircle {
		t.Errorf("root shape = %v, want Circle", graph.MindmapRoot.Shape)
	}
	if len(graph.MindmapRoot.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(graph.MindmapRoot.Children))
	}
	a := graph.MindmapRoot.Children[0]
	if a.Label != "Square" || a.Shape != ir.MindmapSquare {
		t.Errorf("child A = %q/%v, want Square/square", a.Label, a.Shape)
	}
	b := graph.MindmapRoot.Children[1]
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
	graph := out.Graph
	if len(graph.MindmapRoot.Children) != 7 {
		t.Fatalf("children = %d, want 7", len(graph.MindmapRoot.Children))
	}
	expected := []ir.MindmapShape{
		ir.MindmapShapeDefault, ir.MindmapSquare, ir.MindmapRounded, ir.MindmapCircle,
		ir.MindmapBang, ir.MindmapCloud, ir.MindmapHexagon,
	}
	for i, want := range expected {
		got := graph.MindmapRoot.Children[i].Shape
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
	graph := out.Graph
	if graph.MindmapRoot.Children[0].Icon != "fa fa-book" {
		t.Errorf("icon = %q, want %q", graph.MindmapRoot.Children[0].Icon, "fa fa-book")
	}
	if graph.MindmapRoot.Children[1].Class != "urgent" {
		t.Errorf("class = %q, want %q", graph.MindmapRoot.Children[1].Class, "urgent")
	}
}
