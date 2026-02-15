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

func TestComputeLayoutNonNegativeCoordinates(t *testing.T) {
	// All Sugiyama-based layouts must produce non-negative coordinates so that
	// the SVG viewBox "0 0 W H" doesn't clip content.
	tests := []struct {
		name      string
		direction ir.Direction
		nodeCount int
	}{
		{"LR-2", ir.LeftRight, 2},
		{"TD-2", ir.TopDown, 2},
		{"LR-5", ir.LeftRight, 5},
		{"TD-5", ir.TopDown, 5},
		{"BT-3", ir.BottomTop, 3},
		{"RL-3", ir.RightLeft, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := ir.NewGraph()
			g.Kind = ir.Flowchart
			g.Direction = tt.direction

			// Create a chain: A->B->C->...
			ids := make([]string, tt.nodeCount)
			for i := range ids {
				ids[i] = string(rune('A' + i))
				g.EnsureNode(ids[i], nil, nil)
			}
			for i := 0; i < len(ids)-1; i++ {
				g.Edges = append(g.Edges, edge(ids[i], ids[i+1]))
			}

			th := theme.Modern()
			cfg := config.DefaultLayout()
			l := ComputeLayout(g, th, cfg)

			// Every node must have non-negative bounding box.
			for id, n := range l.Nodes {
				left := n.X - n.Width/2
				top := n.Y - n.Height/2
				if left < 0 {
					t.Errorf("node %s left edge = %f, want >= 0", id, left)
				}
				if top < 0 {
					t.Errorf("node %s top edge = %f, want >= 0", id, top)
				}
			}

			// Every edge point and label anchor must be non-negative.
			for i, e := range l.Edges {
				for j, pt := range e.Points {
					if pt[0] < 0 {
						t.Errorf("edge %d point %d X = %f, want >= 0", i, j, pt[0])
					}
					if pt[1] < 0 {
						t.Errorf("edge %d point %d Y = %f, want >= 0", i, j, pt[1])
					}
				}
				if e.LabelAnchor[0] < 0 {
					t.Errorf("edge %d LabelAnchor X = %f, want >= 0", i, e.LabelAnchor[0])
				}
				if e.LabelAnchor[1] < 0 {
					t.Errorf("edge %d LabelAnchor Y = %f, want >= 0", i, e.LabelAnchor[1])
				}
			}
		})
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
