package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseBlockBasic(t *testing.T) {
	input := `block-beta
columns 3
a["A"] b["B"] c["C"]
d["D"]:2 e["E"]`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Block {
		t.Fatalf("Kind = %v, want Block", g.Kind)
	}
	if g.BlockColumns != 3 {
		t.Fatalf("BlockColumns = %d, want 3", g.BlockColumns)
	}
	if len(g.Blocks) != 5 {
		t.Fatalf("Blocks = %d, want 5", len(g.Blocks))
	}
	if g.Blocks[3].Width != 2 {
		t.Errorf("d width = %d, want 2", g.Blocks[3].Width)
	}
}

func TestParseBlockEdges(t *testing.T) {
	input := `block-beta
a["Source"] b["Target"]
a --> b`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
	}
	e := out.Graph.Edges[0]
	if e.From != "a" || e.To != "b" {
		t.Errorf("Edge = %s -> %s", e.From, e.To)
	}
}

func TestParseBlockShapes(t *testing.T) {
	input := `block-beta
a["Square"] b("Rounded") c(("Circle")) d{"Diamond"}`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Blocks) != 4 {
		t.Fatalf("Blocks = %d, want 4", len(out.Graph.Blocks))
	}
	shapes := []ir.NodeShape{ir.Rectangle, ir.RoundRect, ir.Circle, ir.Diamond}
	for i, want := range shapes {
		if out.Graph.Blocks[i].Shape != want {
			t.Errorf("block[%d] shape = %v, want %v", i, out.Graph.Blocks[i].Shape, want)
		}
	}
}

func TestParseBlockEmpty(t *testing.T) {
	input := `block-beta`
	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Blocks) != 0 {
		t.Errorf("Blocks = %d, want 0", len(out.Graph.Blocks))
	}
}
