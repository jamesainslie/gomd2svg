package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestComputeLayoutSimple(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.LeftRight
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.EnsureNode("C", nil, nil)
	g.Edges = []*ir.Edge{edge("A", "B"), edge("B", "C")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Flowchart {
		t.Errorf("Kind = %v, want Flowchart", l.Kind)
	}
	if len(l.Nodes) != 3 {
		t.Errorf("Nodes = %d, want 3", len(l.Nodes))
	}
	if len(l.Edges) != 2 {
		t.Errorf("Edges = %d, want 2", len(l.Edges))
	}
	if l.Width <= 0 {
		t.Errorf("Width = %f, want > 0", l.Width)
	}
	if l.Height <= 0 {
		t.Errorf("Height = %f, want > 0", l.Height)
	}

	// In LR direction, nodes should be positioned left to right
	ax := l.Nodes["A"].X
	bx := l.Nodes["B"].X
	cx := l.Nodes["C"].X
	if ax >= bx || bx >= cx {
		t.Errorf("expected A.x < B.x < C.x, got %f, %f, %f", ax, bx, cx)
	}
}

func TestComputeLayoutTopDown(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.TopDown
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Edges = []*ir.Edge{edge("A", "B")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	// In TD direction, A should be above B (smaller Y)
	ay := l.Nodes["A"].Y
	by := l.Nodes["B"].Y
	if ay >= by {
		t.Errorf("expected A.y < B.y in TopDown, got A.y=%f B.y=%f", ay, by)
	}
}

func TestComputeLayoutEdgeRouting(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.LeftRight
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Edges = []*ir.Edge{edge("A", "B")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if len(l.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(l.Edges))
	}
	e := l.Edges[0]
	if len(e.Points) < 2 {
		t.Errorf("Edge points = %d, want >= 2", len(e.Points))
	}
	if e.From != "A" {
		t.Errorf("Edge.From = %q, want %q", e.From, "A")
	}
	if e.To != "B" {
		t.Errorf("Edge.To = %q, want %q", e.To, "B")
	}
}

func TestComputeLayoutDiamondShape(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.TopDown
	diamondShape := ir.Diamond
	g.EnsureNode("D", nil, &diamondShape)
	g.EnsureNode("A", nil, nil)
	g.Edges = []*ir.Edge{edge("A", "D")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	d := l.Nodes["D"]
	// Diamond nodes should be sized as a square (width == height)
	if d.Width != d.Height {
		t.Errorf("Diamond node Width=%f Height=%f, want square", d.Width, d.Height)
	}
}

func TestComputeLayoutEmptyGraph(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	g.Direction = ir.LeftRight

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d, want 0", len(l.Nodes))
	}
	if len(l.Edges) != 0 {
		t.Errorf("Edges = %d, want 0", len(l.Edges))
	}
}
