package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestBlockGridLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block
	graph.BlockColumns = 3

	for _, id := range []string{"a", "b", "c", "d", "e"} {
		label := id
		graph.EnsureNode(id, &label, nil)
		graph.Blocks = append(graph.Blocks, &ir.BlockDef{ID: id, Label: id, Shape: ir.Rectangle, Width: 1})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeBlockLayout(graph, th, cfg)

	if lay.Kind != ir.Block {
		t.Fatalf("Kind = %v", lay.Kind)
	}
	bd, ok := lay.Diagram.(BlockData)
	if !ok {
		t.Fatal("Diagram is not BlockData")
	}
	if bd.Columns != 3 {
		t.Errorf("Columns = %d, want 3", bd.Columns)
	}
	if len(lay.Nodes) != 5 {
		t.Errorf("Nodes = %d, want 5", len(lay.Nodes))
	}
	ay := lay.Nodes["a"].Y
	by := lay.Nodes["b"].Y
	cy := lay.Nodes["c"].Y
	if ay != by || by != cy {
		t.Errorf("Row 1 Y mismatch: a=%v b=%v c=%v", ay, by, cy)
	}
	dy := lay.Nodes["d"].Y
	if dy == ay {
		t.Error("Row 2 should have different Y than row 1")
	}
}

func TestBlockSpanLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block
	graph.BlockColumns = 3

	aLabel := "A"
	bLabel := "B"
	graph.EnsureNode("a", &aLabel, nil)
	graph.EnsureNode("b", &bLabel, nil)
	graph.Blocks = append(graph.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Shape: ir.Rectangle, Width: 2},
		&ir.BlockDef{ID: "b", Label: "B", Shape: ir.Rectangle, Width: 1},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeBlockLayout(graph, th, cfg)

	if lay.Nodes["a"].Width <= lay.Nodes["b"].Width {
		t.Errorf("a width (%v) should be > b width (%v)", lay.Nodes["a"].Width, lay.Nodes["b"].Width)
	}
}

func TestBlockSugiyamaFallback(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block

	aLabel := "A"
	bLabel := "B"
	graph.EnsureNode("a", &aLabel, nil)
	graph.EnsureNode("b", &bLabel, nil)
	graph.Blocks = append(graph.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Width: 1},
		&ir.BlockDef{ID: "b", Label: "B", Width: 1},
	)
	graph.Edges = append(graph.Edges, &ir.Edge{From: "a", To: "b", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeBlockLayout(graph, th, cfg)

	if len(lay.Edges) != 1 {
		t.Errorf("Edges = %d, want 1", len(lay.Edges))
	}
}

func TestBlockLayoutEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block
	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeBlockLayout(graph, th, cfg)
	if len(lay.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(lay.Nodes))
	}
}
