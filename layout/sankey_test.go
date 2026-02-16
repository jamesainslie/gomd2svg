package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestSankeyLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey
	graph.SankeyLinks = []*ir.SankeyLink{
		{Source: "A", Target: "X", Value: 100},
		{Source: "A", Target: "Y", Value: 200},
		{Source: "B", Target: "X", Value: 150},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", lay.Width, lay.Height)
	}

	sd, ok := lay.Diagram.(SankeyData)
	if !ok {
		t.Fatal("Diagram is not SankeyData")
	}
	if len(sd.Nodes) < 4 {
		t.Errorf("nodes = %d, want >= 4", len(sd.Nodes))
	}
	if len(sd.Links) != 3 {
		t.Errorf("links = %d, want 3", len(sd.Links))
	}
	// Source nodes (A,B) should be in column 0, targets (X,Y) in column 1.
	// A should have X < X's X.
	if sd.Nodes[0].X >= sd.Nodes[2].X {
		t.Error("source A should be left of target X")
	}
}

func TestSankeyLayoutEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", lay.Width, lay.Height)
	}
}

func TestSankeyCollectNodes(t *testing.T) {
	links := []*ir.SankeyLink{
		{Source: "A", Target: "B", Value: 10},
		{Source: "A", Target: "C", Value: 20},
		{Source: "B", Target: "D", Value: 5},
	}

	names, index := sankeyCollectNodes(links)

	if len(names) != 4 {
		t.Fatalf("names = %d, want 4", len(names))
	}
	// First appearance order: A, B, C, D.
	want := []string{"A", "B", "C", "D"}
	for i, w := range want {
		if names[i] != w {
			t.Errorf("names[%d] = %q, want %q", i, names[i], w)
		}
		if index[w] != i {
			t.Errorf("index[%q] = %d, want %d", w, index[w], i)
		}
	}
}

func TestSankeyAssignColumns(t *testing.T) {
	links := []*ir.SankeyLink{
		{Source: "A", Target: "B", Value: 10},
		{Source: "B", Target: "C", Value: 5},
		{Source: "A", Target: "C", Value: 20},
	}

	names, index := sankeyCollectNodes(links)
	columns := sankeyAssignColumns(names, index, links)

	// A -> col 0, B -> col 1, C -> col 2 (longest path A->B->C).
	wantCols := map[string]int{"A": 0, "B": 1, "C": 2}
	for name, wantCol := range wantCols {
		got := columns[index[name]]
		if got != wantCol {
			t.Errorf("column[%q] = %d, want %d", name, got, wantCol)
		}
	}
}

func TestSankeyComputeFlow(t *testing.T) {
	links := []*ir.SankeyLink{
		{Source: "A", Target: "X", Value: 100},
		{Source: "A", Target: "Y", Value: 200},
		{Source: "B", Target: "X", Value: 150},
	}

	names, index := sankeyCollectNodes(links)
	flow := sankeyComputeFlow(names, index, links)

	// A: outflow=300, inflow=0 -> 300.
	if flow[index["A"]] != 300 {
		t.Errorf("flow[A] = %v, want 300", flow[index["A"]])
	}
	// B: outflow=150, inflow=0 -> 150.
	if flow[index["B"]] != 150 {
		t.Errorf("flow[B] = %v, want 150", flow[index["B"]])
	}
	// X: outflow=0, inflow=250 -> 250.
	if flow[index["X"]] != 250 {
		t.Errorf("flow[X] = %v, want 250", flow[index["X"]])
	}
	// Y: outflow=0, inflow=200 -> 200.
	if flow[index["Y"]] != 200 {
		t.Errorf("flow[Y] = %v, want 200", flow[index["Y"]])
	}
}

func TestSankeyLinkPositions(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey
	graph.SankeyLinks = []*ir.SankeyLink{
		{Source: "A", Target: "X", Value: 100},
		{Source: "A", Target: "Y", Value: 200},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeSankeyLayout(graph, th, cfg)

	sd, ok := lay.Diagram.(SankeyData)
	if !ok {
		t.Fatal("expected SankeyData")
	}

	// Links from A should stack: second link's SourceY > first link's SourceY.
	if sd.Links[1].SourceY <= sd.Links[0].SourceY {
		t.Error("second link should be stacked below first at source")
	}

	// Each link width should be > 0.
	for i, link := range sd.Links {
		if link.Width <= 0 {
			t.Errorf("link[%d].Width = %v, want > 0", i, link.Width)
		}
	}
}
