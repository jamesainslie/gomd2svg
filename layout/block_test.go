package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestBlockGridLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	g.BlockColumns = 3

	for _, id := range []string{"a", "b", "c", "d", "e"} {
		label := id
		g.EnsureNode(id, &label, nil)
		g.Blocks = append(g.Blocks, &ir.BlockDef{ID: id, Label: id, Shape: ir.Rectangle, Width: 1})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeBlockLayout(g, th, cfg)

	if l.Kind != ir.Block {
		t.Fatalf("Kind = %v", l.Kind)
	}
	bd, ok := l.Diagram.(BlockData)
	if !ok {
		t.Fatal("Diagram is not BlockData")
	}
	if bd.Columns != 3 {
		t.Errorf("Columns = %d, want 3", bd.Columns)
	}
	if len(l.Nodes) != 5 {
		t.Errorf("Nodes = %d, want 5", len(l.Nodes))
	}
	ay := l.Nodes["a"].Y
	by := l.Nodes["b"].Y
	cy := l.Nodes["c"].Y
	if ay != by || by != cy {
		t.Errorf("Row 1 Y mismatch: a=%v b=%v c=%v", ay, by, cy)
	}
	dy := l.Nodes["d"].Y
	if dy == ay {
		t.Error("Row 2 should have different Y than row 1")
	}
}

func TestBlockSpanLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	g.BlockColumns = 3

	aLabel := "A"
	bLabel := "B"
	g.EnsureNode("a", &aLabel, nil)
	g.EnsureNode("b", &bLabel, nil)
	g.Blocks = append(g.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Shape: ir.Rectangle, Width: 2},
		&ir.BlockDef{ID: "b", Label: "B", Shape: ir.Rectangle, Width: 1},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeBlockLayout(g, th, cfg)

	if l.Nodes["a"].Width <= l.Nodes["b"].Width {
		t.Errorf("a width (%v) should be > b width (%v)", l.Nodes["a"].Width, l.Nodes["b"].Width)
	}
}

func TestBlockSugiyamaFallback(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block

	aLabel := "A"
	bLabel := "B"
	g.EnsureNode("a", &aLabel, nil)
	g.EnsureNode("b", &bLabel, nil)
	g.Blocks = append(g.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Width: 1},
		&ir.BlockDef{ID: "b", Label: "B", Width: 1},
	)
	g.Edges = append(g.Edges, &ir.Edge{From: "a", To: "b", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeBlockLayout(g, th, cfg)

	if len(l.Edges) != 1 {
		t.Errorf("Edges = %d, want 1", len(l.Edges))
	}
}

func TestBlockLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeBlockLayout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
